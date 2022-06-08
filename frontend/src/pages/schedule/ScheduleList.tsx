import * as React from 'react';
import { Schedule, ScheduleType } from '../../types/Schedule';
import { PageLoading } from '../../components/Loading';
import { Page } from '../../components/Page';
import { CreateButton } from '../../components/Button';
import { Column, Table } from '../../components/Table';
import { Script, ScriptType } from '../../types/Script';
import { Link } from 'react-router-dom';
import { EnabledBadge } from '../../components/Badge';
import { ContextMenuItem } from '../../components/ContextMenu';
import { DateLabel } from '../../components/DateLabel';
import { Icon } from '../../components/Icon';
import { DefaultSort } from '../../services/Sort';
import { SchedulePattern } from './SchedulePattern';

export const ScheduleList: React.FC = () => {
    const [loading, setLoading] = React.useState(true);
    const [schedules, setSchedules] = React.useState<ScheduleType[]>();
    const [scripts, setScripts] = React.useState<Map<string, ScriptType>>();

    React.useEffect(() => {
        loadData();
    }, []);

    const loadSchedules = () => {
        return Schedule.List().then(schedules => {
            setSchedules(schedules);
        });
    };

    const loadScripts = () => {
        return Script.List().then(scripts => {
            const m = new Map<string, ScriptType>();
            scripts.forEach(script => {
                m.set(script.ID, script);
            });
            setScripts(m);
        });
    };

    const loadData = () => {
        Promise.all([loadSchedules(), loadScripts()]).then(() => {
            setLoading(false);
        });
    };

    if (loading) {
        return (<PageLoading />);
    }

    const toolbar = (
        <React.Fragment>
            <CreateButton to="/schedules/schedule/" />
        </React.Fragment>
    );

    const tableCols: Column[] = [
        {
            title: 'Name',
            value: (v: ScheduleType) => {
                return (<Link to={'/schedules/schedule/' + v.ID}>{v.Name}</Link>);
            },
            sort: 'Name'
        },
        {
            title: 'Script',
            value: (v: ScheduleType) => {
                return (<span>{scripts.get(v.ScriptID).Name}</span>);
            },
            sort: (asc: boolean, left: ScheduleType, right: ScheduleType) => {
                return DefaultSort(asc, scripts.get(left.ScriptID).Name, scripts.get(right.ScriptID).Name);
            }
        },
        {
            title: 'Frequency',
            value: (v: ScheduleType) => {
                return (<SchedulePattern pattern={v.Pattern} />);
            },
        },
        {
            title: 'Scope',
            value: (v: ScheduleType) => {
                if (v.Scope.GroupIDs && v.Scope.GroupIDs.length > 0) {
                    let unit = 'groups';
                    if (v.Scope.GroupIDs.length == 1) {
                        unit = 'group';
                    }
                    return (<span>{v.Scope.GroupIDs.length} {unit}</span>);
                } else if (v.Scope.HostIDs && v.Scope.HostIDs.length > 0) {
                    let unit = 'hosts';
                    if (v.Scope.HostIDs.length == 1) {
                        unit = 'host';
                    }
                    return (<span>{v.Scope.HostIDs.length} {unit}</span>);
                }
        
                return (<span></span>);
            }
        },
        {
            title: 'Last Run',
            value: (v: ScheduleType) => {
                return (<DateLabel date={v.LastRunTime} />);
            }
        },
        {
            title: 'Status',
            value: (v: ScheduleType) => {
                return (<EnabledBadge value={v.Enabled} />);
            }
        },
    ];

    return (
        <Page title="Schedules" toolbar={toolbar}>
            <Table columns={tableCols} data={schedules} contextMenu={(a: ScheduleType) => ScheduleTableContextMenu(a, loadSchedules)} defaultSort={{ ColumnIdx: 0, Ascending: true }} />
        </Page>
    );
};

const ScheduleTableContextMenu = (schedule: ScheduleType, onReload: () => void): (ContextMenuItem | 'separator')[] => {
    const deleteMenuClick = () => {
        Schedule.DeleteModal(schedule).then(confirmed => {
            if (confirmed) {
                onReload();
            }
        });
    };

    return [
        {
            title: 'Edit',
            icon: (<Icon.Edit />),
            href: '/schedules/schedule/' + schedule.ID + '/edit',
        },
        'separator',
        {
            title: 'Delete',
            icon: (<Icon.Delete />),
            onClick: deleteMenuClick,
        },
    ];
};
