import * as React from 'react';
import { Group, GroupType } from '../../types/Group';
import { useParams, Link } from 'react-router-dom';
import { URLParams } from '../../services/Params';
import { PageLoading } from '../../components/Loading';
import { HostType } from '../../types/Host';
import { ScriptType } from '../../types/Script';
import { Page } from '../../components/Page';
import { Layout } from '../../components/Layout';
import { EditButton, DeleteButton } from '../../components/Button';
import { Card } from '../../components/Card';
import { ListGroup } from '../../components/ListGroup';
import { Icon } from '../../components/Icon';
import { EnvironmentVariableCard } from '../../components/EnvironmentVariableCard';
import { Nothing } from '../../components/Nothing';
import { Redirect } from '../../components/Redirect';
import { ScriptListCard } from '../../components/ScriptListCard';
import { ScheduleType } from '../../types/Schedule';
import { ScheduleListCard } from '../../components/ScheduleListCard';

export const GroupView: React.FC = () => {
    const { id } = useParams() as URLParams;
    const [loading, setLoading] = React.useState(true);
    const [group, setGroup] = React.useState<GroupType>();
    const [hosts, setHosts] = React.useState<HostType[]>();
    const [scripts, setScripts] = React.useState<ScriptType[]>();
    const [schedules, setSchedules] = React.useState<ScheduleType[]>();

    React.useEffect(() => {
        loadData();
    }, []);

    const loadGroup = async () => {
        const group = await Group.Get(id);
        setGroup(group);
    };

    const loadHosts = async () => {
        const hosts = await Group.Hosts(id);
        setHosts(hosts);
    };

    const loadScripts = async () => {
        const scripts = await Group.Scripts(id);
        setScripts(scripts);
    };

    const loadSchedules = async () => {
        const schedules = await Group.Schedules(id);
        setSchedules(schedules);
    };

    const loadData = () => {
        Promise.all([loadGroup(), loadHosts(), loadScripts(), loadSchedules()]).then(() => {
            setLoading(false);
        });
    };

    const deleteClick = () => {
        Group.DeleteModal(group).then(deleted => {
            if (!deleted) {
                return;
            }

            Redirect.To('/groups');
        });
    };

    if (loading) {
        return (<PageLoading />);
    }

    const toolbar = (
        <React.Fragment>
            <EditButton to={'/groups/group/' + group.ID + '/edit'} />
            <DeleteButton onClick={deleteClick} />
        </React.Fragment>
    );

    const breadcrumbs = [
        {
            title: 'Groups',
            href: '/groups'
        },
        {
            title: group.Name
        }
    ];

    return (
        <Page title={breadcrumbs} toolbar={toolbar}>
            <Layout.Row>
                <Layout.Column>
                    <Card.Card className="mb-3">
                        <Card.Header>Host Information</Card.Header>
                        <ListGroup.List>
                            <ListGroup.TextItem title="Name">{group.Name}</ListGroup.TextItem>
                        </ListGroup.List>
                    </Card.Card>
                    <EnvironmentVariableCard variables={group.Environment} className="mb-3" />
                    <ScheduleListCard schedules={schedules} className="mb-3" />
                </Layout.Column>
                <Layout.Column>
                    <Card.Card className="mb-3">
                        <Card.Header>Hosts</Card.Header>
                        <HostListCard hosts={hosts} />
                    </Card.Card>
                    <ScriptListCard scripts={scripts} hostIDs={hosts.map(h => h.ID)} className="mb-3" />
                </Layout.Column>
            </Layout.Row>
        </Page>
    );
};

interface HostListCardProps {
    hosts: HostType[];
}
export const HostListCard: React.FC<HostListCardProps> = (props: HostListCardProps) => {
    if (!props.hosts || props.hosts.length < 1) {
        return (
            <Card.Body>
                <Nothing />
            </Card.Body>
        );
    }
    return (
        <ListGroup.List>
            {
                props.hosts.map((host, index) => {
                    return (
                        <ListGroup.Item key={index}>
                            <Icon.Desktop />
                            <Link to={'/hosts/host/' + host.ID} className="ms-1">{host.Name}</Link>
                        </ListGroup.Item>
                    );
                })
            }
        </ListGroup.List>
    );
};
