import static util.SplunkUtil.SPLUNK_ADMIN_PASSWORD
import static util.SplunkUtil.postToSplunk
import static util.SplunkUtil.tearDownSplunk

import groups.Integration
import io.stackrox.proto.api.v1.AlertServiceOuterClass
import services.AlertService
import services.NetworkBaselineService
import services.ApiTokenService
import spock.lang.Unroll
import util.SplunkUtil
import util.Timer

import java.nio.file.Paths
import java.time.LocalDateTime
import java.util.concurrent.TimeUnit

import org.junit.Rule
import org.junit.rules.Timeout
import org.junit.experimental.categories.Category

import com.jayway.restassured.path.json.JsonPath
import com.jayway.restassured.response.Response

class IntegrationsSplunkViolationsTest extends BaseSpecification {
    @Rule
    @SuppressWarnings(["JUnitPublicProperty"])
    Timeout globalTimeout = new Timeout(1000, TimeUnit.SECONDS)

    private static final String ASSETS_DIR = Paths.get(
            System.getProperty("user.dir"), "artifacts", "splunk-violations-test")
    private static final String PATH_TO_SPLUNK_TA_SPL = Paths.get(ASSETS_DIR,
    "2021-07-23-TA-stackrox-1.2.0-input-validation-patched.spl")
    private static final String PATH_TO_CIM_TA_TGZ = Paths.get(ASSETS_DIR,
    "splunk-common-information-model-cim_4190.tgz")
    private static final String STACKROX_REMOTE_LOCATION = "/tmp/stackrox.spl"
    private static final String CIM_REMOTE_LOCATION = "/tmp/cim.tgz"
    private static final String TEST_NAMESPACE = "qa-splunk-violation"
    private static final String SPLUNK_INPUT_NAME = "stackrox-violations-input"

    def setupSpec() {
        // when using "Analyst" api token to access central Splunk violations endpoint
        // authorisation plugin prevents violations from being returned
        // this leads to no violations being propagated to Splunk
        disableAuthzPlugin()
    }

    private void configureSplunkTA(SplunkUtil.SplunkDeployment splunkDeployment, String centralHost) {
        println "${LocalDateTime.now()} Starting Splunk TA configuration"
        def podName = orchestrator
                .getPods(TEST_NAMESPACE, splunkDeployment.deployment.getName())
                .get(0)
                .getMetadata()
                .getName()
        int port = splunkDeployment.splunkPortForward.getLocalPort()

        println "${LocalDateTime.now()} Copying TA and CIM app files to splunk pod"
        orchestrator.copyFileToPod(PATH_TO_SPLUNK_TA_SPL, TEST_NAMESPACE, podName, STACKROX_REMOTE_LOCATION)
        orchestrator.copyFileToPod(PATH_TO_CIM_TA_TGZ, TEST_NAMESPACE, podName, CIM_REMOTE_LOCATION)
        println "${LocalDateTime.now()} Installing TA"
        postToSplunk(port, "/services/apps/local",
                ["name": STACKROX_REMOTE_LOCATION, "filename": "true"])
        println "${LocalDateTime.now()} Installing CIM app"
        postToSplunk(port, "/services/apps/local",
                ["name": CIM_REMOTE_LOCATION, "filename": "true"])
        // fix minimum free disk space parameter
        // default value is 5Gb and CircleCI free disk space is less than that
        // that can prevent data from being indexed
        orchestrator.execInContainer(splunkDeployment.deployment,
                "sudo /opt/splunk/bin/splunk set minfreemb 200 -auth admin:${SPLUNK_ADMIN_PASSWORD}"
        )
        // Splunk needs to be restarted after TA installation
        postToSplunk(splunkDeployment.splunkPortForward.getLocalPort(), "/services/server/control/restart", [:])

        println("${LocalDateTime.now()} Configuring Stackrox TA")
        def tokenResp = ApiTokenService.generateToken("splunk-token-${splunkDeployment.uid}", "Analyst")
        postToSplunk(port, "/servicesNS/nobody/TA-stackrox/configs/conf-ta_stackrox_settings/additional_parameters",
                ["central_endpoint": "${centralHost}:443",
                 "api_token": tokenResp.getToken(),])
        // create new input to search violations from
        postToSplunk(port, "/servicesNS/nobody/TA-stackrox/data/inputs/stackrox_violations",
                ["name": SPLUNK_INPUT_NAME, "interval": "1", "from_checkpoint": "2000-01-01T00:00:00.000Z"])
    }

    @Unroll
    @Category(Integration)
    def "Verify Splunk violations: StackRox violations reach Splunk TA"() {
        given:
        "Splunk TA is installed and configured, network and process violations triggered"
        orchestrator.deleteNamespace(TEST_NAMESPACE)
        orchestrator.ensureNamespaceExists(TEST_NAMESPACE)
        addStackroxImagePullSecret(TEST_NAMESPACE)
        String centralHost = orchestrator.getServiceIP("central", "stackrox")
        def splunkDeployment = SplunkUtil.createSplunk(orchestrator, TEST_NAMESPACE, false)
        configureSplunkTA(splunkDeployment, centralHost)
        triggerProcessViolation(splunkDeployment)
        triggerNetworkFlowViolation(splunkDeployment, centralHost)

        when:
        "Search for violations in Splunk"
        // Splunk search for violations is volatile for some reason.
        // We added retries to make this test less flaky.
        List<Map<String, String>> results = Collections.emptyList()
        boolean hasNetworkViolation = false
        boolean hasProcessViolation = false
        def port = splunkDeployment.splunkPortForward.getLocalPort()
        for (int i = 0; i < 20; i++) {
            println "Attempt ${i} to get violations from Splunk"
            def searchId = SplunkUtil.createSearch(port, "| from datamodel Alerts.Alerts")
            TimeUnit.SECONDS.sleep(10)
            Response response = SplunkUtil.getSearchResults(port, searchId)
            // We should have at least one violation in the response
            if (response != null) {
                results = response.getBody().jsonPath().getList("results")
                if (!results.isEmpty()) {
                    for (result in results) {
                        hasNetworkViolation |= isNetworkViolation(result)
                        hasProcessViolation |= isProcessViolation(result)
                    }
                    if (hasNetworkViolation && hasProcessViolation) {
                        println "Success!"
                        break
                    }
                }
            }
        }

        then:
        "StackRox violations are in Splunk"
        assert !results.isEmpty()
        assert hasNetworkViolation
        assert hasProcessViolation
        for (result in results) {
            validateCimMappings(result)
        }

        cleanup:
        "remove splunk"
        if (splunkDeployment) {
            tearDownSplunk(orchestrator, splunkDeployment)
        }
        orchestrator.deleteNamespace(TEST_NAMESPACE)
    }

    private static void validateCimMappings(Map<String, String> result) {
        def originalEvent = new JsonPath(result.get("_raw"))
        Map<String, String> violationInfo = originalEvent.getMap("violationInfo") ?: [:]
        Map<String, String> policyInfo = originalEvent.getMap("policyInfo") ?: [:]
        Map<String, String> processInfo = originalEvent.getMap("processInfo") ?: [:]

        assert result.get("app") == "stackrox"
        assert result.get("type") == "alert"
        verifyRequiredResultKey(result, "id", violationInfo.get("violationId"))
        verifyRequiredResultKey(result, "description", violationInfo.get("violationMessage"))
        verifyRequiredResultKey(result, "signature_id", policyInfo.get("policyName"))
        // Note that policyDescription and signature might be absent, i.e. null
        assert result.get("signature") == policyInfo.get("policyDescription")

        // user
        def processUid = processInfo.get("processUid")
        def processGid = processInfo.get("processGid")
        def expectedUser = processUid == null || processGid == null
                ? "unknown" : processUid + ":" + processGid
        verifyRequiredResultKey(result, "user", expectedUser)

        // severity
        String severity = coalesce(extractNestedString(originalEvent, "policyInfo.policySeverity"), "unknown")
                .replace("UNSET_", "unknown_")
                .replace("_SEVERITY", "")
                .toLowerCase()
        assert result.get("severity") == severity

        // dest_type
        String destType = coalesce(
                extractNestedString(originalEvent, "networkFlowInfo.destination.deploymentType"),
                extractNestedString(originalEvent, "networkFlowInfo.destination.entityType")
        )
        assert result.get("dest_type") == destType

        // src_type
        String srcType = coalesce(
                extractNestedString(originalEvent, "networkFlowInfo.source.deploymentType"),
                extractNestedString(originalEvent, "networkFlowInfo.source.entityType"),
                extractNestedString(originalEvent, "deploymentInfo.deploymentType"),
                extractNestedString(originalEvent, "resourceInfo.resourceType")
        )
        verifyRequiredResultKey(result, "src_type", srcType)

        // dest
        String dest = coalesce(
                extractDestOrSrc(originalEvent, "destination"),
                extractNestedString(originalEvent, "networkFlowInfo.destination.name"),
                "unknown")
        assert result.get("dest") == dest

        // src
        String src = coalesce(
                extractDestOrSrc(originalEvent, "source"),
                extractNestedString(originalEvent, "networkFlowInfo.source.name"),
                extractSourceViaDeploymentInfo(originalEvent),
                extractSourceViaResourceInfo(originalEvent)
        )
        verifyRequiredResultKey(result, "src", src)
    }

    private static void verifyRequiredResultKey(Map<String, String> result, String key, String expectedValue) {
        assert Objects.requireNonNull(result.get(key)) == expectedValue
    }

    private static <T> T coalesce(T... args) {
        for (T arg : args) {
            if (arg != null) {
                return arg
            }
        }
        return null
    }

    @SuppressWarnings(["ReturnNullFromCatchBlock"])
    private static String extractNestedString(JsonPath jsonPath, String path) {
        try {
            return jsonPath.getString(path)
        } catch (IllegalArgumentException e) {
            return null
        }
    }

    private static String extractDestOrSrc(JsonPath originalEvent, String prefix) {
        String clusterName = extractNestedString(originalEvent, "deploymentInfo.clusterName")
        if (clusterName == null) {
            return null
        }
        String deploymentNamespace = extractNestedString(originalEvent, "networkFlowInfo.${prefix}.deploymentNamespace")
        if (deploymentNamespace == null) {
            return null
        }
        String deploymentType = extractNestedString(originalEvent, "networkFlowInfo.${prefix}.deploymentType")
        if (deploymentType == null) {
            return null
        }
        String delimiter = deploymentType == "Pod" ? " > " : "/"
        String name = extractNestedString(originalEvent, "networkFlowInfo.${prefix}.name")

        return name == null ? null : "${clusterName}/${deploymentNamespace}${delimiter}${deploymentType}:${name}"
    }

    static String extractSourceViaDeploymentInfo(JsonPath originalEvent) {
        String clusterName = extractNestedString(originalEvent, "deploymentInfo.clusterName")
        if (clusterName == null) {
            return null
        }
        String deploymentNamespace = extractNestedString(originalEvent, "deploymentInfo.deploymentNamespace")
        if (deploymentNamespace == null) {
            return null
        }
        String deploymentType = extractNestedString(originalEvent, "deploymentInfo.deploymentType")
        if (deploymentType == null) {
            return null
        }
        String deploymentName = extractNestedString(originalEvent, "deploymentInfo.deploymentName")
        if (deploymentName == null) {
            return null
        }
        String podId = extractNestedString(originalEvent, "violationInfo.podId")
        String podPart = podId == null ? "" : " > ${podId}"
        String containerName = extractNestedString(originalEvent, "violationInfo.containerName")
        String containerPart = containerName == null ? "" : "/${containerName}"
        String podDescription = deploymentType == "Pod"
            ? " > ${deploymentType}:${deploymentName}"
            : "/${deploymentType}:${deploymentName}${podPart}"
        return "${clusterName}/${deploymentNamespace}${podDescription}${containerPart}"
    }

    @SuppressWarnings(["IfStatementCouldBeTernary"]) // much more readable this way
    static String extractSourceViaResourceInfo(JsonPath originalEvent) {
        String clusterName = extractNestedString(originalEvent, "resourceInfo.clusterName")
        if (clusterName == null) {
            return null
        }
        String namespace = extractNestedString(originalEvent, "resourceInfo.namespace")
        if (namespace == null) {
            return null
        }
        String resourceType = extractNestedString(originalEvent, "resourceInfo.resourceType")
        if (resourceType == null) {
            return null
        }
        String resourceName = extractNestedString(originalEvent, "resourceInfo.name")
        if (resourceName == null) {
            return null
        }

        return "${clusterName}/${namespace}/${resourceType}:${resourceName}"
    }

    def triggerProcessViolation(SplunkUtil.SplunkDeployment splunkDeployment) {
        orchestrator.execInContainer(splunkDeployment.deployment, "curl http://127.0.0.1:10248/")
        assert waitForAlertWithPolicyId("86804b96-e87e-4eae-b56e-1718a8a55763")
    }

    def triggerNetworkFlowViolation(SplunkUtil.SplunkDeployment splunkDeployment, String centralHost) {
        NetworkBaselineService.lockNetworkBaseline(splunkDeployment.getDeployment().deploymentUid)
        orchestrator.execInContainer(splunkDeployment.deployment,
                "for i in `seq 25`; do curl http://${centralHost}:443; sleep 1; done")

        // TODO: this code is flaky; see https://stack-rox.atlassian.net/browse/ROX-7772
        assert waitForAlertWithPolicyId("1b74ffdd-8e67-444c-9814-1c23863c8ccb")
    }

    private boolean waitForAlertWithPolicyId(String policyId) {
        retryUntilTrue({
            AlertService.getViolations(AlertServiceOuterClass.ListAlertsRequest.newBuilder()
                    .setQuery("Namespace:${TEST_NAMESPACE},Violation State:*")
                    .build())
                    .asList()
                    .any { a -> a.getPolicy().getId() == policyId }
            }, 10
        )
    }

    boolean isNetworkViolation(Map<String, String> result) {
        return isViolationOfType(result, "NETWORK_FLOW")
    }

    boolean isProcessViolation(Map<String, String> result) {
        return isViolationOfType(result, "PROCESS_EVENT")
    }

    boolean isViolationOfType(Map<String, String> result, String type) {
        Map<String, String> violationInfo = new JsonPath(result.get("_raw")).getMap("violationInfo") ?: [:]
        return violationInfo.get("violationType") == type
    }

    // returns whether true condition was achieved
    boolean retryUntilTrue(Closure<Boolean> closure, int retries) {
        Timer timer = new Timer(retries, 10)
        while (timer.IsValid()) {
            def result = closure()
            if (result) {
                return true
            }
        }
        return false
    }
}
