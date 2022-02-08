import * as React from 'react';
import { Icon } from './Icon';
import { Page } from './Page';
import '../../css/loading.scss';

export const Loading: React.FC = () => {
    return (
        <div>
            <Icon.Spinner pulse />
            <span className="ms-1 loading-text">
                Loading...
            </span>
        </div>
    );
};

export const PageLoading: React.FC = () => {
    return (
        <Page loading />
    );
};
