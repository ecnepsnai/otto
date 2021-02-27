import * as React from 'react';
import { Host, HostType } from '../../types/Host';
import { Link } from 'react-router-dom';
import { Menu } from '../../components/Menu';
import { Icon } from '../../components/Icon';
import { Table } from '../../components/Table';
import { HeartbeatBadge } from '../../components/Badge';
import { HeartbeatType } from '../../types/Heartbeat';
import { Formatter } from '../../services/Formatter';
import { ClientVersion } from '../../components/ClientVersion';

interface HostListItemProps {
    host: HostType;
    heartbeat: HeartbeatType;
    onReload: () => (void);
}
export class HostListItem extends React.Component<HostListItemProps, unknown> {
    private deleteMenuClick = () => {
        Host.DeleteModal(this.props.host).then(confirmed => {
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
                <td>{ Formatter.ValueOrNothing(this.props.host.GroupIDs.length) }</td>
                <td><HeartbeatBadge heartbeat={this.props.heartbeat}/></td>
                <td><ClientVersion heartbeat={this.props.heartbeat} /></td>
                <Table.Menu>
                    <Menu.Link label="Edit" icon={<Icon.Edit />} to={'/hosts/host/' + this.props.host.ID +'/edit'}/>
                    <Menu.Divider />
                    <Menu.Item label="Delete" icon={<Icon.Delete />} onClick={this.deleteMenuClick}/>
                </Table.Menu>
            </Table.Row>
        );
    }
}
