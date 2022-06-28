import * as React from 'react';
import { Input } from '../../../components/input/Input';
import { Options } from '../../../types/Options';

interface OptionsRotateIDProps {
    defaultValue: Options.RotateID;
    onUpdate: (value: Options.RotateID) => (void);
}
export const OptionsRotateID: React.FC<OptionsRotateIDProps> = (props: OptionsRotateIDProps) => {
    const [value, setValue] = React.useState(props.defaultValue);

    React.useEffect(() => {
        props.onUpdate(value);
    }, [value]);

    const changeEnabled = (Enabled: boolean) => {
        setValue(value => {
            value.Enabled = Enabled;
            return { ...value };
        });
    };

    const changeFrequencyDays = (FrequencyDays: number) => {
        setValue(value => {
            value.FrequencyDays = FrequencyDays;
            return { ...value };
        });
    };

    const content = () => {
        if (!value.Enabled) {
            return null;
        }

        return (<Input.Number label="Rotation Every" append="Days" minimum={1} defaultValue={value.FrequencyDays} onChange={changeFrequencyDays} required />);
    };

    return (
        <React.Fragment>
            <Input.Checkbox
                label="Automatically Rotate Agent Identites"
                defaultValue={value.Enabled}
                helpText="If checked then agent IDs are updated at the frequency specified below."
                onChange={changeEnabled} />
            {content()}
        </React.Fragment>
    );
};
