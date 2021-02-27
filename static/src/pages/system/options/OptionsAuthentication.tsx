import * as React from 'react';
import { Card } from '../../../components/Card';
import { Icon } from '../../../components/Icon';
import { Input } from '../../../components/input/Input';
import { Options } from '../../../types/Options';

interface OptionsAuthenticationProps {
    defaultValue: Options.Authentication;
    onUpdate: (value: Options.Authentication) => (void);
}
interface OptionsAuthenticationState {
    value: Options.Authentication;
}
export class OptionsAuthentication extends React.Component<OptionsAuthenticationProps, OptionsAuthenticationState> {
    constructor(props: OptionsAuthenticationProps) {
        super(props);
        this.state = {
            value: props.defaultValue,
        };
    }

    private changeMaxAgeMinutes = (MaxAgeMinutes: number) => {
        this.setState(state => {
            const options = state.value;
            options.MaxAgeMinutes = MaxAgeMinutes;
            return {
                value: options,
            };
        }, () => {
            this.props.onUpdate(this.state.value);
        });
    }

    private changeSecureOnly = (SecureOnly: boolean) => {
        this.setState(state => {
            const options = state.value;
            options.SecureOnly = SecureOnly;
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
                    <Icon.Label icon={<Icon.Key />} label="Authentication" />
                </Card.Header>
                <Card.Body>
                    <Input.Number
                        label="Session Timeout"
                        append="Minutes"
                        helpText="The number of minutes of inactivity before a session is automatically ended"
                        defaultValue={this.state.value.MaxAgeMinutes}
                        onChange={this.changeMaxAgeMinutes} />
                    <Input.Checkbox
                        label="Require HTTPS"
                        helpText="If checked users must access the Otto web UI using HTTPS."
                        defaultValue={this.state.value.SecureOnly}
                        onChange={this.changeSecureOnly} />
                </Card.Body>
            </Card.Card>
        );
    }
}
