import * as React from 'react';

export interface FileEditProps {}
interface FileEditState {}
export class FileEdit extends React.Component<FileEditProps, FileEditState> {
    constructor(props: FileEditProps) {
        super(props);
        this.state = { };
    }

    render(): JSX.Element {
        return (
            <div></div>
        );
    }
}
