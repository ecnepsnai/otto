import * as React from 'react';
import { Runbook as RunbookAPI, RunbookType } from '../../types/Runbook';
import { useParams, Link, useNavigate } from 'react-router-dom';
import { URLParams } from '../../services/Params';
import { PageLoading } from '../../components/Loading';
import { HostType } from '../../types/Host';
import { ScriptType } from '../../types/Script';
import { Page } from '../../components/Page';
import { Layout } from '../../components/Layout';
import { EditButton, DeleteButton } from '../../components/Button';
import { Card } from '../../components/Card';
import { GroupType } from '../../types/Group';
import { ListGroup } from '../../components/ListGroup';
import { DateLabel } from '../../components/DateLabel';
import { Icon } from '../../components/Icon';
import { Permissions, UserAction } from '../../services/Permissions';

export const RunbookView: React.FC = () => {
    const { id } = useParams() as URLParams;
    const [IsLoading, SetIsLoading] = React.useState<boolean>(true);
    const [Runbook, SetRunbook] = React.useState<RunbookType>();
    const [Scripts, SetScripts] = React.useState<ScriptType[]>();
    const [Hosts, SetHosts] = React.useState<HostType[]>();
    const [Groups, SetGroups] = React.useState<GroupType[]>();
    const navigate = useNavigate();

    React.useEffect(() => {
        loadData();
    }, []);

    const loadRunbook = () => {
        return RunbookAPI.Get(id);
    };

    const loadScripts = () => {
        return RunbookAPI.Scripts(id);
    };

    const loadGroups = () => {
        return RunbookAPI.Groups(id);
    };

    const loadHosts = () => {
        return RunbookAPI.Hosts(id);
    };

    const loadData = () => {
        Promise.all([loadRunbook(), loadScripts(), loadGroups(), loadHosts()]).then(results => {
            SetRunbook(results[0]);
            SetScripts(results[1]);
            SetGroups(results[2]);
            SetHosts(results[3]);
            SetIsLoading(false);
        });
    };

    const deleteClick = () => {
        RunbookAPI.DeleteModal(Runbook).then(deleted => {
            if (!deleted) {
                return;
            }

            navigate('/runbooks');
        });
    };

    const groupsList = () => {
        if (!Groups || Groups.length === 0) {
            return null;
        }

        return (
            <ListGroup.TextItem title="Groups">
                {Groups.map((group, idx) => {
                    return (
                        <div key={idx}>
                            <Icon.LayerGroup />
                            <Link className="ms-1" to={'/groups/group/' + group.ID}>{group.Name}</Link>
                        </div>
                    );
                })}
            </ListGroup.TextItem>
        );
    };

    const hostsList = () => {
        return (
            <Card.Card className='mt-2'>
                <Card.Header><Icon.Label icon={<Icon.Desktop />} label="Hosts" /></Card.Header>
                <ListGroup.List>
                    {Hosts.map((host, idx) => {
                        return (
                            <ListGroup.Item key={idx}>
                                <Link to={'/hosts/host/' + host.ID}>{host.Name}</Link>
                            </ListGroup.Item>
                        );
                    })}
                </ListGroup.List>
            </Card.Card>
        );
    };

    const scriptsList = () => {
        return (
            <Card.Card>
                <Card.Header><Icon.Label icon={<Icon.Scroll />} label="Scripts" /></Card.Header>
                <ListGroup.List>
                    {Scripts.map((script, idx) => {
                        return (
                            <ListGroup.Item key={idx}>
                                <Link to={'/scripts/script/' + script.ID}>{script.Name}</Link>
                            </ListGroup.Item>
                        );
                    })}
                </ListGroup.List>
            </Card.Card>
        );
    };

    if (IsLoading) {
        return (<PageLoading />);
    }

    const toolbar = (
        <React.Fragment>
            <EditButton to={'/runbooks/runbook/' + Runbook.ID + '/edit'} disabled={!Permissions.UserCan(UserAction.ModifyRunbooks)} />
            <DeleteButton onClick={deleteClick} disabled={!Permissions.UserCan(UserAction.ModifyRunbooks)} />
        </React.Fragment>
    );

    const breadcrumbs = [
        {
            title: 'Runbooks',
            href: '/runbooks'
        },
        {
            title: Runbook.Name
        }
    ];

    return (
        <Page title={breadcrumbs} toolbar={toolbar}>
            <Layout.Row>
                <Layout.Column>
                    <Card.Card>
                        <Card.Header>Runbook Information</Card.Header>
                        <ListGroup.List>
                            <ListGroup.TextItem title="Name">{Runbook.Name}</ListGroup.TextItem>
                            <ListGroup.TextItem title="Last Run"><DateLabel date={Runbook.LastRun} /></ListGroup.TextItem>
                            {groupsList()}
                        </ListGroup.List>
                    </Card.Card>
                </Layout.Column>
                <Layout.Column>
                    {scriptsList()}
                    {hostsList()}
                </Layout.Column>
            </Layout.Row>
        </Page>
    );
};
