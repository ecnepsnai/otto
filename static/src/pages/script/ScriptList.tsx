import * as React from 'react';
import { Script } from '../../types/Script';
import { PageLoading } from '../../components/Loading';
import { Page } from '../../components/Page';
import { Buttons, CreateButton } from '../../components/Button';
import { Table } from '../../components/Table';
import { ScriptListItem } from './ScriptListItem';

export interface ScriptListProps {}
interface ScriptListState {
    loading: boolean;
    scripts: Script[];
}
export class ScriptList extends React.Component<ScriptListProps, ScriptListState> {
    constructor(props: ScriptListProps) {
        super(props);
        this.state = {
            loading: true,
            scripts: [],
        };
    }

    componentDidMount(): void {
        this.loadScripts();
    }

    private loadScripts = () => {
        Script.List().then(scripts => {
            this.setState({
                loading: false,
                scripts: scripts,
            });
        });
    }

    render(): JSX.Element {
        if (this.state.loading) { return ( <PageLoading /> ); }

        return (
            <Page title="Scripts">
                <Buttons>
                    <CreateButton to="/scripts/script/" />
                </Buttons>
                <Table.Table>
                    <Table.Head>
                        <Table.Column>Name</Table.Column>
                        <Table.Column>Executable</Table.Column>
                        <Table.Column>Attachments</Table.Column>
                        <Table.MenuColumn />
                    </Table.Head>
                    <Table.Body>
                        {
                            this.state.scripts.map(script => {
                                return <ScriptListItem script={script} key={script.ID} onReload={this.loadScripts}></ScriptListItem>;
                            })
                        }
                    </Table.Body>
                </Table.Table>
            </Page>
        );
    }
}
