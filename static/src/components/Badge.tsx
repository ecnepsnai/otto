import * as React from 'react';
import { HeartbeatType } from '../types/Heartbeat';
import { Icon } from './Icon';
import { Style } from './Style';
import '../../css/badge.scss';

interface BadgeProps {
    color: Style.Palette;
    pill?: boolean;
    outline?: boolean;
    className?: string;
    children?: React.ReactNode;
}
export const Badge: React.FC<BadgeProps> = (props: BadgeProps) => {
    let className = 'badge ';
    className += (props.outline ? 'badge-outline-' : 'badge-') + props.color.toString();
    if (props.pill) {
        className += ' rounded-pill';
    }
    if (props.className) {
        className += ' ' + props.className;
    }
    return (
        <div className={className}>{props.children}</div>
    );
};

interface EnabledBadgeProps {
    value: boolean;
    trueText?: string;
    falseText?: string;
}
export const EnabledBadge: React.FC<EnabledBadgeProps> = (props: EnabledBadgeProps) => {
    const color = (): Style.Palette => {
        if (props.value) {
            return Style.Palette.Success;
        }
        return Style.Palette.Danger;
    };
    const text = (): string => {
        if (props.value) {
            return props.trueText ?? 'Enabled';
        }
        return props.falseText ?? 'Disabled';
    };
    return (
        <Badge color={color()} pill>{text()}</Badge>
    );
};

interface HeartbeatBadgeProps {
    heartbeat: HeartbeatType;
    outline?: boolean;
}
export const HeartbeatBadge: React.FC<HeartbeatBadgeProps> = (props: HeartbeatBadgeProps) => {
    const color = (): Style.Palette => {
        if (!props.heartbeat) {
            return Style.Palette.Secondary;
        }

        if (props.heartbeat.IsReachable) {
            return Style.Palette.Success;
        }
        return Style.Palette.Danger;
    };

    const text = (): string => {
        if (!props.heartbeat) {
            return 'Unknown';
        }

        if (props.heartbeat.IsReachable) {
            return 'Reachable';
        }
        return 'Not Reachable';
    };

    const icon = (): JSX.Element => {
        if (!props.heartbeat) {
            return (<Icon.QuestionCircle />);
        }

        if (props.heartbeat.IsReachable) {
            return (<Icon.CheckCircle />);
        }
        return (<Icon.TimesCircle />);
    };

    return (
        <Badge color={color()} pill outline={props.outline}>
            <Icon.Label icon={icon()} label={text()} />
        </Badge>
    );
};
