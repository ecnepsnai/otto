import * as React from 'react';
import { Style } from './Style';

export namespace Alert {
    interface AlertProps {
        color: Style.Palette;
        onClose?: () => (void);
        children?: React.ReactNode;
    }
    export const Alert: React.FC<AlertProps> = (props: AlertProps) => {
        const closeButton = () => {
            if (!props.onClose) {
                return null;
            }

            return (<button type="button" className="btn-close" data-dismiss="alert" aria-label="Close" onClick={props.onClose}></button>);
        };

        let className = 'alert fade show alert-' + props.color.toString();
        if (props.onClose) {
            className += ' alert-dismissible';
        }
        return (
            <div className={className} role="alert">
                { closeButton() }
                { props.children }
            </div>
        );
    };

    interface CommonAlertProps {
        onClose?: () => (void);
        children?: React.ReactNode;
    }
    export const Primary: React.FC<CommonAlertProps> = (props: CommonAlertProps) => Alert({ color: Style.Palette.Primary, onClose: props.onClose, children: props.children });
    export const Secondary: React.FC<CommonAlertProps> = (props: CommonAlertProps) => Alert({ color: Style.Palette.Secondary, onClose: props.onClose, children: props.children });
    export const Success: React.FC<CommonAlertProps> = (props: CommonAlertProps) => Alert({ color: Style.Palette.Success, onClose: props.onClose, children: props.children });
    export const Danger: React.FC<CommonAlertProps> = (props: CommonAlertProps) => Alert({ color: Style.Palette.Danger, onClose: props.onClose, children: props.children });
    export const Warning: React.FC<CommonAlertProps> = (props: CommonAlertProps) => Alert({ color: Style.Palette.Warning, onClose: props.onClose, children: props.children });
    export const Info: React.FC<CommonAlertProps> = (props: CommonAlertProps) => Alert({ color: Style.Palette.Info, onClose: props.onClose, children: props.children });
    export const Light: React.FC<CommonAlertProps> = (props: CommonAlertProps) => Alert({ color: Style.Palette.Light, onClose: props.onClose, children: props.children });
    export const Dark: React.FC<CommonAlertProps> = (props: CommonAlertProps) => Alert({ color: Style.Palette.Dark, onClose: props.onClose, children: props.children });
}
