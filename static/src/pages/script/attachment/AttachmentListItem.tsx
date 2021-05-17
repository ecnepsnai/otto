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
export const AttachmentListItem: React.FC<AttachmentListItemProps> = (props: AttachmentListItemProps) => {
    const didEditAttachment = (attachment: AttachmentType) => {
        props.didEdit(attachment);
    };

    const editClick = () => {
        GlobalModalFrame.showModal(<AttachmentEdit attachment={props.attachment} didUpdate={didEditAttachment} />);
    };

    const deleteClick = () => {
        Attachment.DeleteModal(props.attachment).then(deleted => {
            if (deleted) {
                props.didDelete();
            }
        });
    };

    return (
        <Table.Row>
            <td>{props.attachment.Path}</td>
            <td>{props.attachment.MimeType}</td>
            <td>{props.attachment.UID + ':' + props.attachment.GID}</td>
            <td>{props.attachment.Mode}</td>
            <td>{Formatter.Bytes(props.attachment.Size)}</td>
            <Table.Menu>
                <Menu.Item label="Edit" icon={<Icon.Edit />} onClick={editClick} />
                <Menu.Anchor label="Download" icon={<Icon.Download />} href={'/api/attachments/attachment/' + props.attachment.ID + '/download'} />
                <Menu.Divider />
                <Menu.Item label="Delete" icon={<Icon.Delete />} onClick={deleteClick} />
            </Table.Menu>
        </Table.Row>
    );
};
