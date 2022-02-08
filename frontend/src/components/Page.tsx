import * as React from 'react';
import { Link } from 'react-router-dom';
import { Icon } from './Icon';
import '../../css/page.scss';

interface PageTitleBreadcrumb {
    title: string;
    href?: string;
}

interface PageProps {
    title?: string | PageTitleBreadcrumb[];
    loading?: boolean;
    toolbar?: JSX.Element;
    children?: React.ReactNode;
}
export const Page: React.FC<PageProps> = (props: PageProps) => {
    if (props.loading) {
        return (
            <div className="page">
                <div className="page-title">
                    <div className="container-fluid">
                        <p><Icon.Label icon={<Icon.Spinner pulse />} label="Please Wait..." /></p>
                    </div>
                </div>
            </div>
        );
    }

    let pageTitle: JSX.Element = undefined;
    if (props.title) {
        if (typeof props.title === 'object') {
            const breadcrumbs = props.title as PageTitleBreadcrumb[];
            pageTitle = (
                <div className="page-title page-title-breadcrumbs">
                    <div className="container-fluid">
                        {
                            breadcrumbs.map((crumb, idx) => {
                                let content = (<span>{crumb.title}</span>);
                                if (crumb.href) {
                                    content = (<span><Link to={crumb.href}>{crumb.title}</Link></span>);
                                }
                                let nextIcon: JSX.Element = undefined;
                                if (idx < breadcrumbs.length - 1) {
                                    nextIcon = (<Icon.ChevronRight />);
                                }
                                return (
                                    <React.Fragment key={idx}>
                                        {content}
                                        {nextIcon}
                                    </React.Fragment>
                                );
                            })
                        }
                    </div>
                </div>
            );
        } else {
            pageTitle = (
                <div className="page-title">
                    <div className="container-fluid">
                        <p>{props.title}</p>
                    </div>
                </div>
            );
        }
    }

    let toolbar: JSX.Element = undefined;
    if (props.toolbar) {
        toolbar = (
            <div className="page-toolbar">
                <div className="container-fluid buttons">
                    {props.toolbar}
                </div>
            </div>
        );
    }

    return (
        <div className="page">
            {pageTitle}
            {toolbar}
            <div className="container-fluid page-content">
                {props.children}
            </div>
        </div>
    );
};
