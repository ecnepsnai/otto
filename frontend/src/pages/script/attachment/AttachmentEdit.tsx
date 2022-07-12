import * as React from 'react';
import { Input } from '../../../components/input/Input';
import { GlobalModalFrame, ModalForm } from '../../../components/Modal';
import { Attachment, AttachmentType } from '../../../types/Attachment';
import { RunAs } from '../../../types/Script';

interface AttachmentEditProps {
    attachment?: AttachmentType;
    didUpdate: (attachment: AttachmentType) => (void);
}
export const AttachmentEdit: React.FC<AttachmentEditProps> = (props: AttachmentEditProps) => {
    const [attachment, setAttachment] = React.useState<AttachmentType>(props.attachment || Attachment.Blank());
    const [file, setFile] = React.useState<File>();

    const saveAttachment = () => {
        if (props.attachment) {
            return editAttachment();
        } else {
            return uploadAttachment();
        }
    };

    const uploadAttachment = () => {
        return Attachment.New(file, attachment).then(attachment => {
            props.didUpdate(attachment);
            GlobalModalFrame.removeModal();
        });
    };

    const editAttachment = () => {
        return Attachment.Save(file, attachment).then(attachment => {
            props.didUpdate(attachment);
            GlobalModalFrame.removeModal();
        });
    };

    const changePath = (Path: string) => {
        setAttachment(attachment => {
            attachment.Path = Path;
            return { ...attachment };
        });
    };

    const changeFile = (file: File) => {
        setFile(file);
    };

    const changeOwner = (Owner: RunAs) => {
        setAttachment(attachment => {
            attachment.Owner = Owner;
            return { ...attachment };
        });
    };

    const changeMode = (Mode: number) => {
        setAttachment(attachment => {
            attachment.Mode = Mode;
            return { ...attachment };
        });
    };

    const changeAfterScript = (AfterScript: boolean) => {
        setAttachment(attachment => {
            attachment.AfterScript = AfterScript;
            return { ...attachment };
        });
    };

    const fileInput = () => {
        const labelText = props.attachment ? 'Replace File' : 'Select File';
        const helpText = props.attachment ? 'Select a new file to replace the existing file, otherwise the file is not changed.' : '';
        return (<Input.FileChooser label={labelText} onChange={changeFile} helpText={helpText} />);
    };

    const title = props.attachment ? 'Edit Attachment' : 'New Attachment';
    return (
        <ModalForm title={title} onSubmit={saveAttachment}>
            { fileInput()}
            <Input.Text type="text" label="File Path" defaultValue={attachment.Path} required onChange={changePath} helpText="The absolute path where the file will be located on hosts. If the parent directory does not exist it will be created with the same owner as the attachment." fixedWidth />
            <Input.RunAsInput inheritLabel="Specify File Owner" label="Owned By" defaultValue={attachment.Owner} onChange={changeOwner} />
            <Input.Number label="Permission / Mode" defaultValue={attachment.Mode} required onChange={changeMode} helpText="The permission value (mode) for the file" />
            <Input.Checkbox label="Upload After Script Execution" defaultValue={attachment.AfterScript} onChange={changeAfterScript} helpText="If checked the file will be uploaded once the script has completed successfully. If unchecked the file will be uploaded before the script is executed." />
        </ModalForm>
    );
};
