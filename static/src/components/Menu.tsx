import * as React from 'react';
import { Link as ReactLink } from 'react-router-dom';
import { Rand } from '../services/Rand';
import { ButtonProps, Button } from './Button';
import { Icon } from './Icon';

export interface DropdownProps {
    label: JSX.Element;
    button: ButtonProps;
}

export class Dropdown extends React.Component<DropdownProps> {
    private id = Rand.ID();
    constructor(props: DropdownProps) {
        super(props);
    }
    render(): JSX.Element {
        return (
            <div className="dropdown">
                <button className={Button.className(this.props.button)} type="button" id={this.id} data-toggle="dropdown" aria-expanded="false">
                    { this.props.label }
                </button>
                <Menu.Menu name={this.id}>
                    { this.props.children }
                </Menu.Menu>
            </div>
        );
    }
}

export namespace Menu {
    export interface MenuProps {
        name: string;
    }

    export class Menu extends React.Component<MenuProps> {
        constructor(props: MenuProps) {
            super(props);
        }
        render(): JSX.Element {
            return (
                <ul className="dropdown-menu" aria-labelledby={this.props.name}>{ this.props.children }</ul>
            );
        }
    }

    export interface ItemProps {
        icon?: JSX.Element;
        label: string;
        onClick: () => (void);
    }

    export class Item extends React.Component<ItemProps, {}> {
        private onClick = (event: React.MouseEvent<HTMLAnchorElement>) => {
            event.preventDefault();
            this.props.onClick();
        }
        render(): JSX.Element {
            return (
                <li><a className="dropdown-item" href="#" onClick={this.onClick}>{ this.props.icon }<span className="ml-1">{ this.props.label }</span></a></li>
            );
        }
    }

    export interface LinkProps {
        icon?: JSX.Element;
        label: string;
        to: string;
    }

    export class Link extends React.Component<LinkProps, {}> {
        render(): JSX.Element {
            return (
                <li>
                    <ReactLink to={this.props.to} className="dropdown-item">
                        <Icon.Label icon={this.props.icon} label={this.props.label} />
                    </ReactLink>
                </li>
            );
        }
    }

    export interface AnchorProps {
        icon?: JSX.Element;
        label: string;
        href: string;
    }

    export class Anchor extends React.Component<AnchorProps, {}> {
        render(): JSX.Element {
            return (
                <li><a className="dropdown-item" href={this.props.href}>{ this.props.icon }<span className="ml-1">{ this.props.label }</span></a></li>
            );
        }
    }

    export class Divider extends React.Component<{}, {}> {
        render(): JSX.Element {
            return (<li><hr className="dropdown-divider" /></li>);
        }
    }
}
