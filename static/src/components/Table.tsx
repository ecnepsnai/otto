import * as React from 'react';
import { Nothing } from './Nothing';
import { ContextMenu, ContextMenuItem, GlobalContextMenuFrame } from './ContextMenu';
import { ContextMenuHandler } from '../services/ContextMenuHandler';
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
            <div className="table-responsive">
                <table className={className}>{props.children}</table>
            </div>
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
        menu?: (ContextMenuItem | 'separator')[];
    }
    export const Row: React.FC<RowProps> = (props: RowProps) => {
        const contextMenuHandler = new ContextMenuHandler((x: number, y: number) => {
            if (!props.menu) {
                return;
            }

            GlobalContextMenuFrame.showMenu(<ContextMenu x={x} y={y} items={props.menu} />);
        });

        let className = 'table-tr';
        if (props.disabled) {
            className += ' disabled';
        }

        return (
            <tr
                className={className}
                onContextMenu={contextMenuHandler.onContextMenu}
                onTouchStart={contextMenuHandler.onTouchStart}
                onTouchEnd={contextMenuHandler.onTouchEnd}
                onTouchMove={contextMenuHandler.onTouchMove}
            >{props.children}</tr>
        );
    };
}
