import * as React from 'react';
import { Group } from '../../types/Group';
import { Link } from 'react-router-dom';
import { MenuItem, MenuLink } from '../../components/Menu';
import { Icon } from '../../components/Icon';
import { Table } from '../../components/Table';

export interface GroupListItemProps { group: Group, hosts: string[], onReload: () => (void), numGroups: number }
export class GroupListItem extends React.Component<GroupListItemProps, {}> {
    private deleteMenuClick = () => {
        this.props.group.DeleteModal().then(confirmed => {
            if (confirmed) {
                this.props.onReload();
            }
        });
    }

    render(): JSX.Element {
        const link = <Link to={'/groups/group/' + this.props.group.ID}>{ this.props.group.Name }</Link>;

        let deleteMenuItem: JSX.Element = null;
        if (this.props.numGroups > 1) {
            deleteMenuItem = (<MenuItem label="Delete" icon={<Icon.Delete />} onClick={this.deleteMenuClick}/>);
        }

        return (
            <Table.Row>
                <td>{ link }</td>
                <td>{ this.props.hosts.length }</td>
                <td>{ this.props.group.ScriptIDs.length }</td>
                <Table.Menu>
                    <MenuLink label="Edit" icon={<Icon.Edit />} to={'/groups/group/' + this.props.group.ID + '/edit'}/>
                    {deleteMenuItem}
                </Table.Menu>
            </Table.Row>
        );
    }
}
