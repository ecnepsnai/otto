import * as React from 'react';
import { Schedule, ScheduleType } from '../../types/Schedule';
import { Link } from 'react-router-dom';
import { Icon } from '../../components/Icon';
import { Table } from '../../components/Table';
import { ScriptType } from '../../types/Script';
import { EnabledBadge } from '../../components/Badge';
import { SchedulePattern } from './SchedulePattern';
import { DateLabel } from '../../components/DateLabel';
import { ContextMenuItem } from '../../components/ContextMenu';

interface ScheduleListItemProps {
    schedule: ScheduleType;
    script: ScriptType;
    onReload: () => (void);
}
export const ScheduleListItem: React.FC<ScheduleListItemProps> = (props: ScheduleListItemProps) => {
    const deleteMenuClick = () => {
        Schedule.DeleteModal(props.schedule).then(confirmed => {
            if (confirmed) {
                props.onReload();
            }
        });
    };

    const enabledOnColumn = () => {
        if (props.schedule.Scope.GroupIDs && props.schedule.Scope.GroupIDs.length > 0) {
            let unit = 'groups';
            if (props.schedule.Scope.GroupIDs.length == 1) {
                unit = 'group';
            }
            return (<td>{props.schedule.Scope.GroupIDs.length} {unit}</td>);
        } else if (props.schedule.Scope.HostIDs && props.schedule.Scope.HostIDs.length > 0) {
            let unit = 'hosts';
            if (props.schedule.Scope.HostIDs.length == 1) {
                unit = 'host';
            }
            return (<td>{props.schedule.Scope.HostIDs.length} {unit}</td>);
        }

        return (<td></td>);
    };

    const link = <Link to={'/schedules/schedule/' + props.schedule.ID}>{props.schedule.Name}</Link>;

    const contextMenu: (ContextMenuItem | 'separator')[] = [
        {
            title: 'Edit',
            icon: (<Icon.Edit />),
            href: '/schedules/schedule/' + props.schedule.ID + '/edit',
        },
        'separator',
        {
            title: 'Delete',
            icon: (<Icon.Delete />),
            onClick: deleteMenuClick,
        },
    ];

    return (
        <Table.Row menu={contextMenu}>
            <td>{link}</td>
            <td>{props.script.Name}</td>
            <td><SchedulePattern pattern={props.schedule.Pattern} /></td>
            {enabledOnColumn()}
            <td><DateLabel date={props.schedule.LastRunTime} /></td>
            <td><EnabledBadge value={props.schedule.Enabled} /></td>
        </Table.Row>
    );
};
