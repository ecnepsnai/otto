import * as React from 'react';
import { Heartbeat } from '../../types/Heartbeat';
import { Host, ScriptEnabledGroup } from '../../types/Host';
import { Group } from '../../types/Group';
import { PageLoading } from '../../components/Loading';
import { Page } from '../../components/Page';
import { Buttons, EditButton, DeleteButton } from '../../components/Button';
import { Layout } from '../../components/Layout';
import { match } from 'react-router-dom';
import { URLParams } from '../../services/Params';
import { Card } from '../../components/Card';
import { ListGroup } from '../../components/ListGroup';
import { EnabledBadge, HeartbeatBadge } from '../../components/Badge';
import { EnvironmentVariableCard } from '../../components/EnvironmentVariableCard';
import { Redirect } from '../../components/Redirect';
import { Schedule } from '../../types/Schedule';
import { GroupListCard } from '../../components/GroupListCard';
import { ScriptListCard } from '../../components/ScriptListCard';
import { ScheduleListCard } from '../../components/ScheduleListCard';
import { CopyButton } from '../../components/CopyButton';
import { DateLabel } from '../../components/DateLabel';
import { ClientVersion } from '../../components/ClientVersion';

export interface HostViewProps { match: match }
interface HostViewState {
    loading: boolean;
    host?: Host;
    heartbeat?: Heartbeat;
    groups?: Group[];
    scripts?: ScriptEnabledGroup[];
    schedules?: Schedule[];
}
export class HostView extends React.Component<HostViewProps, HostViewState> {
    private hostID: string;

    constructor(props: HostViewProps) {
        super(props);
        this.hostID = (this.props.match.params as URLParams).id;
        this.state = {
            loading: true,
        };
    }

    private loadHost = async () => {
        const host = await Host.Get(this.hostID);
        this.setState({
            host: host,
        });

        const heartbeats = await Heartbeat.List();
        heartbeats.forEach(heartbeat => {
            if (heartbeat.Address === host.Address) {
                this.setState({
                    heartbeat: heartbeat,
                });
            }
        });
    }

    private loadScripts = async () => {
        const scripts = await Host.Scripts(this.hostID);
        this.setState({
            scripts: scripts,
        });
    }

    private loadGroups = async () => {
        const groups = await Host.Groups(this.hostID);
        this.setState({
            groups: groups,
        });
    }

    private loadSchedules = async () => {
        const schedules = await Host.Schedules(this.hostID);
        this.setState({
            schedules: schedules,
        });
    }

    private loadData = () => {
        Promise.all([this.loadHost(), this.loadScripts(), this.loadGroups(), this.loadSchedules()]).then(() => {
            this.setState({
                loading: false,
            });
        });
    }

    componentDidMount(): void {
        this.loadData();
    }

    private deleteClick = () => {
        this.state.host.DeleteModal().then(deleted => {
            if (deleted) {
                Redirect.To('/hosts');
            }
        });
    }

    private lastReply = (): JSX.Element => {
        if (!this.state.heartbeat) { return null; }

        return (<ListGroup.TextItem title="Last Heartbeat"><DateLabel date={this.state.heartbeat.LastReply} /></ListGroup.TextItem>);
    }

    private clientVersion = (): JSX.Element => {
        if (!this.state.heartbeat) { return null; }

        return (<ListGroup.TextItem title="Client Version"><ClientVersion heartbeat={this.state.heartbeat} /></ListGroup.TextItem>);
    }

    render(): JSX.Element {
        if (this.state.loading) { return (<PageLoading />); }

        return (
            <Page title="View Host">
                <Layout.Container>
                    <Buttons>
                        <EditButton to={'/hosts/host/' + this.state.host.ID + '/edit'} />
                        <DeleteButton onClick={this.deleteClick} />
                    </Buttons>
                    <Layout.Row>
                        <Layout.Column>
                            <Card.Card className="mb-3">
                                <Card.Header>Host Information</Card.Header>
                                <ListGroup.List>
                                    <ListGroup.TextItem title="Name">{ this.state.host.Name }</ListGroup.TextItem>
                                    <ListGroup.TextItem title="Address">{ this.state.host.Address }:{ this.state.host.Port }</ListGroup.TextItem>
                                    <ListGroup.TextItem title="Enabled"><EnabledBadge value={this.state.host.Enabled} /></ListGroup.TextItem>
                                    <ListGroup.TextItem title="PSK"><code>*****</code> <CopyButton text={this.state.host.PSK} /></ListGroup.TextItem>
                                </ListGroup.List>
                            </Card.Card>
                            <Card.Card className="mb-3">
                                <Card.Header>Otto Client Information</Card.Header>
                                <ListGroup.List>
                                    <ListGroup.TextItem title="Status"><HeartbeatBadge heartbeat={this.state.heartbeat} /></ListGroup.TextItem>
                                    { this.lastReply() }
                                    { this.clientVersion() }
                                </ListGroup.List>
                            </Card.Card>
                            <EnvironmentVariableCard className="mb-3" variables={this.state.host.Environment} />
                            <ScheduleListCard schedules={this.state.schedules} className="mb-3" />
                        </Layout.Column>
                        <Layout.Column>
                            <GroupListCard groups={this.state.groups} className="mb-3"/>
                            <ScriptListCard scripts={this.state.scripts} hostIDs={[this.state.host.ID]} className="mb-3"/>
                        </Layout.Column>
                    </Layout.Row>
                </Layout.Container>
            </Page>
        );
    }
}
