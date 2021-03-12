import * as React from 'react';
import { Input } from '../../../components/input/Input';
import { GlobalModalFrame, ModalForm } from '../../../components/Modal';
import { Attachment, AttachmentType } from '../../../types/Attachment';

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
        return Attachment.Save(props.attachment).then(attachment => {
            props.didUpdate(attachment);
            GlobalModalFrame.removeModal();
        });
    };

    const changePath = (Path: string) => {
        setAttachment(attachment => {
            attachment.Path = Path;
            return {...attachment};
        });
    };

    const changeFile = (file: File) => {
        setFile(file);
    };

    const changeOwner = (UID: number, GID: number) => {
        setAttachment(attachment => {
            attachment.UID = UID;
            attachment.GID = GID;
            return {...attachment};
        });
    };

    const changeMode = (Mode: number) => {
        setAttachment(attachment => {
            attachment.Mode = Mode;
            return {...attachment};
        });
    };

    const fileInput = () => {
        if (props.attachment) {
            return null;
        }
        return (<Input.FileChooser label="Upload File" onChange={changeFile}/>);
    };

    const title = props.attachment ? 'Edit Attachment' : 'New Attachment';
    return (
        <ModalForm title={title} onSubmit={saveAttachment}>
            { fileInput() }
            <Input.Text type="text" label="File Path" defaultValue={attachment.Path} required onChange={changePath} helpText="The absolute path where the file will be located on hosts. If the parent directory does not exist it will be created with the same owner as the attachment." fixedWidth/>
            <Input.IDInput label="Owned By" defaultUID={attachment.UID} defaultGID={attachment.GID} onChange={changeOwner} />
            <Input.Number label="Permission / Mode" defaultValue={attachment.Mode} required onChange={changeMode} helpText="The permission value (mode) for the file" />
        </ModalForm>
    );
};
