package clusters

import (
	"bytes"
	"fmt"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/pkg/env"
	"github.com/stackrox/rox/pkg/images/utils"
	"github.com/stackrox/rox/pkg/logging"
	"github.com/stackrox/rox/pkg/templates"
	"github.com/stackrox/rox/pkg/urlfmt"
	"github.com/stackrox/rox/pkg/version"
	"github.com/stackrox/rox/pkg/zip"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	log = logging.LoggerForModule()
)

// Wrap adds additional functionality to a v1.Cluster.
type Wrap v1.Cluster

// NewDeployer takes in a cluster and returns the cluster implementation
func NewDeployer(c *v1.Cluster) (Deployer, error) {
	dep, ok := deployers[c.Type]
	if !ok {
		return nil, status.Errorf(codes.Unimplemented, "Cluster type %s is not currently implemented", c.Type.String())
	}
	return dep, nil
}

// Deployer is the interface that defines how to get the specific files per orchestrator
type Deployer interface {
	Render(Wrap) ([]*v1.File, error)
}

var deployers = make(map[v1.ClusterType]Deployer)

func executeTemplate(temp *template.Template, fields map[string]interface{}) ([]byte, error) {
	var b []byte
	buf := bytes.NewBuffer(b)
	err := temp.Execute(buf, fields)
	if err != nil {
		log.Errorf("template execution: %s", err)
		return nil, err
	}
	return buf.Bytes(), nil
}

func generateCollectorImage(mainImageName *v1.ImageName, tag string) *v1.ImageName {
	// Populate the tag
	collectorName := &v1.ImageName{
		Tag: tag,
	}
	// Populate Registry
	collectorName.Registry = mainImageName.GetRegistry()
	if mainImageName.GetRegistry() == "stackrox.io" {
		collectorName.Registry = "collector.stackrox.io"
	}
	// Populate Remote
	// This handles the case where there is no namespace. e.g. stackrox.io/collector:latest
	if slashIdx := strings.Index(mainImageName.GetRemote(), "/"); slashIdx == -1 {
		collectorName.Remote = "collector"
	} else {
		collectorName.Remote = mainImageName.GetRemote()[:slashIdx] + "/collector"
	}
	// Populate FullName
	collectorName.FullName = fmt.Sprintf("%s/%s:%s",
		collectorName.GetRegistry(), collectorName.GetRemote(), collectorName.GetTag())
	return collectorName
}

func fieldsFromWrap(c Wrap) (map[string]interface{}, error) {
	mainImageName := utils.GenerateImageFromString(c.MainImage).GetName()
	collectorName := generateCollectorImage(mainImageName, version.GetCollectorVersion())

	mainRegistry, err := urlfmt.FormatURL(mainImageName.GetRegistry(), urlfmt.HTTPS, urlfmt.NoTrailingSlash)
	if err != nil {
		return nil, err
	}
	collectorRegistry, err := urlfmt.FormatURL(collectorName.GetRegistry(), urlfmt.HTTPS, urlfmt.NoTrailingSlash)
	if err != nil {
		return nil, err
	}

	fields := map[string]interface{}{
		"ImageEnv":      env.Image.EnvVar(),
		"Image":         c.MainImage,
		"ImageRegistry": mainRegistry,
		"ImageTag":      mainImageName.GetTag(),

		"PublicEndpointEnv": env.CentralEndpoint.EnvVar(),
		"PublicEndpoint":    c.CentralApiEndpoint,

		"ClusterIDEnv": env.ClusterID.EnvVar(),
		"ClusterID":    c.Id,
		"ClusterName":  c.Name,

		"AdvertisedEndpointEnv": env.AdvertisedEndpoint.EnvVar(),
		"AdvertisedEndpoint":    env.AdvertisedEndpoint.Setting(),

		"RuntimeSupport":    c.RuntimeSupport,
		"CollectorRegistry": collectorRegistry,
		"CollectorImage":    collectorName.GetFullName(),
		"CollectorTag":      version.GetCollectorVersion(),

		"MonitoringEndpoint": c.MonitoringEndpoint,
		"ClusterType":        c.Type.String(),
	}
	return fields, nil
}

func renderFilenames(filenames []string, c map[string]interface{}, staticFilenames ...string) ([]*v1.File, error) {
	var files []*v1.File
	for _, f := range filenames {
		t, err := templates.ReadFileAndTemplate(f)
		if err != nil {
			return nil, err
		}
		d, err := executeTemplate(t, c)
		if err != nil {
			return nil, err
		}
		files = append(files, zip.NewFile(filepath.Base(f), d, path.Ext(f) == ".sh"))
	}
	for _, staticFilename := range staticFilenames {
		f, err := zip.NewFromFile(staticFilename, path.Base(staticFilename), path.Ext(staticFilename) == ".sh")
		if err != nil {
			return nil, err
		}
		files = append(files, f)
	}
	return files, nil
}
