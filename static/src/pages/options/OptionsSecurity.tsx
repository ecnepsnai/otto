import * as React from 'react';
import { Card } from '../../components/Card';
import { Icon } from '../../components/Icon';
import { Checkbox } from '../../components/Form';
import { Options } from '../../types/Options';

export interface OptionsSecurityProps {
    defaultValue: Options.Security;
    onUpdate: (value: Options.Security) => (void);
}
interface OptionsSecurityState {
    value: Options.Security;
}
export class OptionsSecurity extends React.Component<OptionsSecurityProps, OptionsSecurityState> {
    constructor(props: OptionsSecurityProps) {
        super(props);
        this.state = {
            value: props.defaultValue,
        };
    }

    private changeIncludePSKEnv = (IncludePSKEnv: boolean) => {
        this.setState(state => {
            const options = state.value;
            options.IncludePSKEnv = IncludePSKEnv;
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
                    <Icon.Label icon={<Icon.ShieldAlt />} label="Security" />
                </Card.Header>
                <Card.Body>
                    <Checkbox
                        label="Include Client PSK Environment Variable"
                        defaultValue={this.state.value.IncludePSKEnv}
                        helpText="If checked the OTTO_CLIENT_PSK environment variable is included when scripts are run."
                        onChange={this.changeIncludePSKEnv} />
                </Card.Body>
            </Card.Card>
        );
    }
}
