package clair

import (
	"encoding/json"
	"time"

	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/cvss/cvssv2"
	"github.com/stackrox/rox/pkg/cvss/cvssv3"
	"github.com/stackrox/rox/pkg/logging"
	"github.com/stackrox/rox/pkg/protoconv"
	"github.com/stackrox/rox/pkg/scans"
	clairV1 "github.com/stackrox/scanner/api/v1"
	clientMetadata "github.com/stackrox/scanner/pkg/clairify/client/metadata"
	"github.com/stackrox/scanner/pkg/component"
)

const (
	timeFormat = "2006-01-02T15:04Z"
)

var (
	log = logging.LoggerForModule()

	versionFormatsToSource = map[string]storage.SourceType{
		component.GemSourceType.String():               storage.SourceType_RUBY,
		component.JavaSourceType.String():              storage.SourceType_JAVA,
		component.NPMSourceType.String():               storage.SourceType_NODEJS,
		component.PythonSourceType.String():            storage.SourceType_PYTHON,
		component.DotNetCoreRuntimeSourceType.String(): storage.SourceType_DOTNETCORERUNTIME,
	}
)

type metadata struct {
	PublishedOn  string `json:"PublishedDateTime"`
	LastModified string `json:"LastModifiedDateTime"`
	CvssV2       *cvss  `json:"CVSSv2"`
	CvssV3       *cvss  `json:"CVSSv3"`
}

type cvss struct {
	Score               float32
	Vectors             string
	ExploitabilityScore float32
	ImpactScore         float32
}

// ConvertVulnerability converts a clair vulnerability to a proto vulnerability
func ConvertVulnerability(v clairV1.Vulnerability) *storage.EmbeddedVulnerability {
	var vulnMetadataMap interface{}
	var link string
	if metadata, ok := v.Metadata[clientMetadata.NVD]; ok {
		vulnMetadataMap = metadata
		link = scans.GetVulnLink(v.Name)
	} else if metadata, ok := v.Metadata[clientMetadata.RedHat]; ok {
		vulnMetadataMap = metadata
		link = scans.GetRedHatVulnLink(v.Name)
	} else {
		return nil
	}

	if v.Link == "" {
		v.Link = link
	}
	vul := &storage.EmbeddedVulnerability{
		Cve:     v.Name,
		Summary: v.Description,
		Link:    v.Link,
		SetFixedBy: &storage.EmbeddedVulnerability_FixedBy{
			FixedBy: v.FixedBy,
		},
		VulnerabilityType: storage.EmbeddedVulnerability_IMAGE_VULNERABILITY,
	}

	d, err := json.Marshal(vulnMetadataMap)
	if err != nil {
		return vul
	}
	var m metadata
	if err := json.Unmarshal(d, &m); err != nil {
		return vul
	}
	if m.PublishedOn != "" {
		if ts, err := time.Parse(timeFormat, m.PublishedOn); err == nil {
			vul.PublishedOn = protoconv.ConvertTimeToTimestamp(ts)
		}
	}
	if m.LastModified != "" {
		if ts, err := time.Parse(timeFormat, m.LastModified); err == nil {
			vul.LastModified = protoconv.ConvertTimeToTimestamp(ts)
		}
	}

	if m.CvssV2 != nil && m.CvssV2.Vectors != "" {
		if cvssV2, err := cvssv2.ParseCVSSV2(m.CvssV2.Vectors); err == nil {
			cvssV2.ExploitabilityScore = m.CvssV2.ExploitabilityScore
			cvssV2.ImpactScore = m.CvssV2.ImpactScore
			cvssV2.Score = m.CvssV2.Score

			vul.CvssV2 = cvssV2
			// This sets the top level score for use in policies. It will be overwritten if v3 exists
			vul.Cvss = m.CvssV2.Score
			vul.ScoreVersion = storage.EmbeddedVulnerability_V2
			vul.GetCvssV2().Severity = cvssv2.Severity(vul.GetCvss())
		} else {
			log.Error(err)
		}
	}

	if m.CvssV3 != nil && m.CvssV3.Vectors != "" {
		if cvssV3, err := cvssv3.ParseCVSSV3(m.CvssV3.Vectors); err == nil {
			cvssV3.ExploitabilityScore = m.CvssV3.ExploitabilityScore
			cvssV3.ImpactScore = m.CvssV3.ImpactScore
			cvssV3.Score = m.CvssV3.Score

			vul.CvssV3 = cvssV3
			vul.Cvss = m.CvssV3.Score
			vul.ScoreVersion = storage.EmbeddedVulnerability_V3
			vul.GetCvssV3().Severity = cvssv3.Severity(vul.GetCvss())
		} else {
			log.Error(err)
		}
	}
	return vul
}

func convertFeature(feature clairV1.Feature) *storage.EmbeddedImageScanComponent {
	component := &storage.EmbeddedImageScanComponent{
		Name:     feature.Name,
		Version:  feature.Version,
		Location: feature.Location,
	}
	if source, ok := versionFormatsToSource[feature.VersionFormat]; ok {
		component.Source = source
	}
	component.Vulns = make([]*storage.EmbeddedVulnerability, 0, len(feature.Vulnerabilities))
	for _, v := range feature.Vulnerabilities {
		if convertedVuln := ConvertVulnerability(v); convertedVuln != nil {
			component.Vulns = append(component.Vulns, convertedVuln)
		}
	}
	return component
}

func buildSHAToIndexMap(image *storage.Image) map[string]int32 {
	layerSHAToIndex := make(map[string]int32)

	if image.GetMetadata().GetV2() != nil {
		var layerIdx int
		for i, l := range image.GetMetadata().GetV1().GetLayers() {
			if !l.Empty {
				if layerIdx >= len(image.Metadata.LayerShas) {
					log.Error("More layers than expected when correlating V2 instructions to V1 layers")
					break
				}
				sha := image.GetMetadata().LayerShas[layerIdx]
				layerSHAToIndex[sha] = int32(i)
				layerIdx++
			}
		}
	} else {
		// If it's V1 then we should have a 1:1 mapping of layer SHAs to the layerOrdering slice
		for i := range image.GetMetadata().GetV1().GetLayers() {
			if i >= len(image.Metadata.LayerShas) {
				log.Error("More layers than expected when correlating V1 instructions to V1 layers")
				break
			}
			layerSHAToIndex[image.Metadata.LayerShas[i]] = int32(i)
		}
	}
	return layerSHAToIndex
}

// ConvertFeatures converts clair features to proto components
func ConvertFeatures(image *storage.Image, features []clairV1.Feature) (components []*storage.EmbeddedImageScanComponent) {
	layerSHAToIndex := buildSHAToIndexMap(image)

	components = make([]*storage.EmbeddedImageScanComponent, 0, len(features))
	for _, feature := range features {
		convertedComponent := convertFeature(feature)
		if val, ok := layerSHAToIndex[feature.AddedBy]; ok {
			convertedComponent.HasLayerIndex = &storage.EmbeddedImageScanComponent_LayerIndex{
				LayerIndex: val,
			}
		}
		components = append(components, convertedComponent)
	}
	return
}
