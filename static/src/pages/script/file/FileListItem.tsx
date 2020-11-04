import * as React from 'react';
import { File } from '../../../types/File';

export interface FileListItemProps {
    file: File;
}
interface FileListItemState {}
export class FileListItem extends React.Component<FileListItemProps, FileListItemState> {
    constructor(props: FileListItemProps) {
        super(props);
        this.state = { };
    }

    render(): JSX.Element {
        return (
            <div></div>
        );
    }
}
