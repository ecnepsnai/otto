import * as React from 'react';
import { Schedule, ScheduleReport, ScheduleType } from '../../types/Schedule';
import { match, Link } from 'react-router-dom';
import { URLParams } from '../../services/Params';
import { PageLoading } from '../../components/Loading';
import { HostType } from '../../types/Host';
import { ScriptType } from '../../types/Script';
import { Page } from '../../components/Page';
import { Layout } from '../../components/Layout';
import { Buttons, EditButton, DeleteButton } from '../../components/Button';
import { Card } from '../../components/Card';
import { Redirect } from '../../components/Redirect';
import { GroupType } from '../../types/Group';
import { ListGroup } from '../../components/ListGroup';
import { EnabledBadge } from '../../components/Badge';
import { DateLabel } from '../../components/DateLabel';
import { SchedulePattern } from './SchedulePattern';
import { Icon } from '../../components/Icon';
import { Style } from '../../components/Style';
import { Nothing } from '../../components/Nothing';

interface ScheduleViewProps {
    match: match;
}
export const ScheduleView: React.FC<ScheduleViewProps> = (props: ScheduleViewProps) => {
    const [loading, setLoading] = React.useState<boolean>(true);
    const [schedule, setSchedule] = React.useState<ScheduleType>();
    const [reports, setReports] = React.useState<ScheduleReport[]>();
    const [script, setScript] = React.useState<ScriptType>();
    const [hosts, setHosts] = React.useState<HostType[]>();
    const [groups, setGroups] = React.useState<GroupType[]>();

    React.useEffect(() => {
        loadData();
    }, []);

    const loadSchedule = () => {
        const scheduleID = (props.match.params as URLParams).id;
        return Schedule.Get(scheduleID);
    };

    const loadReports = () => {
        const scheduleID = (props.match.params as URLParams).id;
        return Schedule.Reports(scheduleID);
    };

    const loadScript = () => {
        const scheduleID = (props.match.params as URLParams).id;
        return Schedule.Script(scheduleID);
    };

    const loadGroups = () => {
        const scheduleID = (props.match.params as URLParams).id;
        return Schedule.Groups(scheduleID);
    };

    const loadHosts = () => {
        const scheduleID = (props.match.params as URLParams).id;
        return Schedule.Hosts(scheduleID);
    };

    const loadData = () => {
        Promise.all([loadSchedule(), loadHosts(), loadReports(), loadScript()]).then(results => {
            const schedule = results[0];
            setSchedule(schedule);
            setHosts(results[1]);
            setReports(results[2]);
            setScript(results[3]);

            let scopePromise: Promise<GroupType[]> = Promise.resolve([]);
            if (schedule.Scope.GroupIDs && schedule.Scope.GroupIDs.length > 0) {
                scopePromise = loadGroups();
            }
            scopePromise.then(groups => {
                setGroups(groups);
                setLoading(false);
            });
        });
    };

    const deleteClick = () => {
        Schedule.DeleteModal(schedule).then(deleted => {
            if (!deleted) {
                return;
            }

            Redirect.To('/schedules');
        });
    };

    const groupsList = () => {
        if (!groups) {
            return null;
        }

        return (
            <ListGroup.TextItem title="Groups">
                { groups.map((group, idx) => {
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
            <ListGroup.TextItem title="Hosts">
                { hosts.map((host, idx) => {
                    return (
                        <div key={idx}>
                            <Icon.Desktop />
                            <Link className="ms-1" to={'/hosts/host/' + host.ID}>{host.Name}</Link>
                        </div>
                    );
                })}
            </ListGroup.TextItem>
        );
    };

    const historyContent = () => {
        if (reports.length == 0) {
            return (<Card.Body><Nothing /></Card.Body>);
        }

        return (<ListGroup.List>
            { reports.map((report, idx) => {
                return (<ScheduleReportItem report={report} key={idx} />);
            })}
        </ListGroup.List>);
    };

    if (loading) {
        return (<PageLoading />);
    }

    return (
        <Page title="View Schedule">
            <Layout.Container>
                <Buttons>
                    <EditButton to={'/schedules/schedule/' + schedule.ID + '/edit'} />
                    <DeleteButton onClick={deleteClick} />
                </Buttons>
                <Layout.Row>
                    <Layout.Column>
                        <Card.Card>
                            <Card.Header>Schedule Information</Card.Header>
                            <ListGroup.List>
                                <ListGroup.TextItem title="Name">{schedule.Name}</ListGroup.TextItem>
                                <ListGroup.TextItem title="Script"><Link to={'/scripts/script/' + schedule.ScriptID}>{script.Name}</Link></ListGroup.TextItem>
                                <ListGroup.TextItem title="Frequency"><SchedulePattern pattern={schedule.Pattern} /></ListGroup.TextItem>
                                <ListGroup.TextItem title="Last Run"><DateLabel date={schedule.LastRunTime} /></ListGroup.TextItem>
                                <ListGroup.TextItem title="Enabled"><EnabledBadge value={schedule.Enabled} /></ListGroup.TextItem>
                                {groupsList()}
                                {hostsList()}
                            </ListGroup.List>
                        </Card.Card>
                    </Layout.Column>
                    <Layout.Column>
                        <Card.Card>
                            <Card.Header>History</Card.Header>
                            {historyContent()}
                        </Card.Card>
                    </Layout.Column>
                </Layout.Row>
            </Layout.Container>
        </Page>
    );
};

interface ScheduleReportItemProps { report: ScheduleReport; }
const ScheduleReportItem: React.FC<ScheduleReportItemProps> = (props: ScheduleReportItemProps) => {
    let icon = (<Icon.QuestionCircle color={Style.Palette.Primary} />);
    if (props.report.Result == 0) {
        icon = (<Icon.CheckCircle color={Style.Palette.Success} />);
    } else if (props.report.Result == 1) {
        icon = (<Icon.ExclamationTriangle color={Style.Palette.Warning} />);
    } else if (props.report.Result == 2) {
        icon = (<Icon.ExclamationCircle color={Style.Palette.Danger} />);
    }

    return (
        <ListGroup.Item>
            {icon}
            <span className="ms-1">
                <DateLabel date={props.report.Time.Start} /> on {props.report.HostIDs.length} hosts
            </span>
        </ListGroup.Item>
    );
};
