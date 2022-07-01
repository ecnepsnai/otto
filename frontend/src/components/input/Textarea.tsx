import * as React from 'react';
import debounce = require('debounce-promise');
import { FormGroup, ValidationResult } from '../Form';
import { Rand } from '../../services/Rand';
import { InputProps } from './Input';
import '../../../css/form.scss';

interface TextareaProps extends InputProps {
    /**
     * The label that appears above the input
     */
    label: string;
    /**
     * Optional placeholder text for the input
     */
    placeholder?: string;
    /**
     * Event called when the value of the input changed
     */
    onChange: (value: string) => (void);
    /**
     * The default value used for the input
     */
    defaultValue?: string;
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
     * If true then a fixed width font is used
     */
    fixedWidth?: boolean;
    /**
     * Number of rows
     */
    rows?: number;
    /**
     * Optional method to invoke for validating the value of this input.
     * Return a promise that resolves with a validation result.
     *
     * You do not need to validate if a required field has any value, that is done automatically.
     */
    validate?: (value: string) => Promise<ValidationResult>;
}
export const Textarea: React.FC<TextareaProps> = (props: TextareaProps) => {
    const [value, setValue] = React.useState<string>(props.defaultValue);
    const initialValidState: ValidationResult = {
        valid: true
    };
    if (props.required && !props.defaultValue) {
        initialValidState.valid = false;
        initialValidState.invalidMessage = 'A value is required';
    }
    const [valid, setValid] = React.useState<ValidationResult>(initialValidState);
    const [touched, setTouched] = React.useState(false);
    const labelID = Rand.ID();

    React.useEffect(() => {
        props.onChange(value);
    }, [value]);

    const debouncedValidate = debounce(props.validate, 250);
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
                debouncedValidate(value).then(valid => {
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

    const onChange = (event: React.FormEvent<HTMLTextAreaElement>) => {
        const target = event.target as HTMLTextAreaElement;
        validate(target.value).then(valid => {
            setValid(valid);
        });
        setValue(target.value);
    };

    const helpText = () => {
        if (props.helpText) {
            return <div id={labelID + 'help'} className="form-text">{props.helpText}</div>;
        } else {
            return null;
        }
    };

    const validationError = () => {
        if (!valid.invalidMessage || !touched) {
            return null;
        }
        return (<div className="invalid-feedback">{valid.invalidMessage}</div>);
    };

    const input = () => {
        let className = 'form-control';
        if (touched && !valid.valid) {
            className += ' is-invalid';
        }
        if (props.fixedWidth) {
            className += ' fixed-width';
        }
        return (
            <textarea
                className={className}
                id={labelID}
                placeholder={props.placeholder}
                defaultValue={value}
                disabled={props.disabled}
                onChange={onChange}
                onBlur={onBlur}
                data-valid={valid.valid ? 'valid' : 'invalid'}
                rows={props.rows}
            />
        );
    };

    const content = () => {
        return (
            <React.Fragment>
                { input()}
                { validationError()}
            </React.Fragment>
        );
    };

    const requiredFlag = () => {
        if (!props.required) {
            return null;
        }
        return (<span className="form-required">*</span>);
    };

    return (
        <FormGroup thin={props.thin}>
            <label htmlFor={labelID} className="form-label">{props.label} {requiredFlag()}</label>
            { content()}
            { helpText()}
        </FormGroup>
    );
};
