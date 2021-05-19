import * as React from 'react';
import { Rand } from '../../services/Rand';
import { RunAs } from '../../types/Script';
import { ValidationResult, FormGroup } from '../Form';
import { Input } from './Input';

interface RunAsInputProps {
    label: string;
    inheritLabel: string;
    defaultValue: RunAs;
    onChange: (runAs: RunAs) => (void);
}
export const RunAsInput: React.FC<RunAsInputProps> = (props: RunAsInputProps) => {
    const [runAs, setRunAs] = React.useState<RunAs>(props.defaultValue);
    const [valid, setValid] = React.useState<ValidationResult>({ valid: true });
    const [touched, setTouched] = React.useState(false);

    React.useEffect(() => {
        if (!runAs) {
            return;
        }
        props.onChange(runAs);
    }, [runAs]);

    const validate = (value: number): Promise<ValidationResult> => {
        return new Promise((resolve) => {
            if (isNaN(value) || value == null) {
                resolve({
                    valid: false,
                    invalidMessage: 'A numeric value is required',
                });
                return;
            }
            if (!isNaN(0) && value < 0) {
                resolve({
                    valid: false,
                    invalidMessage: 'Value must be at least ' + 0,
                });
                return;
            }
            resolve({ valid: true });
        });
    };

    const onBlur = () => {
        setTouched(true);
    };

    const onChangeInherit = (newValue: boolean) => {
        setRunAs(runAs => {
            runAs.Inherit = !newValue;
            return { ...runAs };
        });
    };

    const onChangeUID = (event: React.FormEvent<HTMLInputElement>) => {
        const target = event.target as HTMLInputElement;
        validate(parseInt(target.value)).then(valid => {
            setValid(valid);
        });
        setRunAs(runAs => {
            runAs.UID = parseInt(target.value);
            return { ...runAs };
        });
    };

    const onChangeGID = (event: React.FormEvent<HTMLInputElement>) => {
        const target = event.target as HTMLInputElement;
        validate(parseInt(target.value)).then(valid => {
            setValid(valid);
        });
        setRunAs(runAs => {
            runAs.GID = parseInt(target.value);
            return { ...runAs };
        });
    };

    const validationError = () => {
        if (!valid.invalidMessage || !touched) {
            return null;
        }
        return (<div className="invalid-feedback force-show">{valid.invalidMessage}</div>);
    };

    const labelID = Rand.ID();
    let className = 'form-control';
    if (touched && !valid) {
        className += ' is-invalid';
    }

    const idinputs = () => {
        if (runAs.Inherit) {
            return null;
        }

        return (<FormGroup>
            <label htmlFor={labelID} className="form-label">{props.label} <span className="form-required">*</span></label>
            <div className="input-group">
                <span className="input-group-text">UID</span>
                <input
                    type="text"
                    inputMode="numeric"
                    pattern="[0-9]*"
                    className={className}
                    id={labelID}
                    placeholder="0"
                    defaultValue={runAs.UID}
                    onChange={onChangeUID}
                    onBlur={onBlur}
                    data-valid={valid.valid ? 'valid' : 'invalid'}
                    required />
                <span className="input-group-text">GID</span>
                <input
                    type="text"
                    inputMode="numeric"
                    pattern="[0-9]*"
                    className={className}
                    placeholder="0"
                    defaultValue={runAs.GID}
                    onChange={onChangeGID}
                    onBlur={onBlur}
                    data-valid={valid.valid ? 'valid' : 'invalid'}
                    required />
            </div>
            {validationError()}
        </FormGroup>);
    };

    return (
        <React.Fragment>
            <Input.Checkbox label={props.inheritLabel} defaultValue={!runAs.Inherit} onChange={onChangeInherit} />
            {idinputs()}
        </React.Fragment>
    );
};