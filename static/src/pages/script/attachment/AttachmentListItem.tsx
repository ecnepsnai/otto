import * as React from 'react';
import { Icon } from '../../../components/Icon';
import { MenuItem } from '../../../components/Menu';
import { GlobalModalFrame } from '../../../components/Modal';
import { Table } from '../../../components/Table';
import { Attachment } from '../../../types/Attachment';
import { AttachmentEdit } from './AttachmentEdit';

export interface AttachmentListItemProps {
    attachment: Attachment;
    didEdit: (attachment: Attachment) => (void);
    didDelete: () => (void);
}
interface AttachmentListItemState {}
export class AttachmentListItem extends React.Component<AttachmentListItemProps, AttachmentListItemState> {
    constructor(props: AttachmentListItemProps) {
        super(props);
        this.state = {};
    }

    private didEditAttachment = (attachment: Attachment) => {
        this.props.didEdit(attachment);
    }

    private editClick = () => {
        GlobalModalFrame.showModal(<AttachmentEdit attachment={this.props.attachment} didUpdate={this.didEditAttachment}/>);
    }

    private deleteClick = () => {
        this.props.attachment.DeleteModal().then(deleted => {
            if (deleted) {
                this.props.didDelete();
            }
        });
    }

    render(): JSX.Element {
        return (
            <Table.Row>
                <td>{ this.props.attachment.Path }</td>
                <td>{ this.props.attachment.UID + ':' + this.props.attachment.GID }</td>
                <td>{ this.props.attachment.Mode }</td>
                <Table.Menu>
                    <MenuItem label="Edit" icon={<Icon.Edit />} onClick={this.editClick}/>
                    <MenuItem label="Delete" icon={<Icon.Delete />} onClick={this.deleteClick}/>
                </Table.Menu>
            </Table.Row>
        );
    }
}
