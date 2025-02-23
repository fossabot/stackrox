import static Services.waitForViolation
import static services.ClusterService.DEFAULT_CLUSTER_NAME

import io.stackrox.proto.api.v1.ApiTokenService.GenerateTokenResponse
import io.stackrox.proto.api.v1.NamespaceServiceOuterClass
import io.stackrox.proto.api.v1.SearchServiceOuterClass as SSOC
import io.stackrox.proto.storage.DeploymentOuterClass

import groups.BAT
import objects.Deployment
import orchestratormanager.OrchestratorTypes
import services.AlertService
import services.ApiTokenService
import services.BaseService
import services.DeploymentService
import services.ImageService
import services.NamespaceService
import services.NetworkGraphService
import services.SearchService
import services.SecretService
import services.SummaryService
import util.NetworkGraphUtil

import util.Env
import spock.lang.IgnoreIf
import org.junit.AssumptionViolatedException
import org.junit.experimental.categories.Category
import spock.lang.Unroll

@Category(BAT)
class SACTest extends BaseSpecification {
    static final private String DEPLOYMENTNGINX_NAMESPACE_QA1 = "sac-deploymentnginx-qa1"
    static final private String NAMESPACE_QA1 = "qa-test1"
    static final private String DEPLOYMENTNGINX_NAMESPACE_QA2 = "sac-deploymentnginx-qa2"
    static final private String NAMESPACE_QA2 = "qa-test2"
    static final private String TESTROLE = "Continuous Integration"
    static final private String SECRETNAME = "sac-secret"
    static final protected String ALLACCESSTOKEN = "allAccessToken"
    static final protected String NOACCESSTOKEN = "noAccess"
    static final protected Deployment DEPLOYMENT_QA1 = new Deployment()
            .setName(DEPLOYMENTNGINX_NAMESPACE_QA1)
            .setImage(TEST_IMAGE)
            .addPort(22, "TCP")
            .addAnnotation("test", "annotation")
            .setEnv(["CLUSTER_NAME": "main"])
            .setNamespace(NAMESPACE_QA1)
            .addLabel("app", "test")
    static final protected Deployment DEPLOYMENT_QA2 = new Deployment()
            .setName(DEPLOYMENTNGINX_NAMESPACE_QA2)
            .setImage(TEST_IMAGE)
            .addPort(22, "TCP")
            .addAnnotation("test", "annotation")
            .setEnv(["CLUSTER_NAME": "main"])
            .setNamespace(NAMESPACE_QA2)
            .addLabel("app", "test")

    static final private List<Deployment> DEPLOYMENTS = [DEPLOYMENT_QA1, DEPLOYMENT_QA2,]

    static final private UNSTABLE_FLOWS = [
            // monitoring doesn't keep a persistent outgoing connection, so we might or might not see this flow.
            "stackrox/monitoring -> INTERNET",
    ] as Set

    // Increase the timeout conditionally based on whether we are running race-detection builds or within OpenShift
    // environments. Both take longer than the default values.
    static final private Integer WAIT_FOR_VIOLATION_TIMEOUT =
            isRaceBuild() ? 600 : ((Env.mustGetOrchestratorType() == OrchestratorTypes.OPENSHIFT) ? 100 : 60)

    static final private Integer WAIT_FOR_RISK_RETRIES =
            isRaceBuild() ? 300 : ((Env.mustGetOrchestratorType() == OrchestratorTypes.OPENSHIFT) ? 50 : 30)

    def setupSpec() {
        // Make sure we scan the image initially to make reprocessing faster.
        def img = Services.scanImage(TEST_IMAGE)
        assert img.hasScan()

        orchestrator.batchCreateDeployments(DEPLOYMENTS)
        for (Deployment deployment : DEPLOYMENTS) {
            assert Services.waitForDeployment(deployment)
        }
        // Make sure each deployment has caused at least one alert
        assert waitForViolation(DEPLOYMENT_QA1.name, "Secure Shell (ssh) Port Exposed",
                WAIT_FOR_VIOLATION_TIMEOUT)
        assert waitForViolation(DEPLOYMENT_QA2.name, "Secure Shell (ssh) Port Exposed",
                WAIT_FOR_VIOLATION_TIMEOUT)

        // Make sure each deployment has a risk score.
        def deployments = DeploymentService.listDeployments()
        deployments.each { DeploymentOuterClass.ListDeployment dep ->
            try {
                withRetry(WAIT_FOR_RISK_RETRIES, 2) {
                    assert DeploymentService.getDeploymentWithRisk(dep.id).hasRisk()
                }
            } catch (Exception e) {
                throw new AssumptionViolatedException("Failed to retrieve risk from deployment ${dep.name}", e)
            }
        }
    }

    def cleanupSpec() {
        BaseService.useBasicAuth()
        for (Deployment deployment : DEPLOYMENTS) {
            orchestrator.deleteDeployment(deployment)
        }
        [NAMESPACE_QA1, NAMESPACE_QA2].forEach {
            ns ->
                orchestrator.deleteNamespace(ns)
                orchestrator.waitForNamespaceDeletion(ns)
        }
    }

    static getAlertCount() {
        return AlertService.getViolations().size()
    }

    static getImageCount() {
        return ImageService.getImages().size()
    }

    static getDeploymentCount() {
        return DeploymentService.listDeploymentsSearch().deploymentsCount
    }

    static getNamespaceCount() {
        return NamespaceService.getNamespaces().size()
    }

    GenerateTokenResponse useToken(String tokenName) {
        GenerateTokenResponse token = ApiTokenService.generateToken(tokenName, TESTROLE)
        BaseService.useApiToken(token.token)
        token
    }

    static getSpecificQuery(String category) {
        def queryString = category + ":*"
        def query = SSOC.RawSearchRequest.newBuilder()
                .setQuery(queryString)
                .build()
        return query
    }

    def createSecret(String namespace) {
        String secID = orchestrator.createSecret(SECRETNAME, namespace)
        SecretService.waitForSecret(secID, 10)
    }

    def deleteSecret(String namespace) {
        orchestrator.deleteSecret(SECRETNAME, namespace)
    }

    Boolean summaryTestShouldSeeNoClustersAndNodes() { true }

    @Unroll
    @IgnoreIf({ Env.CI_JOBNAME.contains("postgres") })
    def "Verify that only namespace #sacResource is visible when using SAC"() {
        when:
        "Create test API token with a built-in role"
        useToken("deployments-access-token")
        then:
        "Call API and verify data returned is within scoped access"
        def result = DeploymentService.listDeployments()
        println result.toString()
        assert result.size() == 1
        assert DeploymentService.getDeploymentWithRisk(result.first().id).hasRisk()
        def resourceNotAllowed = result.find { it.namespace != sacResource }
        assert resourceNotAllowed == null
        cleanup:
        BaseService.useBasicAuth()
        where:
        "Data inputs are: "
        sacResource | _
        NAMESPACE_QA2 | _
    }

    @IgnoreIf({ Env.CI_JOBNAME.contains("postgres") })
    def "Verify GetSummaryCounts using a token without access receives no results"() {
        when:
        "GetSummaryCounts is called using a token without access"
        createSecret(DEPLOYMENT_QA1.namespace)
        useToken(NOACCESSTOKEN)
        def result = SummaryService.getCounts()
        then:
        "Verify GetSumamryCounts returns no results"
        assert result.getNumDeployments() == 0
        assert result.getNumSecrets() == 0
        assert result.getNumNodes() == 0
        assert result.getNumClusters() == 0
        assert result.getNumImages() == 0
        cleanup:
        "Cleanup"
        BaseService.useBasicAuth()
        deleteSecret(DEPLOYMENT_QA1.namespace)
    }

    @IgnoreIf({ Env.CI_JOBNAME.contains("postgres") })
    def "Verify GetSummaryCounts using a token with partial access receives partial results"() {
        when:
        "GetSummaryCounts is called using a token with restricted access"
        createSecret(DEPLOYMENT_QA1.namespace)
        createSecret(DEPLOYMENT_QA2.namespace)
        useToken("getSummaryCountsToken")
        def result = SummaryService.getCounts()
        then:
        "Verify correct counts are returned by GetSummaryCounts"
        assert result.getNumDeployments() == 1
        assert result.getNumSecrets() == orchestrator.getSecretCount(DEPLOYMENT_QA1.namespace)
        if (summaryTestShouldSeeNoClustersAndNodes()) {
            assert result.getNumNodes() == 0
            assert result.getNumClusters() == 0
        }
        assert result.getNumImages() == 1
        cleanup:
        "Cleanup"
        BaseService.useBasicAuth()
        deleteSecret(DEPLOYMENT_QA1.namespace)
        deleteSecret(DEPLOYMENT_QA2.namespace)
    }

    @IgnoreIf({ Env.CI_JOBNAME.contains("postgres") })
    def "Verify GetSummaryCounts using a token with all access receives all results"() {
        when:
        "GetSummaryCounts is called using a token with all access"
        createSecret(DEPLOYMENT_QA1.namespace)
        createSecret(DEPLOYMENT_QA2.namespace)
        useToken(ALLACCESSTOKEN)
        def result = SummaryService.getCounts()
        then:
        "Verify results are returned in each category"
        assert result.getNumDeployments() >= 2
        // These may be created by other tests so it's hard to know the exact number.
        assert result.getNumSecrets() >= 2
        assert result.getNumNodes() > 0
        assert result.getNumClusters() >= 1
        assert result.getNumImages() >= 1
        cleanup:
        "Cleanup"
        BaseService.useBasicAuth()
        deleteSecret(DEPLOYMENT_QA1.namespace)
        deleteSecret(DEPLOYMENT_QA2.namespace)
    }

    @Unroll
    @IgnoreIf({ Env.CI_JOBNAME.contains("postgres") })
    def "Verify alerts count is scoped"() {
        given:
        def query = SSOC.RawQuery.newBuilder().setQuery(
                "Deployment:${DEPLOYMENT_QA1.name},${DEPLOYMENT_QA2.name}"
        ).build()

        when:
        def alertsCount = { String tokenName ->
            BaseService.useBasicAuth()
            useToken(tokenName)
            AlertService.alertClient.countAlerts(query).count
        }

        then:
        assert alertsCount(NOACCESSTOKEN) == 0
        // getSummaryCountsToken has access only to QA1 deployment while
        // ALLACCESSTOKEN has access to QA1 and QA2. Since deployments are identical
        // number of alerts for ALLACCESSTOKEN should be twice of getSummaryCountsToken.
        assert 2 * alertsCount("getSummaryCountsToken") == alertsCount(ALLACCESSTOKEN)

        cleanup:
        BaseService.useBasicAuth()
    }

    @IgnoreIf({ Env.CI_JOBNAME.contains("postgres") })
    def "Verify ListSecrets using a token without access receives no results"() {
        when:
        "ListSecrets is called using a token without view access to Secrets"
        BaseService.useBasicAuth()
        createSecret(DEPLOYMENT_QA1.namespace)
        useToken(NOACCESSTOKEN)
        def result = SecretService.listSecrets()
        then:
        "Verify no secrets are returned by ListSecrets"
        assert result.secretsCount == 0
        cleanup:
        "Cleanup"
        BaseService.useBasicAuth()
        deleteSecret(DEPLOYMENT_QA1.namespace)
    }

    @IgnoreIf({ Env.CI_JOBNAME.contains("postgres") })
    def "Verify ListSecrets using a token with access receives some results"() {
        when:
        "ListSecrets is called using a token with view access to Secrets"
        BaseService.useBasicAuth()
        createSecret(DEPLOYMENT_QA1.namespace)
        createSecret(DEPLOYMENT_QA2.namespace)
        useToken("listSecretsToken")
        def result = SecretService.listSecrets()

        then:
        "Verify no secrets are returned by ListSecrets"
        assert result.secretsCount > 0

        cleanup:
        "Cleanup"
        BaseService.useBasicAuth()
        deleteSecret(DEPLOYMENT_QA1.namespace)
        deleteSecret(DEPLOYMENT_QA2.namespace)
    }

    @Unroll
    @IgnoreIf({ Env.CI_JOBNAME.contains("postgres") })
    def "Verify Search on #category resources using the #tokenName token returns #numResults results"() {
        when:
        "A search is performed using the given token"
        def query = getSpecificQuery(category)
        useToken(tokenName)
        def result = SearchService.search(query)

        then:
        "Verify the specified number of results are returned"
        assert result.resultsCount == numResults

        cleanup:
        "Cleanup"
        BaseService.useBasicAuth()

        where:
        "Data inputs are: "
        tokenName                | category     | numResults
        NOACCESSTOKEN            | "Cluster"    | 0
        "searchDeploymentsToken" | "Deployment" | 1
        "searchImagesToken"      | "Image"      | 1
    }

    @Unroll
    @IgnoreIf({ Env.CI_JOBNAME.contains("postgres") })
    def "Verify Search on #category resources using the #tokenName token returns >= #minReturned results"() {
        when:
        "A search is performed using the given token"
        def query = getSpecificQuery(category)
        useToken(tokenName)
        def result = SearchService.search(query)

        then:
        "Verify >= the specified number of results are returned"
        assert result.resultsCount >= minReturned

        cleanup:
        "Cleanup"
        BaseService.useBasicAuth()
        where:
        "Data inputs are: "
        tokenName           | category     | minReturned
        "searchAlertsToken" | "Deployment" | 1
    }

    @IgnoreIf({ Env.CI_JOBNAME.contains("postgres") })
    def "Verify Search using the allAccessToken returns results for all search categories"() {
        when:
        "A search is performed using the allAccessToken"
        createSecret(DEPLOYMENT_QA1.namespace)
        def query = getSpecificQuery("Cluster")
        useToken(ALLACCESSTOKEN)
        def result = SearchService.search(query)
        then:
        "Verify something was returned for every search category"
        for (SSOC.SearchResponse.Count numResults : result.countsList) {
            // Policies are globally scoped so our cluster-scoped query won't return any
            if (numResults.category == SSOC.SearchCategory.POLICIES) {
                continue
            }
            assert numResults.count > 0
        }
        cleanup:
        "Cleanup"
        BaseService.useBasicAuth()
        deleteSecret(DEPLOYMENT_QA1.namespace)
    }

    @Unroll
    @IgnoreIf({ Env.CI_JOBNAME.contains("postgres") })
    def "Verify Autocomplete on #category resources using the #tokenName token returns #numResults results"() {
        when:
        "Search is called using a token without view access to Deployments"
        def query = getSpecificQuery(category)
        useToken(tokenName)
        def result = SearchService.autocomplete(query)
        then:
        "Verify no results are returned by Search"
        assert result.getValuesCount() == numResults
        cleanup:
        "Cleanup"
        BaseService.useBasicAuth()
        where:
        "Data inputs are: "
        tokenName                | category     | numResults
        NOACCESSTOKEN            | "Deployment" | 0
        NOACCESSTOKEN            | "Image"      | 0
        "searchDeploymentsToken" | "Deployment" | 1
        "searchImagesToken"      | "Image"      | 1
        "searchNamespacesToken"  | "Namespace"  | 1
    }

    @Unroll
    @IgnoreIf({ Env.CI_JOBNAME.contains("postgres") })
    def "Verify Autocomplete on #category resources using the #tokenName token returns >= to #minReturned results"() {
        when:
        "Autocomplete is called using the given token"
        def query = getSpecificQuery(category)
        useToken(tokenName)
        def result = SearchService.autocomplete(query)
        then:
        "Verify exactly the expected number of results are returned"
        assert result.getValuesCount() >= minReturned
        cleanup:
        "Cleanup"
        BaseService.useBasicAuth()
        where:
        "Data inputs are: "
        tokenName      | category     | minReturned
        ALLACCESSTOKEN | "Deployment" | 2
        ALLACCESSTOKEN | "Image"      | 1
    }

    @Unroll
    @IgnoreIf({ Env.CI_JOBNAME.contains("postgres") })
    def "Verify using the #tokenName token with the #service service returns #numReturned results"() {
        when:
        "The service under test is called using the given token"
        useToken(tokenName)
        def result = resultCountFunc()
        then:
        "Verify exactly the expected number of results are returned"
        assert result == numReturned
        cleanup:
        "Cleanup"
        BaseService.useBasicAuth()
        where:
        "Data inputs are: "
        tokenName                | numReturned | resultCountFunc          | service
        NOACCESSTOKEN            | 0           | this.&getDeploymentCount | "Deployment"
        "searchDeploymentsToken" | 1           | this.&getDeploymentCount | "Deployment"
        NOACCESSTOKEN            | 0           | this.&getAlertCount      | "Alert"
        NOACCESSTOKEN            | 0           | this.&getImageCount      | "Image"
        NOACCESSTOKEN            | 0           | this.&getNamespaceCount  | "Namespace"
        "searchNamespacesToken"  | 1           | this.&getNamespaceCount  | "Namespace"
        "searchImagesToken"      | 1           | this.&getImageCount      | "Image"
    }

    @Unroll
    @IgnoreIf({ Env.CI_JOBNAME.contains("postgres") })
    def "Verify using the #tokenName token with the #service service returns >= to #minNumReturned results"() {
        when:
        "The service under test is called using the given token"
        useToken(tokenName)
        def result = resultCountFunc()
        then:
        "Verify greater than or equal to the expected number of results are returned"
        assert result >= minNumReturned
        cleanup:
        "Cleanup"
        BaseService.useBasicAuth()
        where:
        "Data inputs are: "
        tokenName           | minNumReturned | resultCountFunc          | service
        ALLACCESSTOKEN      | 1              | this.&getAlertCount      | "Alert"
        "searchAlertsToken" | 1              | this.&getAlertCount      | "Alert"
        ALLACCESSTOKEN      | 1              | this.&getImageCount      | "Image"
        ALLACCESSTOKEN      | 2              | this.&getDeploymentCount | "Deployment"
        ALLACCESSTOKEN      | 2              | this.&getNamespaceCount  | "Namespace"
    }

    static getNamespaceId(String name) {
        def namespaces = NamespaceService.getNamespaces()
        for (NamespaceServiceOuterClass.Namespace namespace : namespaces) {
            if (namespace.getMetadata().name == name) {
                return namespace.getMetadata().id
            }
        }
        return null
    }

    @Unroll
    @IgnoreIf({ Env.CI_JOBNAME.contains("postgres") })
    def "Verify Namespace service SAC is enforced properly when using the #tokenName token"() {
        when:
        "We try to get one namespace we have access to and one namespace we don't have access to "
        def qa1NamespaceId = getNamespaceId(DEPLOYMENT_QA1.namespace)
        def qa2NamespaceId = getNamespaceId(DEPLOYMENT_QA2.namespace)
        useToken(tokenName)
        def qa1 = NamespaceService.getNamespace(qa1NamespaceId)
        def qa2 = NamespaceService.getNamespace(qa2NamespaceId)
        then:
        "We should get results for the namespace we have access to and null for the namespace we don't have access to"
        // Either the value should be null and it is, else the value is not null
        assert qa1Null && qa1 == null || qa1 != null
        assert qa2Null && qa2 == null || qa2 != null
        cleanup:
        "Cleanup"
        BaseService.useBasicAuth()
        where:
        "Data inputs are: "
        tokenName               | qa1Null | qa2Null
        NOACCESSTOKEN           | true    | true
        "searchNamespacesToken" | false   | true
        ALLACCESSTOKEN          | false   | false
    }

    @Unroll
    @IgnoreIf({ Env.CI_JOBNAME.contains("postgres") })
    def "Verify search with SAC and token #tokenName yields the same number of results as restricted search"() {
        when:
        "Searching for categories ${categories} in namespace ${namespace} with basic auth"
        def restrictedQuery = SSOC.RawSearchRequest.newBuilder()
                .addAllCategories(categories)
                .setQuery("Cluster:${DEFAULT_CLUSTER_NAME}+Namespace:${namespace}")
                .build()
        BaseService.useBasicAuth()
        def restrictedWithBasicAuthCount = SearchService.search(restrictedQuery).resultsCount

        and:
        "Searching for categories ${categories} in namespace ${namespace} with a token with all access"
        useToken(ALLACCESSTOKEN)
        def restrictedWithAllAccessCount = SearchService.search(restrictedQuery).resultsCount

        and:
        "Searching for categories ${categories} in all NS with token ${tokenName} restricted to namespace ${namespace}"
        useToken(tokenName)
        def unrestrictedQuery = SSOC.RawSearchRequest.newBuilder()
                .addAllCategories(categories)
                .setQuery("Cluster:${DEFAULT_CLUSTER_NAME}")
                .build()
        def unrestrictedWithSACCount = SearchService.search(unrestrictedQuery).resultsCount

        then:
        "The number of results should be the same for everything"

        println "With basic auth + restricted query: ${restrictedWithBasicAuthCount}"
        println "With all access token + restricted query: ${restrictedWithAllAccessCount}"
        println "With SAC restricted token + unrestricted query: ${unrestrictedWithSACCount}"

        assert restrictedWithBasicAuthCount == restrictedWithAllAccessCount
        assert restrictedWithAllAccessCount == unrestrictedWithSACCount

        cleanup:
        "Cleanup"
        BaseService.useBasicAuth()

        where:
        "Data inputs are: "
        tokenName                          | namespace      | categories
        NOACCESSTOKEN                      | "non_existent" | [SSOC.SearchCategory.NAMESPACES,
                                                               SSOC.SearchCategory.IMAGES,
                                                               SSOC.SearchCategory.DEPLOYMENTS]
        "kubeSystemDeploymentsImagesToken" | "kube-system"  | [SSOC.SearchCategory.IMAGES]
        "searchNamespacesToken"            | NAMESPACE_QA1  | [SSOC.SearchCategory.NAMESPACES]
        "searchDeploymentsToken"           | NAMESPACE_QA1  | [SSOC.SearchCategory.DEPLOYMENTS]
        "searchDeploymentsImagesToken"     | NAMESPACE_QA1  | [SSOC.SearchCategory.IMAGES]
    }

    @IgnoreIf({ Env.CI_JOBNAME.contains("postgres") })
    def "Verify that SAC has the same effect as query restriction for network flows"() {
        when:
        "Obtaining the network graph for the StackRox namespace with all access"
        BaseService.useBasicAuth()
        def networkGraphWithAllAccess = NetworkGraphService.getNetworkGraph(null, "Namespace:stackrox")
        def allAccessFlows = NetworkGraphUtil.flowStrings(networkGraphWithAllAccess)
        allAccessFlows.removeAll(UNSTABLE_FLOWS)
        println allAccessFlows

        def allAccessFlowsWithoutNeighbors = allAccessFlows.findAll {
            it.matches("(stackrox/.*|INTERNET) -> (stackrox/.*|INTERNET)")
        }
        println allAccessFlowsWithoutNeighbors

        and:
        "Obtaining the network graph for the StackRox namespace with a SAC restricted token"
        useToken("stackroxNetFlowsToken")
        def networkGraphWithSAC = NetworkGraphService.getNetworkGraph(null, "Namespace:stackrox")
        def sacFlows = NetworkGraphUtil.flowStrings(networkGraphWithSAC)
        sacFlows.removeAll(UNSTABLE_FLOWS)
        println sacFlows

        and:
        "Obtaining the network graph for the StackRox namespace with a SAC restricted token and no query"
        def networkGraphWithSACNoQuery = NetworkGraphService.getNetworkGraph()
        def sacFlowsNoQuery = NetworkGraphUtil.flowStrings(networkGraphWithSACNoQuery)
        sacFlowsNoQuery.removeAll(UNSTABLE_FLOWS)
        println sacFlowsNoQuery

        then:
        "Query-restricted and non-restricted flows should be equal under SAC"
        assert sacFlows == sacFlowsNoQuery

        and:
        "The flows should be equal to the flows obtained with all access after removing masked endpoints"
        def sacFlowsFiltered = new HashSet<String>(sacFlows)
        sacFlowsFiltered.removeAll { it.contains("masked deployment") }

        def sacFlowsNoQueryFiltered = new HashSet<String>(sacFlowsNoQuery)
        sacFlowsNoQueryFiltered.removeAll { it.contains("masked deployment") }

        assert allAccessFlowsWithoutNeighbors == sacFlowsFiltered
        assert allAccessFlowsWithoutNeighbors == sacFlowsNoQueryFiltered

        and:
        "The flows obtained with SAC should contain some masked deployments"
        assert sacFlowsFiltered.size() < sacFlows.size()
        assert sacFlowsNoQueryFiltered.size() < sacFlowsNoQuery.size()

        and:
        "The masked deployments should be external to stackrox namespace"
        assert sacFlows.size() - sacFlowsFiltered.size() ==
                allAccessFlows.size() - allAccessFlowsWithoutNeighbors.size()
        assert sacFlowsNoQuery.size() - sacFlowsNoQueryFiltered.size() ==
                allAccessFlows.size() - allAccessFlowsWithoutNeighbors.size()

        cleanup:
        "Cleanup"
        BaseService.useBasicAuth()
    }
}
