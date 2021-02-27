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
export class GroupListItem extends React.Component<GroupListItemProps, unknown> {
    private deleteMenuClick = () => {
        Group.DeleteModal(this.props.group).then(confirmed => {
            if (confirmed) {
                this.props.onReload();
            }
        });
    }

    render(): JSX.Element {
        const link = <Link to={'/groups/group/' + this.props.group.ID}>{ this.props.group.Name }</Link>;

        let deleteMenuItem: JSX.Element = null;
        if (this.props.numGroups > 1) {
            deleteMenuItem = (
                <React.Fragment>
                    <Menu.Divider />
                    <Menu.Item label="Delete" icon={<Icon.Delete />} onClick={this.deleteMenuClick}/>
                </React.Fragment>
            );
        }

        return (
            <Table.Row>
                <td>{ link }</td>
                <td>{ Formatter.ValueOrNothing(this.props.hosts.length) }</td>
                <td>{ Formatter.ValueOrNothing(this.props.group.ScriptIDs.length) }</td>
                <Table.Menu>
                    <Menu.Link label="Edit" icon={<Icon.Edit />} to={'/groups/group/' + this.props.group.ID + '/edit'}/>
                    {deleteMenuItem}
                </Table.Menu>
            </Table.Row>
        );
    }
}
