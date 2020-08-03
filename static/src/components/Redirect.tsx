import * as React from 'react';
import { Redirect as ReactRedirect } from 'react-router-dom';

export class Redirect {
    public static To(url: string): void {
        GlobalRedirectFrame.redirectTo(url);
    }
}

export interface GlobalRedirectFrameProps { }
interface GlobalRedirectFrameState {
    redirect?: JSX.Element;
}

export class GlobalRedirectFrame extends React.Component<GlobalRedirectFrameProps, GlobalRedirectFrameState> {
    constructor(props: GlobalRedirectFrameProps) {
        super(props);
        this.state = {};
        GlobalRedirectFrame.instance = this;
    }

    private static instance: GlobalRedirectFrame;

    public static redirectTo(url: string): void {
        this.instance.setState({
            redirect: ( <ReactRedirect push to={url} /> )
        }, () => {
            setTimeout(() => {
                this.instance.setState({ redirect: undefined });
            }, 20);
        });
    }

    render(): JSX.Element {
        return (
            <div id="global-redirect-frame">
                { this.state.redirect }
            </div>
        );
    }
}
