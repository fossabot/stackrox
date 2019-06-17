import entityTypes from 'constants/entityTypes';
import contextTypes from 'constants/contextTypes';
import { mainPath, urlEntityListTypes, urlEntityTypes } from '../routePaths';
import URLService, { getTypeKeyFromParamValue } from './URLService';

function getMatch(params) {
    return {
        params
    };
}

const location = {
    search: {
        q1: 'v1'
    }
};

function translateTypes(obj) {
    return Object.keys(obj).reduce((prev, key) => {
        const val = getTypeKeyFromParamValue(obj[key]);
        const ret = { ...prev };
        ret[key] = val || obj[key];
        return ret;
    }, {});
}

const baseParams = {
    context: 'compliance'
};

const ENTITY_ENTITY_PARAMS = {
    context: 'compliance',
    pageEntityType: urlEntityTypes.CLUSTER,
    pageEntityId: 'pageEntityId',
    entityListType1: urlEntityListTypes.DEPLOYMENT,
    entityId1: 'entityId1',
    entityType2: urlEntityTypes.NAMESPACE,
    entityId2: 'entityId2'
};

const LIST_ENTITY_PARAMS = {
    context: 'compliance',
    pageEntityListType: urlEntityListTypes.CLUSTER,
    entityId1: 'entityId1',
    entityType2: urlEntityTypes.NAMESPACE,
    entityId2: 'entityId2'
};

const ENTITY_LIST_PARAMS = {
    context: 'compliance',
    pageEntityType: urlEntityTypes.CLUSTER,
    pageEntityId: 'pageEntityId',
    entityListType1: urlEntityListTypes.DEPLOYMENT,
    entityId1: 'entityId1',
    entityListType2: urlEntityListTypes.NAMESPACE,
    entityId2: 'entityId2'
};

const LIST_LIST_PARAMS = {
    context: 'compliance',
    pageEntityListType: urlEntityListTypes.CLUSTER,
    entityId1: 'entityId1',
    entityListType2: urlEntityListTypes.NAMESPACE,
    entityId2: 'entityId2'
};

it('copies and translates params', () => {
    const match = getMatch(ENTITY_ENTITY_PARAMS);
    const url = URLService.getURL(match, location);
    expect(url.urlParams).toEqual(translateTypes(ENTITY_ENTITY_PARAMS));
    expect(url.q).toEqual(location.search);
});

it('sets base params for list path', () => {
    const match = getMatch({ context: 'compliance' });
    const url = URLService.getURL(match, location);
    url.base(entityTypes.CLUSTER);
    expect(url.urlParams.pageEntityListType).toEqual(entityTypes.CLUSTER);
});

it('sets base params for entity path', () => {
    const match = getMatch({ context: 'compliance' });
    const url = URLService.getURL(match, location);
    const pageEntityId = '123';

    url.base(entityTypes.CLUSTER, pageEntityId);
    expect(url.urlParams.pageEntityType).toEqual(entityTypes.CLUSTER);
    expect(url.urlParams.pageEntityId).toEqual(pageEntityId);
});

it('sets base context', () => {
    const match = getMatch({});
    const context = 'testContext';
    const url = URLService.getURL(match, location, context);

    url.base(entityTypes.CLUSTER, '123', context);
    expect(url.urlParams.context).toEqual(context);
});

it('Incrementally pushes list_entity path', () => {
    const match = getMatch(baseParams);
    const url = URLService.getURL(match);
    expect(url.url()).toEqual(`${mainPath}/compliance`);

    url.push(entityTypes.CLUSTER);
    expect(url.urlParams).toEqual({
        context: contextTypes.COMPLIANCE,
        pageEntityListType: entityTypes.CLUSTER
    });
    expect(url.url()).toEqual(`${mainPath}/compliance/clusters`);

    url.push('entityId1');
    expect(url.urlParams).toEqual({
        context: contextTypes.COMPLIANCE,
        pageEntityListType: entityTypes.CLUSTER,
        entityId1: 'entityId1'
    });
    expect(url.url()).toEqual(`${mainPath}/compliance/clusters/entityId1`);

    url.push(entityTypes.DEPLOYMENT);
    expect(url.urlParams).toEqual({
        context: contextTypes.COMPLIANCE,
        pageEntityListType: entityTypes.CLUSTER,
        entityId1: 'entityId1',
        entityListType2: entityTypes.DEPLOYMENT
    });
    expect(url.url()).toEqual(`${mainPath}/compliance/clusters/entityId1/deployments`);

    url.push('entityId2');
    expect(url.urlParams).toEqual({
        context: contextTypes.COMPLIANCE,
        pageEntityListType: entityTypes.CLUSTER,
        entityId1: 'entityId1',
        entityListType2: entityTypes.DEPLOYMENT,
        entityId2: 'entityId2'
    });
    expect(url.url()).toEqual(`${mainPath}/compliance/clusters/entityId1/deployments/entityId2`);
});

it('Incrementally pushes list_list path', () => {
    const match = getMatch(baseParams);
    const url = URLService.getURL(match)
        .push(entityTypes.CLUSTER)
        .push('entityId1')
        .push(entityTypes.DEPLOYMENT, 'entityId2');

    expect(url.urlParams).toEqual({
        context: contextTypes.COMPLIANCE,
        pageEntityListType: entityTypes.CLUSTER,
        entityId1: 'entityId1',
        entityType2: entityTypes.DEPLOYMENT,
        entityId2: 'entityId2'
    });
    expect(url.url()).toEqual(`${mainPath}/compliance/clusters/entityId1/deployment/entityId2`);
});

it('incrementally pushes entity_entity path', () => {
    const match = getMatch(baseParams);
    const url = URLService.getURL(match);

    url.push(entityTypes.CLUSTER, 'pageEntityId');
    expect(url.urlParams).toEqual({
        context: contextTypes.COMPLIANCE,
        pageEntityId: 'pageEntityId',
        pageEntityType: entityTypes.CLUSTER
    });
    expect(url.url()).toEqual(`${mainPath}/compliance/cluster/pageEntityId`);

    url.push(entityTypes.DEPLOYMENT);
    expect(url.urlParams).toEqual({
        context: contextTypes.COMPLIANCE,
        pageEntityId: 'pageEntityId',
        pageEntityType: entityTypes.CLUSTER,
        entityListType1: entityTypes.DEPLOYMENT
    });
    expect(url.url()).toEqual(`${mainPath}/compliance/cluster/pageEntityId/deployments`);

    url.push('entityId1');
    expect(url.urlParams).toEqual({
        context: contextTypes.COMPLIANCE,
        pageEntityId: 'pageEntityId',
        pageEntityType: entityTypes.CLUSTER,
        entityListType1: entityTypes.DEPLOYMENT,
        entityId1: 'entityId1'
    });
    expect(url.url()).toEqual(`${mainPath}/compliance/cluster/pageEntityId/deployments/entityId1`);

    url.push(entityTypes.NAMESPACE);
    expect(url.urlParams).toEqual({
        context: contextTypes.COMPLIANCE,
        pageEntityId: 'pageEntityId',
        pageEntityType: entityTypes.CLUSTER,
        entityListType1: entityTypes.DEPLOYMENT,
        entityId1: 'entityId1',
        entityListType2: entityTypes.NAMESPACE
    });
    expect(url.url()).toEqual(
        `${mainPath}/compliance/cluster/pageEntityId/deployments/entityId1/namespaces`
    );

    url.push('entityId2');
    expect(url.urlParams).toEqual({
        context: contextTypes.COMPLIANCE,
        pageEntityId: 'pageEntityId',
        pageEntityType: entityTypes.CLUSTER,
        entityListType1: entityTypes.DEPLOYMENT,
        entityId1: 'entityId1',
        entityListType2: entityTypes.NAMESPACE,
        entityId2: 'entityId2'
    });
    expect(url.url()).toEqual(
        `${mainPath}/compliance/cluster/pageEntityId/deployments/entityId1/namespaces/entityId2`
    );
});

it('incrementally pushes entity_list path', () => {
    const match = getMatch(baseParams);
    const url = URLService.getURL(match)
        .push(entityTypes.CLUSTER, 'pageEntityId')
        .push(entityTypes.DEPLOYMENT)
        .push('entityId1')
        .push(entityTypes.NAMESPACE, 'entityId2');

    expect(url.urlParams).toEqual({
        context: contextTypes.COMPLIANCE,
        pageEntityId: 'pageEntityId',
        pageEntityType: entityTypes.CLUSTER,
        entityListType1: entityTypes.DEPLOYMENT,
        entityId1: 'entityId1',
        entityType2: entityTypes.NAMESPACE,
        entityId2: 'entityId2'
    });

    expect(url.url()).toEqual(
        `${mainPath}/compliance/cluster/pageEntityId/deployments/entityId1/namespace/entityId2`
    );
});

it('pops entity_entity path', () => {
    const match = getMatch(ENTITY_ENTITY_PARAMS);
    const url = URLService.getURL(match);

    url.pop();
    expect(url.urlParams).toEqual({
        context: contextTypes.COMPLIANCE,
        pageEntityType: entityTypes.CLUSTER,
        pageEntityId: ENTITY_ENTITY_PARAMS.pageEntityId,
        entityListType1: entityTypes.DEPLOYMENT,
        entityId1: ENTITY_ENTITY_PARAMS.entityId1
    });
    expect(url.url()).toEqual(`${mainPath}/compliance/cluster/pageEntityId/deployments/entityId1`);

    url.pop();
    expect(url.urlParams).toEqual({
        context: contextTypes.COMPLIANCE,
        pageEntityId: ENTITY_ENTITY_PARAMS.pageEntityId,
        pageEntityType: entityTypes.CLUSTER,
        entityListType1: entityTypes.DEPLOYMENT
    });
    expect(url.url()).toEqual(`${mainPath}/compliance/cluster/pageEntityId/deployments`);

    url.pop();
    expect(url.urlParams).toEqual({
        context: contextTypes.COMPLIANCE,
        pageEntityId: ENTITY_ENTITY_PARAMS.pageEntityId,
        pageEntityType: entityTypes.CLUSTER
    });
    expect(url.url()).toEqual(`${mainPath}/compliance/cluster/pageEntityId`);

    url.pop();
    expect(url.urlParams).toEqual({
        context: contextTypes.COMPLIANCE
    });
    expect(url.url()).toEqual(`${mainPath}/compliance`);
});

it('pops entity_list path', () => {
    const match = getMatch(ENTITY_LIST_PARAMS);
    const url = URLService.getURL(match);

    url.pop();
    expect(url.urlParams).toEqual({
        context: contextTypes.COMPLIANCE,
        pageEntityType: entityTypes.CLUSTER,
        pageEntityId: ENTITY_ENTITY_PARAMS.pageEntityId,
        entityListType1: entityTypes.DEPLOYMENT,
        entityId1: ENTITY_ENTITY_PARAMS.entityId1,
        entityListType2: entityTypes.NAMESPACE
    });
    expect(url.url()).toEqual(
        `${mainPath}/compliance/cluster/pageEntityId/deployments/entityId1/namespaces`
    );
});

it('pops list_entity path', () => {
    const match = getMatch(LIST_ENTITY_PARAMS);
    const url = URLService.getURL(match);

    expect(url.urlParams).toEqual({
        context: contextTypes.COMPLIANCE,
        pageEntityListType: entityTypes.CLUSTER,
        entityId1: LIST_ENTITY_PARAMS.entityId1,
        entityType2: entityTypes.NAMESPACE,
        entityId2: LIST_ENTITY_PARAMS.entityId2
    });
    expect(url.url()).toEqual(`${mainPath}/compliance/clusters/entityId1/namespace/entityId2`);

    url.pop();
    expect(url.urlParams).toEqual({
        context: contextTypes.COMPLIANCE,
        pageEntityListType: entityTypes.CLUSTER,
        entityId1: LIST_ENTITY_PARAMS.entityId1
    });
    expect(url.url()).toEqual(`${mainPath}/compliance/clusters/entityId1`);

    url.pop();
    expect(url.urlParams).toEqual({
        context: contextTypes.COMPLIANCE,
        pageEntityListType: entityTypes.CLUSTER
    });
    expect(url.url()).toEqual(`${mainPath}/compliance/clusters`);
});

it('pops list_list path', () => {
    const match = getMatch(LIST_LIST_PARAMS);
    const url = URLService.getURL(match);

    url.pop();
    expect(url.urlParams).toEqual({
        context: contextTypes.COMPLIANCE,
        pageEntityListType: entityTypes.CLUSTER,
        entityId1: LIST_ENTITY_PARAMS.entityId1,
        entityListType2: entityTypes.NAMESPACE
    });
    expect(url.url()).toEqual(`${mainPath}/compliance/clusters/entityId1/namespaces`);
});

it('replaces entity_entity path', () => {
    const match = getMatch(baseParams);
    const url = URLService.getURL(match);

    url.base(entityTypes.CLUSTER, 'pageEntityId')
        .push(entityTypes.DEPLOYMENT)
        .push(entityTypes.NAMESPACE);
    expect(url.urlParams).toEqual({
        context: contextTypes.COMPLIANCE,
        pageEntityId: 'pageEntityId',
        pageEntityType: entityTypes.CLUSTER,
        entityListType1: entityTypes.NAMESPACE
    });
    expect(url.url()).toEqual(`${mainPath}/compliance/cluster/pageEntityId/namespaces`);

    url.push('entityId1').push('entityId1-1');
    expect(url.urlParams).toEqual({
        context: contextTypes.COMPLIANCE,
        pageEntityId: 'pageEntityId',
        pageEntityType: entityTypes.CLUSTER,
        entityListType1: entityTypes.NAMESPACE,
        entityId1: 'entityId1-1'
    });
    expect(url.url()).toEqual(`${mainPath}/compliance/cluster/pageEntityId/namespaces/entityId1-1`);

    url.push(entityTypes.DEPLOYMENT, 'entityId2').push('entityId2-1');
    expect(url.urlParams).toEqual({
        context: contextTypes.COMPLIANCE,
        pageEntityId: 'pageEntityId',
        pageEntityType: entityTypes.CLUSTER,
        entityListType1: entityTypes.NAMESPACE,
        entityId1: 'entityId1-1',
        entityType2: entityTypes.DEPLOYMENT,
        entityId2: 'entityId2-1'
    });
    expect(url.url()).toEqual(
        `${mainPath}/compliance/cluster/pageEntityId/namespaces/entityId1-1/deployment/entityId2-1`
    );
});

it('replaces entity_list path', () => {
    const match = getMatch(baseParams);
    const url = URLService.getURL(match);

    url.base(entityTypes.CLUSTER, 'pageEntityId')
        .push(entityTypes.DEPLOYMENT)
        .push(entityTypes.NAMESPACE)
        .push('entityId1')
        .push(entityTypes.DEPLOYMENT)
        .push(entityTypes.NODE);

    expect(url.urlParams).toEqual({
        context: contextTypes.COMPLIANCE,
        pageEntityId: 'pageEntityId',
        pageEntityType: entityTypes.CLUSTER,
        entityListType1: entityTypes.NAMESPACE,
        entityId1: 'entityId1',
        entityListType2: entityTypes.NODE
    });
});

it('replaces list_entity path', () => {
    const match = getMatch(baseParams);
    const url = URLService.getURL(match)
        .push(entityTypes.CLUSTER)
        .push(entityTypes.NAMESPACE);

    expect(url.urlParams).toEqual({
        context: contextTypes.COMPLIANCE,
        pageEntityListType: entityTypes.NAMESPACE
    });

    url.push('entityId1').push('entityId1-1');
    expect(url.urlParams).toEqual({
        context: contextTypes.COMPLIANCE,
        pageEntityListType: entityTypes.NAMESPACE,
        entityId1: 'entityId1-1'
    });

    url.push(entityTypes.DEPLOYMENT).push(entityTypes.NODE);
    expect(url.urlParams).toEqual({
        context: contextTypes.COMPLIANCE,
        pageEntityListType: entityTypes.NAMESPACE,
        entityId1: 'entityId1-1',
        entityListType2: entityTypes.NODE
    });

    url.push('entityId2').push('entityId2-1');

    expect(url.urlParams).toEqual({
        context: contextTypes.COMPLIANCE,
        pageEntityListType: entityTypes.NAMESPACE,
        entityId1: 'entityId1-1',
        entityListType2: entityTypes.NODE,
        entityId2: 'entityId2-1'
    });
});
