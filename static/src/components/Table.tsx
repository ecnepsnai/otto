import * as React from 'react';
import { Nothing } from './Nothing';
import { Dropdown } from './Menu';
import { Icon } from './Icon';
import '../../css/table.scss';

export namespace Table {
    interface TableProps {
        className?: string;
        children?: React.ReactNode;
    }
    export const Table: React.FC<TableProps> = (props: TableProps) => {
        let className = 'table';
        if (props.className) {
            className += ' ' + props.className;
        }
        return (
            <table className={className}>{props.children}</table>
        );
    };

    interface HeadProps {
        children?: React.ReactNode;
    }
    export const Head: React.FC<HeadProps> = (props: HeadProps) => {
        return (<thead className="table-thead"><tr>{props.children}</tr></thead>);
    };

    interface ColumnProps {
        className?: string;
        children?: React.ReactNode;
    }
    export const Column: React.FC<ColumnProps> = (props: ColumnProps) => {
        let className = 'table-th';
        if (props.className) {
            className += ' ' + props.className;
        }
        return (
            <th className={className}>{props.children}</th>
        );
    };
    export const MenuColumn: React.FC<ColumnProps> = (props: ColumnProps) => Column({ className: 'table-th-menu', children: props.children });

    interface BodyProps {
        children?: React.ReactNode;
    }
    export const Body: React.FC<BodyProps> = (props: BodyProps) => {
        let content = props.children;
        if (!React.Children.count(props.children)) {
            content = (<tr><td colSpan={10}><Nothing /></td></tr>);
        }

        return (
            <tbody>{content}</tbody>
        );
    };

    interface RowProps {
        disabled?: boolean;
        children?: React.ReactNode;
    }
    export const Row: React.FC<RowProps> = (props: RowProps) => {
        let className = 'table-tr';
        if (props.disabled) {
            className += ' disabled';
        }
        return (
            <tr className={className}>{props.children}</tr>
        );
    };

    interface MenuProps {
        disabled?: boolean;
        children?: React.ReactNode;
    }
    export const Menu: React.FC<MenuProps> = (props: MenuProps) => {
        return (<td>
            <Dropdown label={<Icon.Bars />}>
                {props.children}
            </Dropdown>
        </td>);
    };
}
