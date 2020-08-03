import * as React from 'react';
import { Group } from '../../types/Group';
import { match, Link } from 'react-router-dom';
import { URLParams } from '../../services/Params';
import { PageLoading } from '../../components/Loading';
import { Host } from '../../types/Host';
import { Script } from '../../types/Script';
import { Page } from '../../components/Page';
import { Layout } from '../../components/Layout';
import { Buttons, EditButton, DeleteButton, SmallPlayButton } from '../../components/Button';
import { Card } from '../../components/Card';
import { ListGroup } from '../../components/ListGroup';
import { Icon } from '../../components/Icon';
import { EnvironmentVariableCard } from '../../components/EnvironmentVariableCard';
import { Nothing } from '../../components/Nothing';
import { Redirect } from '../../components/Redirect';
import { GlobalModalFrame } from '../../components/Modal';
import { RunModal } from '../run/RunModal';
import { Rand } from '../../services/Rand';

export interface GroupViewProps {
    match: match;
}
interface GroupViewState {
    loading: boolean;
    group?: Group;
    hosts?: Host[];
    scripts?: Script[];
}
export class GroupView extends React.Component<GroupViewProps, GroupViewState> {
    private groupID: string;

    constructor(props: GroupViewProps) {
        super(props);
        this.groupID = (this.props.match.params as URLParams).id;
        this.state = {
            loading: true,
        };
    }

    private loadGroup = async () => {
        const group = await Group.Get(this.groupID);
        this.setState({
            group: group,
        });
    }

    private loadHosts = async () => {
        const hosts = await Group.Hosts(this.groupID);
        this.setState({
            hosts: hosts,
        });
    }

    private loadScripts = async () => {
        const scripts = await Group.Scripts(this.groupID);
        this.setState({
            scripts: scripts,
        });
    }

    private loadData = () => {
        Promise.all([this.loadGroup(), this.loadHosts(), this.loadScripts()]).then(() => {
            this.setState({
                loading: false,
            });
        });
    }

    componentDidMount(): void {
        this.loadData();
    }

    private deleteClick = () => {
        this.state.group.DeleteModal().then(deleted => {
            if (!deleted) {
                return;
            }

            Redirect.To('/groups');
        });
    }

    render(): JSX.Element {
        if (this.state.loading) { return (<PageLoading />); }

        return (
            <Page title="View Group">
                <Layout.Container>
                    <Buttons>
                        <EditButton to={'/groups/group/' + this.state.group.ID + '/edit'} />
                        <DeleteButton onClick={this.deleteClick} />
                    </Buttons>
                    <Layout.Row>
                        <Layout.Column>
                            <Card.Card>
                                <Card.Header>Host Information</Card.Header>
                                <ListGroup.List>
                                    <ListGroup.TextItem title="Name">{ this.state.group.Name }</ListGroup.TextItem>
                                </ListGroup.List>
                            </Card.Card>
                        </Layout.Column>
                        <Layout.Column>
                            <Card.Card>
                                <Card.Header>Hosts</Card.Header>
                                <HostListCard hosts={this.state.hosts} />
                            </Card.Card>
                        </Layout.Column>
                    </Layout.Row>
                    <Layout.Row>
                        <Layout.Column>
                            <EnvironmentVariableCard variables={this.state.group.Environment} />
                        </Layout.Column>
                        <Layout.Column>
                            <Card.Card>
                                <Card.Header>Script</Card.Header>
                                <ScriptListCard scripts={this.state.scripts} groupID={this.state.group.ID}/>
                            </Card.Card>
                        </Layout.Column>
                    </Layout.Row>
                </Layout.Container>
            </Page>
        );
    }
}

interface HostListCardProps {
    hosts: Host[];
}
class HostListCard extends React.Component<HostListCardProps, {}> {
    render(): JSX.Element {
        if (!this.props.hosts || this.props.hosts.length < 1) {
            return (
                <Card.Body>
                    <Nothing />
                </Card.Body>
            );
        }
        return (
            <ListGroup.List>
                {
                    this.props.hosts.map((host, index) => {
                        return (
                        <ListGroup.Item key={index}>
                            <Icon.Desktop />
                            <Link to={'/hosts/host/' + host.ID} className="ml-1">{ host.Name }</Link>
                        </ListGroup.Item>
                        );
                    })
                }
            </ListGroup.List>
        );
    }
}

interface ScriptListCardProps {
    groupID: string;
    scripts: Script[];
}
class ScriptListCard extends React.Component<ScriptListCardProps, {}> {
    private runScriptClick = (scriptID: string) => {
        return () => {
            Group.Hosts(this.props.groupID).then(hosts => {
                const hostIDs = hosts.map(host => { return host.ID; });
                GlobalModalFrame.showModal(<RunModal scriptID={scriptID} hostIDs={hostIDs} key={Rand.ID()}/>);
            });
        };
    }

    render(): JSX.Element {
        if (!this.props.scripts || this.props.scripts.length < 1) {
            return (
                <Card.Body>
                    <Nothing />
                </Card.Body>
            );
        }
        return (
            <ListGroup.List>
                {
                    this.props.scripts.map((script, index) => {
                        return (
                        <ListGroup.Item key={index}>
                            <div className="d-flex justify-content-between">
                                <div>
                                    <Icon.Scroll />
                                    <Link to={'/scripts/script/' + script.ID} className="ml-1">{ script.Name }</Link>
                                </div>
                                <div>
                                    <SmallPlayButton onClick={this.runScriptClick(script.ID)} />
                                </div>
                            </div>
                        </ListGroup.Item>
                        );
                    })
                }
            </ListGroup.List>
        );
    }
}
