import { EnforcementAction } from 'Containers/Violations/types/violationTypes';

export const ENFORCEMENT_ACTIONS: Record<EnforcementAction, string> = Object.freeze({
    UNSET_ENFORCEMENT: 'UNSET_ENFORCEMENT',
    SCALE_TO_ZERO_ENFORCEMENT: 'SCALE_TO_ZERO_ENFORCEMENT',
    UNSATISFIABLE_NODE_CONSTRAINT_ENFORCEMENT: 'UNSATISFIABLE_NODE_CONSTRAINT_ENFORCEMENT',
    KILL_POD_ENFORCEMENT: 'KILL_POD_ENFORCEMENT',
    FAIL_BUILD_ENFORCEMENT: 'FAIL_BUILD_ENFORCEMENT',
    FAIL_KUBE_REQUEST_ENFORCEMENT: 'FAIL_KUBE_REQUEST_ENFORCEMENT',
    FAIL_DEPLOYMENT_CREATE_ENFORCEMENT: 'FAIL_DEPLOYMENT_CREATE_ENFORCEMENT',
    FAIL_DEPLOYMENT_UPDATE_ENFORCEMENT: 'FAIL_DEPLOYMENT_UPDATE_ENFORCEMENT',
});

export const ENFORCEMENT_ACTIONS_AS_STRING: Record<EnforcementAction, string> = Object.freeze({
    UNSET_ENFORCEMENT: 'No enforcement',
    SCALE_TO_ZERO_ENFORCEMENT: 'Scale to 0',
    UNSATISFIABLE_NODE_CONSTRAINT_ENFORCEMENT: 'Unsatisfiable Node Constraint "BlockedByStackRox"',
    KILL_POD_ENFORCEMENT: 'Kill Pod',
    FAIL_BUILD_ENFORCEMENT: 'Fail Build',
    FAIL_KUBE_REQUEST_ENFORCEMENT: 'Block Kubernetes Action',
    FAIL_DEPLOYMENT_CREATE_ENFORCEMENT: 'Block Deployment Create',
    FAIL_DEPLOYMENT_UPDATE_ENFORCEMENT: 'Block Deployment Update',
});

export const ENFORCEMENT_ACTIONS_AS_PAST_TENSE: Record<EnforcementAction, string> = Object.freeze({
    UNSET_ENFORCEMENT: 'No enforcement',
    SCALE_TO_ZERO_ENFORCEMENT: 'Scaled to 0',
    UNSATISFIABLE_NODE_CONSTRAINT_ENFORCEMENT:
        'Applied Unsatisfiable Node Constraint "BlockedByStackRox"',
    KILL_POD_ENFORCEMENT: 'Killed Pod',
    FAIL_BUILD_ENFORCEMENT: 'Failed Build',
    FAIL_KUBE_REQUEST_ENFORCEMENT: 'Blocked Kubernetes Action',
    FAIL_DEPLOYMENT_CREATE_ENFORCEMENT: 'Blocked Deployment Create',
    FAIL_DEPLOYMENT_UPDATE_ENFORCEMENT: 'Blocked Deployment Update',
});

export const BLOCKING_ENFORCEMENT_ACTIONS = Object.freeze(
    new Set([
        ENFORCEMENT_ACTIONS.FAIL_KUBE_REQUEST_ENFORCEMENT,
        ENFORCEMENT_ACTIONS.FAIL_DEPLOYMENT_CREATE_ENFORCEMENT,
        ENFORCEMENT_ACTIONS.FAIL_DEPLOYMENT_UPDATE_ENFORCEMENT,
    ])
);
