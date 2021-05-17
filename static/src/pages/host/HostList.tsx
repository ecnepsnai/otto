import * as React from 'react';
import { Host, HostType } from '../../types/Host';
import { PageLoading } from '../../components/Loading';
import { Page } from '../../components/Page';
import { Buttons, CreateButton } from '../../components/Button';
import { Table } from '../../components/Table';
import { HostListItem } from './HostListItem';
import { Heartbeat, HeartbeatType } from '../../types/Heartbeat';

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

    return (
        <Page title="Hosts">
            <Buttons>
                <CreateButton to="/hosts/host/" />
            </Buttons>
            <Table.Table>
                <Table.Head>
                    <Table.Column>Name</Table.Column>
                    <Table.Column>Address</Table.Column>
                    <Table.Column>Groups</Table.Column>
                    <Table.Column>Reachable</Table.Column>
                    <Table.Column>Version</Table.Column>
                    <Table.MenuColumn />
                </Table.Head>
                <Table.Body>
                    {
                        hosts.map(host => {
                            return <HostListItem host={host} heartbeat={heartbeats.get(host.Address)} key={host.ID} onReload={loadData}></HostListItem>;
                        })
                    }
                </Table.Body>
            </Table.Table>
        </Page>
    );
};
