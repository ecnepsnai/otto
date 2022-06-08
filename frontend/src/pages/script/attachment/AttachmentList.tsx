import * as React from 'react';
import { Buttons, AddButton } from '../../../components/Button';
import { ContextMenuItem } from '../../../components/ContextMenu';
import { Icon } from '../../../components/Icon';
import { Loading } from '../../../components/Loading';
import { GlobalModalFrame } from '../../../components/Modal';
import { Column, Table } from '../../../components/Table';
import { Formatter } from '../../../services/Formatter';
import { Attachment, AttachmentType } from '../../../types/Attachment';
import { Script } from '../../../types/Script';
import { AttachmentEdit } from './AttachmentEdit';

interface AttachmentListProps {
    scriptID?: string;
    didUpdateAttachments: (fileIDs: string[]) => void;
}
export const AttachmentList: React.FC<AttachmentListProps> = (props: AttachmentListProps) => {
    const [loading, setLoading] = React.useState(true);
    const [attachments, setAttachments] = React.useState<AttachmentType[]>();

    React.useEffect(() => {
        loadData();
    }, []);

    React.useEffect(() => {
        if (attachments == undefined) {
            return;
        }

        props.didUpdateAttachments(attachments.map(attachment => attachment.ID));
    }, [attachments]);

    const loadData = () => {
        if (props.scriptID) {
            Script.Attachments(props.scriptID).then(attachments => {
                setAttachments(attachments);
                setLoading(false);
            });
        } else {
            setAttachments([]);
            setLoading(false);
        }
    };

    const didAddAttachment = (attachment: Attachment) => {
        setAttachments(attachments => {
            return [...attachments, attachment];
        });
    };

    const didEditAttachment = (attachment: AttachmentType) => {
        return () => {
            setAttachments(attachments => {
                for (let i = 0; i < attachments.length; i++) {
                    if (attachments[i].ID == attachment.ID) {
                        attachments[i] = attachment;
                    }
                }
                return [...attachments];
            });
        };
    };

    const didDeleteAttachment = (attachment: AttachmentType) => {
        return () => {
            setAttachments(attachments => {
                let idx = -1;
                for (let i = 0; i < attachments.length; i++) {
                    if (attachments[i].ID == attachment.ID) {
                        idx = i;
                        break;
                    }
                }
                attachments.splice(idx, 1);
                return [...attachments];
            });
        };
    };

    const createButtonClick = () => {
        GlobalModalFrame.showModal(<AttachmentEdit didUpdate={didAddAttachment} />);
    };

    if (loading) {
        return (<Loading />);
    }

    const tableCols: Column[] = [
        {
            title: 'Path',
            value: 'Path',
            sort: 'Path'
        },
        {
            title: 'Type',
            value: 'MimeType',
            sort: 'MimeType'
        },
        {
            title: 'Owner',
            value: (v: AttachmentType) => {
                if (v.Owner.Inherit) {
                    return (<em>Inherit from Script</em>);
                }
                return (<span>{v.Owner.UID + ':' + v.Owner.GID}</span>);
            },
        },
        {
            title: 'Permissions',
            value: 'Mode',
            sort: 'Mode'
        },
        {
            title: 'Size',
            value: (v: AttachmentType) => {
                return (<span>{Formatter.Bytes(v.Size)}</span>);
            },
            sort: 'Size'
        }
    ];

    return (<div>
        <Buttons>
            <AddButton onClick={createButtonClick} />
        </Buttons>
        <Table columns={tableCols} data={attachments} contextMenu={(a: AttachmentType) => AttachmentTableContextMenu(a, didEditAttachment(a), didDeleteAttachment(a))} defaultSort={{ ColumnIdx: 0, Ascending: true }} />
    </div>);
};

const AttachmentTableContextMenu = (attachment: AttachmentType, didUpate: () => void, didDelete: () => void): (ContextMenuItem | 'separator')[] => {
    const editClick = () => {
        GlobalModalFrame.showModal(<AttachmentEdit attachment={attachment} didUpdate={didUpate} />);
    };

    const deleteClick = () => {
        Attachment.DeleteModal(attachment).then(deleted => {
            if (deleted) {
                didDelete();
            }
        });
    };

    return [
        {
            title: 'Edit',
            icon: (<Icon.Edit />),
            onClick: editClick
        },
        {
            title: 'Download',
            icon: (<Icon.Download />),
            href: '/api/attachments/attachment/' + attachment.ID + '/download'
        },
        'separator',
        {
            title: 'Delete',
            icon: (<Icon.Delete />),
            onClick: deleteClick,
        },
    ];
};
