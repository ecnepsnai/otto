import * as React from 'react';
import { Group } from '../../types/Group';
import { match } from 'react-router-dom';
import { URLParams } from '../../services/Params';
import { PageLoading } from '../../components/Loading';
import { Page } from '../../components/Page';
import { Form, Input } from '../../components/Form';
import { EnvironmentVariableEdit } from '../../components/EnvironmentVariableEdit';
import { ScriptCheckList, HostCheckList } from '../../components/CheckList';
import { Card } from '../../components/Card';
import { Notification } from '../../components/Notification';
import { Redirect } from '../../components/Redirect';
import { Variable } from '../../types/Variable';

export interface GroupEditProps { match: match }
interface GroupEditState {
    loading: boolean;
    group?: Group;
    isNew?: boolean;
    hostIDs?: string[];
}
export class GroupEdit extends React.Component<GroupEditProps, GroupEditState> {
    constructor(props: GroupEditProps) {
        super(props);
        this.state = {
            loading: true,
        };
    }

    componentDidMount(): void {
        this.loadGroup();
    }

    loadGroup(): void {
        const id = (this.props.match.params as URLParams).id;
        if (id == null) {
            this.setState({ isNew: true, group: Group.Blank(), loading: false, hostIDs: [] });
        } else {
            Group.Get(id).then(group => {
                group.Hosts().then(hostIDs => {
                    this.setState({ loading: false, group: group, hostIDs: hostIDs.map(host => { return host.ID; }) });
                });
            });
        }
    }

    private changeName = (Name: string) => {
        this.setState(state => {
            state.group.Name = Name;
            return state;
        });
    }

    private changeEnvironment = (Environment: Variable[]) => {
        this.setState(state => {
            state.group.Environment = Environment;
            return state;
        });
    }

    private changeScriptIDs = (ScriptIDs: string[]) => {
        this.setState(state => {
            state.group.ScriptIDs = ScriptIDs;
            return state;
        });
    }

    private changeHostIDs = (hostIDs: string[]) => {
        this.setState({ hostIDs: hostIDs });
    }

    private formSave = () => {
        let promise: Promise<Group>;
        if (this.state.isNew) {
            promise = Group.New(this.state.group);
        } else {
            promise = this.state.group.Save();
        }

        return promise.then(group => {
            group.SetHosts(this.state.hostIDs).then(() => {
                Notification.success('Group Saved');
                Redirect.To('/groups/group/' + group.ID);
            });
        });
    }

    render(): JSX.Element {
        if (this.state.loading) { return (<PageLoading />); }

        return (
        <Page title={ this.state.isNew ? 'New Group' : 'Edit Group' }>
            <Form showSaveButton onSubmit={this.formSave}>
                <Input
                    label="Name"
                    type="text"
                    defaultValue={this.state.group.Name}
                    onChange={this.changeName}
                    required />
                <Card.Card className="mt-3">
                    <Card.Header>Environment Variables</Card.Header>
                    <Card.Body>
                        <EnvironmentVariableEdit
                            variables={this.state.group.Environment}
                            onChange={this.changeEnvironment} />
                    </Card.Body>
                </Card.Card>
                <Card.Card className="mt-3">
                    <Card.Header>Scripts</Card.Header>
                    <Card.Body>
                        <ScriptCheckList selectedScripts={this.state.group.ScriptIDs} onChange={this.changeScriptIDs}/>
                    </Card.Body>
                </Card.Card>
                <Card.Card className="mt-3">
                    <Card.Header>Hosts</Card.Header>
                    <Card.Body>
                        <HostCheckList selectedHosts={this.state.hostIDs} onChange={this.changeHostIDs}/>
                    </Card.Body>
                </Card.Card>
            </Form>
        </Page>
        );
    }
}
