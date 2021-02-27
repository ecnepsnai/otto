import * as React from 'react';
import { FormGroup } from '../Form';
import { Rand } from '../../services/Rand';
import '../../../css/form.scss';

interface CheckboxProps {
    /**
     * The label that appears above the input
     */
    label: string;
    /**
     * Event called when the value of the input changed
     */
    onChange: (checked: boolean) => (void);
    /**
     * The default value used for the input
     */
    defaultValue?: boolean;
    /**
     * The value used for the input
     */
    checked?: boolean;
    /**
     * Optional help text to appear below this input
     */
    helpText?: string;
    /**
     * Should the input be disabled
     */
    disabled?: boolean;
}
export const Checkbox: React.FC<CheckboxProps> = (props: CheckboxProps) => {
    const [checked, setChecked] = React.useState<boolean>(props.defaultValue);
    const labelID = Rand.ID();

    React.useEffect(() => {
        props.onChange(checked);
    }, [checked]);

    const onChange = (event: React.FormEvent<HTMLInputElement>) => {
        const target = event.target as HTMLInputElement;
        setChecked(target.checked);
    };

    const helpText =() => {
        if (props.helpText) {
            return <div id={labelID + 'help'} className="form-text">{props.helpText}</div>;
        } else {
            return null;
        }
    };

    return (
        <FormGroup className="form-check">
            <input type="checkbox" className="form-check-input" id={labelID} defaultChecked={checked} onChange={onChange} disabled={props.disabled}/>
            <label htmlFor={labelID} className="form-check-label">{props.label}</label>
            { helpText() }
        </FormGroup>
    );
};
