import * as React from 'react';
import { FormGroup } from '../Form';
import { Rand } from '../../services/Rand';
import '../../../css/form.scss';

interface FileChooserProps {
    /**
     * The label that appears above the input
     */
    label: string;
    /**
     * Additional text to show below the input
     */
    helpText?: string;
    /**
     * Event called when a file is selected
     */
    onChange: (file: File) => (void);
}
export const FileChooser: React.FC<FileChooserProps> = (props: FileChooserProps) => {
    const labelID = Rand.ID();

    const didSelectFile = (event: React.ChangeEvent<HTMLInputElement>) => {
        const file = event.target.files[0];
        props.onChange(file);
    };

    const helpText = () => {
        if (props.helpText) {
            return <div id={labelID + 'help'} className="form-text">{props.helpText}</div>;
        } else {
            return null;
        }
    };

    return (
        <FormGroup>
            <label className="form-label" htmlFor={labelID}>{props.label}</label>
            <div className="form-file">
                <input className="form-control" type="file" id={labelID} onChange={didSelectFile} />
            </div>
            { helpText()}
        </FormGroup>
    );
};
