import * as React from 'react';
import { Group, GroupType } from '../../types/Group';
import { Link } from 'react-router-dom';
import { Icon } from '../../components/Icon';
import { Table } from '../../components/Table';
import { Formatter } from '../../services/Formatter';
import { ContextMenuItem } from '../../components/ContextMenu';

interface GroupListItemProps {
    group: GroupType;
    hosts: string[];
    onReload: () => (void);
    numGroups: number;
}
export const GroupListItem: React.FC<GroupListItemProps> = (props: GroupListItemProps) => {
    const deleteMenuClick = () => {
        Group.DeleteModal(props.group).then(confirmed => {
            if (confirmed) {
                props.onReload();
            }
        });
    };

    const link = <Link to={'/groups/group/' + props.group.ID}>{props.group.Name}</Link>;

    const contextMenu: (ContextMenuItem | 'separator')[] = [
        {
            title: 'Edit',
            icon: (<Icon.Edit />),
            href: '/groups/group/' + props.group.ID + '/edit',
        },
    ];
    if (props.numGroups > 1) {
        contextMenu.push('separator');
        contextMenu.push({
            title: 'Delete',
            icon: (<Icon.Delete />),
            onClick: deleteMenuClick,
        });
    }

    return (
        <Table.Row menu={contextMenu}>
            <td>{link}</td>
            <td>{Formatter.ValueOrNothing(props.hosts.length)}</td>
            <td>{Formatter.ValueOrNothing((props.group.ScriptIDs || []).length)}</td>
        </Table.Row>
    );
};
