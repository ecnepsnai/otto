import * as React from 'react';
import { ContextMenuHandler } from '../services/ContextMenuHandler';
import { Rand } from '../services/Rand';
import { DefaultSort } from '../services/Sort';
import { ContextMenu, ContextMenuItem, GlobalContextMenuFrame } from './ContextMenu';
import { Icon } from './Icon';
import { Nothing } from './Nothing';
import { Style } from './Style';
import '../../css/table.scss';

export interface Column {
    /** The title of the header */
    title: string;
    /** The value for the cell of this column for each row of data.
     * Can either be a function that returns JSX, or a string which defined the property for a simple object. */
    value: (string | ((v: unknown) => JSX.Element));
    /** The sort option for this column. Can either be a property of the object to do a basic comparison on,
     * or a function that will be called with each comparable object for a custom sort.
     * If no value is provided then sorting is disabled for this column. */
    sort?: (string | ((asc: boolean, left: unknown, right: unknown) => number));
}

export interface SortProps {
    ColumnIdx: number;
    Ascending?: boolean;
}

interface TableProps {
    /** The columns of this table. Cannot be changed once the table is rendered. */
    columns: Column[];
    /** The data of the table. Can be changed after the table is rendered. */
    data: unknown[];
    /** The context menu for each row of the table. Called while the table is being rendered. */
    contextMenu?: (v: unknown) => (ContextMenuItem | 'separator')[];
    /** The default sort settings for the table. */
    defaultSort?: SortProps;
    /** A menu to appear above the table */
    menu?: JSX.Element;
}
export const Table: React.FC<TableProps> = (props: TableProps) => {
    const [VirtualData, SetVirtualData] = React.useState<unknown[]>([]);
    const [Sort, SetSort] = React.useState<SortProps>(props.defaultSort);
    const [TableKey, SetTableKey] = React.useState(Rand.ID());

    React.useEffect(() => {
        SetVirtualData(props.data);
    }, []);

    React.useEffect(() => {
        SetVirtualData(props.data);
        SetSort(s => {
            if (props.defaultSort === s) {
                return s;
            }
            return null;
        });
    }, [props.data]);

    React.useEffect(() => {
        if (!Sort) {
            return;
        }

        SetVirtualData(d => {
            return d.sort((a: unknown, b: unknown) => {
                const col = props.columns[Sort.ColumnIdx];
                if (typeof col.sort === 'string') {
                    return DefaultSort(Sort.Ascending, unknownGet(a, col.sort), unknownGet(b, col.sort));
                }

                return col.sort(Sort.Ascending, a, b);
            });
        });

        SetTableKey(Rand.ID());
    }, [Sort]);

    const headClick = (idx: number) => {
        return () => {
            SetSort(s => {
                if (s && s.ColumnIdx === idx) {
                    s.Ascending = !s.Ascending;
                    return { ...s };
                }

                return {
                    ColumnIdx: idx,
                    Ascending: true
                };
            });
        };
    };

    const thead = () => {
        return (
            <thead>
                <tr>
                    {
                        props.columns.map((col, idx) => {
                            if (!col.sort) {
                                return (<th key={idx} className="table-th">{col.title}</th>);
                            }

                            let title = (<span>{col.title}</span>);
                            if (Sort && Sort.ColumnIdx === idx) {
                                const icon = Sort.Ascending ? (<Icon.CaretUp />) : (<Icon.CaretDown />);
                                title = (<span>{col.title}<span className="ms-1">{icon}</span></span>);
                            }

                            return (<th key={idx} onClick={headClick(idx)} className="table-th sortable">{title}</th>);
                        })
                    }
                </tr>
            </thead>
        );
    };
    const tbody = () => {
        if (!VirtualData || VirtualData.length === 0) {
            return (
                <tbody>
                    <tr>
                        <td colSpan={props.columns.length}><Nothing /></td>
                    </tr>
                </tbody>
            );
        }

        return (
            <tbody>
                {
                    VirtualData.map((obj, idx) => {
                        return (<TableData key={idx} columns={props.columns} data={obj} contextMenu={props.contextMenu ? props.contextMenu(obj) : null} />);
                    })
                }
            </tbody>
        );
    };

    const table = () => {
        return (
            <div className="table-responsive">
                <table className="table" key={TableKey}>
                    {thead()}
                    {tbody()}
                </table>
            </div>
        );
    };

    if (props.menu) {
        return (
            <div className="table-wrapper">
                <div className="table-wrapper-menu">{props.menu}</div>
                {table()}
            </div>
        );
    }

    return (
        <React.Fragment>{table()}</React.Fragment>
    );
};

interface TableDataProps {
    columns: Column[];
    data: unknown;
    contextMenu?: (ContextMenuItem | 'separator')[];
}
const TableData: React.FC<TableDataProps> = (props: TableDataProps) => {
    const contextMenuHandler = new ContextMenuHandler((x: number, y: number) => {
        if (!props.contextMenu) {
            return;
        }

        GlobalContextMenuFrame.showMenu(<ContextMenu x={x} y={y} items={props.contextMenu} />);
    });

    const value = (column: Column) => {
        if (typeof column.value === 'string') {
            return (
                <React.Fragment>{unknownGet(props.data, column.value) as React.ReactNode}</React.Fragment>
            );
        }

        let element: JSX.Element;
        try {
            element = column.value(props.data);
        } catch (ex) {
            console.error('Error rendering cell value', { column: column, data: props.data, error: ex });
            element = (<Icon.Label icon={<Icon.ExclamationCircle color={Style.Palette.Danger} />} label="##ERROR##" />);
        }

        return element;
    };

    return (
        <tr
            className="table-tr"
            onContextMenu={(contextMenuHandler ? contextMenuHandler.onContextMenu : null)}
            onTouchStart={(contextMenuHandler ? contextMenuHandler.onTouchStart : null)}
            onTouchEnd={(contextMenuHandler ? contextMenuHandler.onTouchEnd : null)}
            onTouchMove={(contextMenuHandler ? contextMenuHandler.onTouchMove : null)}>
            {
                props.columns.map((col, colIdx) => {
                    return (<td key={colIdx}>
                        {value(col)}
                    </td>);
                })
            }
        </tr>
    );
};

// eslint-disable-next-line @typescript-eslint/no-explicit-any
const unknownGet = (obj: any, key: string): unknown => {
    return obj[key];
};
