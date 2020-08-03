import * as React from 'react';
import { Heartbeat } from '../../types/Heartbeat';
import { Host, ScriptEnabledGroup } from '../../types/Host';
import { Group } from '../../types/Group';
import { PageLoading } from '../../components/Loading';
import { Page } from '../../components/Page';
import { Buttons, EditButton, DeleteButton, SmallPlayButton } from '../../components/Button';
import { Layout } from '../../components/Layout';
import { match, Link } from 'react-router-dom';
import { URLParams } from '../../services/Params';
import { Card } from '../../components/Card';
import { ListGroup } from '../../components/ListGroup';
import { EnabledBadge } from '../../components/Badge';
import { EnvironmentVariableCard } from '../../components/EnvironmentVariableCard';
import { Icon } from '../../components/Icon';
import { Redirect } from '../../components/Redirect';
import { GlobalModalFrame } from '../../components/Modal';
import { RunModal } from '../run/RunModal';
import { Rand } from '../../services/Rand';

export interface HostViewProps { match: match }
interface HostViewState {
    loading: boolean;
    host?: Host;
    groups?: Group[];
    scripts?: ScriptEnabledGroup[];
    heartbeat?: Heartbeat;
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

    private loadData = () => {
        Promise.all([this.loadHost(), this.loadScripts(), this.loadGroups()]).then(() => {
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

    private runScriptClick = (scriptID: string) => {
        return () => {
            GlobalModalFrame.showModal(<RunModal scriptID={scriptID} hostIDs={[this.state.host.ID]} key={Rand.ID()}/>);
        };
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
                            <Card.Card>
                                <Card.Header>Host Information</Card.Header>
                                <ListGroup.List>
                                    <ListGroup.TextItem title="Name">{ this.state.host.Name }</ListGroup.TextItem>
                                    <ListGroup.TextItem title="Address">{ this.state.host.Address }:{ this.state.host.Port }</ListGroup.TextItem>
                                    <ListGroup.TextItem title="Enabled"><EnabledBadge value={this.state.host.Enabled} /></ListGroup.TextItem>
                                </ListGroup.List>
                            </Card.Card>
                        </Layout.Column>
                        <Layout.Column>
                            <Card.Card>
                                <Card.Header>Groups</Card.Header>
                                <ListGroup.List>
                                    {
                                        this.state.groups.map((group, index) => {
                                            return (
                                            <ListGroup.Item key={index}>
                                                <Icon.LayerGroup />
                                                <Link to={'/groups/group/' + group.ID} className="ml-1">{ group.Name }</Link>
                                            </ListGroup.Item>
                                            );
                                        })
                                    }
                                </ListGroup.List>
                            </Card.Card>
                        </Layout.Column>
                    </Layout.Row>
                    <Layout.Row>
                        <Layout.Column>
                            <EnvironmentVariableCard variables={this.state.host.Environment} />
                        </Layout.Column>
                        <Layout.Column>
                            <Card.Card>
                                <Card.Header>Script</Card.Header>
                                <ListGroup.List>
                                    {
                                        this.state.scripts.map((script, index) => {
                                            return (
                                            <ListGroup.Item key={index}>
                                                <div className="d-flex justify-content-between">
                                                    <div>
                                                        <Icon.Scroll />
                                                        <Link to={'/scripts/script/' + script.ScriptID} className="ml-1">{ script.ScriptName }</Link>
                                                    </div>
                                                    <div>
                                                        <SmallPlayButton onClick={this.runScriptClick(script.ScriptID)} />
                                                    </div>
                                                </div>
                                            </ListGroup.Item>
                                            );
                                        })
                                    }
                                </ListGroup.List>
                            </Card.Card>
                        </Layout.Column>
                    </Layout.Row>
                </Layout.Container>
            </Page>
        );
    }
}
