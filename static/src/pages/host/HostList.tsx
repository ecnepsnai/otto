import * as React from 'react';
import { Host } from '../../types/Host';
import { PageLoading } from '../../components/Loading';
import { Page } from '../../components/Page';
import { Buttons, CreateButton } from '../../components/Button';
import { Table } from '../../components/Table';
import { HostListItem } from './HostListItem';
import { Heartbeat } from '../../types/Heartbeat';

export interface HostListProps {}
interface HostListState {
    loading: boolean;
    hosts: Host[];
    heartbeats?: Map<string, Heartbeat>;
}
export class HostList extends React.Component<HostListProps, HostListState> {
    constructor(props: HostListProps) {
        super(props);
        this.state = {
            loading: true,
            hosts: [],
        };
    }

    componentDidMount(): void {
        this.loadData();
    }

    private loadHosts = () => {
        return Host.List().then(hosts => {
            this.setState({
                hosts: hosts,
            });
        });
    }

    private loadHeartbeats = () => {
        return Heartbeat.List().then(heartbeats => {
            this.setState({
                heartbeats: heartbeats,
            });
        });
    }

    private loadData = () => {
        Promise.all([this.loadHosts(), this.loadHeartbeats()]).then(() => {
            this.setState({ loading: false });
        });
    }

    render(): JSX.Element {
        if (this.state.loading) { return ( <PageLoading /> ); }

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
                        <Table.MenuColumn />
                    </Table.Head>
                    <Table.Body>
                        {
                            this.state.hosts.map(host => {
                                return <HostListItem host={host} heartbeat={this.state.heartbeats.get(host.Address)} key={host.ID} onReload={this.loadData}></HostListItem>;
                            })
                        }
                    </Table.Body>
                </Table.Table>
            </Page>
        );
    }
}
