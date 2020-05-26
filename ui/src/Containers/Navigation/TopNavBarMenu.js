import React from 'react';
import PropTypes from 'prop-types';
import { withRouter } from 'react-router-dom';
import { connect } from 'react-redux';
import { createStructuredSelector } from 'reselect';
import { MoreHorizontal } from 'react-feather';

import { selectors } from 'reducers';
import { actions as authActions } from 'reducers/auth';
import Menu from 'Components/Menu';
import Avatar from 'Components/Avatar';
import getUserAttributeMap from 'utils/userDataUtils';
import { isBackendFeatureFlagEnabled, knownBackendFlags } from 'utils/featureFlags';

const topNavMenuBtnClass =
    'no-underline text-base-600 hover:bg-base-200 items-center cursor-pointer';

function TopNavBarMenu({ logout, shouldHaveReadPermission, userData, featureFlags }) {
    const options = [{ label: 'Logout', onClick: () => logout() }];

    if (shouldHaveReadPermission('Licenses')) {
        options.unshift({ label: 'Manage Product License', link: '/main/license' });
    }

    let buttonIcon = <MoreHorizontal className="mx-4 h-4 w-4 pointer-events-none" />;
    let buttonText = null;
    const buttonTextClassName = 'border rounded-full mx-3 p-3 text-xl border-base-400';

    if (isBackendFeatureFlagEnabled(featureFlags, knownBackendFlags.ROX_CURRENT_USER_INFO, false)) {
        const { userInfo, userAttributes } = userData;
        if (userAttributes) {
            const userAttributeMap = getUserAttributeMap(userAttributes);
            const { name, email, username } = userAttributeMap;
            const userDisplayedName = name || username;
            const menuOptionComponent = (
                <div className="flex flex-col pl-2">
                    <div className="font-700">{userDisplayedName}</div>
                    {email && <div className="lowercase text-base-500 italic pt-px">{email}</div>}
                    <div className="pt-1">
                        <span className="font-700 pr-2">Roles ({userInfo.roles.length}):</span>
                        <span>{userInfo.roles[0].name}</span>
                    </div>
                </div>
            );
            options.unshift({ component: menuOptionComponent, link: '/main/user' });
            buttonIcon = <Avatar name={userDisplayedName} className="mx-3" />;
            buttonText = null;
        }
    }

    return (
        <Menu
            className={`${topNavMenuBtnClass} border-l border-base-400`}
            menuClassName="min-w-48"
            buttonIcon={buttonIcon}
            buttonText={buttonText}
            buttonTextClassName={buttonTextClassName}
            button
            hideCaret
            options={options}
        />
    );
}

TopNavBarMenu.propTypes = {
    logout: PropTypes.func.isRequired,
    shouldHaveReadPermission: PropTypes.func.isRequired,
    userData: PropTypes.shape({
        userInfo: PropTypes.shape({
            roles: PropTypes.arrayOf(
                PropTypes.shape({
                    name: PropTypes.string,
                })
            ),
        }),
        userAttributes: PropTypes.arrayOf(PropTypes.shape({})),
    }).isRequired,
    featureFlags: PropTypes.arrayOf(
        PropTypes.shape({
            envVar: PropTypes.string.isRequired,
            enabled: PropTypes.bool.isRequired,
        })
    ).isRequired,
};

const mapStateToProps = createStructuredSelector({
    shouldHaveReadPermission: selectors.shouldHaveReadPermission,
    userData: selectors.getCurrentUser,
    featureFlags: selectors.getFeatureFlags,
});

const mapDispatchToProps = (dispatch) => ({
    logout: () => dispatch(authActions.logout()),
});

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(TopNavBarMenu));
