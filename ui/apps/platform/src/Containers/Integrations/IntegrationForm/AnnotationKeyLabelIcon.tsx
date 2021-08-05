import React, { ReactElement } from 'react';
import { Popover } from '@patternfly/react-core';
import { HelpIcon } from '@patternfly/react-icons';

function AnnotationKeyLabelIcon(): ReactElement {
    return (
        <Popover
            showClose={false}
            bodyContent={
                <div>
                    Using an annotation key, you can define an audience to notify about policy
                    violations associated with any given deployment or namespace. If the deployment
                    has the annotation, its value overrides the default; otherwise, if the namespace
                    has the annotation, its value overrides the default.
                </div>
            }
        >
            <button
                type="button"
                aria-label="More info for annotation field"
                // eslint-disable-next-line react/jsx-no-bind
                onClick={(e) => e.preventDefault()}
                aria-describedby="simple-form-name-01"
                className="pf-c-form__group-label-help"
            >
                <HelpIcon noVerticalAlign />
            </button>
        </Popover>
    );
}

export default AnnotationKeyLabelIcon;
