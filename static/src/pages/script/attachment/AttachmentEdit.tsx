import * as React from 'react';
import { FileBrowser, Input, NumberInput } from '../../../components/Form';
import { GlobalModalFrame, ModalForm } from '../../../components/Modal';
import { Attachment } from '../../../types/Attachment';

export interface AttachmentEditProps {
    attachment?: Attachment;
    didUpdate: (attachment: Attachment) => (void);
}
interface AttachmentEditState {
    attachment: Attachment;
    file?: File;
    loading: boolean;
}
export class AttachmentEdit extends React.Component<AttachmentEditProps, AttachmentEditState> {
    constructor(props: AttachmentEditProps) {
        super(props);
        const attachment = this.props.attachment ?? Attachment.Blank();
        this.state = {
            attachment: attachment,
            loading: false,
        };
    }

    private saveAttachment = () => {
        if (this.props.attachment) {
            return this.editAttachment();
        } else {
            return this.uploadAttachment();
        }
    }

    private uploadAttachment = () => {
        this.setState({ loading: true });
        return Attachment.New(this.state.file, this.state.attachment).then(attachment => {
            this.props.didUpdate(attachment);
            GlobalModalFrame.removeModal();
        }, () => {
            this.setState({ loading: false });
        });
    }

    private editAttachment = () => {
        this.setState({ loading: false });
        return this.props.attachment.Save().then(attachment => {
            this.props.didUpdate(attachment);
            GlobalModalFrame.removeModal();
        }, () => {
            this.setState({ loading: false });
        });
    }

    private changePath = (Path: string) => {
        this.setState(state => {
            state.attachment.Path = Path;
            return state;
        });
    }

    private changeFile = (file: File) => {
        this.setState({ file: file });
    }

    private changeUID = (UID: number) => {
        this.setState(state => {
            state.attachment.UID = UID;
            return state;
        });
    }

    private changeGID = (GID: number) => {
        this.setState(state => {
            state.attachment.GID = GID;
            return state;
        });
    }

    private changeMode = (Mode: number) => {
        this.setState(state => {
            state.attachment.Mode = Mode;
            return state;
        });
    }

    private fileInput = () => {
        if (this.props.attachment) { return null; }
        return (<FileBrowser label="Upload File" onChange={this.changeFile}/>);
    }

    render(): JSX.Element {
        const title = this.props.attachment ? 'Edit Attachment' : 'New Attachment';
        return (
            <ModalForm title={title} onSubmit={this.saveAttachment}>
                { this.fileInput() }
                <Input type="text" label="File Path" defaultValue={this.state.attachment.Path} required onChange={this.changePath} helpText="The absolute path where the file will be located on hosts" fixedWidth/>
                <NumberInput label="Owner User ID" defaultValue={this.state.attachment.UID} required onChange={this.changeUID} helpText="The user ID (UID) that this file will be owned by" />
                <NumberInput label="Owner Group ID" defaultValue={this.state.attachment.GID} required onChange={this.changeGID} helpText="The group ID (GID) that this file will be owned by" />
                <NumberInput label="Permission / Mode" defaultValue={this.state.attachment.Mode} required onChange={this.changeMode} helpText="The permission value (Mode) for the file" />
            </ModalForm>
        );
    }
}
