import * as React from 'react';
import '../../css/table.scss';
import { Icon } from './Icon';
import { Dropdown } from './Menu';
import { Nothing } from './Nothing';
import { Style } from './Style';

export namespace Table {
    export interface TableProps { className?: string; }
    export class Table extends React.Component<TableProps, {}> {
        render(): JSX.Element {
            let className = 'table';
            if (this.props.className) {
                className += ' ' + this.props.className;
            }
            return (
                <div className="table-responsive">
                    <table className={className}>{this.props.children}</table>
                </div>
            );
        }
    }

    export class Head extends React.Component<{}, {}> {
        render(): JSX.Element {
            return (
                <thead className="table-thead"><tr>{this.props.children}</tr></thead>
            );
        }
    }

    export interface ColumnProps { className?: string; }
    export class Column extends React.Component<ColumnProps, {}> {
        render(): JSX.Element {
            let className = 'table-th';
            if (this.props.className) {
                className += ' ' + this.props.className;
            }
            return (
                <th className={className}>{this.props.children}</th>
            );
        }
    }

    export class MenuColumn extends React.Component<{}, {}> {
        render(): JSX.Element {
            return (
                <Column className=" table-th-menu"></Column>
            );
        }
    }

    export class Body extends React.Component<{}, {}> {
        render(): JSX.Element {
            let content = this.props.children;
            if (!React.Children.count(this.props.children)) {
                content = ( <Nothing /> );
            }

            return (
                <tbody>{ content }</tbody>
            );
        }
    }

    export interface RowProps { disabled?: boolean; }
    export class Row extends React.Component<RowProps, {}> {
        render(): JSX.Element {
            let className = 'table-tr';
            if (this.props.disabled) {
                className += ' disabled';
            }
            return (
                <tr className={className}>{this.props.children}</tr>
            );
        }
    }

    export interface MenuProps { disabled?: boolean; }
    export class Menu extends React.Component<MenuProps, {}> {
        render(): JSX.Element {
            const buttonProps = {
                color: Style.Palette.Secondary,
                outline: true,
                size: Style.Size.XS,
            };
            return (<td>
                <Dropdown label={<Icon.Bars />} button={buttonProps}>
                    {this.props.children}
                </Dropdown>
            </td>);
        }
    }
}
