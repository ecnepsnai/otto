import * as React from 'react';
import { Group } from '../../types/Group';
import { match, Link } from 'react-router-dom';
import { URLParams } from '../../services/Params';
import { PageLoading } from '../../components/Loading';
import { Host } from '../../types/Host';
import { Script } from '../../types/Script';
import { Page } from '../../components/Page';
import { Layout } from '../../components/Layout';
import { Buttons, EditButton, DeleteButton } from '../../components/Button';
import { Card } from '../../components/Card';
import { ListGroup } from '../../components/ListGroup';
import { Icon } from '../../components/Icon';
import { EnvironmentVariableCard } from '../../components/EnvironmentVariableCard';
import { Nothing } from '../../components/Nothing';
import { Redirect } from '../../components/Redirect';
import { ScriptListCard } from '../../components/ScriptListCard';
import { Schedule } from '../../types/Schedule';
import { ScheduleListCard } from '../../components/ScheduleListCard';

export interface GroupViewProps {
    match: match;
}
interface GroupViewState {
    loading: boolean;
    group?: Group;
    hosts?: Host[];
    scripts?: Script[];
    schedules?: Schedule[];
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

    private loadSchedules = async () => {
        const schedules = await Group.Schedules(this.groupID);
        this.setState({
            schedules: schedules,
        });
    }

    private loadData = () => {
        Promise.all([this.loadGroup(), this.loadHosts(), this.loadScripts(), this.loadSchedules()]).then(() => {
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
                            <Card.Card className="mb-3">
                                <Card.Header>Host Information</Card.Header>
                                <ListGroup.List>
                                    <ListGroup.TextItem title="Name">{ this.state.group.Name }</ListGroup.TextItem>
                                </ListGroup.List>
                            </Card.Card>
                            <EnvironmentVariableCard variables={this.state.group.Environment} className="mb-3" />
                            <ScheduleListCard schedules={this.state.schedules} className="mb-3" />
                        </Layout.Column>
                        <Layout.Column>
                            <Card.Card className="mb-3">
                                <Card.Header>Hosts</Card.Header>
                                <HostListCard hosts={this.state.hosts} />
                            </Card.Card>
                            <ScriptListCard scripts={this.state.scripts} hostIDs={this.state.hosts.map(h => h.ID)} className="mb-3" />
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
                            <Link to={'/hosts/host/' + host.ID} className="ms-1">{ host.Name }</Link>
                        </ListGroup.Item>
                        );
                    })
                }
            </ListGroup.List>
        );
    }
}
