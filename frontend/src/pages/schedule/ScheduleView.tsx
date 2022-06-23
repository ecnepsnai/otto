import * as React from 'react';
import { Schedule, ScheduleReport, ScheduleType } from '../../types/Schedule';
import { useParams, Link, useNavigate } from 'react-router-dom';
import { URLParams } from '../../services/Params';
import { Loading, PageLoading } from '../../components/Loading';
import { Host, HostType } from '../../types/Host';
import { ScriptType } from '../../types/Script';
import { Page } from '../../components/Page';
import { Layout } from '../../components/Layout';
import { EditButton, DeleteButton, Button } from '../../components/Button';
import { Card } from '../../components/Card';
import { GroupType } from '../../types/Group';
import { ListGroup } from '../../components/ListGroup';
import { EnabledBadge } from '../../components/Badge';
import { DateLabel } from '../../components/DateLabel';
import { SchedulePattern } from './SchedulePattern';
import { Icon } from '../../components/Icon';
import { Style } from '../../components/Style';
import { Nothing } from '../../components/Nothing';
import { GlobalModalFrame, Modal } from '../../components/Modal';

export const ScheduleView: React.FC = () => {
    const { id } = useParams() as URLParams;
    const [loading, setLoading] = React.useState<boolean>(true);
    const [schedule, setSchedule] = React.useState<ScheduleType>();
    const [reports, setReports] = React.useState<ScheduleReport[]>();
    const [script, setScript] = React.useState<ScriptType>();
    const [hosts, setHosts] = React.useState<HostType[]>();
    const [groups, setGroups] = React.useState<GroupType[]>();
    const navigate = useNavigate();

    React.useEffect(() => {
        loadData();
    }, []);

    const loadSchedule = () => {
        return Schedule.Get(id);
    };

    const loadReports = () => {
        return Schedule.Reports(id);
    };

    const loadScript = () => {
        return Schedule.Script(id);
    };

    const loadGroups = () => {
        return Schedule.Groups(id);
    };

    const loadHosts = () => {
        return Schedule.Hosts(id);
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

            navigate('/schedules');
        });
    };

    const groupsList = () => {
        if (!groups || groups.length === 0) {
            return null;
        }

        return (
            <ListGroup.TextItem title="Groups">
                {groups.map((group, idx) => {
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
                {hosts.map((host, idx) => {
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

    if (loading) {
        return (<PageLoading />);
    }

    const toolbar = (
        <React.Fragment>
            <EditButton to={'/schedules/schedule/' + schedule.ID + '/edit'} />
            <DeleteButton onClick={deleteClick} />
        </React.Fragment>
    );

    const breadcrumbs = [
        {
            title: 'Schedules',
            href: '/schedules'
        },
        {
            title: schedule.Name
        }
    ];

    return (
        <Page title={breadcrumbs} toolbar={toolbar}>
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
                        <ScheduleReportList reports={reports} />
                    </Card.Card>
                </Layout.Column>
            </Layout.Row>
        </Page>
    );
};

interface ScheduleReportListProps { reports: ScheduleReport[]; }
const ScheduleReportList: React.FC<ScheduleReportListProps> = (props: ScheduleReportListProps) => {
    const [ShownReports, SetShownReports] = React.useState<ScheduleReport[]>(props.reports.slice(0, Math.min(5, props.reports.length)));

    if (props.reports.length == 0) {
        return (
            <React.Fragment>
                <Card.Header>History</Card.Header>
                <Card.Body>
                    <Nothing />
                </Card.Body>
            </React.Fragment>
        );
    }

    const showMoreClick = () => {
        SetShownReports(reports => {
            return props.reports.slice(0, Math.min(reports.length + 20, props.reports.length));
        });
    };

    const showMoreDiabled = () => {
        return ShownReports.length >= props.reports.length;
    };

    return (
        <React.Fragment>
            <Card.Header>History</Card.Header>
            <ListGroup.List>
                {ShownReports.map((report, idx) => {
                    return (<ScheduleReportItem report={report} key={idx} />);
                })}
            </ListGroup.List>
            <Card.Footer>
                <Button color={Style.Palette.Primary} onClick={showMoreClick} disabled={showMoreDiabled()}><Icon.Label icon={<Icon.Plus />} label="Show More" /></Button>
            </Card.Footer>
        </React.Fragment>
    );
};


interface ScheduleReportItemProps { report: ScheduleReport; }
const ScheduleReportItem: React.FC<ScheduleReportItemProps> = (props: ScheduleReportItemProps) => {
    let resultIcon = (<Icon.QuestionCircle color={Style.Palette.Primary} />);
    if (props.report.Result == 0) {
        resultIcon = (<Icon.CheckCircle color={Style.Palette.Success} />);
    } else if (props.report.Result == 1) {
        resultIcon = (<Icon.ExclamationTriangle color={Style.Palette.Warning} />);
    } else if (props.report.Result == 2) {
        resultIcon = (<Icon.ExclamationCircle color={Style.Palette.Danger} />);
    }

    const linkClick = () => {
        GlobalModalFrame.showModal(<ScheduleReportDetails report={props.report} />);
    };

    const link = () => {
        const h = props.report.HostIDs.length == 1 ? 'host' : 'hosts';

        return (<Link onClick={linkClick} to="#">{props.report.HostIDs.length} {h}</Link>);
    };

    return (
        <ListGroup.Item>
            {resultIcon}
            <span className="ms-1">
                <DateLabel date={props.report.Time.Start} /> on {link()}
            </span>
        </ListGroup.Item>
    );
};

interface ScheduleReportDetailsProps { report: ScheduleReport; }
const ScheduleReportDetails: React.FC<ScheduleReportDetailsProps> = (props: ScheduleReportDetailsProps) => {
    const [IsLoading, SetIsLoading] = React.useState(true);
    const [Hosts, SetHosts] = React.useState<{ [id: string]: string }>();

    React.useEffect(() => {
        Host.List().then(hosts => {
            SetHosts(() => {
                const hostMap: { [id: string]: string } = {};
                hosts.forEach(host => {
                    hostMap[host.ID] = host.Name;
                });
                return hostMap;
            });
            SetIsLoading(false);
        });
    }, []);

    if (IsLoading) {
        return (
            <Modal title="Schedule Report">
                <Loading />
            </Modal>
        );
    }

    const hostEntry = (hostID: string) => {
        const resultIcon = props.report.HostResult[hostID] == 0 ? (<Icon.CheckCircle color={Style.Palette.Success} />) : (<Icon.ExclamationCircle color={Style.Palette.Danger} />);

        return (
            <ListGroup.Item key={hostID}>
                <Icon.Label icon={resultIcon} label={Hosts[hostID]} />
            </ListGroup.Item>
        );
    };

    return (
        <Modal title="Schedule Report">
            <Layout.Row>
                <Layout.Column>
                    <Card.Card>
                        <Card.Header>Run Information</Card.Header>
                        <ListGroup.List>
                            <ListGroup.TextItem title="Started"><DateLabel date={props.report.Time.Start} /></ListGroup.TextItem>
                            <ListGroup.TextItem title="Finished"><DateLabel date={props.report.Time.Finished} /></ListGroup.TextItem>
                            <ListGroup.TextItem title="Elapsed">{props.report.Time.ElapsedSeconds} seconds</ListGroup.TextItem>
                        </ListGroup.List>
                    </Card.Card>
                    <Card.Card>
                        <Card.Header>Hosts</Card.Header>
                        <ListGroup.List>
                            {
                                props.report.HostIDs.map(hostID => {
                                    return hostEntry(hostID);
                                })
                            }
                        </ListGroup.List>
                    </Card.Card>
                </Layout.Column>
            </Layout.Row>
        </Modal>
    );
};
