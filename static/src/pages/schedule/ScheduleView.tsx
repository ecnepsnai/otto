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
interface ScheduleViewState {
    loading: boolean;
    schedule?: ScheduleType;
    reports?: ScheduleReport[];
    script?: ScriptType;
    hosts?: HostType[];
    groups?: GroupType[];
}
export class ScheduleView extends React.Component<ScheduleViewProps, ScheduleViewState> {
    private scheduleID: string;

    constructor(props: ScheduleViewProps) {
        super(props);
        this.scheduleID = (this.props.match.params as URLParams).id;
        this.state = {
            loading: true,
        };
    }

    private loadSchedule = async () => {
        const schedule = await Schedule.Get(this.scheduleID);
        this.setState({
            schedule: schedule,
        });
    }

    private loadReports = async () => {
        const reports = await Schedule.Reports(this.scheduleID);
        this.setState({
            reports: reports,
        });
    }

    private loadScript = async () => {
        const script = await Schedule.Script(this.scheduleID);
        this.setState({
            script: script,
        });
    }

    private loadGroups = async () => {
        const groups = await Schedule.Groups(this.scheduleID);
        this.setState({
            groups: groups,
        });
    }

    private loadHosts = async () => {
        const hosts = await Schedule.Hosts(this.scheduleID);
        this.setState({
            hosts: hosts,
        });
    }

    private loadData = () => {
        Promise.all([this.loadSchedule(), this.loadHosts(), this.loadReports(), this.loadScript()]).then(() => {
            let scopePromise: Promise<void> = Promise.resolve();
            if (this.state.schedule.Scope.GroupIDs && this.state.schedule.Scope.GroupIDs.length > 0) {
                scopePromise = this.loadGroups();
            }
            scopePromise.then(() => {
                this.setState({
                    loading: false,
                });
            });
        });
    }

    componentDidMount(): void {
        this.loadData();
    }

    private deleteClick = () => {
        Schedule.DeleteModal(this.state.schedule).then(deleted => {
            if (!deleted) {
                return;
            }

            Redirect.To('/schedules');
        });
    }

    private groupsList = () => {
        if (!this.state.groups) {
            return null;
        }

        return (
            <ListGroup.TextItem title="Groups">
                { this.state.groups.map((group, idx) => {
                    return (
                        <div key={idx}>
                            <Icon.LayerGroup />
                            <Link className="ms-1" to={'/groups/group/' + group.ID}>{group.Name}</Link>
                        </div>
                    );
                }) }
            </ListGroup.TextItem>
        );
    }

    private hostsList = () => {
        return (
            <ListGroup.TextItem title="Hosts">
                { this.state.hosts.map((host, idx) => {
                    return (
                        <div key={idx}>
                            <Icon.Desktop />
                            <Link className="ms-1" to={'/hosts/host/' + host.ID}>{host.Name}</Link>
                        </div>
                    );
                }) }
            </ListGroup.TextItem>
        );
    }

    private historyContent = () => {
        if (this.state.reports.length == 0) {
            return (<Card.Body><Nothing /></Card.Body>);
        }

        return (<ListGroup.List>
            { this.state.reports.map((report, idx) => {
                return (<ScheduleReportItem report={report} key={idx}/>);
            })}
        </ListGroup.List>);
    }

    render(): JSX.Element {
        if (this.state.loading) {
            return (<PageLoading />);
        }

        return (
            <Page title="View Schedule">
                <Layout.Container>
                    <Buttons>
                        <EditButton to={'/schedules/schedule/' + this.state.schedule.ID + '/edit'} />
                        <DeleteButton onClick={this.deleteClick} />
                    </Buttons>
                    <Layout.Row>
                        <Layout.Column>
                            <Card.Card>
                                <Card.Header>Schedule Information</Card.Header>
                                <ListGroup.List>
                                    <ListGroup.TextItem title="Name">{ this.state.schedule.Name }</ListGroup.TextItem>
                                    <ListGroup.TextItem title="Script"><Link to={'/scripts/script/' + this.state.schedule.ScriptID}>{this.state.script.Name}</Link></ListGroup.TextItem>
                                    <ListGroup.TextItem title="Frequency"><SchedulePattern pattern={this.state.schedule.Pattern} /></ListGroup.TextItem>
                                    <ListGroup.TextItem title="Last Run"><DateLabel date={this.state.schedule.LastRunTime} /></ListGroup.TextItem>
                                    <ListGroup.TextItem title="Enabled"><EnabledBadge value={this.state.schedule.Enabled} /></ListGroup.TextItem>
                                    { this.groupsList() }
                                    { this.hostsList() }
                                </ListGroup.List>
                            </Card.Card>
                        </Layout.Column>
                        <Layout.Column>
                            <Card.Card>
                                <Card.Header>History</Card.Header>
                                {this.historyContent()}
                            </Card.Card>
                        </Layout.Column>
                    </Layout.Row>
                </Layout.Container>
            </Page>
        );
    }
}

interface ScheduleReportItemProps { report: ScheduleReport; }
class ScheduleReportItem extends React.Component<ScheduleReportItemProps, unknown> {
    render(): JSX.Element {
        let icon = (<Icon.QuestionCircle color={Style.Palette.Primary} />);
        if (this.props.report.Result == 0) {
            icon = (<Icon.CheckCircle color={Style.Palette.Success} />);
        } else if (this.props.report.Result == 1) {
            icon = (<Icon.ExclamationTriangle color={Style.Palette.Warning} />);
        } else if (this.props.report.Result == 2) {
            icon = (<Icon.ExclamationCircle color={Style.Palette.Danger} />);
        }

        return (
            <ListGroup.Item>
                {icon}
                <span className="ms-1">
                    <DateLabel date={this.props.report.Time.Start} /> on {this.props.report.HostIDs.length} hosts
                </span>
            </ListGroup.Item>
        );
    }
}
