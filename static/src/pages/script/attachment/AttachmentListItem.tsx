import * as React from 'react';
import { Icon } from '../../../components/Icon';
import { Menu } from '../../../components/Menu';
import { GlobalModalFrame } from '../../../components/Modal';
import { Table } from '../../../components/Table';
import { Formatter } from '../../../services/Formatter';
import { Attachment, AttachmentType } from '../../../types/Attachment';
import { AttachmentEdit } from './AttachmentEdit';

interface AttachmentListItemProps {
    attachment: AttachmentType;
    didEdit: (attachment: AttachmentType) => (void);
    didDelete: () => (void);
}
export class AttachmentListItem extends React.Component<AttachmentListItemProps, unknown> {
    private didEditAttachment = (attachment: AttachmentType) => {
        this.props.didEdit(attachment);
    }

    private editClick = () => {
        GlobalModalFrame.showModal(<AttachmentEdit attachment={this.props.attachment} didUpdate={this.didEditAttachment}/>);
    }

    private deleteClick = () => {
        Attachment.DeleteModal(this.props.attachment).then(deleted => {
            if (deleted) {
                this.props.didDelete();
            }
        });
    }

    render(): JSX.Element {
        return (
            <Table.Row>
                <td>{ this.props.attachment.Path }</td>
                <td>{ this.props.attachment.MimeType }</td>
                <td>{ this.props.attachment.UID + ':' + this.props.attachment.GID }</td>
                <td>{ this.props.attachment.Mode }</td>
                <td>{ Formatter.Bytes(this.props.attachment.Size) }</td>
                <Table.Menu>
                    <Menu.Item label="Edit" icon={<Icon.Edit />} onClick={this.editClick}/>
                    <Menu.Anchor label="Download" icon={<Icon.Download />} href={'/api/attachments/attachment/' + this.props.attachment.ID + '/download'} />
                    <Menu.Divider />
                    <Menu.Item label="Delete" icon={<Icon.Delete />} onClick={this.deleteClick}/>
                </Table.Menu>
            </Table.Row>
        );
    }
}
