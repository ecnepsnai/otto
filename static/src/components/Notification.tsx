import * as React from 'react';
import { Style } from './Style';
import { Rand } from '../services/Rand';
import '../../css/notification.scss';
import { Alert } from './Alert';

export interface NotificationProps {
    message: string;
    color: Style.Palette;
    id: string;
}

export class Notification extends React.Component<NotificationProps, {}> {
    private onClose = () => {
        GlobalNotificationFrame.removeNotification(this.props.id);
    }
    componentDidMount(): void {
        const duration = Math.max(this.props.message.length * 100, 1250);
        setTimeout(() => {
            this.onClose();
        }, duration);
    }
    render(): JSX.Element {
        return (
            <div className="notification">
                <Alert color={this.props.color} onClose={this.onClose}>
                    { this.props.message }
                </Alert>
            </div>
        );
    }

    /**
     * Post a green 'Success' notification
     * @param message The notification message
     */
    public static success(message: string): void {
        this.post(message, Style.Palette.Success);
    }

    /**
     * Post a blue 'Informational' notification
     * @param message The notification message
     */
    public static information(message: string): void {
        this.post(message, Style.Palette.Info);
    }

    /**
     * Post a yellow 'Warning' notification
     * @param message The notification message
     */
    public static warning(message: string): void {
        this.post(message, Style.Palette.Warning);
    }

    /**
     * Post a red 'Error' notification
     * @param message The notification message
     */
    public static error(message: string): void {
        this.post(message, Style.Palette.Danger);
    }

    private static post(message: string, color: Style.Palette) {
        const id = Rand.ID();
        GlobalNotificationFrame.addNotification(
            <Notification message={message} color={color} key={id} id={id}/>
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
