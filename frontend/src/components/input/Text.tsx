import * as React from 'react';
import debounce = require('debounce-promise');
import { FormGroup, ValidationResult } from '../Form';
import { Rand } from '../../services/Rand';
import { InputProps } from './Input';
import '../../../css/form.scss';

interface TextProps extends InputProps {
    /**
     * The label that appears above the input
     */
    label?: string;
    /**
     * The value used in the type attribute on the input node
     */
    type: 'text' | 'password' | 'email' | 'search';
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
     * Text label to appear before the input
     */
    prepend?: JSX.Element | string;
    /**
     * Text label to appear after the input
     */
    append?: JSX.Element | string;
    /**
     * If true then a fixed width font is used
     */
    fixedWidth?: boolean;
    /**
     * Optional method to invoke for validating the value of this input.
     * Return a promise that resolves with a validation result.
     *
     * You do not need to validate if a required field has any value, that is done automatically.
     */
    validate?: (value: string) => Promise<ValidationResult>;
}
export const Text: React.FC<TextProps> = (props: TextProps) => {
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

    const onChange = (event: React.FormEvent<HTMLInputElement>) => {
        const target = event.target as HTMLInputElement;
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
            <input
                type={props.type}
                className={className}
                id={labelID}
                placeholder={props.placeholder}
                defaultValue={value}
                disabled={props.disabled}
                onChange={onChange}
                onBlur={onBlur}
                data-valid={valid.valid ? 'valid' : 'invalid'}
            />
        );
    };

    const content = () => {
        if (!props.prepend && !props.append) {
            return (
                <React.Fragment>
                    {input()}
                    {validationError()}
                </React.Fragment>
            );
        }

        let prepend: JSX.Element = null;
        if (props.prepend) {
            prepend = (<span className="input-group-text">{props.prepend}</span>);
        }
        let append: JSX.Element = null;
        if (props.append) {
            append = (<span className="input-group-text">{props.append}</span>);
        }

        return (
            <div className="input-group">
                {prepend}
                {input()}
                {append}
                {validationError()}
            </div>
        );
    };

    const label = () => {
        if (!props.label) {
            return null;
        }

        const requiredFlag = () => {
            if (!props.required) {
                return null;
            }
            return (<span className="form-required">*</span>);
        };

        return (<label htmlFor={labelID} className="form-label">{props.label} {requiredFlag()}</label>);
    };

    return (
        <FormGroup thin={props.thin}>
            { label() }
            { content() }
            { helpText() }
        </FormGroup>
    );
};
