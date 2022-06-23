import * as React from 'react';
import { RunAs, Script, ScriptType } from '../../types/Script';
import { useParams, useNavigate } from 'react-router-dom';
import { URLParams } from '../../services/Params';
import { PageLoading } from '../../components/Loading';
import { Page } from '../../components/Page';
import { Input } from '../../components/input/Input';
import { Form } from '../../components/Form';
import { EnvironmentVariableEdit } from '../../components/EnvironmentVariableEdit';
import { GroupCheckList } from '../../components/CheckList';
import { Card } from '../../components/Card';
import { Notification } from '../../components/Notification';
import { Variable } from '../../types/Variable';
import { AttachmentList } from './attachment/AttachmentList';

export const ScriptEdit: React.FC = () => {
    const { id } = useParams() as URLParams;
    const [loading, setLoading] = React.useState<boolean>(true);
    const [script, setScript] = React.useState<ScriptType>();
    const [isNew, setIsNew] = React.useState<boolean>();
    const [groupIDs, setGroupIDs] = React.useState<string[]>();
    const navigate = useNavigate();

    React.useEffect(() => {
        loadScript();
    }, []);

    const loadScript = () => {
        if (id == null) {
            setScript(Script.Blank());
            setIsNew(true);
            setGroupIDs([]);
            setLoading(false);
        } else {
            Script.Get(id).then(script => {
                Script.Groups(script.ID).then(selectedGroups => {
                    setScript(script);
                    setIsNew(false);
                    setGroupIDs(selectedGroups.map(group => group.ID));
                    setLoading(false);
                });
            });
        }
    };

    const changeName = (Name: string) => {
        setScript(script => {
            script.Name = Name;
            return { ...script };
        });
    };

    const changeEnvironment = (Environment: Variable[]) => {
        setScript(script => {
            script.Environment = Environment;
            return { ...script };
        });
    };

    const changeGroupIDs = (GroupIDs: string[]) => {
        setGroupIDs(GroupIDs);
    };

    const changeRunAs = (runAs: RunAs) => {
        setScript(script => {
            script.RunAs = runAs;
            return { ...script };
        });
    };

    const changeWorkingDirectory = (WorkingDirectory: string) => {
        setScript(script => {
            script.WorkingDirectory = WorkingDirectory;
            return { ...script };
        });
    };

    const changeAfterExecution = (AfterExecution: string) => {
        setScript(script => {
            script.AfterExecution = AfterExecution;
            return { ...script };
        });
    };

    const changeExecutable = (Executable: string) => {
        setScript(script => {
            script.Executable = Executable;
            return { ...script };
        });
    };

    const changeScript = (Script: string) => {
        setScript(script => {
            script.Script = Script;
            return { ...script };
        });
    };

    const changeAttachments = (AttachmentIDs: string[]) => {
        setScript(script => {
            script.AttachmentIDs = AttachmentIDs;
            return { ...script };
        });
    };

    const formSave = () => {
        let promise: Promise<ScriptType>;
        if (isNew) {
            promise = Script.New(script);
        } else {
            promise = Script.Save(script);
        }

        return promise.then(script => {
            Script.SetGroups(script.ID, groupIDs).then(() => {
                Notification.success('Script Saved');
                navigate('/scripts/script/' + script.ID);
            });
        });
    };

    if (loading) {
        return (<PageLoading />);
    }

    const breadcrumbs = [
        {
            title: 'Scripts',
            href: '/scripts',
        },
        {
            title: 'New Script'
        }
    ];
    if (!isNew) {
        breadcrumbs[1] = {
            title: script.Name,
            href: '/scripts/script/' + script.ID
        };
        breadcrumbs.push({
            title: 'Edit'
        });
    }

    return (
        <Page title={breadcrumbs}>
            <Form showSaveButton onSubmit={formSave}>
                <Input.Text
                    label="Name"
                    type="text"
                    defaultValue={script.Name}
                    onChange={changeName}
                    required />
                <Input.RunAsInput inheritLabel="Run As Specific User" label="Run As" defaultValue={script.RunAs} onChange={changeRunAs} />
                <Input.Text
                    label="Working Directory"
                    type="text"
                    defaultValue={script.WorkingDirectory}
                    onChange={changeWorkingDirectory}
                    helpText="Optional directory that the script should run in."
                    fixedWidth />
                <Input.Select
                    label="After Script Execution"
                    defaultValue={script.AfterExecution}
                    onChange={changeAfterExecution}>
                    <option value="">Do Nothing</option>
                    <option value="exit_client">Stop the Otto Client</option>
                    <option value="reboot">Reboot the Host</option>
                    <option value="shutdown">Shutdown the Host</option>
                    <option value="update_psk">Rotate the Client PSK</option>
                </Input.Select>
                <Input.Text
                    label="Executable"
                    type="text"
                    defaultValue={script.Executable}
                    onChange={changeExecutable}
                    fixedWidth
                    required />
                <Card.Card className="mt-3">
                    <Card.Header>Environment Variables</Card.Header>
                    <Card.Body>
                        <EnvironmentVariableEdit
                            variables={script.Environment || []}
                            onChange={changeEnvironment} />
                    </Card.Body>
                </Card.Card>
                <Card.Card className="mt-3">
                    <Card.Header>Groups</Card.Header>
                    <Card.Body>
                        <GroupCheckList selectedGroups={groupIDs} onChange={changeGroupIDs} />
                    </Card.Body>
                </Card.Card>
                <Card.Card className="mt-3">
                    <Card.Header>Attachments</Card.Header>
                    <Card.Body>
                        <AttachmentList scriptID={script.ID} didUpdateAttachments={changeAttachments} />
                    </Card.Body>
                </Card.Card>
                <hr />
                <Input.Textarea
                    label="Script"
                    defaultValue={script.Script}
                    onChange={changeScript}
                    rows={10}
                    fixedWidth
                    required />
            </Form>
        </Page>
    );
};
