package translation

import (
	"context"
	"fmt"

	// Required for the usage of go:embed below.
	_ "embed"

	"github.com/pkg/errors"
	central "github.com/stackrox/rox/operator/api/central/v1alpha1"
	"github.com/stackrox/rox/operator/pkg/values/translation"
	"github.com/stackrox/rox/pkg/helmutil"
	"github.com/stackrox/rox/pkg/utils"
	"helm.sh/helm/v3/pkg/chartutil"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
)

var (
	//go:embed base-values.yaml
	baseValuesYAML []byte
)

// Translator translates and enriches helm values
type Translator struct {
	Client kubernetes.Interface
}

// Translate translates and enriches helm values
func (t Translator) Translate(ctx context.Context, u *unstructured.Unstructured) (chartutil.Values, error) {
	baseValues, err := chartutil.ReadValues(baseValuesYAML)
	utils.CrashOnError(err) // ensured through unit test that this doesn't happen.

	c := central.Central{}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, &c)
	if err != nil {
		return nil, err
	}

	valsFromCR, err := translate(ctx, t.Client, c)
	if err != nil {
		return nil, err
	}

	imageOverrideVals, err := imageOverrides.ToValues()
	if err != nil {
		return nil, errors.Wrap(err, "computing image override values")
	}
	return helmutil.CoalesceTables(baseValues, imageOverrideVals, valsFromCR), nil
}

// translate translates a Central CR into helm values.
func translate(ctx context.Context, clientSet kubernetes.Interface, c central.Central) (chartutil.Values, error) {
	v := translation.NewValuesBuilder()

	v.AddAllFrom(translation.GetImagePullSecrets(c.Spec.ImagePullSecrets))
	v.AddAllFrom(getEnv(ctx, clientSet, c.Namespace, c.Spec.Egress))
	v.AddAllFrom(translation.GetTLSValues(c.Spec.TLS))

	customize := translation.NewValuesBuilder()
	customize.AddAllFrom(translation.GetCustomize(c.Spec.Customize))

	if c.Spec.Central != nil {
		v.AddChild("central", getCentralComponentValues(ctx, clientSet, c.Namespace, c.Spec.Central))
	}

	if c.Spec.Scanner != nil {
		v.AddChild("scanner", getScannerComponentValues(c.Spec.Scanner))
	}

	v.AddChild("customize", &customize)

	return v.Build()
}

func getEnv(ctx context.Context, clientSet kubernetes.Interface, namespace string, egress *central.Egress) *translation.ValuesBuilder {
	env := translation.NewValuesBuilder()
	if egress != nil {
		if egress.ConnectivityPolicy != nil {
			switch *egress.ConnectivityPolicy {
			case central.ConnectivityOnline:
				env.SetBoolValue("offlineMode", false)
			case central.ConnectivityOffline:
				env.SetBoolValue("offlineMode", true)
			default:
				return env.SetError(fmt.Errorf("invalid spec.egress.connectivityPolicy %q", *egress.ConnectivityPolicy))
			}
		}
		env.AddAllFrom(translation.NewBuilderFromSecret(ctx, clientSet, namespace, egress.ProxyConfigSecret, map[string]string{"config.yaml": "proxyConfig"}, "spec.egress.proxyConfigSecret"))
	}
	ret := translation.NewValuesBuilder()
	ret.AddChild("env", &env)
	return &ret
}

func getCentralComponentValues(ctx context.Context, clientSet kubernetes.Interface, namespace string, c *central.CentralComponentSpec) *translation.ValuesBuilder {
	cv := translation.NewValuesBuilder()

	cv.AddChild(translation.ResourcesKey, translation.GetResources(c.Resources))

	if c.DefaultTLSSecret != nil {
		cv.SetMap("defaultTLS", map[string]interface{}{"reference": c.DefaultTLSSecret.Name})
	}

	cv.SetStringMap("nodeSelector", c.NodeSelector)

	// TODO(ROX-7147): design CentralEndpointSpec, see central_types.go

	if c.Persistence != nil {
		persistence := translation.NewValuesBuilder()
		persistence.SetString("hostPath", c.Persistence.HostPath)
		if c.Persistence.PersistentVolumeClaim != nil {
			pvc := translation.NewValuesBuilder()
			pvc.SetString("claimName", c.Persistence.PersistentVolumeClaim.ClaimName)
			if c.Persistence.PersistentVolumeClaim.CreateClaim != nil {
				switch *c.Persistence.PersistentVolumeClaim.CreateClaim {
				case central.ClaimCreate:
					pvc.SetBoolValue("createClaim", true)
				case central.ClaimReuse:
					pvc.SetBoolValue("createClaim", false)
				default:
					return cv.SetError(fmt.Errorf("invalid spec.central.persistence.persistentVolumeClaim.createClaim %q", *c.Persistence.PersistentVolumeClaim.CreateClaim))
				}
			}
			// TODO(ROX-7149): more details TBD, values files are inconsistent and require more investigation and template reading
			persistence.AddChild("persistentVolumeClaim", &pvc)
		}
		cv.AddChild("persistence", &persistence)
	}

	if c.Exposure != nil {
		exposure := translation.NewValuesBuilder()
		if c.Exposure.LoadBalancer != nil {
			lb := translation.NewValuesBuilder()
			lb.SetBool("enabled", c.Exposure.LoadBalancer.Enabled)
			lb.SetInt32("port", c.Exposure.LoadBalancer.Port)
			lb.SetString("ip", c.Exposure.LoadBalancer.IP)
			exposure.AddChild("loadBalancer", &lb)
		}
		if c.Exposure.NodePort != nil {
			np := translation.NewValuesBuilder()
			np.SetBool("enabled", c.Exposure.NodePort.Enabled)
			np.SetInt32("port", c.Exposure.NodePort.Port)
			exposure.AddChild("nodePort", &np)
		}
		if c.Exposure.Route != nil {
			route := translation.NewValuesBuilder()
			route.SetBool("enabled", c.Exposure.Route.Enabled)
			exposure.AddChild("route", &route)
		}
		cv.AddChild("exposure", &exposure)
	}
	return &cv
}

func getScannerComponentValues(s *central.ScannerComponentSpec) *translation.ValuesBuilder {
	sv := translation.NewValuesBuilder()

	if s.ScannerComponent != nil {
		switch *s.ScannerComponent {
		case central.ScannerComponentDisabled:
			sv.SetBoolValue("disable", true)
		case central.ScannerComponentEnabled:
			sv.SetBoolValue("disable", false)
		default:
			return sv.SetError(fmt.Errorf("invalid spec.scanner.scannerComponent %q", *s.ScannerComponent))
		}
	}

	if s.GetAnalyzer().GetScaling() != nil {
		scaling := s.GetAnalyzer().GetScaling()
		sv.SetInt32("replicas", scaling.Replicas)

		autoscaling := translation.NewValuesBuilder()
		if scaling.AutoScaling != nil {
			switch *scaling.AutoScaling {
			case central.ScannerAutoScalingDisabled:
				autoscaling.SetBoolValue("disable", true)
			case central.ScannerAutoScalingEnabled:
				autoscaling.SetBoolValue("disable", false)
			default:
				return autoscaling.SetError(fmt.Errorf("invalid spec.scanner.replicas.autoScaling %q", *scaling.AutoScaling))
			}
		}
		autoscaling.SetInt32("minReplicas", scaling.MinReplicas)
		autoscaling.SetInt32("maxReplicas", scaling.MaxReplicas)
		sv.AddChild("autoscaling", &autoscaling)
	}

	if s.GetAnalyzer() != nil {
		sv.SetStringMap("nodeSelector", s.GetAnalyzer().NodeSelector)
		sv.AddChild(translation.ResourcesKey, translation.GetResources(s.GetAnalyzer().Resources))
	}

	if s.ScannerDB != nil {
		sv.SetStringMap("dbNodeSelector", s.ScannerDB.NodeSelector)
		sv.AddChild("dbResources", translation.GetResources(s.ScannerDB.Resources))
	}

	return &sv
}
