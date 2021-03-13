import * as React from 'react';
import { Rand } from '../services/Rand';
import { Button } from './Button';
import { FormGroup } from './Form';
import { Icon } from './Icon';
import { Style } from './Style';

/**
 * Describes the properties for a multi-input
 */
interface MultiInputProps {
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
    /** If a single value is required */
    required?: boolean;
}

/**
 * A "multi" input, which is a list of text fields with add/remove buttons that accepts an array of strings
 * as its model
 */
export const MultiInput: React.FC<MultiInputProps> = (props: MultiInputProps) => {
    const [values, setValues] = React.useState(props.defaultValue);
    const labelID = Rand.ID();
    const [valid, setValid] = React.useState('valid');

    React.useEffect(() => {
        if (!values || values.length == 0) {
            setValues(['']);
            if (props.required) {
                setValid('invalid');
            }
        }
    }, []);

    React.useEffect(() => {
        return () => {
            let vl = 'valid';
            if (props.required && values.length <= 0) {
                vl = 'invalid';
            }
            setValid(vl);
            props.onChange(values);
        };
    }, [values]);

    const onChange = (value: string, index: number) => {
        setValues(values => {
            values[index] = value;
            values = values.filter(v => {
                return v && v.length > 0;
            });
            return values;
        });
    };

    const showAddButton = (): boolean => {
        if (props.max > 0) {
            return values.length < props.max;
        }

        return true;
    };

    const showRemoveButton = (index: number): boolean => {
        return index > 0;
    };

    const addButtonClicked = (index: number) => {
        setValues(v => {
            v.splice(index + 1, 0, '');
            return [...v];
        });
    };

    const removeButtonClicked = (index: number) => {
        setValues(v => {
            v.splice(index, 1);
            return [...v];
        });
    };

    const helpText = (): JSX.Element => {
        if (props.helpText) {
            return <div id={labelID + 'help'} className="form-text">{props.helpText}</div>;
        }

        return null;
    };

    return (
        <div data-valid={valid}>
            <FormGroup>
                <label htmlFor={labelID} className="form-label">{props.label}</label>
                {
                    values.map((value, idx) => {
                        return <InputGroup
                            key={idx}
                            value={value}
                            index={idx}
                            placeholder={props.placeholder}
                            onChange={onChange}
                            showAddButton={showAddButton}
                            showRemoveButton={showRemoveButton}
                            addButtonClicked={addButtonClicked}
                            removeButtonClicked={removeButtonClicked}/>;
                    })
                }
                { helpText() }
            </FormGroup>
        </div>
    );
};

interface FieldProps {
    value: string;
    onChange: (value: string) => void;
    placeholder?: string;
}
const Field: React.FC<FieldProps> = (props: FieldProps) => {
    const onChange = (event: React.FormEvent<HTMLInputElement>) => {
        const target = event.target as HTMLInputElement;
        props.onChange(target.value);
    };

    return (<input className="form-control" value={props.value} placeholder={props.placeholder} onChange={onChange}/>);
};

interface ButtonProps {
    onClick: () => void;
}
const AddButton: React.FC<ButtonProps> = (props: ButtonProps) => {
    return (<Button color={Style.Palette.Secondary} outline onClick={props.onClick}><Icon.Plus /></Button>);
};
const RemoveButton: React.FC<ButtonProps> = (props: ButtonProps) => {
    return (<Button color={Style.Palette.Secondary} outline onClick={props.onClick}><Icon.Minus /></Button>);
};

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
const InputGroup: React.FC<InputGroupProps> = (props: InputGroupProps) => {
    const onChange = (value: string) => {
        props.onChange(value, props.index);
    };

    const addButtonClicked = () => {
        props.addButtonClicked(props.index);
    };

    const removeButtonClicked = () => {
        props.removeButtonClicked(props.index);
    };

    const addButton = (): JSX.Element => {
        if (props.showAddButton(props.index)) {
            return <AddButton onClick={addButtonClicked}/>;
        }

        return null;
    };

    const removeButton = (): JSX.Element => {
        if (props.showRemoveButton(props.index)) {
            return <RemoveButton onClick={removeButtonClicked}/>;
        }

        return null;
    };

    let inputClass = '';
    if (props.index > 0) {
        inputClass = 'mt-2';
    }
    return (<div className={'input-group ' + inputClass}>
        <Field value={props.value} onChange={onChange} placeholder={props.placeholder}/>
        { removeButton() }
        { addButton() }
    </div>);
};
