import * as React from 'react';
import '../../css/page.scss';

interface PageProps {
    title?: string;
    header?: JSX.Element;
    children?: React.ReactNode;
}
export const Page: React.FC<PageProps> = (props: PageProps) => {
    let content: JSX.Element;
    if (props.header) {
        content = props.header;
    } else {
        content = (<p>{ props.title }</p>);
    }

    return (
        <div className="page">
            <div className="page-title">
                <div className="container">
                    { content }
                </div>
            </div>
            <div className="container">
                { props.children }
            </div>
        </div>
    );
};
