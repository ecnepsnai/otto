import * as React from 'react';
import { Buttons, AddButton } from '../../../components/Button';
import { Loading } from '../../../components/Loading';
import { GlobalModalFrame } from '../../../components/Modal';
import { Table } from '../../../components/Table';
import { Attachment, AttachmentType } from '../../../types/Attachment';
import { Script } from '../../../types/Script';
import { AttachmentEdit } from './AttachmentEdit';
import { AttachmentListItem } from './AttachmentListItem';

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

    const didEditAttachment = (idx: number) => {
        return (attachment: Attachment) => {
            setAttachments(attachments => {
                attachments[idx] = attachment;
                return [...attachments];
            });
        };
    };

    const didDeleteAttachment = (idx: number) => {
        return () => {
            setAttachments(attachments => {
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

    return (<div>
        <Buttons>
            <AddButton onClick={createButtonClick} />
        </Buttons>
        <Table.Table>
            <Table.Head>
                <Table.Column>Path</Table.Column>
                <Table.Column>Type</Table.Column>
                <Table.Column>Owner</Table.Column>
                <Table.Column>Permission</Table.Column>
                <Table.Column>Size</Table.Column>
                <Table.MenuColumn />
            </Table.Head>
            <Table.Body>
                {
                    attachments.map((attachment, idx) => {
                        return (<AttachmentListItem attachment={attachment} key={idx} didEdit={didEditAttachment(idx)} didDelete={didDeleteAttachment(idx)}/>);
                    })
                }
            </Table.Body>
        </Table.Table>
    </div>);
};
