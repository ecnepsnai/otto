import * as React from 'react';
import { Host, HostType } from '../../types/Host';
import { Link } from 'react-router-dom';
import { Icon } from '../../components/Icon';
import { Table } from '../../components/Table';
import { HeartbeatBadge } from '../../components/Badge';
import { HeartbeatType } from '../../types/Heartbeat';
import { Formatter } from '../../services/Formatter';
import { ClientVersion } from '../../components/ClientVersion';
import { ContextMenuItem } from '../../components/ContextMenu';

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

    const contextMenu: (ContextMenuItem | 'separator')[] = [
        {
            title: 'Edit',
            icon: (<Icon.Edit />),
            href: '/hosts/host/' + props.host.ID + '/edit',
        },
        'separator',
        {
            title: 'Delete',
            icon: (<Icon.Delete />),
            onClick: deleteMenuClick,
        },
    ];

    const link = <Link to={'/hosts/host/' + props.host.ID}>{props.host.Name}</Link>;
    return (
        <Table.Row disabled={!props.host.Enabled} menu={contextMenu}>
            <td>{link}</td>
            <td>{props.host.Address}</td>
            <td>{Formatter.ValueOrNothing(props.host.GroupIDs.length)}</td>
            <td><HeartbeatBadge heartbeat={props.heartbeat} /></td>
            <td><ClientVersion heartbeat={props.heartbeat} /></td>
        </Table.Row>
    );
};
