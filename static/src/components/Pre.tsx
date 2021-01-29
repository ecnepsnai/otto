import * as React from 'react';
import '../../css/pre.scss';

export class Pre extends React.Component<{}, {}> {
    render(): JSX.Element {
        return (
            <div className="pre-wrapper"><pre>{ this.props.children }</pre></div>
        );
    }
}
