import * as React from 'react';
import { Input } from '../../../components/input/Input';
import { Options } from '../../../types/Options';

interface OptionsNetworkProps {
    defaultValue: Options.Network;
    onUpdate: (value: Options.Network) => (void);
}
export const OptionsNetwork: React.FC<OptionsNetworkProps> = (props: OptionsNetworkProps) => {
    const [value, setValue] = React.useState(props.defaultValue);

    React.useEffect(() => {
        props.onUpdate(value);
    }, [value]);

    const changeForceIPVersion = (ForceIPVersion: string) => {
        setValue(value => {
            value.ForceIPVersion = ForceIPVersion;
            return { ...value };
        });
    };

    const changeTimeout = (Timeout: number) => {
        setValue(value => {
            value.Timeout = Timeout;
            return { ...value };
        });
    };

    const changeHeartbeatFrequency = (HeartbeatFrequency: number) => {
        setValue(value => {
            value.HeartbeatFrequency = HeartbeatFrequency;
            return { ...value };
        });
    };

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
        <div>
            <Input.Radio
                label="IP Version"
                choices={radioChoices}
                defaultValue={value.ForceIPVersion}
                onChange={changeForceIPVersion} />
            <Input.Number
                label="Timeout"
                append="Seconds"
                helpText="The maximum number of seconds Otto will wait while trying to connect to a agent"
                defaultValue={value.Timeout}
                onChange={changeTimeout} />
            <Input.Number
                label="Heartbeat Interval"
                append="Minutes"
                helpText="The frequency (in minutes) to check the reachability of all Otto hosts"
                defaultValue={value.HeartbeatFrequency}
                onChange={changeHeartbeatFrequency} />
        </div>
    );
};
