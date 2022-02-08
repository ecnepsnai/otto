import * as React from 'react';
import { Options } from '../../../types/Options';
import { OptionsRotateID } from './OptionsRotateID';

interface OptionsSecurityProps {
    defaultValue: Options.Security;
    onUpdate: (value: Options.Security) => (void);
}
export const OptionsSecurity: React.FC<OptionsSecurityProps> = (props: OptionsSecurityProps) => {
    const [value, setValue] = React.useState(props.defaultValue);

    React.useEffect(() => {
        props.onUpdate(value);
    }, [value]);

    const changeRotateID = (RotateID: Options.RotateID) => {
        setValue(value => {
            value.RotateID = RotateID;
            return { ...value };
        });
    };

    return (
        <div>
            <OptionsRotateID defaultValue={value.RotateID} onUpdate={changeRotateID} />
        </div>
    );
};
