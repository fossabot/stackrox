import { gql } from '@apollo/client';

export const CLUSTER_LIST_FRAGMENT = gql`
    fragment clusterFields on Cluster {
        id
        name
        vulnCounter {
            all {
                fixable
                total
            }
            critical {
                fixable
                total
            }
            important {
                fixable
                total
            }
            moderate {
                fixable
                total
            }
            low {
                fixable
                total
            }
        }
        status {
            orchestratorMetadata {
                version
            }
        }
        # createdAt
        namespaceCount
        deploymentCount
        nodeCount
        # policyCount(query: $policyQuery) # see https://stack-rox.atlassian.net/browse/ROX-4080
        policyStatus(query: $policyQuery) {
            status
        }
        latestViolation(query: $policyQuery)
        priority
    }
`;

export const VULN_CVE_ONLY_FRAGMENT = gql`
    fragment cveFields on EmbeddedVulnerability {
        id
        cve
        cvss
        severity
        scoreVersion
        summary
        fixedByVersion
        isFixable(query: $scopeQuery)
    }
`;

export const VULN_CVE_LIST_FRAGMENT = gql`
    fragment cveFields on EmbeddedVulnerability {
        id: cve
        cve
        cvss
        vulnerabilityTypes
        scoreVersion
        envImpact
        impactScore
        summary
        severity
        fixedByVersion
        isFixable(query: $scopeQuery)
        createdAt
        discoveredAtImage(query: $scopeQuery)
        publishedOn
        deploymentCount(query: $query)
        imageCount(query: $query)
        componentCount(query: $query)
        activeState(query: $query) {
            state
            activeContexts {
                containerName
            }
        }
    }
`;

export const DEPLOYMENT_LIST_FRAGMENT = gql`
    fragment deploymentFields on Deployment {
        id
        name
        vulnCounter {
            all {
                total
                fixable
            }
            low {
                total
                fixable
            }
            moderate {
                total
                fixable
            }
            important {
                total
                fixable
            }
            critical {
                total
                fixable
            }
        }
        deployAlerts {
            policy {
                id
            }
            time
        }
        # policyCount(query: $policyQuery) # see https://stack-rox.atlassian.net/browse/ROX-4080
        # failingPolicyCount(query: $policyQuery) # see https://stack-rox.atlassian.net/browse/ROX-4080
        policyStatus(query: $policyQuery)
        clusterName
        clusterId
        namespace
        namespaceId
        imageCount
        latestViolation(query: $policyQuery)
        priority
        images {
            scan {
                scanTime
            }
        }
    }
`;

export const NODE_LIST_FRAGMENT = gql`
    fragment nodeFields on Node {
        id
        name
        vulnCounter {
            all {
                total
                fixable
            }
            low {
                total
                fixable
            }
            moderate {
                total
                fixable
            }
            important {
                total
                fixable
            }
            critical {
                total
                fixable
            }
        }
        topVuln {
            cvss
            scoreVersion
        }
        notes
        scan {
            scanTime
            notes
        }
        osImage
        containerRuntimeVersion
        clusterName
        clusterId
        joinedAt
        priority
    }
`;

export const IMAGE_LIST_FRAGMENT = gql`
    fragment imageFields on Image {
        id
        name {
            fullName
        }
        watchStatus
        deploymentCount(query: $query)
        priority
        topVuln {
            cvss
            scoreVersion
        }
        metadata {
            v1 {
                created
            }
        }
        componentCount(query: $query)
        notes
        scan {
            scanTime
            operatingSystem
            notes
        }
        vulnCounter {
            all {
                total
                fixable
            }
            low {
                total
                fixable
            }
            moderate {
                total
                fixable
            }
            important {
                total
                fixable
            }
            critical {
                total
                fixable
            }
        }
    }
`;

export const VULN_COMPONENT_LIST_FRAGMENT = gql`
    fragment componentFields on EmbeddedImageScanComponent {
        id
        name
        version
        location
        source
        fixedIn
        vulnCounter {
            all {
                total
                fixable
            }
            low {
                total
                fixable
            }
            moderate {
                total
                fixable
            }
            important {
                total
                fixable
            }
            critical {
                total
                fixable
            }
        }
        topVuln {
            cvss
            scoreVersion
        }
        imageCount(query: $query)
        deploymentCount(query: $query)
        nodeCount(query: $query)
        priority
    }
`;

export const VULN_COMPONENT_ACTIVE_STATUS_LIST_FRAGMENT = gql`
    fragment componentFields on EmbeddedImageScanComponent {
        id
        name
        version
        location
        source
        fixedIn
        vulnCounter {
            all {
                total
                fixable
            }
            low {
                total
                fixable
            }
            moderate {
                total
                fixable
            }
            important {
                total
                fixable
            }
            critical {
                total
                fixable
            }
        }
        topVuln {
            cvss
            scoreVersion
        }
        activeState(query: $scopeQuery) {
            state
            activeContexts {
                containerName
            }
        }
        imageCount(query: $query)
        deploymentCount(query: $query)
        nodeCount(query: $query)
        priority
    }
`;

export const NAMESPACE_LIST_FRAGMENT = gql`
    fragment namespaceFields on Namespace {
        metadata {
            id
            clusterName
            clusterId
            priority
            name
        }
        vulnCounter {
            all {
                fixable
                total
            }
            critical {
                fixable
                total
            }
            important {
                fixable
                total
            }
            moderate {
                fixable
                total
            }
            low {
                fixable
                total
            }
        }
        deploymentCount: numDeployments # numDeployments is pre-calculated in namespace resolver
        imageCount(query: $query)
        # policyCount(query: $policyQuery) # see https://stack-rox.atlassian.net/browse/ROX-4080
        policyStatusOnly(query: $policyQuery)
        latestViolation(query: $policyQuery)
    }
`;

export const POLICY_LIST_FRAGMENT_CORE = gql`
    fragment corePolicyFields on Policy {
        id
        disabled
        notifiers
        name
        description
        lastUpdated
        severity
        lifecycleStages
        enforcementActions
        isDefault
    }
`;

export const UNSCOPED_POLICY_LIST_FRAGMENT = gql`
    fragment unscopedPolicyFields on Policy {
        ...corePolicyFields
        deploymentCount: failingDeploymentCount(query: $scopeQuery) # field changed to failingDeploymentCount to improve performance
        latestViolation
        policyStatus
    }
    ${POLICY_LIST_FRAGMENT_CORE}
`;

export const POLICY_LIST_FRAGMENT = gql`
    fragment policyFields on Policy {
        ...corePolicyFields
        deploymentCount: failingDeploymentCount(query: $scopeQuery) # field changed to failingDeploymentCount to improve performance
        latestViolation
        policyStatus
    }
    ${POLICY_LIST_FRAGMENT_CORE}
`;

export const POLICY_ENTITY_ALL_FIELDS_FRAGMENT = gql`
    fragment policyFields on Policy {
        id
        name
        description
        disabled
        rationale
        remediation
        severity
        policyStatus
        categories
        lastUpdated
        enforcementActions
        lifecycleStages
        isDefault
        policySections {
            sectionName
            policyGroups {
                fieldName
                values {
                    value
                }
            }
        }
        scope {
            cluster
            label {
                key
                value
            }
            namespace
        }
        exclusions {
            deployment {
                name
                scope {
                    cluster
                    label {
                        key
                        value
                    }
                    namespace
                }
            }
            expiration
            image {
                name
            }
            name
        }
    }
`;
