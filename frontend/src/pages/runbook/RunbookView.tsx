import * as React from 'react';
import { Runbook, RunbookReport, RunbookType } from '../../types/Runbook';
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
import { RunbookPattern } from './RunbookPattern';
import { Icon } from '../../components/Icon';
import { Style } from '../../components/Style';
import { Nothing } from '../../components/Nothing';
import { GlobalModalFrame, Modal } from '../../components/Modal';
import { Permissions, UserAction } from '../../services/Permissions';

export const RunbookView: React.FC = () => {
    const { id } = useParams() as URLParams;
    const [loading, setLoading] = React.useState<boolean>(true);
    const [runbook, setRunbook] = React.useState<RunbookType>();
    const [reports, setReports] = React.useState<RunbookReport[]>();
    const [script, setScript] = React.useState<ScriptType>();
    const [hosts, setHosts] = React.useState<HostType[]>();
    const [groups, setGroups] = React.useState<GroupType[]>();
    const navigate = useNavigate();

    React.useEffect(() => {
        loadData();
    }, []);

    const loadRunbook = () => {
        return Runbook.Get(id);
    };

    const loadReports = () => {
        return Runbook.Reports(id);
    };

    const loadScript = () => {
        return Runbook.Script(id);
    };

    const loadGroups = () => {
        return Runbook.Groups(id);
    };

    const loadHosts = () => {
        return Runbook.Hosts(id);
    };

    const loadData = () => {
        Promise.all([loadRunbook(), loadHosts(), loadReports(), loadScript()]).then(results => {
            const runbook = results[0];
            setRunbook(runbook);
            setHosts(results[1]);
            setReports(results[2]);
            setScript(results[3]);

            let scopePromise: Promise<GroupType[]> = Promise.resolve([]);
            if (runbook.Scope.GroupIDs && runbook.Scope.GroupIDs.length > 0) {
                scopePromise = loadGroups();
            }
            scopePromise.then(groups => {
                setGroups(groups);
                setLoading(false);
            });
        });
    };

    const deleteClick = () => {
        Runbook.DeleteModal(runbook).then(deleted => {
            if (!deleted) {
                return;
            }

            navigate('/runbooks');
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
            <EditButton to={'/runbooks/runbook/' + runbook.ID + '/edit'} disabled={!Permissions.UserCan(UserAction.ModifyRunbooks)} />
            <DeleteButton onClick={deleteClick} disabled={!Permissions.UserCan(UserAction.ModifyRunbooks)} />
        </React.Fragment>
    );

    const breadcrumbs = [
        {
            title: 'Runbooks',
            href: '/runbooks'
        },
        {
            title: runbook.Name
        }
    ];

    return (
        <Page title={breadcrumbs} toolbar={toolbar}>
            <Layout.Row>
                <Layout.Column>
                    <Card.Card>
                        <Card.Header>Runbook Information</Card.Header>
                        <ListGroup.List>
                            <ListGroup.TextItem title="Name">{runbook.Name}</ListGroup.TextItem>
                            <ListGroup.TextItem title="Script"><Link to={'/scripts/script/' + runbook.ScriptID}>{script.Name}</Link></ListGroup.TextItem>
                            <ListGroup.TextItem title="Frequency"><RunbookPattern pattern={runbook.Pattern} /></ListGroup.TextItem>
                            <ListGroup.TextItem title="Last Run"><DateLabel date={runbook.LastRunTime} /></ListGroup.TextItem>
                            <ListGroup.TextItem title="Enabled"><EnabledBadge value={runbook.Enabled} /></ListGroup.TextItem>
                            {groupsList()}
                            {hostsList()}
                        </ListGroup.List>
                    </Card.Card>
                </Layout.Column>
                <Layout.Column>
                    <Card.Card>
                        <RunbookReportList reports={reports} />
                    </Card.Card>
                </Layout.Column>
            </Layout.Row>
        </Page>
    );
};

interface RunbookReportListProps { reports: RunbookReport[]; }
const RunbookReportList: React.FC<RunbookReportListProps> = (props: RunbookReportListProps) => {
    const [ShownReports, SetShownReports] = React.useState<RunbookReport[]>(props.reports.slice(0, Math.min(5, props.reports.length)));

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
                    return (<RunbookReportItem report={report} key={idx} />);
                })}
            </ListGroup.List>
            <Card.Footer>
                <Button color={Style.Palette.Primary} onClick={showMoreClick} disabled={showMoreDiabled()}><Icon.Label icon={<Icon.Plus />} label="Show More" /></Button>
            </Card.Footer>
        </React.Fragment>
    );
};


interface RunbookReportItemProps { report: RunbookReport; }
const RunbookReportItem: React.FC<RunbookReportItemProps> = (props: RunbookReportItemProps) => {
    let resultIcon = (<Icon.QuestionCircle color={Style.Palette.Primary} />);
    if (props.report.Result == 0) {
        resultIcon = (<Icon.CheckCircle color={Style.Palette.Success} />);
    } else if (props.report.Result == 1) {
        resultIcon = (<Icon.ExclamationTriangle color={Style.Palette.Warning} />);
    } else if (props.report.Result == 2) {
        resultIcon = (<Icon.ExclamationCircle color={Style.Palette.Danger} />);
    }

    const linkClick = () => {
        GlobalModalFrame.showModal(<RunbookReportDetails report={props.report} />);
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

interface RunbookReportDetailsProps { report: RunbookReport; }
const RunbookReportDetails: React.FC<RunbookReportDetailsProps> = (props: RunbookReportDetailsProps) => {
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
            <Modal title="Runbook Report">
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
        <Modal title="Runbook Report">
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
