import * as React from 'react';
import { Script, ScriptType, ScriptEnabledHost } from '../../types/Script';
import { match, Link } from 'react-router-dom';
import { PageLoading } from '../../components/Loading';
import { Page } from '../../components/Page';
import { Layout } from '../../components/Layout';
import { Card } from '../../components/Card';
import { URLParams } from '../../services/Params';
import { Buttons, EditButton, DeleteButton, Button, SmallPlayButton, ButtonAnchor } from '../../components/Button';
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
import { AttachmentType } from '../../types/Attachment';
import { Formatter } from '../../services/Formatter';
import { Nothing } from '../../components/Nothing';
import { ScheduleType } from '../../types/Schedule';
import { ScheduleListCard } from '../../components/ScheduleListCard';

interface ScriptViewProps {
    match: match;
}
export const ScriptView: React.FC<ScriptViewProps> = (props: ScriptViewProps) => {
    const [loading, setLoading] = React.useState<boolean>(true);
    const [script, setScript] = React.useState<ScriptType>();
    const [hosts, setHosts] = React.useState<ScriptEnabledHost[]>();
    const [attachments, setAttachments] = React.useState<AttachmentType[]>();
    const [schedules, setSchedules] = React.useState<ScheduleType[]>();

    React.useEffect(() => {
        loadData();
    }, []);

    const loadHosts = async () => {
        const scriptID = (props.match.params as URLParams).id;
        setHosts(await Script.Hosts(scriptID));
    };

    const loadScript = async () => {
        const scriptID = (props.match.params as URLParams).id;
        setScript(await Script.Get(scriptID));
    };

    const loadAttachments = async () => {
        const scriptID = (props.match.params as URLParams).id;
        setAttachments(await Script.Attachments(scriptID));
    };

    const loadSchedules = async () => {
        const scriptID = (props.match.params as URLParams).id;
        setSchedules(await Script.Schedules(scriptID));
    };

    const loadData = () => {
        Promise.all([loadHosts(), loadScript(), loadAttachments(), loadSchedules()]).then(() => {
            setLoading(false);
        });
    };

    const deleteClick = () => {
        Script.DeleteModal(script).then(deleted => {
            if (!deleted) {
                return;
            }

            Redirect.To('/scripts');
        });
    };

    const executeClick = () => {
        GlobalModalFrame.showModal(<RunModal scriptID={script.ID} key={Rand.ID()}/>);
    };

    const runScriptGroupClick = (groupID: string) => {
        return () => {
            Group.Hosts(groupID).then(hosts => {
                const hostIDs = hosts.map(host => host.ID);
                GlobalModalFrame.showModal(<RunModal scriptID={script.ID} hostIDs={hostIDs} key={Rand.ID()}/>);
            });
        };
    };

    const runScriptHostClick = (hostID: string) => {
        return () => {
            GlobalModalFrame.showModal(<RunModal scriptID={script.ID} hostIDs={[hostID]} key={Rand.ID()}/>);
        };
    };

    const runAs = () => {
        if (script.RunAs.Inherit) {
            return null;
        }

        return (<ListGroup.TextItem title="Run As">User: {script.RunAs.UID} Group: {script.RunAs.GID}</ListGroup.TextItem>);
    };

    const attachmentList = () => {
        if (!attachments || attachments.length == 0) {
            return (<Card.Body><Nothing /></Card.Body>);
        }

        return (<ListGroup.List>{
            attachments.map((attachment, idx) => {
                return (<ListGroup.Item key={idx}>
                    <div className="d-flex justify-content-between">
                        <span>
                            <Icon.Label icon={<Icon.Paperclip />} label={attachment.Name} />
                            <span className="text-muted ms-1">
                                {Formatter.Bytes(attachment.Size)}
                            </span>
                        </span>
                        <ButtonAnchor href={'/api/attachments/attachment/' + attachment.ID + '/download'} color={Style.Palette.Secondary} outline size={Style.Size.XS} download><Icon.Download /></ButtonAnchor>
                    </div>
                </ListGroup.Item>);
            })
        }</ListGroup.List>);
    };

    if (loading) {
        return (<PageLoading />);
    }

    return (
        <Page title="View Script">
            <Buttons>
                <EditButton to={'/scripts/script/' + script.ID + '/edit'} />
                <DeleteButton onClick={deleteClick} />
                <Button color={Style.Palette.Success} outline onClick={executeClick}><Icon.Label icon={<Icon.PlayCircle />} label="Run Script" /></Button>
            </Buttons>
            <Layout.Row>
                <Layout.Column>
                    <Card.Card className="mb-3">
                        <Card.Header>Script Details</Card.Header>
                        <ListGroup.List>
                            <ListGroup.TextItem title="Name">{script.Name}</ListGroup.TextItem>
                            { runAs() }
                            <ListGroup.TextItem title="Working Directory">{script.WorkingDirectory}</ListGroup.TextItem>
                            <ListGroup.TextItem title="Executable">{script.Executable}</ListGroup.TextItem>
                            <ListGroup.TextItem title="Status"><EnabledBadge value={script.Enabled}/></ListGroup.TextItem>
                        </ListGroup.List>
                    </Card.Card>
                    <EnvironmentVariableCard variables={script.Environment} className="mb-3"/>
                    <Card.Card className="mb-3">
                        <Card.Header>Attachments</Card.Header>
                        {attachmentList()}
                    </Card.Card>
                    <ScheduleListCard schedules={schedules} className="mb-3" />
                </Layout.Column>
                <Layout.Column>
                    <Card.Card className="mb-3">
                        <Card.Header>Enabled on Hosts</Card.Header>
                        <ListGroup.List>
                            {
                                hosts.map((host, index) => {
                                    return (<ListGroup.Item key={index}>
                                        <div className="d-flex justify-content-between">
                                            <div>
                                                <Icon.LayerGroup />
                                                <Link to={'/groups/group/' + host.GroupID} className="ms-1">{ host.GroupName }</Link>
                                            </div>
                                            <div>
                                                <SmallPlayButton onClick={runScriptGroupClick(host.GroupID)} />
                                            </div>
                                        </div>
                                        <div className="d-flex justify-content-between">
                                            <div>
                                                <Icon.Descendant />
                                                <Icon.Desktop />
                                                <Link to={'/hosts/host/' + host.HostID} className="ms-1">{ host.HostName }</Link>
                                            </div>
                                            <div>
                                                <SmallPlayButton onClick={runScriptHostClick(host.HostID)} />
                                            </div>
                                        </div>
                                    </ListGroup.Item>);
                                })
                            }
                        </ListGroup.List>
                    </Card.Card>
                    <Card.Card className="mb-3">
                        <Card.Header>Script</Card.Header>
                        <Card.Body>
                            <Pre>{script.Script}</Pre>
                        </Card.Body>
                    </Card.Card>
                </Layout.Column>
            </Layout.Row>
        </Page>
    );
};
