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

    return (
        <FormGroup>
            <label className="form-label" htmlFor={labelID}>{props.label}</label>
            <div className="form-file">
                <input className="form-control" type="file" id={labelID} onChange={didSelectFile} />
            </div>
        </FormGroup>
    );
};
