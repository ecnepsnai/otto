import * as React from 'react';
import { Host, HostType } from '../../types/Host';
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
import { RandomPSK } from '../../components/RandomPSK';

interface HostEditProps { match: match }
interface HostEditState {
    loading: boolean;
    host?: HostType;
    isNew?: boolean;
    useHostName?: boolean;
}
export class HostEdit extends React.Component<HostEditProps, HostEditState> {
    constructor(props: HostEditProps) {
        super(props);
        this.state = {
            loading: true,
        };
    }

    componentDidMount(): void {
        this.loadHost();
    }

    loadHost(): void {
        const id = (this.props.match.params as URLParams).id;
        if (id == null) {
            this.setState({
                isNew: true,
                host: Host.Blank(),
                loading: false,
                useHostName: true,
            });
        } else {
            Host.Get(id).then(host => {
                this.setState({
                    loading: false,
                    host: host,
                    useHostName: host.Name == host.Address,
                });
            });
        }
    }

    private changeName = (Name: string) => {
        this.setState(state => {
            state.host.Name = Name;
            if (state.useHostName) {
                state.host.Address = Name;
            }
            return state;
        });
    }

    private changeAddress = (Address: string) => {
        this.setState(state => {
            state.host.Address = Address;
            return state;
        });
    }

    private changePort = (Port: number) => {
        this.setState(state => {
            state.host.Port = Port;
            return state;
        });
    }

    private changePSK = (PSK: string) => {
        this.setState(state => {
            state.host.PSK = PSK;
            return state;
        });
    }

    private enabledCheckbox = () => {
        if (this.state.isNew) {
            return null;
        }

        return (
            <Input.Checkbox
                label="Enabled"
                helpText=""
                defaultValue={this.state.host.Enabled}
                onChange={this.changeEnabled} />
        );
    }

    private changeEnabled = (Enabled: boolean) => {
        this.setState(state => {
            state.host.Enabled = Enabled;
            return state;
        });
    }

    private changeEnvironment = (Environment: Variable[]) => {
        this.setState(state => {
            state.host.Environment = Environment;
            return state;
        });
    }

    private changeGroupIDs = (groupIDs: string[]) => {
        this.setState(state => {
            state.host.GroupIDs = groupIDs;
            return state;
        });
    }

    private formSave = () => {
        let promise: Promise<HostType>;
        if (this.state.isNew) {
            promise = Host.New(this.state.host);
        } else {
            promise = Host.Save(this.state.host);
        }

        return promise.then(host => {
            Notification.success('Host Saved');
            Redirect.To('/hosts/host/' + host.ID);
        });
    }

    private changeUseHostName = (useHostName: boolean) => {
        this.setState({ useHostName: useHostName });
    }

    private addressInput = () => {
        if (this.state.useHostName) {
            return null;
        }

        return (
            <Input.Text
                label="Address"
                type="text"
                defaultValue={this.state.host.Address}
                onChange={this.changeAddress}
                required />
        );
    }

    render(): JSX.Element {
        if (this.state.loading) {
            return (<PageLoading />);
        }

        return (
            <Page title={ this.state.isNew ? 'New Host' : 'Edit Host' }>
                <Form showSaveButton onSubmit={this.formSave}>
                    <Input.Text
                        label="Name"
                        type="text"
                        defaultValue={this.state.host.Name}
                        onChange={this.changeName}
                        required />
                    <Input.Checkbox label="Connect to host using this name" defaultValue={this.state.useHostName} onChange={this.changeUseHostName} />
                    { this.addressInput() }
                    <Input.Number
                        label="Port"
                        defaultValue={this.state.host.Port}
                        onChange={this.changePort}
                        required />
                    <Input.Text
                        label="Pre-Shared Key"
                        type="password"
                        defaultValue={this.state.host.PSK}
                        onChange={this.changePSK}
                        required />
                    <RandomPSK newPSK={this.changePSK} />
                    { this.enabledCheckbox() }
                    <Card.Card className="mt-3">
                        <Card.Header>Environment Variables</Card.Header>
                        <Card.Body>
                            <EnvironmentVariableEdit
                                variables={this.state.host.Environment}
                                onChange={this.changeEnvironment} />
                        </Card.Body>
                    </Card.Card>
                    <Card.Card className="mt-3">
                        <Card.Header>Groups</Card.Header>
                        <Card.Body>
                            <GroupCheckList selectedGroups={this.state.host.GroupIDs} onChange={this.changeGroupIDs}/>
                        </Card.Body>
                    </Card.Card>
                </Form>
            </Page>
        );
    }
}
