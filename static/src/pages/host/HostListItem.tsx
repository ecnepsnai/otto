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
export const HostListItem: React.FC<HostListItemProps> = (props: HostListItemProps) => {
    const deleteMenuClick = () => {
        Host.DeleteModal(props.host).then(confirmed => {
            if (confirmed) {
                props.onReload();
            }
        });
    };

    const link = <Link to={'/hosts/host/' + props.host.ID}>{props.host.Name}</Link>;
    return (
        <Table.Row disabled={!props.host.Enabled}>
            <td>{link}</td>
            <td>{props.host.Address}</td>
            <td>{Formatter.ValueOrNothing(props.host.GroupIDs.length)}</td>
            <td><HeartbeatBadge heartbeat={props.heartbeat} /></td>
            <td><ClientVersion heartbeat={props.heartbeat} /></td>
            <Table.Menu>
                <Menu.Link label="Edit" icon={<Icon.Edit />} to={'/hosts/host/' + props.host.ID + '/edit'} />
                <Menu.Divider />
                <Menu.Item label="Delete" icon={<Icon.Delete />} onClick={deleteMenuClick} />
            </Table.Menu>
        </Table.Row>
    );
};
