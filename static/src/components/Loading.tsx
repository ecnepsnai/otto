import * as React from 'react';
import { Icon } from './Icon';
import '../../css/loading.scss';
import { Page } from './Page';

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
        <Page header={<Loading />} />
    );
};
