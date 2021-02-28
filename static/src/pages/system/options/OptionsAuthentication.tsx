import * as React from 'react';
import { Card } from '../../../components/Card';
import { Icon } from '../../../components/Icon';
import { Input } from '../../../components/input/Input';
import { Options } from '../../../types/Options';

interface OptionsAuthenticationProps {
    defaultValue: Options.Authentication;
    onUpdate: (value: Options.Authentication) => (void);
}
export const OptionsAuthentication: React.FC<OptionsAuthenticationProps> = (props: OptionsAuthenticationProps) => {
    const [value, setValue] = React.useState(props.defaultValue);

    React.useEffect(() => {
        props.onUpdate(value);
    }, [value]);

    const changeMaxAgeMinutes = (MaxAgeMinutes: number) => {
        setValue(value => {
            value.MaxAgeMinutes = MaxAgeMinutes;
            return {...value};
        });
    };

    const changeSecureOnly = (SecureOnly: boolean) => {
        setValue(value => {
            value.SecureOnly = SecureOnly;
            return {...value};
        });
    };

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
                    defaultValue={value.MaxAgeMinutes}
                    onChange={changeMaxAgeMinutes} />
                <Input.Checkbox
                    label="Require HTTPS"
                    helpText="If checked users must access the Otto web UI using HTTPS."
                    defaultValue={value.SecureOnly}
                    onChange={changeSecureOnly} />
            </Card.Body>
        </Card.Card>
    );
};
