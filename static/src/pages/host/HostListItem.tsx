import * as React from 'react';

import { Host } from '../../types/Host';
import { Link } from 'react-router-dom';
import { MenuItem, MenuLink } from '../../components/Menu';
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
        return (
            <Table.Row disabled={!this.props.host.Enabled}>
                <td>{ link }</td>
                <td>{ this.props.host.Address }</td>
                <td><HeartbeatBadge heartbeat={this.props.heartbeat} /></td>
                <Table.Menu>
                    <MenuLink label="Edit" icon={<Icon.Edit />} to={'/hosts/host/' + this.props.host.ID +'/edit'}/>
                    <MenuItem label="Delete" icon={<Icon.Delete />} onClick={this.deleteMenuClick}/>
                </Table.Menu>
            </Table.Row>
        );
    }
}
