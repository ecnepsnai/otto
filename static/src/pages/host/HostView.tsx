import * as React from 'react';
import { Heartbeat, HeartbeatType } from '../../types/Heartbeat';
import { Host, HostType, ScriptEnabledGroup } from '../../types/Host';
import { GroupType } from '../../types/Group';
import { PageLoading } from '../../components/Loading';
import { Page } from '../../components/Page';
import { EditButton, DeleteButton } from '../../components/Button';
import { Layout } from '../../components/Layout';
import { match } from 'react-router-dom';
import { URLParams } from '../../services/Params';
import { Card } from '../../components/Card';
import { ListGroup } from '../../components/ListGroup';
import { EnabledBadge } from '../../components/Badge';
import { EnvironmentVariableCard } from '../../components/EnvironmentVariableCard';
import { Redirect } from '../../components/Redirect';
import { ScheduleType } from '../../types/Schedule';
import { GroupListCard } from '../../components/GroupListCard';
import { ScriptListCard } from '../../components/ScriptListCard';
import { ScheduleListCard } from '../../components/ScheduleListCard';
import { HostHeartbeat } from './HostHeartbeat';

interface HostViewProps {
    match: match;
}
export const HostView: React.FC<HostViewProps> = (props: HostViewProps) => {
    const [loading, setLoading] = React.useState(true);
    const [host, setHost] = React.useState<HostType>();
    const [heartbeat, setHeartbeat] = React.useState<HeartbeatType>();
    const [groups, setGroups] = React.useState<GroupType[]>();
    const [scripts, setScripts] = React.useState<ScriptEnabledGroup[]>();
    const [schedules, setSchedules] = React.useState<ScheduleType[]>();

    React.useEffect(() => {
        loadData();
    }, []);

    const loadHost = async () => {
        const hostID = (props.match.params as URLParams).id;
        const host = await Host.Get(hostID);
        setHost(host);

        const heartbeats = await Heartbeat.List();
        heartbeats.forEach(heartbeat => {
            if (heartbeat.Address === host.Address) {
                setHeartbeat(heartbeat);
            }
        });
    };

    const loadScripts = async () => {
        const hostID = (props.match.params as URLParams).id;
        setScripts(await Host.Scripts(hostID));
    };

    const loadGroups = async () => {
        const hostID = (props.match.params as URLParams).id;
        const groups = await Host.Groups(hostID);
        setGroups(groups);
    };

    const loadSchedules = async () => {
        const hostID = (props.match.params as URLParams).id;
        setSchedules(await Host.Schedules(hostID));
    };

    const loadData = () => {
        Promise.all([loadHost(), loadScripts(), loadGroups(), loadSchedules()]).then(() => {
            setLoading(false);
        });
    };

    const deleteClick = () => {
        Host.DeleteModal(host).then(deleted => {
            if (deleted) {
                Redirect.To('/hosts');
            }
        });
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
                        </ListGroup.List>
                    </Card.Card>
                    <HostHeartbeat host={host} defaultHeartbeat={heartbeat} />
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
