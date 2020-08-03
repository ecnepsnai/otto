import * as React from 'react';
import { Icon } from './Icon';
import '../../css/loading.scss';
import { Page } from './Page';

export class Loading extends React.Component<{}, {}> {
    render(): JSX.Element {
        return (
            <div>
                <Icon.Spinner pulse />
                <span className="ml-1 loading-text">
                    Loading...
                </span>
            </div>
        );
    }
}

export class PageLoading extends React.Component<{}, {}> {
    render(): JSX.Element {
        return (
            <Page header={<Loading />} />
        );
    }
}
