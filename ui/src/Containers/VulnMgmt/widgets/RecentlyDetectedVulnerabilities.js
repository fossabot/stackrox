import React, { useContext } from 'react';
import PropTypes from 'prop-types';
import { useQuery } from 'react-apollo';
import gql from 'graphql-tag';
import sortBy from 'lodash/sortBy';

import entityTypes from 'constants/entityTypes';
import queryService from 'modules/queryService';
import workflowStateContext from 'Containers/workflowStateContext';
import ViewAllButton from 'Components/ViewAllButton';
import Loader from 'Components/Loader';
import Widget from 'Components/Widget';
import NumberedList from 'Components/NumberedList';
import { getVulnerabilityChips } from 'utils/vulnerabilityUtils';
import NoResultsMessage from 'Components/NoResultsMessage';
import { cveSortFields } from 'constants/sortFields';

export const RECENTLY_DETECTED_VULNERABILITIES = gql`
    query recentlyDetectedVulnerabilities(
        $query: String
        $scopeQuery: String
        $pagination: Pagination
    ) {
        results: vulnerabilities(query: $query, pagination: $pagination) {
            id: cve
            cve
            cvss
            scoreVersion
            deploymentCount
            imageCount
            isFixable(query: $scopeQuery)
            envImpact
            createdAt
            summary
        }
    }
`;

const processData = (data, workflowState, limit) => {
    let results = data && data.results && data.results.filter(datum => datum.createdAt);
    results = sortBy(results, ['createdAt', 'cvss', 'envImpact'])
        .slice(-limit)
        .reverse(); // @TODO: filter on the client side until we have pagination on Vulnerabilities

    // @TODO: remove JSX generation from processing data and into Numbered List function
    return getVulnerabilityChips(workflowState, results);
};

const RecentlyDetectedVulnerabilities = ({ entityContext, search, limit }) => {
    const entityContextObject = queryService.entityContextToQueryObject(entityContext); // deals with BE inconsistency

    const queryObject = {
        ...entityContextObject,
        ...search,
        [cveSortFields.CVE_TYPE]: 'IMAGE_CVE'
    }; // Combine entity context and search
    const query = queryService.objectToWhereClause(queryObject); // get final gql query string

    const { loading, data = {} } = useQuery(RECENTLY_DETECTED_VULNERABILITIES, {
        variables: {
            query,
            scopeQuery: queryService.objectToWhereClause(entityContextObject),
            pagination: queryService.getPagination(
                {
                    id: cveSortFields.CVE_CREATED_TIME,
                    desc: true
                },
                0,
                limit
            )
        }
    });

    let content = <Loader />;

    const workflowState = useContext(workflowStateContext);
    if (!loading) {
        const processedData = processData(data, workflowState, limit);

        if (!processedData || processedData.length === 0) {
            content = (
                <NoResultsMessage message="No vulnerabilities found" className="p-6" icon="info" />
            );
        } else {
            content = (
                <div className="w-full">
                    <NumberedList data={processedData} />
                </div>
            );
        }
    }

    const viewAllURL = workflowState
        .pushList(entityTypes.CVE)
        .setSort([{ id: cveSortFields.CVE_CREATED_TIME, desc: true }])
        .toUrl();

    return (
        <Widget
            className="h-full pdf-page"
            header="Recently Detected Vulnerabilities"
            headerComponents={<ViewAllButton url={viewAllURL} />}
        >
            {content}
        </Widget>
    );
};

RecentlyDetectedVulnerabilities.propTypes = {
    entityContext: PropTypes.shape({}),
    search: PropTypes.shape({}),
    limit: PropTypes.number
};

RecentlyDetectedVulnerabilities.defaultProps = {
    entityContext: {},
    search: {},
    limit: 5
};

export default RecentlyDetectedVulnerabilities;
