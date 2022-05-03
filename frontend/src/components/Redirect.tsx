import * as React from 'react';
import { Redirect as ReactRedirect } from 'react-router-dom';

export class Redirect {
    public static To(url: string): void {
        GlobalRedirectFrame.redirectTo(url);
    }
}

interface GlobalRedirectFrameState {
    redirect?: string;
}
export class GlobalRedirectFrame extends React.Component<unknown, GlobalRedirectFrameState> {
    constructor(props: unknown) {
        super(props);
        this.state = {};
        GlobalRedirectFrame.instance = this;
    }

    private static instance: GlobalRedirectFrame;

    public static redirectTo(url: string): void {
        this.instance.setState({
            redirect: url
        }, () => {
            setTimeout(() => {
                this.instance.setState({ redirect: undefined });
            }, 20);
        });
    }

    render(): JSX.Element {
        if (this.state.redirect) {
            return (
                <div id="global-redirect-frame">
                    <ReactRedirect to={this.state.redirect} push />
                </div>
            );
        }

        return (<div id="global-redirect-frame"></div>);
    }
}
