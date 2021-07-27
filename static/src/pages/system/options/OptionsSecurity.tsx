import * as React from 'react';
import { Input } from '../../../components/input/Input';
import { Options } from '../../../types/Options';
import { OptionsRotatePSK } from './OptionsRotatePSK';

interface OptionsSecurityProps {
    defaultValue: Options.Security;
    onUpdate: (value: Options.Security) => (void);
}
export const OptionsSecurity: React.FC<OptionsSecurityProps> = (props: OptionsSecurityProps) => {
    const [value, setValue] = React.useState(props.defaultValue);

    React.useEffect(() => {
        props.onUpdate(value);
    }, [value]);

    const changeIncludePSKEnv = (IncludePSKEnv: boolean) => {
        setValue(value => {
            value.IncludePSKEnv = IncludePSKEnv;
            return { ...value };
        });
    };

    const changeRotatePSK = (RotatePSK: Options.RotatePSK) => {
        setValue(value => {
            value.RotatePSK = RotatePSK;
            return { ...value };
        });
    };

    return (
        <div>
            <Input.Checkbox
                label="Include Client PSK Environment Variable"
                defaultValue={value.IncludePSKEnv}
                helpText="If checked the OTTO_CLIENT_PSK environment variable is included when scripts are run."
                onChange={changeIncludePSKEnv} />
            <hr />
            <OptionsRotatePSK defaultValue={value.RotatePSK} onUpdate={changeRotatePSK} />
        </div>
    );
};
