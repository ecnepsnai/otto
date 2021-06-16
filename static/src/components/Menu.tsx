import * as React from 'react';
import { Link as ReactLink } from 'react-router-dom';
import { Rand } from '../services/Rand';
import { Icon } from './Icon';

interface DropdownProps {
    label: JSX.Element;
    children?: React.ReactNode;
}
export const Dropdown: React.FC<DropdownProps> = (props: DropdownProps) => {
    const id = Rand.ID();
    return (
        <span className="dropdown">
            <button className="btn btn-outline-secondary btn-xs" type="button" role="button" id={id} data-bs-toggle="dropdown" aria-expanded="false">
                {props.label}
            </button>
            <Menu.Menu name={id}>
                {props.children}
            </Menu.Menu>
        </span>
    );
};

export namespace Menu {
    interface MenuProps {
        name: string;
        children?: React.ReactNode;
    }
    export const Menu: React.FC<MenuProps> = (props: MenuProps) => {
        return (
            <ul className="dropdown-menu" aria-labelledby={props.name}>{props.children}</ul>
        );
    };

    interface ItemProps {
        icon?: JSX.Element;
        label: string;
        onClick: () => (void);
    }
    export const Item: React.FC<ItemProps> = (props: ItemProps) => {
        const onClick = (event: React.MouseEvent<HTMLAnchorElement>) => {
            event.preventDefault();
            props.onClick();
        };
        return (
            <li><a className="dropdown-item" href="#" onClick={onClick}>{props.icon}<span className="ms-1">{props.label}</span></a></li>
        );
    };

    interface LinkProps {
        icon?: JSX.Element;
        label: string;
        to: string;
    }
    export const Link: React.FC<LinkProps> = (props: LinkProps) => {
        return (
            <li>
                <ReactLink to={props.to} className="dropdown-item">
                    <Icon.Label icon={props.icon} label={props.label} />
                </ReactLink>
            </li>
        );
    };

    interface AnchorProps {
        icon?: JSX.Element;
        label: string;
        href: string;
    }
    export const Anchor: React.FC<AnchorProps> = (props: AnchorProps) => {
        return (
            <li><a className="dropdown-item" href={props.href}>{props.icon}<span className="ms-1">{props.label}</span></a></li>
        );
    };

    export const Divider: React.FC = () => {
        return (<li><hr className="dropdown-divider" /></li>);
    };
}
