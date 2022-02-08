import * as React from 'react';
import { Script, ScriptType } from '../../types/Script';
import { PageLoading } from '../../components/Loading';
import { Page } from '../../components/Page';
import { CreateButton } from '../../components/Button';
import { Table } from '../../components/Table';
import { ScriptListItem } from './ScriptListItem';

export const ScriptList: React.FC = () => {
    const [loading, setLoading] = React.useState<boolean>(true);
    const [scripts, setScripts] = React.useState<ScriptType[]>([]);

    React.useEffect(() => {
        loadScripts();
    }, []);

    const loadScripts = () => {
        Script.List().then(scripts => {
            setLoading(false);
            setScripts(scripts);
        });
    };

    if (loading) {
        return (<PageLoading />);
    }

    const toolbar = (
        <React.Fragment>
            <CreateButton to="/scripts/script/" />
        </React.Fragment>
    );

    return (
        <Page title="Scripts" toolbar={toolbar}>
            <Table.Table>
                <Table.Head>
                    <Table.Column>Name</Table.Column>
                    <Table.Column>Executable</Table.Column>
                    <Table.Column>Attachments</Table.Column>
                </Table.Head>
                <Table.Body>
                    {
                        scripts.map(script => {
                            return <ScriptListItem script={script} key={script.ID} onReload={loadScripts}></ScriptListItem>;
                        })
                    }
                </Table.Body>
            </Table.Table>
        </Page>
    );
};
