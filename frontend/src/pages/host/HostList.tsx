import * as React from 'react';
import { Host, HostType } from '../../types/Host';
import { PageLoading } from '../../components/Loading';
import { Page } from '../../components/Page';
import { CreateButton } from '../../components/Button';
import { Column, Table } from '../../components/Table';
import { Heartbeat, HeartbeatType } from '../../types/Heartbeat';
import { Link } from 'react-router-dom';
import { Formatter } from '../../services/Formatter';
import { DefaultSort } from '../../services/Sort';
import { HostTrust } from './HostTrust';
import { HeartbeatBadge } from '../../components/Badge';
import { AgentVersion } from '../../components/AgentVersion';
import { ContextMenuItem } from '../../components/ContextMenu';
import { Icon } from '../../components/Icon';

export const HostList: React.FC = () => {
    const [loading, setLoading] = React.useState(true);
    const [hosts, setHosts] = React.useState<HostType[]>();
    const [heartbeats, setHeartbeats] = React.useState<Map<string, HeartbeatType>>();

    React.useEffect(() => {
        loadData();
    }, []);

    const loadHosts = () => {
        return Host.List().then(hosts => {
            setHosts(hosts);
        });
    };

    const loadHeartbeats = () => {
        return Heartbeat.List().then(heartbeats => {
            setHeartbeats(heartbeats);
        });
    };

    const loadData = () => {
        Promise.all([loadHosts(), loadHeartbeats()]).then(() => {
            setLoading(false);
        });
    };

    if (loading) {
        return (<PageLoading />);
    }

    const toolbar = (
        <React.Fragment>
            <CreateButton to="/hosts/host/" />
        </React.Fragment>
    );

    const tableCols: Column[] = [
        {
            title: 'Name',
            value: (v: HostType) => {
                return (<Link to={'/hosts/host/' + v.ID}>{v.Name}</Link>);
            },
            sort: 'Name'
        },
        {
            title: 'Address',
            value: 'Address',
            sort: 'Address'
        },
        {
            title: 'Groups',
            value: (v: HostType) => {
                return (<span>{Formatter.ValueOrNothing(v.GroupIDs.length)}</span>);
            },
            sort: (asc: boolean, left: HostType, right: HostType) => {
                return DefaultSort(asc, (left.GroupIDs || []).length, (right.GroupIDs || []).length);
            }
        },
        {
            title: 'Trust',
            value: (v: HostType) => {
                return (<HostTrust host={v} badgeOnly outline />);
            },
        },
        {
            title: 'Status',
            value: (v: HostType) => {
                return (<HeartbeatBadge heartbeat={heartbeats.get(v.Address)} outline />);
            },
        },
        {
            title: 'Version',
            value: (v: HostType) => {
                return (<AgentVersion heartbeat={heartbeats.get(v.Address)} />);
            },
        },
    ];

    return (
        <Page title="Hosts" toolbar={toolbar}>
            <Table columns={tableCols} data={hosts} contextMenu={(a: HostType) => HostTableContextMenu(a, loadHosts)} defaultSort={{ ColumnIdx: 0, Ascending: true }} />
        </Page>
    );
};

const HostTableContextMenu = (host: HostType, onReload: () => (void)): (ContextMenuItem | 'separator')[] => {
    const deleteMenuClick = () => {
        Host.DeleteModal(host).then(confirmed => {
            if (confirmed) {
                onReload();
            }
        });
    };

    return [
        {
            title: 'Edit',
            icon: (<Icon.Edit />),
            href: '/hosts/host/' + host.ID + '/edit',
        },
        'separator',
        {
            title: 'Delete',
            icon: (<Icon.Delete />),
            onClick: deleteMenuClick,
        },
    ];
};
