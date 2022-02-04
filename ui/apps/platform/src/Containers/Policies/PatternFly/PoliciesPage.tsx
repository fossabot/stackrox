import React from 'react';
import { useLocation, useParams } from 'react-router-dom';

import usePermissions from 'hooks/usePermissions';

import { parsePoliciesSearchString } from './policies.utils';
import PoliciesTablePage from './Table/PoliciesTablePage';
import PolicyPage from './PolicyPage';

function PoliciesPage() {
    /*
     * Examples of urls for PolicyPage:
     * /main/policies/:policyId
     * /main/policies/:policyId?action=edit
     * /main/policies?action=create
     *
     * Examples of urls for PolicyTablePage:
     * /main/policies
     * /main/policies?s[Lifecycle Stage]=BUILD
     * /main/policies?s[Lifecycle Stage]=BUILD&s[Lifecycle State]=DEPLOY
     * /main/policies?s[Lifecycle State]=RUNTIME&s[Severity]=CRITICAL_SEVERITY
     */
    const { search } = useLocation();
    const { pageAction, searchFilter } = parsePoliciesSearchString(search);
    const { policyId } = useParams();

    const { hasReadAccess, hasReadWriteAccess } = usePermissions();
    const hasReadAccessForPolicy = hasReadAccess('Policy');
    const hasWriteAccessForPolicy = hasReadWriteAccess('Policy');

    if (!hasReadAccessForPolicy) {
        return <div>TODO</div>;
    }

    if (pageAction || policyId) {
        return (
            <PolicyPage
                hasWriteAccessForPolicy={hasWriteAccessForPolicy}
                pageAction={pageAction}
                policyId={policyId}
            />
        );
    }

    return (
        <PoliciesTablePage
            hasWriteAccessForPolicy={hasWriteAccessForPolicy}
            searchFilter={searchFilter}
        />
    );
}

export default PoliciesPage;
