import * as React from 'react';
import '../../css/title.scss';

export interface PageProps {
    title?: string;
    header?: JSX.Element;
}

export class Page extends React.Component<PageProps, {}> {
    componentDidMount(): void {
        if (this.props.title) {
            document.title = this.props.title + ' — Otto';
        } else {
            document.title = 'Otto';
        }
    }
    render(): JSX.Element {
        let content: JSX.Element;
        if (this.props.header) {
            content = this.props.header;
        } else {
            content = (<p>{ this.props.title }</p>);
        }

        return (
            <div className="page">
                <div className="page-title">
                    <div className="container">
                        { content }
                    </div>
                </div>
                <div className="container">
                    { this.props.children }
                </div>
            </div>
        );
    }
}
