import * as React from 'react';
import { Script } from '../../types/Script';
import { match } from 'react-router-dom';
import { URLParams } from '../../services/Params';
import { PageLoading } from '../../components/Loading';
import { Page } from '../../components/Page';
import { Form, Input, Checkbox, NumberInput, Select, Textarea } from '../../components/Form';
import { EnvironmentVariableEdit } from '../../components/EnvironmentVariableEdit';
import { GroupCheckList } from '../../components/CheckList';
import { Card } from '../../components/Card';
import { Notification } from '../../components/Notification';
import { Redirect } from '../../components/Redirect';

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
            this.setState({ isNew: true, script: Script.Blank(), loading: false });
        } else {
            Script.Get(id).then(script => {
                this.setState({ loading: false, script: script });
            });
        }
    }

    private changeName = (Name: string) => {
        this.setState(state => {
            state.script.Name = Name;
            return state;
        });
    }

    private changeEnvironment = (Environment: {[id: string]: string}) => {
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

    private changeEnabled = (Enabled: boolean) => {
        this.setState(state => {
            state.script.Enabled = Enabled;
            return state;
        });
    }

    private changeUID = (UID: number) => {
        this.setState(state => {
            state.script.UID = UID;
            return state;
        });
    }

    private changeGID = (GID: number) => {
        this.setState(state => {
            state.script.GID = GID;
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

    private formSave = () => {
        let promise: Promise<Script>;
        if (this.state.isNew) {
            promise = Script.New(this.state.script);
        } else {
            promise = this.state.script.Save();
        }

        promise.then(script => {
            script.SetGroups(this.state.groupIDs).then(() => {
                Notification.success('Script Saved', script.Name);
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
                <NumberInput
                    label="User ID"
                    defaultValue={this.state.script.UID}
                    onChange={this.changeUID}
                    required />
                <NumberInput
                    label="Group ID"
                    defaultValue={this.state.script.GID}
                    onChange={this.changeGID}
                    required />
                <Input
                    label="Working Directory"
                    type="text"
                    defaultValue={this.state.script.WorkingDirectory}
                    onChange={this.changeWorkingDirectory}
                    helpText="Optional directory that the script should run in." />
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
