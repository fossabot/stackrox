package resolvers

import (
	"context"
	"net/url"
	"strings"
	"time"

	"github.com/graph-gophers/graphql-go"
	"github.com/pkg/errors"
	"github.com/stackrox/rox/central/graphql/resolvers/loaders"
	"github.com/stackrox/rox/central/metrics"
	"github.com/stackrox/rox/central/policy/matcher"
	riskDS "github.com/stackrox/rox/central/risk/datastore"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/features"
	"github.com/stackrox/rox/pkg/k8srbac"
	pkgMetrics "github.com/stackrox/rox/pkg/metrics"
	"github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/search/scoped"
	"github.com/stackrox/rox/pkg/set"
	"github.com/stackrox/rox/pkg/utils"
)

func init() {
	schema := getBuilder()
	utils.Must(
		schema.AddType("PolicyStatus", []string{"status: String!", "failingPolicies: [Policy!]!"}),
		schema.AddQuery("clusters(query: String, pagination: Pagination): [Cluster!]!"),
		schema.AddQuery("clusterCount(query: String): Int!"),
		schema.AddQuery("cluster(id: ID!): Cluster"),
		schema.AddExtraResolver("Cluster", `alerts(query: String, pagination: Pagination): [Alert!]!`),
		schema.AddExtraResolver("Cluster", `alertCount(query: String): Int!`),
		schema.AddExtraResolver("Cluster", `latestViolation(query: String): Time`),
		schema.AddExtraResolver("Cluster", `failingPolicyCounter(query: String): PolicyCounter`),
		schema.AddExtraResolver("Cluster", `deployments(query: String, pagination: Pagination): [Deployment!]!`),
		schema.AddExtraResolver("Cluster", `deploymentCount(query: String): Int!`),
		schema.AddExtraResolver("Cluster", `nodes(query: String, pagination: Pagination): [Node!]!`),
		schema.AddExtraResolver("Cluster", `nodeCount(query: String): Int!`),
		schema.AddExtraResolver("Cluster", `node(node: ID!): Node`),
		schema.AddExtraResolver("Cluster", `namespaces(query: String, pagination: Pagination): [Namespace!]!`),
		schema.AddExtraResolver("Cluster", `namespace(name: String!): Namespace`),
		schema.AddExtraResolver("Cluster", `namespaceCount(query: String): Int!`),
		schema.AddExtraResolver("Cluster", "complianceResults(query: String): [ControlResult!]!"),
		schema.AddExtraResolver("Cluster", `k8sRoles(query: String, pagination: Pagination): [K8SRole!]!`),
		schema.AddExtraResolver("Cluster", `k8sRole(role: ID!): K8SRole`),
		schema.AddExtraResolver("Cluster", `k8sRoleCount(query: String): Int!`),
		schema.AddExtraResolver("Cluster", `serviceAccounts(query: String, pagination: Pagination): [ServiceAccount!]!`),
		schema.AddExtraResolver("Cluster", `serviceAccount(sa: ID!): ServiceAccount`),
		schema.AddExtraResolver("Cluster", `serviceAccountCount(query: String): Int!`),
		schema.AddExtraResolver("Cluster", `subjects(query: String, pagination: Pagination): [Subject!]!`),
		schema.AddExtraResolver("Cluster", `subject(name: String!): Subject`),
		schema.AddExtraResolver("Cluster", `subjectCount(query: String): Int!`),
		schema.AddExtraResolver("Cluster", `images(query: String, pagination: Pagination): [Image!]!`),
		schema.AddExtraResolver("Cluster", `imageCount(query: String): Int!`),
		schema.AddExtraResolver("Cluster", `components(query: String, pagination: Pagination): [EmbeddedImageScanComponent!]!`),
		schema.AddExtraResolver("Cluster", `componentCount(query: String): Int!`),
		schema.AddExtraResolver("Cluster", `vulns(query: String, scopeQuery: String, pagination: Pagination): [EmbeddedVulnerability!]!`),
		schema.AddExtraResolver("Cluster", `vulnCount(query: String): Int!`),
		schema.AddExtraResolver("Cluster", `vulnCounter(query: String): VulnerabilityCounter!`),
		schema.AddExtraResolver("Cluster", `nodeVulnerabilities(query: String, scopeQuery: String, pagination: Pagination): [NodeVulnerability!]!`),
		schema.AddExtraResolver("Cluster", `nodeVulnerabilityCount(query: String): Int!`),
		schema.AddExtraResolver("Cluster", `nodeVulnerabilityCounter(query: String): VulnerabilityCounter!`),
		schema.AddExtraResolver("Cluster", `k8sVulns(query: String, pagination: Pagination): [EmbeddedVulnerability!]!`),
		schema.AddExtraResolver("Cluster", `k8sVulnCount(query: String): Int!`),
		schema.AddExtraResolver("Cluster", `istioVulns(query: String, pagination: Pagination): [EmbeddedVulnerability!]!`),
		schema.AddExtraResolver("Cluster", `istioVulnCount(query: String): Int!`),
		schema.AddExtraResolver("Cluster", `openShiftVulns(query: String, pagination: Pagination): [EmbeddedVulnerability!]!`),
		schema.AddExtraResolver("Cluster", `openShiftVulnCount(query: String): Int!`),
		schema.AddExtraResolver("Cluster", `policies(query: String, pagination: Pagination): [Policy!]!`),
		schema.AddExtraResolver("Cluster", `policyCount(query: String): Int!`),
		schema.AddExtraResolver("Cluster", `policyStatus(query: String): PolicyStatus!`),
		schema.AddExtraResolver("Cluster", `secrets(query: String, pagination: Pagination): [Secret!]!`),
		schema.AddExtraResolver("Cluster", `secretCount(query: String): Int!`),
		schema.AddExtraResolver("Cluster", `controlStatus(query: String): String!`),
		schema.AddExtraResolver("Cluster", `controls(query: String): [ComplianceControl!]!`),
		schema.AddExtraResolver("Cluster", `failingControls(query: String): [ComplianceControl!]!`),
		schema.AddExtraResolver("Cluster", `passingControls(query: String): [ComplianceControl!]!`),
		schema.AddExtraResolver("Cluster", `complianceControlCount(query: String): ComplianceControlCount!`),
		schema.AddExtraResolver("Cluster", `risk: Risk`),
		schema.AddExtraResolver("Cluster", `isGKECluster: Boolean!`),
		schema.AddExtraResolver("Cluster", `isOpenShiftCluster: Boolean!`),
		schema.AddExtraResolver("Cluster", `unusedVarSink(query: String): Int`),
		schema.AddExtraResolver("Cluster", `istioEnabled: Boolean!`),
		schema.AddExtraResolver("Cluster", "plottedVulns(query: String): PlottedVulnerabilities!"),
		schema.AddExtraResolver("OrchestratorMetadata", `openshiftVersion: String!`),
	)
}

func (resolver *clusterResolver) getClusterRawQuery() string {
	return search.NewQueryBuilder().AddExactMatches(search.ClusterID, resolver.data.GetId()).Query()
}

func (resolver *clusterResolver) getClusterQuery() *v1.Query {
	return search.NewQueryBuilder().AddExactMatches(search.ClusterID, resolver.data.GetId()).ProtoQuery()
}

func (resolver *clusterResolver) getClusterConjunctionQuery(q *v1.Query) (*v1.Query, error) {
	pagination := q.GetPagination()
	q.Pagination = nil

	q, err := search.AddAsConjunction(resolver.getClusterQuery(), q)
	if err != nil {
		return nil, err
	}

	q.Pagination = pagination
	return q, nil
}

// Cluster returns a GraphQL resolver for the given cluster
func (resolver *Resolver) Cluster(ctx context.Context, args struct{ graphql.ID }) (*clusterResolver, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Root, "Cluster")
	if err := readClusters(ctx); err != nil {
		return nil, err
	}
	return resolver.wrapCluster(resolver.ClusterDataStore.GetCluster(ctx, string(args.ID)))
}

// Clusters returns GraphQL resolvers for all clusters
func (resolver *Resolver) Clusters(ctx context.Context, args PaginatedQuery) ([]*clusterResolver, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Root, "Clusters")
	if err := readClusters(ctx); err != nil {
		return nil, err
	}
	query, err := args.AsV1QueryOrEmpty()
	if err != nil {
		return nil, err
	}

	return resolver.wrapClusters(resolver.ClusterDataStore.SearchRawClusters(ctx, query))
}

// ClusterCount returns count of all clusters across infrastructure
func (resolver *Resolver) ClusterCount(ctx context.Context, args RawQuery) (int32, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Root, "ClusterCount")
	if err := readClusters(ctx); err != nil {
		return 0, err
	}
	q, err := args.AsV1QueryOrEmpty()
	if err != nil {
		return 0, err
	}
	results, err := resolver.ClusterDataStore.Search(ctx, q)
	if err != nil {
		return 0, err
	}
	return int32(len(results)), nil
}

// Alerts returns GraphQL resolvers for all alerts on this cluster
func (resolver *clusterResolver) Alerts(ctx context.Context, args PaginatedQuery) ([]*alertResolver, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Cluster, "Alerts")

	if err := readAlerts(ctx); err != nil {
		return nil, err // could return nil, nil to prevent errors from propagating.
	}

	query := search.AddRawQueriesAsConjunction(args.String(), resolver.getClusterRawQuery())

	return resolver.root.Violations(ctx, PaginatedQuery{Query: &query, Pagination: args.Pagination})
}

func (resolver *clusterResolver) AlertCount(ctx context.Context, args RawQuery) (int32, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Cluster, "AlertCount")
	if err := readAlerts(ctx); err != nil {
		return 0, err
	}

	query := search.AddRawQueriesAsConjunction(args.String(), resolver.getClusterRawQuery())

	return resolver.root.ViolationCount(ctx, RawQuery{Query: &query})
}

// FailingPolicyCounter returns a policy counter for all the failed policies.
func (resolver *clusterResolver) FailingPolicyCounter(ctx context.Context, args RawQuery) (*PolicyCounterResolver, error) {
	if err := readAlerts(ctx); err != nil {
		return nil, err
	}

	q, err := args.AsV1QueryOrEmpty()
	if err != nil {
		return nil, err
	}

	q, err = search.AddAsConjunction(q, resolver.getClusterQuery())
	if err != nil {
		return nil, err
	}

	alerts, err := resolver.root.ViolationsDataStore.SearchListAlerts(ctx, q)
	if err != nil {
		return nil, nil
	}
	return mapListAlertsToPolicySeverityCount(alerts), nil
}

// Deployments returns GraphQL resolvers for all deployments in this cluster
func (resolver *clusterResolver) Deployments(ctx context.Context, args PaginatedQuery) ([]*deploymentResolver, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Cluster, "Deployments")

	if err := readDeployments(ctx); err != nil {
		return nil, err
	}

	query := search.AddRawQueriesAsConjunction(args.String(), resolver.getClusterRawQuery())

	return resolver.root.Deployments(ctx, PaginatedQuery{Query: &query, Pagination: args.Pagination})
}

// DeploymentCount returns count of all deployments in this cluster
func (resolver *clusterResolver) DeploymentCount(ctx context.Context, args RawQuery) (int32, error) {
	if err := readDeployments(ctx); err != nil {
		return 0, err
	}

	query := search.AddRawQueriesAsConjunction(args.String(), resolver.getClusterRawQuery())

	return resolver.root.DeploymentCount(ctx, RawQuery{Query: &query})
}

// Nodes returns all nodes on the cluster
func (resolver *clusterResolver) Nodes(ctx context.Context, args PaginatedQuery) ([]*nodeResolver, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Cluster, "Nodes")

	if err := readNodes(ctx); err != nil {
		return nil, err
	}

	query := search.AddRawQueriesAsConjunction(args.String(), resolver.getClusterRawQuery())

	return resolver.root.Nodes(ctx, PaginatedQuery{Query: &query, Pagination: args.Pagination})
}

// NodeCount returns count of all nodes on the cluster
func (resolver *clusterResolver) NodeCount(ctx context.Context, args RawQuery) (int32, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Cluster, "NodeCount")

	if err := readNodes(ctx); err != nil {
		return 0, err
	}

	query := search.AddRawQueriesAsConjunction(args.String(), resolver.getClusterRawQuery())

	return resolver.root.NodeCount(ctx, RawQuery{Query: &query})
}

// Node returns a given node on a cluster
func (resolver *clusterResolver) Node(ctx context.Context, args struct{ Node graphql.ID }) (*nodeResolver, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Cluster, "Node")

	if err := readNodes(ctx); err != nil {
		return nil, err
	}
	store, err := resolver.root.NodeGlobalDataStore.GetClusterNodeStore(ctx, resolver.data.GetId(), false)
	if err != nil {
		return nil, err
	}
	node, err := store.GetNode(string(args.Node))
	return resolver.root.wrapNode(node, node != nil, err)
}

// Namespaces returns the namespaces in a cluster.
func (resolver *clusterResolver) Namespaces(ctx context.Context, args PaginatedQuery) ([]*namespaceResolver, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Cluster, "Namespaces")

	if err := readNamespaces(ctx); err != nil {
		return nil, err
	}

	query := search.AddRawQueriesAsConjunction(args.String(), resolver.getClusterRawQuery())

	return resolver.root.Namespaces(ctx, PaginatedQuery{Query: &query, Pagination: args.Pagination})
}

// Namespace returns a given namespace in a cluster.
func (resolver *clusterResolver) Namespace(ctx context.Context, args struct{ Name string }) (*namespaceResolver, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Cluster, "Namespace")

	return resolver.root.NamespaceByClusterIDAndName(ctx, clusterIDAndNameQuery{
		ClusterID: graphql.ID(resolver.data.GetId()),
		Name:      args.Name,
	})
}

// NamespaceCount returns counts of namespaces on a cluster.
func (resolver *clusterResolver) NamespaceCount(ctx context.Context, args RawQuery) (int32, error) {
	if err := readNamespaces(ctx); err != nil {
		return 0, err
	}

	query := search.AddRawQueriesAsConjunction(args.String(), resolver.getClusterRawQuery())

	return resolver.root.NamespaceCount(ctx, RawQuery{Query: &query})
}

func (resolver *clusterResolver) ComplianceResults(ctx context.Context, args RawQuery) ([]*controlResultResolver, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Cluster, "ComplianceResults")

	if err := readCompliance(ctx); err != nil {
		return nil, err
	}

	runResults, err := resolver.root.ComplianceAggregator.GetResultsWithEvidence(ctx, args.String())
	if err != nil {
		return nil, err
	}
	output := newBulkControlResults()
	output.addClusterData(resolver.root, runResults, nil)
	output.addDeploymentData(resolver.root, runResults, nil)
	output.addNodeData(resolver.root, runResults, nil)
	return *output, nil
}

// K8sRoles returns GraphQL resolvers for all k8s roles
func (resolver *clusterResolver) K8sRoles(ctx context.Context, args PaginatedQuery) ([]*k8SRoleResolver, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Cluster, "K8sRoles")

	if err := readK8sRoles(ctx); err != nil {
		return nil, err
	}

	query := search.AddRawQueriesAsConjunction(args.String(), resolver.getClusterRawQuery())

	return resolver.root.K8sRoles(ctx, PaginatedQuery{Query: &query, Pagination: args.Pagination})
}

// K8sRoleCount returns count of K8s roles in this cluster
func (resolver *clusterResolver) K8sRoleCount(ctx context.Context, args RawQuery) (int32, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Cluster, "K8sRoleCount")

	if err := readK8sRoles(ctx); err != nil {
		return 0, err
	}

	q, err := args.AsV1QueryOrEmpty()
	if err != nil {
		return 0, err
	}

	q, err = search.AddAsConjunction(resolver.getClusterQuery(), q)
	if err != nil {
		return 0, err
	}

	results, err := resolver.root.K8sRoleStore.Search(ctx, q)
	if err != nil {
		return 0, err
	}
	return int32(len(results)), nil
}

// K8sRole returns clusterResolver GraphQL resolver for a given k8s role
func (resolver *clusterResolver) K8sRole(ctx context.Context, args struct{ Role graphql.ID }) (*k8SRoleResolver, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Cluster, "K8sRole")

	if err := readK8sRoles(ctx); err != nil {
		return nil, err
	}

	q := search.NewQueryBuilder().AddExactMatches(search.ClusterID, resolver.data.GetId()).
		AddExactMatches(search.RoleID, string(args.Role)).ProtoQuery()

	roles, err := resolver.root.K8sRoleStore.SearchRawRoles(ctx, q)

	if err != nil {
		return nil, err
	}

	if len(roles) == 0 {
		return nil, nil
	}

	return resolver.root.wrapK8SRole(roles[0], true, nil)
}

// ServiceAccounts returns GraphQL resolvers for all service accounts in this cluster
func (resolver *clusterResolver) ServiceAccounts(ctx context.Context, args PaginatedQuery) ([]*serviceAccountResolver, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Cluster, "ServiceAccounts")

	if err := readServiceAccounts(ctx); err != nil {
		return nil, err
	}

	query := search.AddRawQueriesAsConjunction(args.String(), resolver.getClusterRawQuery())

	return resolver.root.ServiceAccounts(ctx, PaginatedQuery{Query: &query, Pagination: args.Pagination})
}

// ServiceAccountCount returns count of Service Accounts in this cluster
func (resolver *clusterResolver) ServiceAccountCount(ctx context.Context, args RawQuery) (int32, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Cluster, "ServiceAccountCount")

	if err := readServiceAccounts(ctx); err != nil {
		return 0, err
	}

	q, err := args.AsV1QueryOrEmpty()
	if err != nil {
		return 0, err
	}

	q, err = search.AddAsConjunction(resolver.getClusterQuery(), q)
	if err != nil {
		return 0, err
	}

	results, err := resolver.root.ServiceAccountsDataStore.Search(ctx, q)
	if err != nil {
		return 0, err
	}
	return int32(len(results)), nil
}

// ServiceAccount returns clusterResolver GraphQL resolver for a given service account
func (resolver *clusterResolver) ServiceAccount(ctx context.Context, args struct{ Sa graphql.ID }) (*serviceAccountResolver, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Cluster, "ServiceAccount")

	if err := readK8sRoles(ctx); err != nil {
		return nil, err
	}

	q := search.NewQueryBuilder().AddExactMatches(search.ClusterID, resolver.data.GetId()).
		AddExactMatches(search.RoleID, string(args.Sa)).ProtoQuery()

	serviceAccounts, err := resolver.root.ServiceAccountsDataStore.SearchRawServiceAccounts(ctx, q)

	if err != nil {
		return nil, err
	}

	if len(serviceAccounts) == 0 {
		return nil, nil
	}

	return resolver.root.wrapServiceAccount(serviceAccounts[0], true, nil)
}

// Subjects returns GraphQL resolvers for all subjects in this cluster
func (resolver *clusterResolver) Subjects(ctx context.Context, args PaginatedQuery) ([]*subjectResolver, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Cluster, "Subjects")

	q, err := args.AsV1QueryOrEmpty()
	if err != nil {
		return nil, err
	}

	pagination := q.GetPagination()
	q.Pagination = nil

	subjectResolvers, err := resolver.root.wrapSubjects(resolver.getSubjects(ctx, q))
	if err != nil {
		return nil, err
	}

	resolvers, err := paginationWrapper{
		pv: pagination,
	}.paginate(subjectResolvers, nil)
	return resolvers.([]*subjectResolver), err
}

// SubjectCount returns count of Users and Groups in this cluster
func (resolver *clusterResolver) SubjectCount(ctx context.Context, args RawQuery) (int32, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Cluster, "SubjectCount")

	q, err := args.AsV1QueryOrEmpty()
	if err != nil {
		return 0, err
	}

	subjects, err := resolver.getSubjects(ctx, q)
	if err != nil {
		return 0, err
	}
	return int32(len(subjects)), nil
}

// Subject returns clusterResolver GraphQL resolver for a given subject
func (resolver *clusterResolver) Subject(ctx context.Context, args struct{ Name string }) (*subjectResolver, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Cluster, "Subject")

	subjectName, err := url.QueryUnescape(args.Name)
	if err != nil {
		return nil, err
	}
	q := search.NewQueryBuilder().
		AddExactMatches(search.ClusterID, resolver.data.GetId()).
		AddExactMatches(search.SubjectName, subjectName).
		AddExactMatches(search.SubjectKind, storage.SubjectKind_GROUP.String(), storage.SubjectKind_USER.String()).
		ProtoQuery()

	bindings, err := resolver.getRoleBindings(ctx, q)
	if err != nil {
		return nil, err
	}
	if len(bindings) == 0 {
		log.Errorf("Subject: %q not found on Cluster: %q", subjectName, resolver.data.GetName())
		return nil, nil
	}
	return resolver.root.wrapSubject(k8srbac.GetSubject(subjectName, bindings))
}

func (resolver *clusterResolver) Images(ctx context.Context, args PaginatedQuery) ([]*imageResolver, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Cluster, "Images")

	query := search.AddRawQueriesAsConjunction(args.String(), resolver.getClusterRawQuery())

	return resolver.root.Images(ctx, PaginatedQuery{Query: &query, Pagination: args.Pagination})
}

func (resolver *clusterResolver) ImageCount(ctx context.Context, args RawQuery) (int32, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Cluster, "ImageCount")

	query := search.AddRawQueriesAsConjunction(args.String(), resolver.getClusterRawQuery())

	return resolver.root.ImageCount(ctx, RawQuery{Query: &query})
}

func (resolver *clusterResolver) Components(ctx context.Context, args PaginatedQuery) ([]ComponentResolver, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Cluster, "Components")

	query := search.AddRawQueriesAsConjunction(args.String(), resolver.getClusterRawQuery())

	return resolver.root.Components(ctx, PaginatedQuery{Query: &query, Pagination: args.Pagination})
}

func (resolver *clusterResolver) ComponentCount(ctx context.Context, args RawQuery) (int32, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Cluster, "ComponentCount")

	query := search.AddRawQueriesAsConjunction(args.String(), resolver.getClusterRawQuery())

	return resolver.root.ComponentCount(ctx, RawQuery{Query: &query})
}

func (resolver *clusterResolver) Vulns(ctx context.Context, args PaginatedQuery) ([]VulnerabilityResolver, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Cluster, "Vulns")

	query := search.AddRawQueriesAsConjunction(args.String(), resolver.getClusterRawQuery())

	return resolver.root.Vulnerabilities(ctx, PaginatedQuery{Query: &query, Pagination: args.Pagination})
}

func (resolver *clusterResolver) VulnCount(ctx context.Context, args RawQuery) (int32, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Cluster, "VulnCount")

	query := search.AddRawQueriesAsConjunction(args.String(), resolver.getClusterRawQuery())

	return resolver.root.VulnerabilityCount(ctx, RawQuery{Query: &query})
}

func (resolver *clusterResolver) VulnCounter(ctx context.Context, args RawQuery) (*VulnerabilityCounterResolver, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Cluster, "VulnCounter")

	query := search.AddRawQueriesAsConjunction(args.String(), resolver.getClusterRawQuery())

	return resolver.root.VulnCounter(ctx, RawQuery{Query: &query})
}

func (resolver *clusterResolver) NodeVulnerabilities(ctx context.Context, args PaginatedQuery) ([]NodeVulnerabilityResolver, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Cluster, "NodeVulnerabilities")

	if !features.PostgresDatastore.Enabled() {
		query := search.AddRawQueriesAsConjunction(args.String(), resolver.getClusterRawQuery())
		return resolver.root.NodeVulnerabilities(ctx, PaginatedQuery{Query: &query, Pagination: args.Pagination})
	}
	// TODO : Add postgres support
	return nil, errors.New("Sub-resolver NodeVulnerabilities in clusterResolver does not support postgres yet")
}

func (resolver *clusterResolver) NodeVulnerabilityCount(ctx context.Context, args RawQuery) (int32, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Cluster, "NodeVulnerabilityCount")

	if !features.PostgresDatastore.Enabled() {
		query := search.AddRawQueriesAsConjunction(args.String(), resolver.getClusterRawQuery())
		return resolver.root.NodeVulnerabilityCount(ctx, RawQuery{Query: &query})
	}
	// TODO : Add postgres support
	return 0, errors.New("Sub-resolver NodeVulnerabilityCount in clusterResolver does not support postgres yet")
}

func (resolver *clusterResolver) NodeVulnerabilityCounter(ctx context.Context, args RawQuery) (*VulnerabilityCounterResolver, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Cluster, "NodeVulnerabilityCounter")

	if !features.PostgresDatastore.Enabled() {
		query := search.AddRawQueriesAsConjunction(args.String(), resolver.getClusterRawQuery())
		return resolver.root.NodeVulnCounter(ctx, RawQuery{Query: &query})
	}
	// TODO : Add postgres support
	return nil, errors.New("Sub-resolver NodeVulnerabilityCounter in clusterResolver does not support postgres yet")
}

func (resolver *clusterResolver) K8sVulns(ctx context.Context, args PaginatedQuery) ([]VulnerabilityResolver, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Cluster, "K8sVulns")

	query := search.AddRawQueriesAsConjunction(args.String(), resolver.getClusterRawQuery())

	return resolver.root.K8sVulnerabilities(ctx, PaginatedQuery{Query: &query, Pagination: args.Pagination})
}

func (resolver *clusterResolver) K8sVulnCount(ctx context.Context, args RawQuery) (int32, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Cluster, "K8sVulnCount")

	vulns, err := resolver.K8sVulns(ctx, PaginatedQuery{Query: args.Query})
	if err != nil {
		return 0, err
	}
	return int32(len(vulns)), nil
}

func (resolver *clusterResolver) OpenShiftVulns(ctx context.Context, args PaginatedQuery) ([]VulnerabilityResolver, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Cluster, "OpenShiftVulns")

	query := search.AddRawQueriesAsConjunction(args.String(), resolver.getClusterRawQuery())

	return resolver.root.OpenShiftVulnerabilities(ctx, PaginatedQuery{Query: &query, Pagination: args.Pagination})
}

func (resolver *clusterResolver) OpenShiftVulnCount(ctx context.Context, args RawQuery) (int32, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Cluster, "OpenShiftVulnCount")

	vulns, err := resolver.OpenShiftVulns(ctx, PaginatedQuery{Query: args.Query})
	if err != nil {
		return 0, err
	}
	return int32(len(vulns)), nil
}

func (resolver *clusterResolver) IstioVulns(ctx context.Context, args PaginatedQuery) ([]VulnerabilityResolver, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Cluster, "IstioVulns")

	query := search.AddRawQueriesAsConjunction(args.String(), resolver.getClusterRawQuery())

	return resolver.root.IstioVulnerabilities(ctx, PaginatedQuery{Query: &query, Pagination: args.Pagination})
}

func (resolver *clusterResolver) IstioVulnCount(ctx context.Context, args RawQuery) (int32, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Cluster, "IstioVulnCount")

	vulns, err := resolver.IstioVulns(ctx, PaginatedQuery{Query: args.Query})
	if err != nil {
		return 0, err
	}
	return int32(len(vulns)), nil
}

func (resolver *clusterResolver) Policies(ctx context.Context, args PaginatedQuery) ([]*policyResolver, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Cluster, "Policies")

	if err := readPolicies(ctx); err != nil {
		return nil, err
	}

	q, err := args.AsV1QueryOrEmpty()
	if err != nil {
		return nil, err
	}

	// remove pagination from query since we want to paginate the final result
	pagination := q.GetPagination()
	q.Pagination = &v1.QueryPagination{
		SortOptions: pagination.GetSortOptions(),
	}

	policyResolvers, err := resolver.root.wrapPolicies(resolver.getApplicablePolicies(ctx, q))
	if err != nil {
		return nil, err
	}
	for _, policyResolver := range policyResolvers {
		policyResolver.ctx = scoped.Context(ctx, scoped.Scope{
			Level: v1.SearchCategory_CLUSTERS,
			ID:    resolver.data.GetId(),
		})
	}

	resolvers, err := paginationWrapper{
		pv: pagination,
	}.paginate(policyResolvers, nil)
	return resolvers.([]*policyResolver), err
}

func (resolver *clusterResolver) getApplicablePolicies(ctx context.Context, q *v1.Query) ([]*storage.Policy, error) {
	policyLoader, err := loaders.GetPolicyLoader(ctx)
	if err != nil {
		return nil, err
	}

	policies, err := policyLoader.FromQuery(ctx, q)
	if err != nil {
		return nil, err
	}

	namespaces, err := resolver.root.NamespaceDataStore.SearchNamespaces(ctx,
		search.NewQueryBuilder().AddExactMatches(search.ClusterID, resolver.data.GetId()).ProtoQuery())
	if err != nil {
		return nil, err
	}

	applicable, _ := matcher.NewClusterMatcher(resolver.data, namespaces).FilterApplicablePolicies(policies)
	return applicable, nil
}

func (resolver *clusterResolver) PolicyCount(ctx context.Context, args RawQuery) (int32, error) {
	if err := readPolicies(ctx); err != nil {
		return 0, err
	}

	q, err := args.AsV1QueryOrEmpty()
	if err != nil {
		return 0, err
	}

	policies, err := resolver.getApplicablePolicies(ctx, q)
	if err != nil {
		return 0, err
	}

	return int32(len(policies)), nil
}

// PolicyStatus returns true if there is no policy violation for this cluster
func (resolver *clusterResolver) PolicyStatus(ctx context.Context, args RawQuery) (*policyStatusResolver, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Cluster, "PolicyStatus")

	if err := readAlerts(ctx); err != nil {
		return nil, err
	}

	query, err := args.AsV1QueryOrEmpty()
	if err != nil {
		return nil, err
	}

	alerts, err := resolver.getActiveDeployAlerts(ctx, query)
	if err != nil {
		return nil, err
	}

	scopedCtx := scoped.Context(ctx, scoped.Scope{
		Level: v1.SearchCategory_CLUSTERS,
		ID:    resolver.data.GetId(),
	})

	if len(alerts) == 0 {
		return &policyStatusResolver{scopedCtx, resolver.root, "pass", nil}, nil
	}

	policyIDs := set.NewStringSet()
	for _, alert := range alerts {
		policyIDs.Add(alert.GetPolicy().GetId())
	}

	return &policyStatusResolver{scopedCtx, resolver.root, "fail", policyIDs.AsSlice()}, nil
}

func (resolver *clusterResolver) Secrets(ctx context.Context, args PaginatedQuery) ([]*secretResolver, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Cluster, "Secrets")

	query := search.AddRawQueriesAsConjunction(args.String(), resolver.getClusterRawQuery())

	return resolver.root.Secrets(ctx, PaginatedQuery{Query: &query, Pagination: args.Pagination})
}

func (resolver *clusterResolver) SecretCount(ctx context.Context, args RawQuery) (int32, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Cluster, "SecretCount")

	query, err := args.AsV1QueryOrEmpty()
	if err != nil {
		return 0, err
	}

	query, err = search.AddAsConjunction(resolver.getClusterQuery(), query)
	if err != nil {
		return 0, err
	}

	result, err := resolver.root.SecretsDataStore.Search(ctx, query)
	if err != nil {
		return 0, err
	}
	return int32(len(result)), nil
}

func (resolver *clusterResolver) ControlStatus(ctx context.Context, args RawQuery) (string, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Cluster, "ControlStatus")

	if err := readCompliance(ctx); err != nil {
		return "Fail", err
	}
	r, err := resolver.getLastSuccessfulComplianceRunResult(ctx, []storage.ComplianceAggregation_Scope{storage.ComplianceAggregation_CLUSTER}, args)
	if err != nil || r == nil {
		return "Fail", err
	}
	if len(r) != 1 {
		return "Fail", errors.Errorf("unexpected number of results: expected: 1, actual: %d", len(r))
	}
	return getControlStatusFromAggregationResult(r[0]), nil
}

func (resolver *clusterResolver) Controls(ctx context.Context, args RawQuery) ([]*complianceControlResolver, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Cluster, "Controls")

	if err := readCompliance(ctx); err != nil {
		return nil, err
	}
	rs, err := resolver.getLastSuccessfulComplianceRunResult(ctx, []storage.ComplianceAggregation_Scope{storage.ComplianceAggregation_CLUSTER, storage.ComplianceAggregation_CONTROL}, args)
	if err != nil || rs == nil {
		return nil, err
	}
	resolvers, err := resolver.root.wrapComplianceControls(getComplianceControlsFromAggregationResults(rs, any, resolver.root.ComplianceStandardStore))
	if err != nil {
		return nil, err
	}
	return resolvers, nil
}

func (resolver *clusterResolver) PassingControls(ctx context.Context, args RawQuery) ([]*complianceControlResolver, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Cluster, "PassingControls")

	if err := readCompliance(ctx); err != nil {
		return nil, err
	}
	rs, err := resolver.getLastSuccessfulComplianceRunResult(ctx, []storage.ComplianceAggregation_Scope{storage.ComplianceAggregation_CLUSTER, storage.ComplianceAggregation_CONTROL}, args)
	if err != nil || rs == nil {
		return nil, err
	}
	resolvers, err := resolver.root.wrapComplianceControls(getComplianceControlsFromAggregationResults(rs, passing, resolver.root.ComplianceStandardStore))
	if err != nil {
		return nil, err
	}
	return resolvers, nil
}

func (resolver *clusterResolver) FailingControls(ctx context.Context, args RawQuery) ([]*complianceControlResolver, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Cluster, "FailingControls")

	if err := readCompliance(ctx); err != nil {
		return nil, err
	}
	rs, err := resolver.getLastSuccessfulComplianceRunResult(ctx, []storage.ComplianceAggregation_Scope{storage.ComplianceAggregation_CLUSTER, storage.ComplianceAggregation_CONTROL}, args)
	if err != nil || rs == nil {
		return nil, err
	}
	resolvers, err := resolver.root.wrapComplianceControls(getComplianceControlsFromAggregationResults(rs, failing, resolver.root.ComplianceStandardStore))
	if err != nil {
		return nil, err
	}
	return resolvers, nil
}

func (resolver *clusterResolver) ComplianceControlCount(ctx context.Context, args RawQuery) (*complianceControlCountResolver, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Cluster, "ComplianceControlCount")
	if err := readCompliance(ctx); err != nil {
		return nil, err
	}
	results, err := resolver.getLastSuccessfulComplianceRunResult(ctx, []storage.ComplianceAggregation_Scope{storage.ComplianceAggregation_CLUSTER, storage.ComplianceAggregation_CONTROL}, args)
	if err != nil {
		return nil, err
	}
	if results == nil {
		return &complianceControlCountResolver{}, nil
	}
	return getComplianceControlCountFromAggregationResults(results), nil
}

func (resolver *clusterResolver) getLastSuccessfulComplianceRunResult(ctx context.Context, scope []storage.ComplianceAggregation_Scope, args RawQuery) ([]*storage.ComplianceAggregation_Result, error) {
	if err := readCompliance(ctx); err != nil {
		return nil, err
	}
	standardIDs, err := getStandardIDs(ctx, resolver.root.ComplianceStandardStore)
	if err != nil {
		return nil, err
	}
	hasComplianceSuccessfullyRun, err := resolver.root.ComplianceDataStore.IsComplianceRunSuccessfulOnCluster(ctx, resolver.data.GetId(), standardIDs)
	if err != nil || !hasComplianceSuccessfullyRun {
		return nil, err
	}
	query, err := search.NewQueryBuilder().AddExactMatches(search.ClusterID, resolver.data.GetId()).RawQuery()
	if err != nil {
		return nil, err
	}
	if args.Query != nil {
		query = strings.Join([]string{query, *(args.Query)}, "+")
	}
	r, _, _, err := resolver.root.ComplianceAggregator.Aggregate(ctx, query, scope, storage.ComplianceAggregation_CONTROL)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (resolver *clusterResolver) getActiveDeployAlerts(ctx context.Context, q *v1.Query) ([]*storage.ListAlert, error) {
	cluster := resolver.data

	return resolver.root.ViolationsDataStore.SearchListAlerts(ctx,
		search.ConjunctionQuery(q,
			search.NewQueryBuilder().AddExactMatches(search.ClusterID, cluster.GetId()).
				AddStrings(search.ViolationState, storage.ViolationState_ACTIVE.String()).
				AddStrings(search.LifecycleStage, storage.LifecycleStage_DEPLOY.String()).ProtoQuery()))
}

func (resolver *clusterResolver) Risk(ctx context.Context) (*riskResolver, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Cluster, "Risk")
	if err := readRisks(ctx); err != nil {
		return nil, err
	}
	return resolver.root.wrapRisk(resolver.getClusterRisk(ctx))
}

func (resolver *clusterResolver) getClusterRisk(ctx context.Context) (*storage.Risk, bool, error) {
	cluster := resolver.data

	riskQuery := search.NewQueryBuilder().
		AddExactMatches(search.ClusterID, cluster.GetId()).
		AddExactMatches(search.RiskSubjectType, storage.RiskSubjectType_DEPLOYMENT.String()).
		ProtoQuery()

	risks, err := resolver.root.RiskDataStore.SearchRawRisks(ctx, riskQuery)
	if err != nil {
		return nil, false, err
	}

	risks = filterDeploymentRisksOnScope(ctx, risks...)
	scrubRiskFactors(risks...)
	aggregateRiskScore := getAggregateRiskScore(risks...)
	if aggregateRiskScore == float32(0.0) {
		return nil, false, nil
	}

	risk := &storage.Risk{
		Score: aggregateRiskScore,
		Subject: &storage.RiskSubject{
			Id:   cluster.GetId(),
			Type: storage.RiskSubjectType_CLUSTER,
		},
	}

	id, err := riskDS.GetID(risk.GetSubject().GetId(), risk.GetSubject().GetType())
	if err != nil {
		return nil, false, err
	}
	risk.Id = id

	return risk, true, nil
}

func (resolver *clusterResolver) IsGKECluster() (bool, error) {
	version := resolver.data.GetStatus().GetOrchestratorMetadata().GetVersion()
	ok := resolver.root.cveMatcher.IsGKEVersion(version)
	return ok, nil
}

func (resolver *clusterResolver) IsOpenShiftCluster() (bool, error) {
	metadata := resolver.data.GetStatus().GetOrchestratorMetadata()
	return metadata.GetIsOpenshift() != nil, nil
}

func (resolver *clusterResolver) IstioEnabled(ctx context.Context) (bool, error) {
	res, err := resolver.root.NamespaceDataStore.Search(ctx, search.NewQueryBuilder().AddExactMatches(search.Namespace, "istio-system").ProtoQuery())
	if err != nil {
		return false, err
	}
	return len(res) > 0, nil
}

func (resolver *clusterResolver) LatestViolation(ctx context.Context, args RawQuery) (*graphql.Time, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Cluster, "Latest Violation")

	q, err := args.AsV1QueryOrEmpty()
	if err != nil {
		return nil, err
	}

	q, err = resolver.getClusterConjunctionQuery(q)
	if err != nil {
		return nil, err
	}

	return getLatestViolationTime(ctx, resolver.root, q)
}

func (resolver *clusterResolver) PlottedVulns(ctx context.Context, args RawQuery) (*PlottedVulnerabilitiesResolver, error) {
	query := search.AddRawQueriesAsConjunction(args.String(), resolver.getClusterRawQuery())
	return newPlottedVulnerabilitiesResolver(ctx, resolver.root, RawQuery{Query: &query})
}

func (resolver *clusterResolver) UnusedVarSink(ctx context.Context, args RawQuery) *int32 {
	return nil
}

func (resolver *orchestratorMetadataResolver) OpenShiftVersion() (string, error) {
	return resolver.data.GetOpenshiftVersion(), nil
}
