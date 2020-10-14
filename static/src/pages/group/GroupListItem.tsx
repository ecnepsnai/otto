import * as React from 'react';
import { Group } from '../../types/Group';
import { Link } from 'react-router-dom';
import { Dropdown, MenuItem, MenuLink } from '../../components/Menu';
import { Style } from '../../components/Style';
import { Icon } from '../../components/Icon';
import { Table } from '../../components/Table';

export interface GroupListItemProps { group: Group, hosts: string[], onReload: () => (void); }
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
        const dropdownLabel = <Icon.Bars />;
        const buttonProps = {
            color: Style.Palette.Secondary,
            outline: true,
            size: Style.Size.XS,
        };
        return (
            <Table.Row>
                <td>{ link }</td>
                <td>{ this.props.hosts.length }</td>
                <td>{ this.props.group.ScriptIDs.length }</td>
                <td>
                    <Dropdown label={dropdownLabel} button={buttonProps}>
                        <MenuLink label="Edit" icon={<Icon.Edit />} to={'/groups/group/' + this.props.group.ID + '/edit'}/>
                        <MenuItem label="Delete" icon={<Icon.Delete />} onClick={this.deleteMenuClick}/>
                    </Dropdown>
                </td>
            </Table.Row>
        );
    }
}
