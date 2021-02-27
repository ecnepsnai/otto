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
interface ScriptEditState {
    loading: boolean;
    script?: ScriptType;
    isNew?: boolean;
    groupIDs?: string[];
}
export class ScriptEdit extends React.Component<ScriptEditProps, ScriptEditState> {
    constructor(props: ScriptEditProps) {
        super(props);
        this.state = {
            loading: true,
        };
    }

    componentDidMount(): void {
        this.loadScript();
    }

    loadScript(): void {
        const id = (this.props.match.params as URLParams).id;
        if (id == null) {
            this.setState({
                isNew: true,
                script: Script.Blank(),
                loading: false,
                groupIDs: [],
            });
        } else {
            Script.Get(id).then(script => {
                Script.Groups(script.ID).then(selectedGroups => {
                    this.setState({
                        loading: false,
                        script: script,
                        groupIDs: selectedGroups.map(group => group.ID)
                    });
                });
            });
        }
    }

    private changeName = (Name: string) => {
        this.setState(state => {
            state.script.Name = Name;
            return state;
        });
    }

    private changeEnvironment = (Environment: Variable[]) => {
        this.setState(state => {
            state.script.Environment = Environment;
            return state;
        });
    }

    private changeGroupIDs = (groupIDs: string[]) => {
        this.setState({ groupIDs: groupIDs });
    }

    private enabledCheckbox = () => {
        if (this.state.isNew) {
            return null;
        }

        return (
            <Input.Checkbox
                label="Enabled"
                defaultValue={this.state.script.Enabled}
                onChange={this.changeEnabled} />
        );
    }

    private changeRunAsInherit = (DontInherit: boolean) => {
        this.setState(state => {
            state.script.RunAs.Inherit = !DontInherit;
            return state;
        });
    }

    private runAs = () => {
        if (this.state.script.RunAs.Inherit) {
            return null;
        }

        return (<Input.IDInput
            label="Run Script As"
            defaultUID={this.state.script.RunAs.UID}
            defaultGID={this.state.script.RunAs.GID}
            onChange={this.changeID} />);
    }

    private changeEnabled = (Enabled: boolean) => {
        this.setState(state => {
            state.script.Enabled = Enabled;
            return state;
        });
    }

    private changeID = (UID: number, GID: number) => {
        this.setState(state => {
            state.script.RunAs.UID = UID;
            state.script.RunAs.GID = GID;
            return state;
        });
    }

    private changeWorkingDirectory = (WorkingDirectory: string) => {
        this.setState(state => {
            state.script.WorkingDirectory = WorkingDirectory;
            return state;
        });
    }

    private changeAfterExecution = (AfterExecution: string) => {
        this.setState(state => {
            state.script.AfterExecution = AfterExecution;
            return state;
        });
    }

    private changeExecutable = (Executable: string) => {
        this.setState(state => {
            state.script.Executable = Executable;
            return state;
        });
    }

    private changeScript = (Script: string) => {
        this.setState(state => {
            state.script.Script = Script;
            return state;
        });
    }

    private changeAttachments = (AttachmentIDs: string[]) => {
        this.setState(state => {
            state.script.AttachmentIDs = AttachmentIDs;
            return state;
        });
    }

    private formSave = () => {
        let promise: Promise<ScriptType>;
        if (this.state.isNew) {
            promise = Script.New(this.state.script);
        } else {
            promise = Script.Save(this.state.script);
        }

        return promise.then(script => {
            Script.SetGroups(this.state.script.ID, this.state.groupIDs).then(() => {
                Notification.success('Script Saved');
                Redirect.To('/scripts/script/' + script.ID);
            });
        });
    }

    render(): JSX.Element {
        if (this.state.loading) {
            return (<PageLoading />);
        }

        return (
            <Page title={ this.state.isNew ? 'New Script' : 'Edit Script' }>
                <Form showSaveButton onSubmit={this.formSave}>
                    <Input.Text
                        label="Name"
                        type="text"
                        defaultValue={this.state.script.Name}
                        onChange={this.changeName}
                        required />
                    { this.enabledCheckbox() }
                    <Input.Checkbox label="Run As Specific User" defaultValue={!this.state.script.RunAs.Inherit} onChange={this.changeRunAsInherit} />
                    { this.runAs() }
                    <Input.Text
                        label="Working Directory"
                        type="text"
                        defaultValue={this.state.script.WorkingDirectory}
                        onChange={this.changeWorkingDirectory}
                        helpText="Optional directory that the script should run in."
                        fixedWidth />
                    <Input.Select
                        label="After Script Execution"
                        defaultValue={this.state.script.AfterExecution}
                        onChange={this.changeAfterExecution}>
                        <option value="">Do Nothing</option>
                        <option value="exit_client">Stop the Otto Client</option>
                        <option value="reboot">Reboot the Host</option>
                        <option value="shutdown">Shutdown the Host</option>
                    </Input.Select>
                    <Input.Text
                        label="Executable"
                        type="text"
                        defaultValue={this.state.script.Executable}
                        onChange={this.changeExecutable}
                        fixedWidth
                        required />
                    <Card.Card className="mt-3">
                        <Card.Header>Environment Variables</Card.Header>
                        <Card.Body>
                            <EnvironmentVariableEdit
                                variables={this.state.script.Environment}
                                onChange={this.changeEnvironment} />
                        </Card.Body>
                    </Card.Card>
                    <Card.Card className="mt-3">
                        <Card.Header>Groups</Card.Header>
                        <Card.Body>
                            <GroupCheckList selectedGroups={this.state.groupIDs} onChange={this.changeGroupIDs}/>
                        </Card.Body>
                    </Card.Card>
                    <Card.Card className="mt-3">
                        <Card.Header>Attachments</Card.Header>
                        <Card.Body>
                            <AttachmentList scriptID={this.state.script.ID} didUpdateAttachments={this.changeAttachments}/>
                        </Card.Body>
                    </Card.Card>
                    <hr/>
                    <Input.Textarea
                        label="Script"
                        defaultValue={this.state.script.Script}
                        onChange={this.changeScript}
                        rows={10}
                        fixedWidth
                        required />
                </Form>
            </Page>
        );
    }
}
