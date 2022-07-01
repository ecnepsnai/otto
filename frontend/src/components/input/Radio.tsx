import * as React from 'react';
import { FormGroup } from '../Form';
import { Rand } from '../../services/Rand';
import { InputProps } from './Input';
import '../../../css/form.scss';

export interface RadioChoice {
    value: string | number;
    label: string;
}
interface RadioProps extends InputProps {
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
    defaultValue?: string | number;
    /**
     * Called when a new value is selected
     */
    onChange: (value: string | number) => (void);
    /**
     * If toggle buttons should be used instead of classic radio controls
     */
    buttons?: boolean;
}
export const Radio: React.FC<RadioProps> = (props: RadioProps) => {
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
    const [value, setValue] = React.useState<string | number>(props.defaultValue);

    React.useEffect(() => {
        props.onChange(value);
    }, [value]);

    const onChange = (event: React.FormEvent<HTMLInputElement>) => {
        const target = event.target as HTMLInputElement;
        if (target.checked) {
            setValue(target.value);
        }
    };

    const input = () => {
        return (
            <React.Fragment>
                {
                    props.choices.map(choice => {
                        const labelID = Rand.ID();
                        return (
                            <div className="form-check" key={labelID}>
                                <input className="form-check-input" type="radio" name={labelID} id={labelID} value={choice.value} checked={value === choice.value} onChange={onChange} />
                                <label className="form-check-label" htmlFor={labelID}>
                                    {choice.label}
                                </label>
                            </div>
                        );
                    })
                }
            </React.Fragment>
        );
    };

    const buttons = () => {
        return (
            <div>
                <div className="btn-group">
                    {
                        props.choices.map(choice => {
                            const labelID = Rand.ID();
                            return (
                                <React.Fragment key={labelID}>
                                    <input type="radio" className="btn-check" name={labelID} id={labelID} value={choice.value} checked={value === choice.value} onChange={onChange} />
                                    <label className="btn btn-secondary btn-sm" htmlFor={labelID}>{choice.label}</label>
                                </React.Fragment>
                            );
                        })
                    }
                </div>
            </div>
        );
    };

    let content: JSX.Element;
    if (props.buttons) {
        content = buttons();
    } else {
        content = input();
    }

    return (
        <FormGroup thin={props.thin}>
            <label className="form-label">{props.label}</label>
            {content}
        </FormGroup>
    );
};
