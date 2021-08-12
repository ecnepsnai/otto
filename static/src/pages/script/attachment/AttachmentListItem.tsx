import * as React from 'react';
import { ContextMenuItem } from '../../../components/ContextMenu';
import { Icon } from '../../../components/Icon';
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

    const owner = () => {
        if (props.attachment.Owner.Inherit) {
            return (<em>Inherit from Script</em>);
        }
        return (<span>{props.attachment.Owner.UID + ':' + props.attachment.Owner.GID}</span>);
    };

    const contextMenu: (ContextMenuItem | 'separator')[] = [
        {
            title: 'Edit',
            icon: (<Icon.Edit />),
            onClick: editClick
        },
        {
            title: 'Download',
            icon: (<Icon.Download />),
            href: '/api/attachments/attachment/' + props.attachment.ID + '/download'
        },
        'separator',
        {
            title: 'Delete',
            icon: (<Icon.Delete />),
            onClick: deleteClick,
        },
    ];

    return (
        <Table.Row menu={contextMenu}>
            <td>{props.attachment.Path}</td>
            <td>{props.attachment.MimeType}</td>
            <td>{owner()}</td>
            <td>{props.attachment.Mode}</td>
            <td>{Formatter.Bytes(props.attachment.Size)}</td>
        </Table.Row>
    );
};
