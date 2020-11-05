import * as React from 'react';
import '../../css/pre.scss';

export class Pre extends React.Component<{}, {}> {
    render(): JSX.Element {
        return (
            <pre>{ this.props.children }</pre>
        );
    }
}
