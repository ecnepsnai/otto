import * as React from 'react';
import { Card } from '../../../components/Card';
import { Icon } from '../../../components/Icon';
import { Input } from '../../../components/input/Input';
import { Options } from '../../../types/Options';
import { EnvironmentVariableEdit } from '../../../components/EnvironmentVariableEdit';
import { Variable } from '../../../types/Variable';
import { Alert } from '../../../components/Alert';

interface OptionsGeneralProps {
    defaultValue: Options.General;
    onUpdate: (value: Options.General) => (void);
}
interface OptionsGeneralState {
    value: Options.General;
}
export class OptionsGeneral extends React.Component<OptionsGeneralProps, OptionsGeneralState> {
    constructor(props: OptionsGeneralProps) {
        super(props);
        this.state = {
            value: props.defaultValue,
        };
    }

    private originWarning = () => {
        let origin = location.origin;
        if (!origin.endsWith('/')) {
            origin = origin + '/';
        }
        let serverURL = this.props.defaultValue.ServerURL;
        if (!serverURL.endsWith('/')) {
            serverURL = serverURL + '/';
        }

        if (origin === serverURL) {
            return null;
        }

        return (<Alert.Danger>
            The configured server URL is different than the URL you are using to access this page.
        </Alert.Danger>);
    }

    private changeServerURL = (ServerURL: string) => {
        this.setState(state => {
            const options = state.value;
            options.ServerURL = ServerURL;
            return {
                value: options,
            };
        }, () => {
            this.props.onUpdate(this.state.value);
        });
    }

    private changeGlobalEnvironment = (GlobalEnvironment: Variable[]) => {
        this.setState(state => {
            const options = state.value;
            options.GlobalEnvironment = GlobalEnvironment;
            return {
                value: options,
            };
        }, () => {
            this.props.onUpdate(this.state.value);
        });
    }

    render(): JSX.Element {
        return (
            <Card.Card>
                <Card.Header>
                    <Icon.Label icon={<Icon.Wrench />} label="General" />
                </Card.Header>
                <Card.Body>
                    <Input.Text
                        type="text"
                        label="Otto Server URL"
                        placeholder="https://otto.example.com/"
                        helpText="The absolute URL (Including protocol) where this otto server is accessed from"
                        defaultValue={this.state.value.ServerURL}
                        onChange={this.changeServerURL} />
                    { this.originWarning() }
                    <label className="form-label">Global Environment Variables</label>
                    <div>
                        <EnvironmentVariableEdit
                            variables={this.state.value.GlobalEnvironment}
                            onChange={this.changeGlobalEnvironment} />
                    </div>
                </Card.Body>
            </Card.Card>
        );
    }
}
