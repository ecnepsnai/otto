import * as React from 'react';

import { Style } from './Style';

export interface AlertProps {
    color: Style.Palette;
    onClose?: () => (void);
}
export class Alert extends React.Component<AlertProps, {}> {
    private closeButton = () => {
        if (!this.props.onClose) {
            return null;
        }

        return (
        <button type="button" className="btn-close" data-dismiss="alert" aria-label="Close" onClick={this.props.onClose}></button>
        );
    }
    render(): JSX.Element {
        let className = 'alert fade show alert-' + this.props.color.toString();
        if (this.props.onClose) {
            className += ' alert-dismissible';
        }
        return (
            <div className={className} role="alert">
                { this.closeButton() }
                { this.props.children }
            </div>
        );
    }

    public static primary(content: string): JSX.Element {
        return <Alert color={Style.Palette.Primary}>{ content }</Alert>;
    }
    public static secondary(content: string): JSX.Element {
        return <Alert color={Style.Palette.Secondary}>{ content }</Alert>;
    }
    public static success(content: string): JSX.Element {
        return <Alert color={Style.Palette.Success}>{ content }</Alert>;
    }
    public static danger(content: string): JSX.Element {
        return <Alert color={Style.Palette.Danger}>{ content }</Alert>;
    }
    public static warning(content: string): JSX.Element {
        return <Alert color={Style.Palette.Warning}>{ content }</Alert>;
    }
    public static info(content: string): JSX.Element {
        return <Alert color={Style.Palette.Info}>{ content }</Alert>;
    }
    public static light(content: string): JSX.Element {
        return <Alert color={Style.Palette.Light}>{ content }</Alert>;
    }
    public static dark(content: string): JSX.Element {
        return <Alert color={Style.Palette.Dark}>{ content }</Alert>;
    }
}
