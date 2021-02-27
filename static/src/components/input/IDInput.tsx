import * as React from 'react';
import { Rand } from '../../services/Rand';
import { ValidationResult, FormGroup } from '../Form';

interface IDInputProps {
    label: string;
    defaultUID: number;
    defaultGID: number;
    onChange: (uid: number, gid: number) => (void);
}
export const IDInput: React.FC<IDInputProps> = (props: IDInputProps) => {
    const [uid, setUID] = React.useState<string>((props.defaultUID || 0).toString());
    const [gid, setGID] = React.useState<string>((props.defaultGID || 0).toString());
    const [valid, setValid] = React.useState<ValidationResult>({valid: true});
    const [touched, setTouched] = React.useState(false);

    React.useEffect(() => {
        const fuid = parseFloat(uid);
        if (isNaN(fuid)) {
            return;
        }
        const fgid = parseFloat(gid);
        if (isNaN(fgid)) {
            return;
        }

        props.onChange(fuid, fgid);
    }, [uid, gid]);

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

    const onChangeUID = (event: React.FormEvent<HTMLInputElement>) => {
        const target = event.target as HTMLInputElement;
        validate(parseInt(target.value)).then(valid => {
            setValid(valid);
        });
        setUID(target.value);
    };

    const onChangeGID = (event: React.FormEvent<HTMLInputElement>) => {
        const target = event.target as HTMLInputElement;
        validate(parseInt(target.value)).then(valid => {
            setValid(valid);
        });
        setGID(target.value);
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
                defaultValue={uid}
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
                defaultValue={gid}
                onChange={onChangeGID}
                onBlur={onBlur}
                data-valid={valid.valid ? 'valid' : 'invalid'}
                required />
        </div>
        {validationError()}
    </FormGroup>);
};