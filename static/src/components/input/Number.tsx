import * as React from 'react';
import debounce = require('debounce-promise');
import { FormGroup, ValidationResult } from '../Form';
import { Rand } from '../../services/Rand';
import '../../../css/form.scss';

interface NumberProps {
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
    onChange: (value: number) => (void);
    /**
     * The default value used for the input
     */
    defaultValue: number;
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
     * The minimum valid value for this input. Will be validated automatically if specified.
     */
    minimum?: number;
    /**
     * The maximum valid value for this input. Will be validated automatically if specified.
     */
    maximum?: number;
    /**
     * Text label to appear before the input
     */
    prepend?: string;
    /**
     * Text label to appear after the input
     */
    append?: string;
    /**
     * Optional method to invoke for validating the value of this input.
     * Return a promise that resolves with a validation result.
     *
     * You do not need to validate if a required field has any value, that is done automatically.
     */
    validate?: (value: number) => Promise<ValidationResult>;
}
export const Number: React.FC<NumberProps> = (props: NumberProps) => {
    const [value, setValue] = React.useState<string>(props.defaultValue ? props.defaultValue.toString() : '');
    const initialValidState: ValidationResult = {
        valid: true
    };
    if (props.required && props.defaultValue == null) {
        initialValidState.valid = false;
        initialValidState.invalidMessage = 'A numeric value is required';
    }
    const [valid, setValid] = React.useState<ValidationResult>(initialValidState);
    const [touched, setTouched] = React.useState(false);
    const labelID = Rand.ID();

    React.useEffect(() => {
        props.onChange(parseInt(value));
    }, [value]);

    const debouncedValidate = debounce(props.validate, 250);
    const validate = (value: number): Promise<ValidationResult> => {
        return new Promise((resolve) => {
            if (props.required && (isNaN(value) || value == null)) {
                resolve({
                    valid: false,
                    invalidMessage: 'A numeric value is required',
                });
                return;
            }
            if (!isNaN(props.minimum) && value < props.minimum) {
                resolve({
                    valid: false,
                    invalidMessage: 'Value must be at least ' + props.minimum,
                });
                return;
            }
            if (!isNaN(props.maximum) && value > props.maximum) {
                resolve({
                    valid: false,
                    invalidMessage: 'Value must be less than ' + props.maximum,
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
        validate(parseInt(target.value)).then(valid => {
            setValid(valid);
        });
        setValue(target.value);
        props.onChange(parseFloat(target.value));
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
        let defaultValue = '';
        if (!isNaN(props.defaultValue)) {
            defaultValue = props.defaultValue.toString();
        }
        let className = 'form-control';
        if (touched && !valid.valid) {
            className += ' is-invalid';
        }
        return (
            <input
                type="text"
                inputMode="numeric"
                pattern="[0-9]*"
                className={className}
                id={labelID}
                placeholder={props.placeholder}
                defaultValue={defaultValue}
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
                    { input() }
                    { validationError() }
                </React.Fragment>
            );
        }

        let prepend: JSX.Element = null;
        if (props.prepend) {
            prepend = ( <span className="input-group-text">{props.prepend}</span> );
        }
        let append: JSX.Element = null;
        if (props.append) {
            append = ( <span className="input-group-text">{props.append}</span> );
        }

        return (
            <div className="input-group">
                {prepend}
                { input() }
                {append}
                { validationError() }
            </div>
        );
    };

    const requiredFlag = () => {
        if (!props.required) {
            return null; 
        }
        return (<span className="form-required">*</span>);
    };

    return (
        <FormGroup>
            <label htmlFor={labelID} className="form-label">{props.label} {requiredFlag()}</label>
            { content() }
            { helpText() }
        </FormGroup>
    );
};
