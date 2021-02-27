import * as React from 'react';
import { Page } from '../../../components/Page';
import { Options } from '../../../types/Options';
import { StateManager } from '../../../services/StateManager';
import { Form } from '../../../components/Form';
import { OptionsGeneral } from './OptionsGeneral';
import { OptionsAuthentication } from './OptionsAuthentication';
import { OptionsNetwork } from './OptionsNetwork';
import { Notification } from '../../../components/Notification';
import { OptionsSecurity } from './OptionsSecurity';

interface SystemOptionsState {
    loading?: boolean,
    options: Options.OttoOptions,
}
export class SystemOptions extends React.Component<unknown, SystemOptionsState> {
    constructor(props: unknown) {
        super(props);
        this.state = {
            options: StateManager.Current().Options
        };
    }

    private changeGeneral = (value: Options.General) => {
        this.setState(state => {
            const options = state.options;
            options.General = value;
            return { options: options };
        });
    }

    private changeAuthentication = (value: Options.Authentication) => {
        this.setState(state => {
            const options = state.options;
            options.Authentication = value;
            return { options: options };
        });
    }

    private changeNetwork = (value: Options.Network) => {
        this.setState(state => {
            const options = state.options;
            options.Network = value;
            return { options: options };
        });
    }

    private changeSecurity = (value: Options.Security) => {
        this.setState(state => {
            const options = state.options;
            options.Security = value;
            return { options: options };
        });
    }

    private onSubmit = () => {
        this.setState({ loading: true });
        return Options.Options.Save(this.state.options).then(() => {
            Notification.success('Options Saved');
            this.setState({ loading: false });
        });
    }

    render(): JSX.Element {
        return (
            <Page title="Options">
                <Form className="cards" showSaveButton onSubmit={this.onSubmit}>
                    <OptionsGeneral defaultValue={this.state.options.General} onUpdate={this.changeGeneral}/>
                    <OptionsAuthentication defaultValue={this.state.options.Authentication} onUpdate={this.changeAuthentication}/>
                    <OptionsNetwork defaultValue={this.state.options.Network} onUpdate={this.changeNetwork}/>
                    <OptionsSecurity defaultValue={this.state.options.Security} onUpdate={this.changeSecurity}/>
                    <div className="mb-2"></div>
                </Form>
            </Page>
        );
    }
}
