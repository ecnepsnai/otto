import * as React from 'react';
import { Card } from '../../../components/Card';
import { Icon } from '../../../components/Icon';
import { Input } from '../../../components/input/Input';
import { Options } from '../../../types/Options';

interface OptionsNetworkProps {
    defaultValue: Options.Network;
    onUpdate: (value: Options.Network) => (void);
}
interface OptionsNetworkState {
    value: Options.Network;
}
export class OptionsNetwork extends React.Component<OptionsNetworkProps, OptionsNetworkState> {
    constructor(props: OptionsNetworkProps) {
        super(props);
        this.state = {
            value: props.defaultValue,
        };
    }

    private changeForceIPVersion = (ForceIPVersion: string) => {
        this.setState(state => {
            const options = state.value;
            options.ForceIPVersion = ForceIPVersion;
            return {
                value: options,
            };
        }, () => {
            this.props.onUpdate(this.state.value);
        });
    }

    private changeTimeout = (Timeout: number) => {
        this.setState(state => {
            const options = state.value;
            options.Timeout = Timeout;
            return {
                value: options,
            };
        }, () => {
            this.props.onUpdate(this.state.value);
        });
    }

    private changeHeartbeatFrequency = (HeartbeatFrequency: number) => {
        this.setState(state => {
            const options = state.value;
            options.HeartbeatFrequency = HeartbeatFrequency;
            return {
                value: options,
            };
        }, () => {
            this.props.onUpdate(this.state.value);
        });
    }

    render(): JSX.Element {
        const radioChoices = [
            {
                value: 'auto',
                label: 'Automatic'
            },
            {
                value: 'ipv4',
                label: 'IPv4'
            },
            {
                value: 'ipv6',
                label: 'IPv6'
            }
        ];

        return (
            <Card.Card>
                <Card.Header>
                    <Icon.Label icon={<Icon.NetworkWired />} label="Network" />
                </Card.Header>
                <Card.Body>
                    <Input.Radio
                        label="IP Version"
                        choices={radioChoices}
                        defaultValue={this.state.value.ForceIPVersion}
                        onChange={this.changeForceIPVersion}/>
                    <Input.Number
                        label="Timeout"
                        append="Seconds"
                        helpText="The maximum number of seconds Otto will wait while trying to connect to a client"
                        defaultValue={this.state.value.Timeout}
                        onChange={this.changeTimeout} />
                    <Input.Number
                        label="Heartbeat Interval"
                        append="Minutes"
                        helpText="The frequency (in minutes) to check the reachability of all Otto hosts"
                        defaultValue={this.state.value.HeartbeatFrequency}
                        onChange={this.changeHeartbeatFrequency} />
                </Card.Body>
            </Card.Card>
        );
    }
}
