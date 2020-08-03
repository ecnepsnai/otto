import * as React from 'react';
import { Style } from './Style';
import { Bootstrap, BSModule } from '../services/Bootstrap';
import { Rand } from '../services/Rand';
import '../../css/notify.scss';
import { Icon } from './Icon';

export interface NotificationProps {
    title: string;
    body: string;
    color: Style.Palette;
    id?: string;
}

interface NotificationState {
    id: string;
    bsToast?: BSModule;
}

export class Notification extends React.Component<NotificationProps, NotificationState> {
    constructor(props: NotificationProps) {
        super(props);
        this.state = { id: props.id ?? Rand.ID() };
    }
    componentDidMount(): void {
        const element = document.getElementById(this.state.id);
        const bsm = Bootstrap.Toast(element, { animation: true, autohide: true, delay: 2000});
        bsm.show();
        const id = this.props.id;
        element.addEventListener('hidden.bs.toast', function () {
            GlobalNotificationFrame.removeNotification(id);
        });
        this.setState({ bsToast: bsm });
    }
    render(): JSX.Element {
        return (
            <div className="toast" id={this.state.id} role="alert" aria-live="assertive" aria-atomic="true">
                <div className="toast-header">
                    <NotificationIcon color={this.props.color}></NotificationIcon>
                    <strong className="mr-auto">{ this.props.title }</strong>
                    <button type="button" className="ml-2 mb-1 close" data-dismiss="toast" aria-label="Close">
                        <span aria-hidden="true">&times;</span>
                    </button>
                </div>
                <div className="toast-body">{ this.props.body }</div>
            </div>
        );
    }

    /**
     * Post a green 'Success' notification
     * @param title A short title for the notification
     * @param body The body message for the notification
     */
    public static success(title: string, body: string): void {
        this.post(title, body, Style.Palette.Success);
    }

    /**
     * Post a blue 'Informational' notification
     * @param title A short title for the notification
     * @param body The body message for the notification
     */
    public static information(title: string, body: string): void {
        this.post(title, body, Style.Palette.Info);
    }

    /**
     * Post a yellow 'Warning' notification
     * @param title A short title for the notification
     * @param body The body message for the notification
     */
    public static warning(title: string, body: string): void {
        this.post(title, body, Style.Palette.Warning);
    }

    /**
     * Post a red 'Error' notification
     * @param title A short title for the notification
     * @param body The body message for the notification
     */
    public static error(title: string, body: string): void {
        this.post(title, body, Style.Palette.Danger);
    }

    private static post(title: string, body: string, color: Style.Palette) {
        const id = Rand.ID();
        GlobalNotificationFrame.addNotification(
            <Notification title={title} body={body} color={color} key={id} id={id}/>
        );
    }
}

interface NotificationIconProps {
    color: Style.Palette;
}
class NotificationIcon extends React.Component<NotificationIconProps, {}> {
    render(): JSX.Element {
        const className = 'notification-icon bg-' + this.props.color.toString();
        let icon: JSX.Element;

        if (this.props.color == Style.Palette.Success) {
            icon = ( <Icon.CheckCircle color={Style.Palette.Light} /> );
        } else if (this.props.color == Style.Palette.Info) {
            icon = ( <Icon.InfoCircle color={Style.Palette.Light} /> );
        } else if (this.props.color == Style.Palette.Warning) {
            icon = ( <Icon.ExclamationCircle color={Style.Palette.Light} /> );
        } else if (this.props.color == Style.Palette.Danger) {
            icon = ( <Icon.TimesCircle color={Style.Palette.Light} /> );
        }

        return (
            <div className={className}>{ icon }</div>
        );
    }
}

export interface GlobalNotificationFrameProps { }
interface GlobalNotificationFrameState {
    notifications: JSX.Element[];
}

export class GlobalNotificationFrame extends React.Component<GlobalNotificationFrameProps, GlobalNotificationFrameState> {
    constructor(props: GlobalNotificationFrameProps) {
        super(props);
        this.state = { notifications: [] };
        GlobalNotificationFrame.instance = this;
    }

    private static instance: GlobalNotificationFrame;

    public static addNotification(notification: JSX.Element): void {
        this.instance.setState(state => {
            const notifications = state.notifications;
            notifications.push(notification);
            return { notifications: notifications};
        });
    }

    public static removeNotification(key: string): void {
        this.instance.setState(state => {
            const notifications = state.notifications;
            return { notifications: notifications.filter(n => {
                return n.key !== key;
            })};
        });
    }

    render(): JSX.Element {
        return (
            <div id="global-notification-frame">
                {
                    this.state.notifications
                }
            </div>
        );
    }
}
