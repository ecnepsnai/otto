import * as React from 'react';
import { Script, ScriptEnabledHost } from '../../types/Script';
import { match, Link } from 'react-router-dom';
import { PageLoading } from '../../components/Loading';
import { Page } from '../../components/Page';
import { Layout } from '../../components/Layout';
import { Card } from '../../components/Card';
import { URLParams } from '../../services/Params';
import { Buttons, EditButton, DeleteButton, Button, SmallPlayButton } from '../../components/Button';
import { ListGroup } from '../../components/ListGroup';
import { EnabledBadge } from '../../components/Badge';
import { Icon } from '../../components/Icon';
import { EnvironmentVariableCard } from '../../components/EnvironmentVariableCard';
import { Redirect } from '../../components/Redirect';
import { Style } from '../../components/Style';
import { GlobalModalFrame } from '../../components/Modal';
import { RunModal } from '../run/RunModal';
import { Rand } from '../../services/Rand';
import { Group } from '../../types/Group';
import { Pre } from '../../components/Pre';
import { Attachment } from '../../types/Attachment';
import { Formatter } from '../../services/Formatter';
import { Nothing } from '../../components/Nothing';
import { Schedule } from '../../types/Schedule';
import { ScheduleListCard } from '../../components/ScheduleListCard';

export interface ScriptViewProps { match: match; }
interface ScriptViewState {
    loading: boolean;
    script?: Script;
    hosts?: ScriptEnabledHost[];
    attachments?: Attachment[];
    schedules?: Schedule[];
}
export class ScriptView extends React.Component<ScriptViewProps, ScriptViewState> {
    private scriptID: string;

    constructor(props: ScriptViewProps) {
        super(props);
        this.scriptID = (this.props.match.params as URLParams).id;
        this.state = {
            loading: true,
        };
    }

    private loadHosts = async () => {
        const hosts = await Script.Hosts(this.scriptID);
        this.setState({
            hosts: hosts,
        });
    }

    private loadScript = async () => {
        const script = await Script.Get(this.scriptID);
        this.setState({
            script: script,
        });
    }

    private loadAttachments = async () => {
        const attachments = await Script.Attachments(this.scriptID);
        this.setState({
            attachments: attachments,
        });
    }

    private loadSchedules = async () => {
        const schedules = await Script.Schedules(this.scriptID);
        this.setState({
            schedules: schedules,
        });
    }

    private loadData = () => {
        Promise.all([this.loadHosts(), this.loadScript(), this.loadAttachments(), this.loadSchedules()]).then(() => {
            this.setState({
                loading: false,
            });
        });
    }

    componentDidMount(): void {
        this.loadData();
    }

    private deleteClick = () => {
        this.state.script.DeleteModal().then(deleted => {
            if (!deleted) {
                return;
            }

            Redirect.To('/scripts');
        });
    }

    private executeClick = () => {
        GlobalModalFrame.showModal(<RunModal scriptID={this.state.script.ID} key={Rand.ID()}/>);
    }

    private runScriptGroupClick = (groupID: string) => {
        return () => {
            Group.Hosts(groupID).then(hosts => {
                const hostIDs = hosts.map(host => { return host.ID; });
                GlobalModalFrame.showModal(<RunModal scriptID={this.state.script.ID} hostIDs={hostIDs} key={Rand.ID()}/>);
            });
        };
    }

    private runScriptHostClick = (hostID: string) => {
        return () => {
            GlobalModalFrame.showModal(<RunModal scriptID={this.state.script.ID} hostIDs={[hostID]} key={Rand.ID()}/>);
        };
    }

    private runAs = () => {
        if (this.state.script.RunAs.Inherit) { return null; }

        return (<ListGroup.TextItem title="Run As">User: {this.state.script.RunAs.UID} Group: {this.state.script.RunAs.GID}</ListGroup.TextItem>);
    }

    private attachmentList = () => {
        if (!this.state.attachments || this.state.attachments.length == 0) {
            return (<Card.Body><Nothing /></Card.Body>);
        }

        return (<ListGroup.List>{
            this.state.attachments.map((attachment, idx) => {
                return (<ListGroup.Item key={idx}>
                    <div className="d-flex justify-content-between">
                        <span>
                            <Icon.Label icon={<Icon.Paperclip />} label={attachment.Name} />
                            <span className="text-muted ms-1">
                                {Formatter.Bytes(attachment.Size)}
                            </span>
                        </span>
                        <a href={'/api/attachments/attachment/' + attachment.ID + '/download'} className={Button.className({ color: Style.Palette.Secondary, outline: true, size: Style.Size.XS })} download>
                            <Icon.Download />
                        </a>
                    </div>
                </ListGroup.Item>);
            })
        }</ListGroup.List>);
    }

    render(): JSX.Element {
        if (this.state.loading) { return (<PageLoading />); }

        return (
            <Page title="View Script">
                <Buttons>
                    <EditButton to={'/scripts/script/' + this.state.script.ID + '/edit'} />
                    <DeleteButton onClick={this.deleteClick} />
                    <Button color={Style.Palette.Success} outline onClick={this.executeClick}><Icon.Label icon={<Icon.PlayCircle />} label="Run Script" /></Button>
                </Buttons>
                <Layout.Row>
                    <Layout.Column>
                        <Card.Card className="mb-3">
                            <Card.Header>Script Details</Card.Header>
                            <ListGroup.List>
                                <ListGroup.TextItem title="Name">{this.state.script.Name}</ListGroup.TextItem>
                                { this.runAs() }
                                <ListGroup.TextItem title="Working Directory">{this.state.script.WorkingDirectory}</ListGroup.TextItem>
                                <ListGroup.TextItem title="Executable">{this.state.script.Executable}</ListGroup.TextItem>
                                <ListGroup.TextItem title="Status"><EnabledBadge value={this.state.script.Enabled}/></ListGroup.TextItem>
                            </ListGroup.List>
                        </Card.Card>
                        <EnvironmentVariableCard variables={this.state.script.Environment} className="mb-3"/>
                        <Card.Card className="mb-3">
                            <Card.Header>Attachments</Card.Header>
                            {this.attachmentList()}
                        </Card.Card>
                        <ScheduleListCard schedules={this.state.schedules} className="mb-3" />
                    </Layout.Column>
                    <Layout.Column>
                        <Card.Card className="mb-3">
                            <Card.Header>Enabled on Hosts</Card.Header>
                            <ListGroup.List>
                                {
                                    this.state.hosts.map((host, index) => {
                                        return (
                                        <ListGroup.Item key={index}>
                                            <div className="d-flex justify-content-between">
                                                <div>
                                                    <Icon.LayerGroup />
                                                    <Link to={'/groups/group/' + host.GroupID} className="ms-1">{ host.GroupName }</Link>
                                                </div>
                                                <div>
                                                    <SmallPlayButton onClick={this.runScriptGroupClick(host.GroupID)} />
                                                </div>
                                            </div>
                                            <div className="d-flex justify-content-between">
                                                <div>
                                                    <Icon.Descendant />
                                                    <Icon.Desktop />
                                                    <Link to={'/hosts/host/' + host.HostID} className="ms-1">{ host.HostName }</Link>
                                                </div>
                                                <div>
                                                    <SmallPlayButton onClick={this.runScriptHostClick(host.HostID)} />
                                                </div>
                                            </div>
                                        </ListGroup.Item>
                                        );
                                    })
                                }
                            </ListGroup.List>
                        </Card.Card>
                        <Card.Card className="mb-3">
                            <Card.Header>Script</Card.Header>
                            <Card.Body>
                                <Pre>{this.state.script.Script}</Pre>
                            </Card.Body>
                        </Card.Card>
                    </Layout.Column>
                </Layout.Row>
            </Page>
        );
    }
}
