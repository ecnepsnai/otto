import * as React from 'react';
import { Rand } from '../services/Rand';
import { Style } from './Style';
import { Button } from './Button';
import { Icon } from './Icon';
import debounce = require('debounce-promise');
import '../../css/form.scss';

export interface ValidationResult {
    valid: boolean;
    invalidMessage?: string;
}

export interface FormProps {
    /**
     * Should a save button be appended to the bottom of the form
     */
    showSaveButton?: boolean;
    /**
     * Event fired when the form is submitted
     */
    onSubmit?: () => void;
    /**
     * Optional class to add
     */
    className?: string;
    /**
     * If true the submit button is disabled
     */
    loading?: boolean;
}

interface FormState {
    invalid?: boolean;
}

/**
 * A form is a collection of inputs and a submit button
 */
export class Form extends React.Component<FormProps, FormState> {
    private domRef: React.RefObject<HTMLFormElement>;

    constructor(props: FormProps) {
        super(props);
        this.domRef = React.createRef();
        this.state = {};
    }

    /**
     * Performs validation on the form
     * @returns true if the form is valid, false if invalid
     */
    public validateForm = (): boolean => {
        const elemn = this.domRef.current;
        const invalidNodes = elemn.querySelectorAll('[data-valid="invalid"]');
        if (invalidNodes.length > 0) {
            this.setState({ invalid: true });
            return false;
        }
        return true;
    }

    private onClick = () => {
        this.submitForm();
    }

    private onSubmut = (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault();
        this.submitForm();
    }

    private submitForm = () => {
        this.validateForm();

        if (this.props.onSubmit) {
            this.props.onSubmit();
        }
    }

    private saveButton = () => {
        if (!this.props.showSaveButton) {
            return null;
        }

        let content = (<Icon.Label icon={<Icon.CheckCircle/>} label="Apply"/>);
        if (this.props.loading) {
            content = (<Icon.Label icon={<Icon.Spinner pulse/>} label="Loading..."/>);
        }

        return (
            <div className="mt-3">
                <Button color={Style.Palette.Primary} size={Style.Size.M} onClick={this.onClick} disabled={this.props.loading}>
                    { content }
                </Button>
            </div>
        );
    }

    private error = () => {
        if (!this.state.invalid) { return null; }

        return (
            <div className="mt-2">
                <Icon.Label icon={<Icon.TimesCircle color={Style.Palette.Danger}/>} label="Correct Errors Before Contuining" />
            </div>
        );
    }

    render(): JSX.Element {
        const className = this.props.className || '';
        return (
        <form onSubmit={this.onSubmut} ref={this.domRef} className={className}>
            <fieldset>{ this.props.children }</fieldset>
            { this.saveButton() }
            { this.error() }
        </form>
        );
    }
}

export interface FormGroupProps { className?: string }
export class FormGroup extends React.Component<FormGroupProps, {}> {
    render(): JSX.Element {
        return ( <div className={ (this.props.className ?? '') + ' mb-3'}>{ this.props.children }</div> );
    }
}

export interface InputProps {
    /**
     * The label that appears above the input
     */
    label: string;
    /**
     * The value used in the type attribute on the input node
     */
    type: string;
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
    validate?: (value: string) => Promise<ValidationResult>;
    /**
     * If true then a fixed width font is used
     */
    fixedWidth?: boolean;
}

interface InputState { value: string; labelID: string; valid: ValidationResult; touched: boolean; }

/**
 * An input node for regular text inputs. This component is only sutible for text or password types.
 */
export class Input extends React.Component<InputProps, InputState> {
    constructor(props: InputProps) {
        super(props);
        const initialValidState: ValidationResult = {
            valid: true
        };
        if (props.required && !props.defaultValue) {
            initialValidState.valid = false;
            initialValidState.invalidMessage = 'A value is required';
        }
        this.state = { value: '', labelID: Rand.ID(), valid: initialValidState, touched: false, };
    }
    private debouncedValidate = debounce(this.props.validate, 250);
    private validate = (value: string): Promise<ValidationResult> => {
        return new Promise((resolve) => {
            if (this.props.required && value == '') {
                resolve({
                    valid: false,
                    invalidMessage: 'A value is required'
                });
                return;
            }
            if (this.props.validate) {
                this.debouncedValidate(value).then(valid => {
                    resolve(valid);
                });
                return;
            }
            resolve({ valid: true });
        });
    }
    private onBlur = () => {
        this.setState({ touched: true });
    }
    private onChange = (event: React.FormEvent<HTMLInputElement>) => {
        const target = event.target as HTMLInputElement;
        this.validate(target.value).then(valid => {
            this.setState({ valid: valid });
        });
        this.setState({ value: target.value });
        this.props.onChange(target.value);
    }
    private helpText() {
        if (this.props.helpText) {
            return <div id={this.state.labelID + 'help'} className="form-text">{this.props.helpText}</div>;
        } else {
            return null;
        }
    }
    private validationError() {
        if (!this.state.valid.invalidMessage || !this.state.touched) { return null; }
        return (<div className="invalid-feedback">{this.state.valid.invalidMessage}</div>);
    }
    private input = () => {
        let className = 'form-control';
        if (this.state.touched && !this.state.valid.valid) {
            className += ' is-invalid';
        }
        if (this.props.fixedWidth) {
            className += ' fixed-width';
        }
        return (
            <input
                type={this.props.type}
                className={className}
                id={this.state.labelID}
                placeholder={this.props.placeholder}
                defaultValue={this.props.defaultValue}
                disabled={this.props.disabled}
                onChange={this.onChange}
                onBlur={this.onBlur}
                data-valid={this.state.valid.valid ? 'valid' : 'invalid'}
            />
        );
    }
    private content = () => {
        if (!this.props.prepend && !this.props.append) {
            return (
                <React.Fragment>
                    { this.input() }
                    { this.validationError() }
                </React.Fragment>
            );
        }

        let prepend: JSX.Element = null;
        if (this.props.prepend) {
            prepend = ( <span className="input-group-text">{this.props.prepend}</span> );
        }
        let append: JSX.Element = null;
        if (this.props.append) {
            append = ( <span className="input-group-text">{this.props.append}</span> );
        }

        return (
        <div className="input-group">
            {prepend}
            { this.input() }
            {append}
            { this.validationError() }
        </div>
        );
    };
    private requiredFlag = () => {
        if (!this.props.required) { return null; }
        return (<span className="form-required">*</span>);
    }
    render(): JSX.Element {
        return (
            <FormGroup>
                <label htmlFor={this.state.labelID} className="form-label">{this.props.label} {this.requiredFlag()}</label>
                { this.content() }
                { this.helpText() }
            </FormGroup>
        );
    }
}

export interface NumberInputProps {
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

interface NumberInputState { value: string; labelID: string; valid: ValidationResult; touched: boolean; }

/**
 * An input node for number type inputs. Sutible for integer and floating points, but not sutible
 * for hexedecimal or scientific-notation values.
 */
export class NumberInput extends React.Component<NumberInputProps, NumberInputState> {
    constructor(props: NumberInputProps) {
        super(props);
        const initialValidState: ValidationResult = {
            valid: true
        };
        if (props.required && props.defaultValue == null) {
            initialValidState.valid = false;
            initialValidState.invalidMessage = 'A numeric value is required';
        }
        this.state = { value: '', labelID: Rand.ID(), valid: initialValidState, touched: false, };
    }
    private debouncedValidate = debounce(this.props.validate, 250);
    private validate = (value: number): Promise<ValidationResult> => {
        return new Promise((resolve) => {
            if (this.props.required && (isNaN(value) || value == null)) {
                resolve({
                    valid: false,
                    invalidMessage: 'A numeric value is required',
                });
                return;
            }
            if (!isNaN(this.props.minimum) && value < this.props.minimum) {
                resolve({
                    valid: false,
                    invalidMessage: 'Value must be at least ' + this.props.minimum,
                });
                return;
            }
            if (!isNaN(this.props.maximum) && value > this.props.maximum) {
                resolve({
                    valid: false,
                    invalidMessage: 'Value must be less than ' + this.props.maximum,
                });
                return;
            }
            if (this.props.validate) {
                this.debouncedValidate(value).then(valid => {
                    resolve(valid);
                });
                return;
            }
            resolve({ valid: true });
        });
    }
    private onBlur = () => {
        this.setState({ touched: true });
    }
    private onChange = (event: React.FormEvent<HTMLInputElement>) => {
        const target = event.target as HTMLInputElement;
        this.validate(parseInt(target.value)).then(valid => {
            this.setState({ valid: valid });
        });
        this.setState({ value: target.value });
        this.props.onChange(parseFloat(target.value));
    }
    private helpText() {
        if (this.props.helpText) {
            return <div id={this.state.labelID + 'help'} className="form-text">{this.props.helpText}</div>;
        } else {
            return null;
        }
    }
    private validationError() {
        if (!this.state.valid.invalidMessage || !this.state.touched) { return null; }
        return (<div className="invalid-feedback">{this.state.valid.invalidMessage}</div>);
    }
    private input = () => {
        let defaultValue = '';
        if (!isNaN(this.props.defaultValue)) {
            defaultValue = this.props.defaultValue.toString();
        }
        let className = 'form-control';
        if (this.state.touched && !this.state.valid.valid) {
            className += ' is-invalid';
        }
        return (
            <input
                type="text"
                inputMode="numeric"
                pattern="[0-9]*"
                className={className}
                id={this.state.labelID}
                placeholder={this.props.placeholder}
                defaultValue={defaultValue}
                disabled={this.props.disabled}
                onChange={this.onChange}
                onBlur={this.onBlur}
                data-valid={this.state.valid.valid ? 'valid' : 'invalid'}
            />
        );
    }
    private content = () => {
        if (!this.props.prepend && !this.props.append) {
            return (
                <React.Fragment>
                    { this.input() }
                    { this.validationError() }
                </React.Fragment>
            );
        }

        let prepend: JSX.Element = null;
        if (this.props.prepend) {
            prepend = ( <span className="input-group-text">{this.props.prepend}</span> );
        }
        let append: JSX.Element = null;
        if (this.props.append) {
            append = ( <span className="input-group-text">{this.props.append}</span> );
        }

        return (
        <div className="input-group">
            {prepend}
            { this.input() }
            {append}
            { this.validationError() }
        </div>
        );
    };
    private requiredFlag = () => {
        if (!this.props.required) { return null; }
        return (<span className="form-required">*</span>);
    }
    render(): JSX.Element {
        return (
            <FormGroup>
                <label htmlFor={this.state.labelID} className="form-label">{this.props.label} {this.requiredFlag()}</label>
                { this.content() }
                { this.helpText() }
            </FormGroup>
        );
    }
}

export interface SelectProps {
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
}

interface SelectState { value: string; labelID: string; valid: ValidationResult; touched: boolean; }

/**
 * A dropdown, or <select> input, where the user picks a single option from a list
 */
export class Select extends React.Component<SelectProps, SelectState> {
    constructor(props: SelectProps) {
        super(props);
        const initialValidState: ValidationResult = {
            valid: true
        };
        if (props.required && (props.defaultValue == null || props.defaultValue == '')) {
            initialValidState.valid = false;
            initialValidState.invalidMessage = 'A selection is required';
        }
        this.state = { value: props.defaultValue, labelID: Rand.ID(), valid: initialValidState, touched: false };
    }
    private onChange = (event: React.FormEvent<HTMLSelectElement>) => {
        const target = event.target as HTMLSelectElement;
        this.validate(target.value).then(valid => {
            this.setState({ valid: valid });
        });
        this.setState({ value: target.value });
        this.props.onChange(target.value);
    }
    private helpText() {
        if (this.props.helpText) {
            return <div id={this.state.labelID + 'help'} className="form-text">{this.props.helpText}</div>;
        } else {
            return null;
        }
    }
    private requiredFlag = () => {
        if (!this.props.required) { return null; }
        return (<span className="form-required">*</span>);
    }
    private defaultSelection = () => {
        if (this.props.required && this.state.value) {
            return null;
        }

        return (<option selected>Select One...</option>);
    }
    private validate = (value: string): Promise<ValidationResult> => {
        return new Promise((resolve) => {
            if (this.props.required && value == '') {
                resolve({
                    valid: false,
                    invalidMessage: 'A value is required'
                });
                return;
            }
            if (this.props.validate) {
                this.props.validate(value).then(valid => {
                    resolve(valid);
                });
                return;
            }
            resolve({ valid: true });
        });
    }
    private onBlur = () => {
        this.setState({ touched: true });
    }
    private validationError() {
        if (!this.state.valid.invalidMessage || !this.state.touched) { return null; }
        return (<div className="invalid-feedback">{this.state.valid.invalidMessage}</div>);
    }
    render(): JSX.Element {
        let className = 'form-select';
        if (this.state.touched && !this.state.valid.valid) {
            className += ' is-invalid';
        }
        return (
            <FormGroup>
                <label htmlFor={this.state.labelID} className="form-label">{this.props.label} {this.requiredFlag()}</label>
                <select
                    defaultValue={this.props.defaultValue}
                    className={className}
                    id={this.state.labelID}
                    onChange={this.onChange}
                    disabled={this.props.disabled}
                    onBlur={this.onBlur}
                    data-valid={this.state.valid.valid ? 'valid' : 'invalid'}>
                        { this.defaultSelection() }
                        { this.props.children }
                </select>
                { this.validationError() }
                { this.helpText() }
            </FormGroup>
        );
    }
}

/**
 * Describes the properties for a checkbox
 */
export interface CheckboxProps {
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

interface CheckboxState { checked: boolean; labelID: string; }

/**
 * A checkbox
 */
export class Checkbox extends React.Component<CheckboxProps, CheckboxState> {
    constructor(props: CheckboxProps) {
        super(props);
        this.state = { checked: false, labelID: Rand.ID() };
    }
    private onChange = (event: React.FormEvent<HTMLInputElement>) => {
        const target = event.target as HTMLInputElement;
        this.setState({ checked: target.checked });
        this.props.onChange(target.checked);
    }
    private helpText() {
        if (this.props.helpText) {
            return <div id={this.state.labelID + 'help'} className="form-text">{this.props.helpText}</div>;
        } else {
            return null;
        }
    }
    render(): JSX.Element {
        return (
            <FormGroup className="form-check">
                <input type="checkbox" className="form-check-input" id={this.state.labelID} checked={this.props.checked} defaultChecked={this.props.defaultValue} onChange={this.onChange} disabled={this.props.disabled}/>
                <label htmlFor={this.state.labelID} className="form-check-label">{this.props.label}</label>
                { this.helpText() }
            </FormGroup>
        );
    }
}

/**
 * Describes the properties for an <textarea> node
 */
export interface TextareaProps {
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
    /**
     * If true then a fixed width font is used
     */
    fixedWidth?: boolean;
    /**
     * Number of rows
     */
    rows?: number;
}

interface TextareaState { value: string; labelID: string; valid: ValidationResult; touched: boolean; }

/**
 * A multi-line textarea field
 */
export class Textarea extends React.Component<TextareaProps, TextareaState> {
    constructor(props: TextareaProps) {
        super(props);
        const initialValidState: ValidationResult = {
            valid: true
        };
        if (props.required && !props.defaultValue) {
            initialValidState.valid = false;
            initialValidState.invalidMessage = 'A value is required';
        }
        this.state = { value: '', labelID: Rand.ID(), valid: initialValidState, touched: false, };
    }
    private debouncedValidate = debounce(this.props.validate, 250);
    private validate = (value: string): Promise<ValidationResult> => {
        return new Promise((resolve) => {
            if (this.props.required && value == '') {
                resolve({
                    valid: false,
                    invalidMessage: 'A value is required'
                });
                return;
            }
            if (this.props.validate) {
                this.debouncedValidate(value).then(valid => {
                    resolve(valid);
                });
                return;
            }
            resolve({ valid: true });
        });
    }
    private onBlur = () => {
        this.setState({ touched: true });
    }
    private onChange = (event: React.FormEvent<HTMLTextAreaElement>) => {
        const target = event.target as HTMLTextAreaElement;
        this.validate(target.value).then(valid => {
            this.setState({ valid: valid });
        });
        this.setState({ value: target.value });
        this.props.onChange(target.value);
    }
    private helpText() {
        if (this.props.helpText) {
            return <div id={this.state.labelID + 'help'} className="form-text">{this.props.helpText}</div>;
        } else {
            return null;
        }
    }
    private validationError() {
        if (!this.state.valid.invalidMessage || !this.state.touched) { return null; }
        return (<div className="invalid-feedback">{this.state.valid.invalidMessage}</div>);
    }
    private requiredFlag = () => {
        if (!this.props.required) { return null; }
        return (<span className="form-required">*</span>);
    }
    render(): JSX.Element {
        let className = 'form-control';
        if (this.state.touched && !this.state.valid) {
            className += ' is-invalid';
        }
        if (this.props.fixedWidth) {
            className += ' fixed-width';
        }
        return (
            <FormGroup>
                <label htmlFor={this.state.labelID} className="form-label">{this.props.label} {this.requiredFlag()}</label>
                <textarea
                    className={className}
                    id={this.state.labelID}
                    placeholder={this.props.placeholder}
                    defaultValue={this.props.defaultValue}
                    onChange={this.onChange}
                    disabled={this.props.disabled}
                    onBlur={this.onBlur}
                    data-valid={this.state.valid ? 'valid' : 'invalid'}
                    rows={this.props.rows}
                />
                { this.validationError() }
                { this.helpText() }
            </FormGroup>
        );
    }
}

/**
 * Describes the properties for a choice in a group of radio buttons
 */
export interface RadioChoice {
    value: string|number;
    label: string;
}
/**
 * Describes the properties for radio input
 */
export interface RadioProps {
    /**
     * The label that appears above the input
     */
    label: string;
    /**
     * The choices for the input
     */
    choices: RadioChoice[];
    /**
     * The default value to be selected
     */
    defaultValue?: string|number;
    /**
     * Called when a new value is selected
     */
    onChange: (value: string|number) => (void);
    /**
     * If toggle buttons should be used instead of classic radio controls
     */
    buttons?: boolean;
}
interface RadioState {
    value: string|number;
}
/**
 * A group of radio buttons for selecting a single choice from a list
 */
export class Radio extends React.Component<RadioProps, RadioState> {
    constructor(props: RadioProps) {
        super(props);

        if (props.defaultValue != undefined) {
            let found = false;
            props.choices.forEach(choice => {
                if (choice.value === props.defaultValue) {
                    found = true;
                }
            });
            if (!found) {
                throw new Error('default value not a valid choice');
            }
        }

        this.state = {
            value: props.defaultValue,
        };
    }

    componentDidUpdate(props: RadioProps): void {
        if (props.defaultValue !== this.props.defaultValue) {
            this.setState({ value: this.props.defaultValue }, () => {
                this.props.onChange(this.props.defaultValue);
            });
        }
    }

    private onChange = (event: React.FormEvent<HTMLInputElement>) => {
        const target = event.target as HTMLInputElement;
        if (target.checked) {
            this.setState({ value: target.value }, () => {
                this.props.onChange(target.value);
            });
        }
    }

    private input = () => {
        return (
            <React.Fragment>
                {
                    this.props.choices.map(choice => {
                        const labelID = Rand.ID();
                        return (
                            <div className="form-check" key={labelID}>
                                <input className="form-check-input" type="radio" name={labelID} id={labelID} value={choice.value} checked={this.state.value===choice.value} onChange={this.onChange}/>
                                <label className="form-check-label" htmlFor={labelID}>
                                    {choice.label}
                                </label>
                            </div>
                        );
                    })
                }
            </React.Fragment>
        );
    }

    private buttons = () => {
        return (
            <div>
                <div className="btn-group">
                    {
                        this.props.choices.map(choice => {
                            const labelID = Rand.ID();
                            return (
                                <React.Fragment key={labelID}>
                                    <input type="radio" className="btn-check" name={labelID} id={labelID} value={choice.value} checked={this.state.value===choice.value} onChange={this.onChange} />
                                    <label className="btn btn-secondary btn-sm" htmlFor={labelID}>{choice.label}</label>
                                </React.Fragment>
                            );
                        })
                    }
                </div>
            </div>
        );
    }

    render(): JSX.Element {
        let content: JSX.Element;
        if (this.props.buttons) {
            content = this.buttons();
        } else {
            content = this.input();
        }

        return (
            <div className="mb-3">
                <label className="form-label">{this.props.label}</label>
                {content}
            </div>
        );
    }
}

export interface FileBrowserProps {
    /**
     * The label that appears above the input
     */
    label: string;
    /**
     * Event called when a file is selected
     */
    onChange: (file: File) => (void);
}

interface FileBrowserState {
    fileName?: string;
}

export class FileBrowser extends React.Component<FileBrowserProps, FileBrowserState> {
    constructor(props: FileBrowserProps) {
        super(props);
        this.state = {};
    }

    private didSelectFile = (event: React.ChangeEvent<HTMLInputElement>) => {
        const file = event.target.files[0];
        this.props.onChange(file);
        this.setState({ fileName: file.name });
    }

    render(): JSX.Element {
        const fileLabel = this.state.fileName || 'Choose file...';
        return (
            <div className="mb-3">
                <label className="form-label">{this.props.label}</label>
                <div className="form-file">
                    <input type="file" className="form-file-input" id="customFile" onChange={this.didSelectFile}/>
                    <label className="form-file-label" htmlFor="customFile">
                        <span className="form-file-text">{fileLabel}</span>
                        <span className="form-file-button">Browse</span>
                    </label>
                </div>
            </div>
        );
    }
}
