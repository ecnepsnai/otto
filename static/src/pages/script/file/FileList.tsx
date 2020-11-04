import * as React from 'react';
import { Buttons, CreateButton } from '../../../components/Button';
import { Loading } from '../../../components/Loading';
import { Table } from '../../../components/Table';
import { File } from '../../../types/File';
import { Script } from '../../../types/Script';
import { FileListItem } from './FileListItem';

export interface FileListProps {
    scriptID?: string;
    didUpdateFiles: (fileIDs: string[]) => void;
}
interface FileListState {
    loading?: boolean;
    files?: File[];
}
export class FileList extends React.Component<FileListProps, FileListState> {
    constructor(props: FileListProps) {
        super(props);
        this.state = { loading: true };
    }

    private loadData = () => {
        if (this.props.scriptID) {
            Script.Files(this.props.scriptID).then(files => {
                this.setState({ loading: false, files: files });
            });
        } else {
            this.setState({ loading: false, files: [] });
        }
    }

    componentDidMount(): void {
        this.loadData();
    }

    render(): JSX.Element {
        if (this.state.loading) { return (<Loading />); }

        return (<div>
            <Buttons>
                <CreateButton to="/scripts/script/" />
            </Buttons>
            <Table.Table>
                <Table.Head>
                    <Table.Column>Path</Table.Column>
                    <Table.Column>Owner</Table.Column>
                    <Table.MenuColumn />
                </Table.Head>
                <Table.Body>
                    {
                        this.state.files.map((file, idx) => {
                            return (<FileListItem file={file} key={idx} />);
                        })
                    }
                </Table.Body>
            </Table.Table>
        </div>);
    }
}
