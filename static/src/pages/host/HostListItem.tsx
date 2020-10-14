import * as React from 'react';

import { Host } from '../../types/Host';
import { Link } from 'react-router-dom';
import { Dropdown, MenuItem, MenuLink } from '../../components/Menu';
import { Style } from '../../components/Style';
import { Icon } from '../../components/Icon';
import { Table } from '../../components/Table';
import { HeartbeatBadge } from '../../components/Badge';
import { Heartbeat } from '../../types/Heartbeat';

export interface HostListItemProps { host: Host, heartbeat: Heartbeat, onReload: () => (void); }
export class HostListItem extends React.Component<HostListItemProps, {}> {
    private deleteMenuClick = () => {
        this.props.host.DeleteModal().then(confirmed => {
            if (confirmed) {
                this.props.onReload();
            }
        });
    }

    render(): JSX.Element {
        const link = <Link to={'/hosts/host/' + this.props.host.ID}>{ this.props.host.Name }</Link>;
        const dropdownLabel = <Icon.Bars />;
        const buttonProps = {
            color: Style.Palette.Secondary,
            outline: true,
            size: Style.Size.XS,
        };
        return (
            <Table.Row disabled={!this.props.host.Enabled}>
                <td>{ link }</td>
                <td>{ this.props.host.Address }</td>
                <td><HeartbeatBadge heartbeat={this.props.heartbeat} /></td>
                <td>
                    <Dropdown label={dropdownLabel} button={buttonProps}>
                        <MenuLink label="Edit" icon={<Icon.Edit />} to={'/hosts/host/' + this.props.host.ID +'/edit'}/>
                        <MenuItem label="Delete" icon={<Icon.Delete />} onClick={this.deleteMenuClick}/>
                    </Dropdown>
                </td>
            </Table.Row>
        );
    }
}
