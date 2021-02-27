import * as React from 'react';
import { Script, ScriptType } from '../../types/Script';
import { match } from 'react-router-dom';
import { URLParams } from '../../services/Params';
import { PageLoading } from '../../components/Loading';
import { Page } from '../../components/Page';
import { Input } from '../../components/input/Input';
import { Form } from '../../components/Form';
import { EnvironmentVariableEdit } from '../../components/EnvironmentVariableEdit';
import { GroupCheckList } from '../../components/CheckList';
import { Card } from '../../components/Card';
import { Notification } from '../../components/Notification';
import { Redirect } from '../../components/Redirect';
import { Variable } from '../../types/Variable';
import { AttachmentList } from './attachment/AttachmentList';

interface ScriptEditProps {
    match: match
}
export const ScriptEdit: React.FC<ScriptEditProps> = (props: ScriptEditProps) => {
    const [loading, setLoading] = React.useState<boolean>(true);
    const [script, setScript] = React.useState<ScriptType>();
    const [isNew, setIsNew] = React.useState<boolean>();
    const [groupIDs, setGroupIDs] = React.useState<string[]>();

    React.useEffect(() => {
        loadScript();
    }, []);

    const loadScript = () => {
        const id = (props.match.params as URLParams).id;
        if (id == null) {
            setLoading(true);
            setScript(Script.Blank());
            setIsNew(true);
            setGroupIDs([]);
        } else {
            Script.Get(id).then(script => {
                Script.Groups(script.ID).then(selectedGroups => {
                    setLoading(false);
                    setScript(script);
                    setIsNew(false);
                    setGroupIDs(selectedGroups.map(group => group.ID));
                });
            });
        }
    };

    const changeName = (Name: string) => {
        setScript(script => {
            script.Name = Name;
            return {...script};
        });
    };

    const changeEnvironment = (Environment: Variable[]) => {
        setScript(script => {
            script.Environment = Environment;
            return {...script};
        });
    };

    const changeGroupIDs = (GroupIDs: string[]) => {
        setGroupIDs(GroupIDs);
    };

    const enabledCheckbox = () => {
        if (isNew) {
            return null;
        }

        return (
            <Input.Checkbox
                label="Enabled"
                defaultValue={script.Enabled}
                onChange={changeEnabled} />
        );
    };

    const changeRunAsInherit = (DontInherit: boolean) => {
        setScript(script => {
            script.RunAs.Inherit = !DontInherit;
            return {...script};
        });
    };

    const runAs = () => {
        if (script.RunAs.Inherit) {
            return null;
        }

        return (<Input.IDInput
            label="Run Script As"
            defaultUID={script.RunAs.UID}
            defaultGID={script.RunAs.GID}
            onChange={changeID} />);
    };

    const changeEnabled = (Enabled: boolean) => {
        setScript(script => {
            script.Enabled = Enabled;
            return {...script};
        });
    };

    const changeID = (UID: number, GID: number) => {
        setScript(script => {
            script.RunAs.UID = UID;
            script.RunAs.GID = GID;
            return {...script};
        });
    };

    const changeWorkingDirectory = (WorkingDirectory: string) => {
        setScript(script => {
            script.WorkingDirectory = WorkingDirectory;
            return {...script};
        });
    };

    const changeAfterExecution = (AfterExecution: string) => {
        setScript(script => {
            script.AfterExecution = AfterExecution;
            return {...script};
        });
    };

    const changeExecutable = (Executable: string) => {
        setScript(script => {
            script.Executable = Executable;
            return {...script};
        });
    };

    const changeScript = (Script: string) => {
        setScript(script => {
            script.Script = Script;
            return {...script};
        });
    };

    const changeAttachments = (AttachmentIDs: string[]) => {
        setScript(script => {
            script.AttachmentIDs = AttachmentIDs;
            return {...script};
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
                Redirect.To('/scripts/script/' + script.ID);
            });
        });
    };

    if (loading) {
        return (<PageLoading />);
    }

    return (
        <Page title={ isNew ? 'New Script' : 'Edit Script' }>
            <Form showSaveButton onSubmit={formSave}>
                <Input.Text
                    label="Name"
                    type="text"
                    defaultValue={script.Name}
                    onChange={changeName}
                    required />
                { enabledCheckbox() }
                <Input.Checkbox label="Run As Specific User" defaultValue={!script.RunAs.Inherit} onChange={changeRunAsInherit} />
                { runAs() }
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
                            variables={script.Environment}
                            onChange={changeEnvironment} />
                    </Card.Body>
                </Card.Card>
                <Card.Card className="mt-3">
                    <Card.Header>Groups</Card.Header>
                    <Card.Body>
                        <GroupCheckList selectedGroups={groupIDs} onChange={changeGroupIDs}/>
                    </Card.Body>
                </Card.Card>
                <Card.Card className="mt-3">
                    <Card.Header>Attachments</Card.Header>
                    <Card.Body>
                        <AttachmentList scriptID={script.ID} didUpdateAttachments={changeAttachments}/>
                    </Card.Body>
                </Card.Card>
                <hr/>
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
