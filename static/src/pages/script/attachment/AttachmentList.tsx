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
interface AttachmentListState {
    loading?: boolean;
    attachments?: AttachmentType[];
}
export class AttachmentList extends React.Component<AttachmentListProps, AttachmentListState> {
    constructor(props: AttachmentListProps) {
        super(props);
        this.state = { loading: true };
    }

    private loadData = () => {
        if (this.props.scriptID) {
            Script.Attachments(this.props.scriptID).then(attachments => {
                this.setState({ loading: false, attachments: attachments });
            });
        } else {
            this.setState({ loading: false, attachments: [] });
        }
    }

    componentDidMount(): void {
        this.loadData();
    }

    private didUpdateAttachments = () => {
        const ids = this.state.attachments.map(attachment => {
            return attachment.ID;
        });
        this.props.didUpdateAttachments(ids);
    }

    private didAddAttachment = (attachment: Attachment) => {
        this.setState(state => {
            state.attachments.push(attachment);
            return state;
        }, () => {
            this.didUpdateAttachments();
        });
    }

    private didEditAttachment = (idx: number) => {
        return (attachment: Attachment) => {
            this.setState(state => {
                state.attachments[idx] = attachment;
                return state;
            });
        };
    }

    private didDeleteAttachment = (idx: number) => {
        return () => {
            this.setState(state => {
                state.attachments.splice(idx, 1);
                return state;
            }, () => {
                this.didUpdateAttachments();
            });
        };
    }

    private createButtonClick = () => {
        GlobalModalFrame.showModal(<AttachmentEdit didUpdate={this.didAddAttachment} />);
    }

    render(): JSX.Element {
        if (this.state.loading) {
            return (<Loading />);
        }

        return (<div>
            <Buttons>
                <AddButton onClick={this.createButtonClick} />
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
                        this.state.attachments.map((attachment, idx) => {
                            return (<AttachmentListItem attachment={attachment} key={idx} didEdit={this.didEditAttachment(idx)} didDelete={this.didDeleteAttachment(idx)}/>);
                        })
                    }
                </Table.Body>
            </Table.Table>
        </div>);
    }
}
