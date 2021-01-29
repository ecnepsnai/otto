import * as React from 'react';
import { Script } from '../../types/Script';
import { match } from 'react-router-dom';
import { URLParams } from '../../services/Params';
import { PageLoading } from '../../components/Loading';
import { Page } from '../../components/Page';
import { Form, Input, Checkbox, Select, Textarea, IDInput } from '../../components/Form';
import { EnvironmentVariableEdit } from '../../components/EnvironmentVariableEdit';
import { GroupCheckList } from '../../components/CheckList';
import { Card } from '../../components/Card';
import { Notification } from '../../components/Notification';
import { Redirect } from '../../components/Redirect';
import { Variable } from '../../types/Variable';
import { AttachmentList } from './attachment/AttachmentList';

export interface ScriptEditProps { match: match }
interface ScriptEditState {
    loading: boolean;
    script?: Script;
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
                script.Groups().then(selectedGroups => {
                    this.setState({
                        loading: false,
                        script: script,
                        groupIDs: selectedGroups.map(group => { return group.ID; })
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
        if (this.state.isNew) { return null; }

        return (
            <Checkbox
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
        if (this.state.script.RunAs.Inherit) { return null; }

        return (<IDInput
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
        let promise: Promise<Script>;
        if (this.state.isNew) {
            promise = Script.New(this.state.script);
        } else {
            promise = this.state.script.Save();
        }

        return promise.then(script => {
            script.SetGroups(this.state.groupIDs).then(() => {
                Notification.success('Script Saved');
                Redirect.To('/scripts/script/' + script.ID);
            });
        });
    }

    render(): JSX.Element {
        if (this.state.loading) { return (<PageLoading />); }

        return (
        <Page title={ this.state.isNew ? 'New Script' : 'Edit Script' }>
            <Form showSaveButton onSubmit={this.formSave}>
                <Input
                    label="Name"
                    type="text"
                    defaultValue={this.state.script.Name}
                    onChange={this.changeName}
                    required />
                { this.enabledCheckbox() }
                <Checkbox label="Run As Specific User" defaultValue={!this.state.script.RunAs.Inherit} onChange={this.changeRunAsInherit} />
                { this.runAs() }
                <Input
                    label="Working Directory"
                    type="text"
                    defaultValue={this.state.script.WorkingDirectory}
                    onChange={this.changeWorkingDirectory}
                    helpText="Optional directory that the script should run in."
                    fixedWidth />
                <Select
                    label="After Script Execution"
                    defaultValue={this.state.script.AfterExecution}
                    onChange={this.changeAfterExecution}>
                        <option value="">Do Nothing</option>
                        <option value="exit_client">Stop the Otto Client</option>
                        <option value="reboot">Reboot the Host</option>
                        <option value="shutdown">Shutdown the Host</option>
                </Select>
                <Input
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
                <Textarea
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
