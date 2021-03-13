import * as React from 'react';
import { Group, GroupType } from '../../types/Group';
import { Link } from 'react-router-dom';
import { Menu } from '../../components/Menu';
import { Icon } from '../../components/Icon';
import { Table } from '../../components/Table';
import { Formatter } from '../../services/Formatter';

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

    const link = <Link to={'/groups/group/' + props.group.ID}>{ props.group.Name }</Link>;

    let deleteMenuItem: JSX.Element = null;
    if (props.numGroups > 1) {
        deleteMenuItem = (
            <React.Fragment>
                <Menu.Divider />
                <Menu.Item label="Delete" icon={<Icon.Delete />} onClick={deleteMenuClick}/>
            </React.Fragment>
        );
    }

    return (
        <Table.Row>
            <td>{ link }</td>
            <td>{ Formatter.ValueOrNothing(props.hosts.length) }</td>
            <td>{ Formatter.ValueOrNothing((props.group.ScriptIDs || []).length) }</td>
            <Table.Menu>
                <Menu.Link label="Edit" icon={<Icon.Edit />} to={'/groups/group/' + props.group.ID + '/edit'}/>
                {deleteMenuItem}
            </Table.Menu>
        </Table.Row>
    );
};
