import * as React from 'react';
import { Rand } from '../services/Rand';
import { ButtonProps, Button } from './Button';

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
                <Menu name={this.id}>
                    { this.props.children }
                </Menu>
            </div>
        );
    }
}

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

export interface MenuItemProps {
    icon?: JSX.Element;
    label: string;
    onClick: () => (void);
}

export class MenuItem extends React.Component<MenuItemProps, {}> {
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
