import * as React from 'react';
import { Rand } from '../services/Rand';
import { FormGroup } from './Form';
import { Icon } from './Icon';

/**
 * Describes the properties for a multi-input
 */
export interface MultiInputProps {
    /** Label for the input */
    label: string;
    /** Optional placeholder for empty inputs */
    placeholder?: string;
    /** Event for when the model value changes */
    onChange: (values: string[]) => (void);
    /** The default value */
    defaultValue: string[];
    /** Optional help text to show below the input */
    helpText?: string;
    /** The maximum number of values that can be added */
    max?: number;
}

interface MultiInputState { values: string[]; labelID: string; }

/**
 * A "multi" input, which is a list of text fields with add/remove buttons that accepts an array of strings
 * as its model
 */
export class MultiInput extends React.Component<MultiInputProps, MultiInputState> {
    constructor(props: MultiInputProps) {
        super(props);

        let values = this.props.defaultValue;
        if (!values || values.length == 0) {
            values = [''];
        }

        this.state = { values: values, labelID: Rand.ID() };
    }

    private onChange = (value: string, index: number) => {
        const values = this.state.values;
        values[index] = value;
        this.setState({ values: values });
        this.props.onChange(values);
    }

    private showAddButton = (): boolean => {
        if (this.props.max > 0) {
            return this.state.values.length < this.props.max;
        }

        return true;
    }

    private showRemoveButton = (index: number): boolean => {
        return index > 0;
    }

    private addButtonClicked = (index: number) => {
        this.setState(state => {
            const values = state.values;
            values.splice(index + 1, 0, '');
            return { values: values };
        });
    }

    private removeButtonClicked = (index: number) => {
        this.setState(state => {
            const values = state.values;
            values.splice(index, 1);
            return { values: values };
        });
    }

    private helpText(): JSX.Element {
        if (this.props.helpText) {
            return <div id={this.state.labelID + 'help'} className="form-text">{this.props.helpText}</div>;
        } else {
            return null;
        }
    }
    render(): JSX.Element {
        return (
            <FormGroup>
                <label htmlFor={this.state.labelID} className="form-label">{this.props.label}</label>
                {
                    this.state.values.map((value, idx) => {
                        return <InputGroup
                            key={idx}
                            value={value}
                            index={idx}
                            placeholder={this.props.placeholder}
                            onChange={this.onChange}
                            showAddButton={this.showAddButton}
                            showRemoveButton={this.showRemoveButton}
                            addButtonClicked={this.addButtonClicked}
                            removeButtonClicked={this.removeButtonClicked}
                        ></InputGroup>;
                    })
                }
                { this.helpText() }
            </FormGroup>
        );
    }
}

interface FieldProps {
    value: string;
    onChange: (value: string) => void;
    placeholder?: string;
}
class Field extends React.Component<FieldProps, {}> {
    private onChange = (event: React.FormEvent<HTMLInputElement>) => {
        const target = event.target as HTMLInputElement;
        this.props.onChange(target.value);
    }
    render(): JSX.Element {
        return <input className="form-control" value={this.props.value} placeholder={this.props.placeholder} onChange={this.onChange}/>;
    }
}

interface AddButtonProps {
    onClick: () => void;
}
class AddButton extends React.Component<AddButtonProps, {}> {
    private onClick = () => {
        this.props.onClick();
    }
    render(): JSX.Element {
        return <button className="btn btn-outline-secondary" type="button" onClick={this.onClick}><Icon.Plus /></button>;
    }
}

interface RemoveButtonProps {
    onClick: () => void;
}
class RemoveButton extends React.Component<RemoveButtonProps, {}> {
    private onClick = () => {
        this.props.onClick();
    }
    render(): JSX.Element {
        return <button className="btn btn-outline-secondary" type="button" onClick={this.onClick}><Icon.Minus /></button>;
    }
}

interface InputGroupProps {
    value: string;
    index: number;
    placeholder?: string;
    onChange: (value: string, index: number) => void;
    showAddButton: (index: number) => boolean;
    showRemoveButton: (index: number) => boolean;
    addButtonClicked: (index: number) => void;
    removeButtonClicked: (index: number) => void;
}
class InputGroup extends React.Component<InputGroupProps, {}> {
    private onChange = (value: string) => {
        this.props.onChange(value, this.props.index);
    }
    private addButtonClicked = () => {
        this.props.addButtonClicked(this.props.index);
    }
    private removeButtonClicked = () => {
        this.props.removeButtonClicked(this.props.index);
    }
    private addButton(): JSX.Element {
        if (this.props.showAddButton(this.props.index)) {
            return <AddButton onClick={this.addButtonClicked}/>;
        }

        return null;
    }
    private removeButton(): JSX.Element {
        if (this.props.showRemoveButton(this.props.index)) {
            return <RemoveButton onClick={this.removeButtonClicked}/>;
        }

        return null;
    }
    render(): JSX.Element {
        let inputClass = '';
        if (this.props.index > 0) {
            inputClass = 'mt-2';
        }
        return <div className={'input-group ' + inputClass}>
            <Field value={this.props.value} onChange={this.onChange} placeholder={this.props.placeholder}/>
            { this.removeButton() }
            { this.addButton() }
        </div>;
    }
}
