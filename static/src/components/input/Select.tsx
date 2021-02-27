import * as React from 'react';
import { FormGroup, ValidationResult } from '../Form';
import { Rand } from '../../services/Rand';
import '../../../css/form.scss';

interface SelectProps {
    /**
     * The label that appears above the input
     */
    label: string;
    /**
     * Event called when the value of the input changed
     */
    onChange: (value: string) => (void);
    /**
     * The default value used for the input
     */
    defaultValue: string;
    /**
     * If true a value is required for this input
     */
    required?: boolean;
    /**
     * Optional help text to appear below this input
     */
    helpText?: string;
    /**
     * Should the input be disabled
     */
    disabled?: boolean;
    /**
     * Optional method to invoke for validating the value of this input.
     * Return a promise that resolves with a validation result.
     *
     * You do not need to validate if a required field has any value, that is done automatically.
     */
    validate?: (value: string) => Promise<ValidationResult>;

    children?: React.ReactNode;
}
export const Select: React.FC<SelectProps> = (props: SelectProps) => {
    const [value, setValue] = React.useState(props.defaultValue);
    const initialValidState: ValidationResult = {
        valid: true
    };
    if (props.required && (props.defaultValue == null || props.defaultValue == '')) {
        initialValidState.valid = false;
        initialValidState.invalidMessage = 'A selection is required';
    }
    const [valid, setValid] = React.useState(initialValidState);
    const [touched, setTouched] = React.useState(false);
    const labelID = Rand.ID();

    React.useEffect(() => {
        props.onChange(value);
    }, [value]);

    const onChange = (event: React.FormEvent<HTMLSelectElement>) => {
        const target = event.target as HTMLSelectElement;
        validate(target.value).then(valid => {
            setValid(valid);
        });
        setValue(target.value);
        props.onChange(target.value);
    };

    const helpText = () => {
        if (props.helpText) {
            return <div id={labelID + 'help'} className="form-text">{props.helpText}</div>;
        } else {
            return null;
        }
    };

    const requiredFlag = () => {
        if (!props.required) {
            return null; 
        }
        return (<span className="form-required">*</span>);
    };

    const defaultSelection = () => {
        if (props.required && value) {
            return null;
        }

        return (<option selected>Select One...</option>);
    };

    const validate = (value: string): Promise<ValidationResult> => {
        return new Promise((resolve) => {
            if (props.required && value == '') {
                resolve({
                    valid: false,
                    invalidMessage: 'A value is required'
                });
                return;
            }
            if (props.validate) {
                props.validate(value).then(valid => {
                    resolve(valid);
                });
                return;
            }
            resolve({ valid: true });
        });
    };

    const onBlur = () => {
        setTouched(true);
    };

    const validationError = () => {
        if (!valid.invalidMessage || !touched) {
            return null; 
        }
        return (<div className="invalid-feedback">{valid.invalidMessage}</div>);
    };

    let className = 'form-select';
    if (touched && !valid.valid) {
        className += ' is-invalid';
    }
    return (
        <FormGroup>
            <label htmlFor={labelID} className="form-label">{props.label} {requiredFlag()}</label>
            <select
                defaultValue={props.defaultValue}
                className={className}
                id={labelID}
                onChange={onChange}
                disabled={props.disabled}
                onBlur={onBlur}
                data-valid={valid.valid ? 'valid' : 'invalid'}>
                { defaultSelection() }
                { props.children }
            </select>
            { validationError() }
            { helpText() }
        </FormGroup>
    );
};
