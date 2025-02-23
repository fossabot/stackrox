// Package resources lists all resource types used by Central.
package resources

import (
	"sort"

	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/auth/permissions"
	"github.com/stackrox/rox/pkg/features"
)

// All resource types that we want to define (for the purposes of enforcing
// API permissions) must be defined here.
//
// Description for each type and the meaning of the respective Read and Write
// operations is available in
//     "ui/apps/platform/src/Containers/AccessControl/PermissionSets/ResourceDescription.tsx"
//
// UI defines possible values for resource type in
//     "ui/apps/platform/src/types/roleResources.ts"
//
// Each time you touch the list below, you likely need to update both
// aforementioned files.
//
// KEEP THE FOLLOWING LIST SORTED IN LEXICOGRAPHIC ORDER (case-sensitive).
var (
	// Access is the new resource grouping all access related resources.
	Access = newResourceMetadata("Access", permissions.GlobalScope)
	// Administration is the new resource grouping all administration-like resources.
	Administration = newResourceMetadata("Administration", permissions.GlobalScope)
	Alert          = newResourceMetadata("Alert", permissions.NamespaceScope)
	// SAC check is not performed directly on CVE resource. It exists here for postgres sac generation to pass.
	CVE        = newResourceMetadata("CVE", permissions.NamespaceScope)
	Cluster    = newResourceMetadata("Cluster", permissions.ClusterScope)
	Compliance = newResourceMetadata("Compliance", permissions.GlobalScope)
	Deployment = newResourceMetadata("Deployment", permissions.NamespaceScope)
	// DeploymentExtension is the new resource grouping all deployment extending resources.
	DeploymentExtension = newResourceMetadata("DeploymentExtension", permissions.NamespaceScope)
	Detection           = newResourceMetadata("Detection", permissions.GlobalScope)
	Image               = newResourceMetadata("Image", permissions.NamespaceScope)
	// SAC check is not performed directly on ImageComponent resource. It exists here for postgres sac generation to pass.
	ImageComponent = newResourceMetadata("ImageComponent", permissions.NamespaceScope)
	// Integration is the new  resource grouping all integration resources.
	Integration                      = newResourceMetadata("Integration", permissions.GlobalScope)
	K8sRole                          = newResourceMetadata("K8sRole", permissions.NamespaceScope)
	K8sRoleBinding                   = newResourceMetadata("K8sRoleBinding", permissions.NamespaceScope)
	K8sSubject                       = newResourceMetadata("K8sSubject", permissions.NamespaceScope)
	Namespace                        = newResourceMetadata("Namespace", permissions.NamespaceScope)
	NetworkGraph                     = newResourceMetadata("NetworkGraph", permissions.NamespaceScope)
	NetworkPolicy                    = newResourceMetadata("NetworkPolicy", permissions.NamespaceScope)
	Node                             = newResourceMetadata("Node", permissions.ClusterScope)
	Policy                           = newResourceMetadata("Policy", permissions.GlobalScope)
	Secret                           = newResourceMetadata("Secret", permissions.NamespaceScope)
	ServiceAccount                   = newResourceMetadata("ServiceAccount", permissions.NamespaceScope)
	VulnerabilityManagementApprovals = newResourceMetadata("VulnerabilityManagementApprovals",
		permissions.GlobalScope)
	VulnerabilityManagementRequests = newResourceMetadata("VulnerabilityManagementRequests",
		permissions.GlobalScope)
	VulnerabilityReports = newResourceMetadata("VulnerabilityReports", permissions.GlobalScope)
	WatchedImage         = newResourceMetadata("WatchedImage", permissions.GlobalScope)

	// Deprecated resources.

	// Deprecated: AllComments is deprecated, use Administration.
	AllComments = newDeprecatedResourceMetadata("AllComments", permissions.GlobalScope,
		Administration)
	// Deprecated: APIToken is deprecated, use Integration.
	APIToken = newDeprecatedResourceMetadata("APIToken", permissions.GlobalScope, Integration)
	// Deprecated: AuthPlugin is deprecated, use Access.
	AuthPlugin = newDeprecatedResourceMetadata("AuthPlugin", permissions.GlobalScope, Access)
	// Deprecated: AuthProvider is deprecated, use Access.
	AuthProvider = newDeprecatedResourceMetadata("AuthProvider", permissions.GlobalScope,
		Access)
	// Deprecated: BackupPlugins is deprecated, use Integration.
	BackupPlugins = newDeprecatedResourceMetadata("BackupPlugins", permissions.GlobalScope,
		Integration)
	// Deprecated: ComplianceRuns is deprecated, use Compliance.
	ComplianceRuns = newDeprecatedResourceMetadata("ComplianceRuns", permissions.ClusterScope,
		Compliance)
	// Deprecated: ComplianceRunSchedule is deprecated, use Administration.
	ComplianceRunSchedule = newDeprecatedResourceMetadata("ComplianceRunSchedule",
		permissions.GlobalScope, Administration)
	// Deprecated: Config is deprecated, use Administration.
	Config = newDeprecatedResourceMetadata("Config", permissions.GlobalScope,
		Administration)
	// Deprecated: DebugLogs is deprecated, use Administration.
	DebugLogs = newDeprecatedResourceMetadata("DebugLogs", permissions.GlobalScope,
		Administration)
	// Deprecated: Group is deprecated, use Access.
	Group = newDeprecatedResourceMetadata("Group", permissions.GlobalScope, Access)
	// Deprecated: ImageIntegration is deprecated, use Integration.
	ImageIntegration = newDeprecatedResourceMetadata("ImageIntegration",
		permissions.GlobalScope, Integration)
	// Deprecated: Indicator is deprecated, use DeploymentExtension.
	Indicator = newDeprecatedResourceMetadata("Indicator", permissions.NamespaceScope,
		DeploymentExtension)
	// Deprecated: Licenses is deprecated, use Access.
	Licenses = newDeprecatedResourceMetadata("Licenses", permissions.GlobalScope, Access)
	// Deprecated: NetworkBaseline is deprecated, use DeploymentExtension.
	NetworkBaseline = newDeprecatedResourceMetadata("NetworkBaseline",
		permissions.NamespaceScope, DeploymentExtension)
	// Deprecated: NetworkGraphConfig is deprecated, use Administration.
	NetworkGraphConfig = newDeprecatedResourceMetadata("NetworkGraphConfig",
		permissions.GlobalScope, Administration)
	// Deprecated: Notifier is deprecated, use Integration.
	Notifier = newDeprecatedResourceMetadata("Notifier", permissions.GlobalScope,
		Integration)
	// Deprecated: ProbeUpload is deprecated, use Administration.
	ProbeUpload = newDeprecatedResourceMetadata("ProbeUpload", permissions.GlobalScope,
		Administration)
	// Deprecated: ProcessWhitelist is deprecated, use DeploymentExtension.
	ProcessWhitelist = newDeprecatedResourceMetadata("ProcessWhitelist",
		permissions.NamespaceScope, DeploymentExtension)
	// Deprecated: Risk is deprecated, use DeploymentExtension.
	Risk = newDeprecatedResourceMetadata("Risk", permissions.NamespaceScope,
		DeploymentExtension)
	// Deprecated: Role is deprecated, use Access.
	Role = newDeprecatedResourceMetadata("Role", permissions.GlobalScope, Access)
	// Deprecated: ScannerBundle is deprecated, use Administration.
	ScannerBundle = newDeprecatedResourceMetadata("ScannerBundle",
		permissions.GlobalScope, Administration)
	// Deprecated: ScannerDefinitions is deprecated, use Administration.
	ScannerDefinitions = newDeprecatedResourceMetadata("ScannerDefinitions",
		permissions.GlobalScope, Administration)
	// Deprecated: SensorUpgradeConfig is deprecated, use Administration.
	SensorUpgradeConfig = newDeprecatedResourceMetadata("SensorUpgradeConfig",
		permissions.GlobalScope, Administration)
	// Deprecated: ServiceIdentity is deprecated, use Administration.
	ServiceIdentity = newDeprecatedResourceMetadata("ServiceIdentity",
		permissions.GlobalScope, Administration)
	// Deprecated: SignatureIntegration is deprecated, use Integration.
	SignatureIntegration = newDeprecatedResourceMetadataWithFeatureFlag("SignatureIntegration",
		permissions.GlobalScope, Integration, features.ImageSignatureVerification)
	// Deprecated: User is deprecated, use Access.
	User = newDeprecatedResourceMetadata("User", permissions.GlobalScope, Access)

	// Internal Resources.
	ComplianceOperator = newInternalResourceMetadata("ComplianceOperator", permissions.GlobalScope)

	resourceToMetadata         = make(map[permissions.Resource]permissions.ResourceMetadata)
	disabledResourceToMetadata = make(map[permissions.Resource]permissions.ResourceMetadata)
)

func newResourceMetadata(name permissions.Resource, scope permissions.ResourceScope) permissions.ResourceMetadata {
	md := permissions.ResourceMetadata{
		Resource: name,
		Scope:    scope,
	}
	resourceToMetadata[name] = md
	return md
}

func newDeprecatedResourceMetadata(name permissions.Resource, scope permissions.ResourceScope,
	replacingResourceMD permissions.ResourceMetadata) permissions.ResourceMetadata {
	md := permissions.ResourceMetadata{
		Resource:          name,
		Scope:             scope,
		ReplacingResource: &replacingResourceMD,
	}
	resourceToMetadata[name] = md
	return md
}

/*
Commented for now, uncomment in case you need to register a resource guarded behind a feature flag.
func newResourceMetadataWithFeatureFlag(name permissions.Resource, scope permissions.ResourceScope,
	flag features.FeatureFlag) permissions.ResourceMetadata {
	md := permissions.ResourceMetadata{
		Resource:          name,
		Scope:             scope,
	}
	if flag.Enabled() {
		resourceToMetadata[name] = md
	} else {
		disabledResourceToMetadata[name] = md
	}
	return md
}
*/

func newDeprecatedResourceMetadataWithFeatureFlag(name permissions.Resource, scope permissions.ResourceScope,
	replacingResourceMD permissions.ResourceMetadata, flag features.FeatureFlag) permissions.ResourceMetadata {
	md := permissions.ResourceMetadata{
		Resource:          name,
		Scope:             scope,
		ReplacingResource: &replacingResourceMD,
	}
	if flag.Enabled() {
		resourceToMetadata[name] = md
	} else {
		disabledResourceToMetadata[name] = md
	}
	return md
}

func newInternalResourceMetadata(name permissions.Resource,
	scope permissions.ResourceScope) permissions.ResourceMetadata {
	return permissions.ResourceMetadata{
		Resource: name,
		Scope:    scope,
	}
}

// ListAll returns a list of all resources.
func ListAll() []permissions.Resource {
	resources := make([]permissions.Resource, 0, len(resourceToMetadata))
	for _, metadata := range ListAllMetadata() {
		resources = append(resources, metadata.Resource)
	}
	return resources
}

// ListAllMetadata returns a list of all resource metadata.
func ListAllMetadata() []permissions.ResourceMetadata {
	metadatas := make([]permissions.ResourceMetadata, 0, len(resourceToMetadata))
	for _, metadata := range resourceToMetadata {
		metadatas = append(metadatas, metadata)
	}
	sort.SliceStable(metadatas, func(i, j int) bool {
		return string(metadatas[i].Resource) < string(metadatas[j].Resource)
	})
	return metadatas
}

// ListAllDisabledMetadata returns a list of all resource metadata that are currently disable by feature flag.
func ListAllDisabledMetadata() []permissions.ResourceMetadata {
	metadatas := make([]permissions.ResourceMetadata, 0, len(disabledResourceToMetadata))
	for _, metadata := range disabledResourceToMetadata {
		metadatas = append(metadatas, metadata)
	}
	sort.SliceStable(metadatas, func(i, j int) bool {
		return string(metadatas[i].Resource) < string(metadatas[j].Resource)
	})
	return metadatas
}

// AllResourcesViewPermissions returns a slice containing view permissions for all resource types.
func AllResourcesViewPermissions() []permissions.ResourceWithAccess {
	metadatas := ListAllMetadata()
	result := make([]permissions.ResourceWithAccess, len(metadatas))
	for i, metadata := range metadatas {
		result[i] = permissions.ResourceWithAccess{
			// We want to ensure access to *all* resources, so when using SAC, always perform legacy auth (= enforcement
			// at the global scope) even for cluster- or namespace-scoped resources.
			Resource: permissions.WithLegacyAuthForSAC(metadata, true),
			Access:   storage.Access_READ_ACCESS,
		}
	}
	return result
}

// AllResourcesModifyPermissions returns a slice containing write permissions for all resource types.
func AllResourcesModifyPermissions() []permissions.ResourceWithAccess {
	metadatas := ListAllMetadata()
	result := make([]permissions.ResourceWithAccess, len(metadatas))
	for i, metadata := range metadatas {
		result[i] = permissions.ResourceWithAccess{
			// We want to ensure access to *all* resources, so when using SAC, always perform legacy auth (= enforcement
			// at the global scope) even for cluster- or namespace-scoped resources.
			Resource: permissions.WithLegacyAuthForSAC(metadata, true),
			Access:   storage.Access_READ_WRITE_ACCESS,
		}
	}
	return result
}

// MetadataForResource returns the metadata for the given resource. If the resource is unknown, metadata for this
// resource with global scope is returned.
func MetadataForResource(res permissions.Resource) (permissions.ResourceMetadata, bool) {
	md, found := resourceToMetadata[res]
	if !found {
		md.Resource = res
		md.Scope = permissions.GlobalScope
	}
	return md, found
}
