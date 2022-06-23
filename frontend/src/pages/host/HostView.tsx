import * as React from 'react';
import { Heartbeat, HeartbeatType } from '../../types/Heartbeat';
import { Host, HostType, ScriptEnabledGroup } from '../../types/Host';
import { GroupType } from '../../types/Group';
import { PageLoading } from '../../components/Loading';
import { Page } from '../../components/Page';
import { EditButton, DeleteButton } from '../../components/Button';
import { Layout } from '../../components/Layout';
import { useParams, useNavigate } from 'react-router-dom';
import { URLParams } from '../../services/Params';
import { Card } from '../../components/Card';
import { ListGroup } from '../../components/ListGroup';
import { EnabledBadge } from '../../components/Badge';
import { EnvironmentVariableCard } from '../../components/EnvironmentVariableCard';
import { ScheduleType } from '../../types/Schedule';
import { GroupListCard } from '../../components/GroupListCard';
import { ScriptListCard } from '../../components/ScriptListCard';
import { ScheduleListCard } from '../../components/ScheduleListCard';
import { HostHeartbeat } from './HostHeartbeat';
import { HostTrust } from './HostTrust';

export const HostView: React.FC = () => {
    const { id } = useParams() as URLParams;
    const [loading, setLoading] = React.useState(true);
    const [host, setHost] = React.useState<HostType>();
    const [heartbeat, setHeartbeat] = React.useState<HeartbeatType>();
    const [groups, setGroups] = React.useState<GroupType[]>();
    const [scripts, setScripts] = React.useState<ScriptEnabledGroup[]>();
    const [schedules, setSchedules] = React.useState<ScheduleType[]>();
    const navigate = useNavigate();

    React.useEffect(() => {
        loadData();
    }, []);

    const loadHost = async () => {
        const host = await Host.Get(id);
        setHost(host);

        const heartbeats = await Heartbeat.List();
        heartbeats.forEach(heartbeat => {
            if (heartbeat.Address === host.Address) {
                setHeartbeat(heartbeat);
            }
        });
    };

    const loadScripts = async () => {
        setScripts(await Host.Scripts(id));
    };

    const loadGroups = async () => {
        const groups = await Host.Groups(id);
        setGroups(groups);
    };

    const loadSchedules = async () => {
        setSchedules(await Host.Schedules(id));
    };

    const loadData = () => {
        Promise.all([loadHost(), loadScripts(), loadGroups(), loadSchedules()]).then(() => {
            setLoading(false);
        });
    };

    const deleteClick = () => {
        Host.DeleteModal(host).then(deleted => {
            if (deleted) {
                navigate('/hosts');
            }
        });
    };

    const didHeartbeat = () => {
        loadHost();
    };

    if (loading) {
        return (<PageLoading />);
    }

    const toolbar = (
        <React.Fragment>
            <EditButton to={'/hosts/host/' + host.ID + '/edit'} />
            <DeleteButton onClick={deleteClick} />
        </React.Fragment>
    );

    const breadcrumbs = [
        {
            title: 'Hosts',
            href: '/hosts'
        },
        {
            title: host.Name
        }
    ];

    return (
        <Page title={breadcrumbs} toolbar={toolbar}>
            <Layout.Row>
                <Layout.Column>
                    <Card.Card className="mb-3">
                        <Card.Header>Host Information</Card.Header>
                        <ListGroup.List>
                            <ListGroup.TextItem title="Name">{host.Name}</ListGroup.TextItem>
                            <ListGroup.TextItem title="Address">{host.Address}:{host.Port}</ListGroup.TextItem>
                            <ListGroup.TextItem title="Status"><EnabledBadge value={host.Enabled} /></ListGroup.TextItem>
                            <ListGroup.TextItem title="Trust"><HostTrust host={host} onReload={loadHost} /></ListGroup.TextItem>
                        </ListGroup.List>
                    </Card.Card>
                    <HostHeartbeat host={host} defaultHeartbeat={heartbeat} didUpdate={didHeartbeat} />
                    <EnvironmentVariableCard className="mb-3" variables={host.Environment} />
                    <ScheduleListCard schedules={schedules} className="mb-3" />
                </Layout.Column>
                <Layout.Column>
                    <GroupListCard groups={groups} className="mb-3" />
                    <ScriptListCard scripts={scripts} hostIDs={[host.ID]} className="mb-3" />
                </Layout.Column>
            </Layout.Row>
        </Page>
    );
};
