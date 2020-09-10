import * as React from 'react';
import { Schedule } from '../../types/Schedule';
import { match } from 'react-router-dom';
import { URLParams } from '../../services/Params';
import { PageLoading } from '../../components/Loading';
import { Page } from '../../components/Page';
import { Form, Select, Input, Radio, RadioChoice, Checkbox } from '../../components/Form';
import { Notification } from '../../components/Notification';
import { Redirect } from '../../components/Redirect';
import { Script } from '../../types/Script';
import { Group } from '../../types/Group';
import { Host } from '../../types/Host';
import { GroupCheckList, HostCheckList } from '../../components/CheckList';
import { Card } from '../../components/Card';

export interface ScheduleEditProps { match: match }
interface ScheduleEditState {
    loading: boolean;
    schedule?: Schedule;
    isNew?: boolean;
    RunOn: string;
    patternTemplate?: string;
    scripts?: Script[];
    groups?: Group[];
    hosts?: Host[];
}
export class ScheduleEdit extends React.Component<ScheduleEditProps, ScheduleEditState> {
    constructor(props: ScheduleEditProps) {
        super(props);
        this.state = {
            RunOn: 'groups',
            loading: true,
        };
    }

    componentDidMount(): void {
        this.loadData();
    }

    private loadSchedule = () => {
        const id = (this.props.match.params as URLParams).id;
        if (id == null) {
            return Promise.resolve(Schedule.Blank());
        } else {
            return Schedule.Get(id);
        }
    }

    private loadData = () => {
        Promise.all([ this.loadSchedule(), Script.List(), Group.List(), Host.List() ]).then(results => {
            const isNew = results[0].ID == undefined;
            const schedule = results[0];
            const scripts = results[1];
            const groups = results[2];
            const hosts = results[3];
            let patternTemplate = '';
            if (isNew) {
                schedule.ScriptID = scripts[0].ID;
                schedule.Pattern = '0 * * * *';
                patternTemplate = '0 * * * *';
            } else {
                switch (schedule.Pattern) {
                    case '0 * * * *':
                    case '0 */4 * * *':
                    case '0 0 * * *':
                    case '0 0 * * 1':
                        patternTemplate = schedule.Pattern;
                        break;
                    default:
                        patternTemplate = 'custom';
                        break;
                }
            }

            this.setState({
                loading: false,
                schedule: schedule,
                isNew: isNew,
                scripts: scripts,
                groups: groups,
                hosts: hosts,
                patternTemplate: patternTemplate,
            });
        });
    }

    private changeScriptID = (ScriptID: string) => {
        this.setState(state => {
            const schedule = state.schedule;
            schedule.ScriptID = ScriptID;
            return { schedule: schedule };
        });
    }

    private changePatternTemplate = (pattern: string) => {
        this.setState(state => {
            const schedule = state.schedule;
            if (pattern !== 'custom') {
                schedule.Pattern = pattern;
            } else {
                schedule.Pattern = '';
            }
            return { schedule: schedule, patternTemplate: pattern };
        });
    }

    private changeEnabled = (Enabled: boolean) => {
        this.setState(state => {
            const schedule = state.schedule;
            schedule.Enabled = Enabled;
            return { schedule: schedule };
        });
    }

    private enabledCheckbox = () => {
        if (this.state.isNew) { return null; }

        return (
            <Checkbox
                label="Enabled"
                defaultValue={this.state.schedule.Enabled}
                onChange={this.changeEnabled} />
        );
    }

    private cronPatternInput = () => {
        if (this.state.patternTemplate !== 'custom') {
            return null;
        }

        return (
            <Input
                label="Frequency Expression"
                type="text"
                helpText="Cron expression"
                defaultValue={this.state.schedule.Pattern}
                onChange={this.changePattern}
                required />
        );
    }

    private changePattern = (Pattern: string) => {
        this.setState(state => {
            const schedule = state.schedule;
            schedule.Pattern = Pattern;
            return { schedule: schedule };
        });
    }

    private changeRunOn = (RunOn: string) => {
        this.setState(state => {
            const schedule = state.schedule;
            schedule.Scope.HostIDs = [];
            schedule.Scope.GroupIDs = [];
            return { schedule: schedule, RunOn: RunOn };
        });
    }

    private changeHostIDs = (HostIDs: string[]) => {
        this.setState(state => {
            const schedule = state.schedule;
            schedule.Scope.HostIDs = HostIDs;
            return { schedule: schedule };
        });
    }

    private changeGroupIDs = (GroupIDs: string[]) => {
        this.setState(state => {
            const schedule = state.schedule;
            schedule.Scope.GroupIDs = GroupIDs;
            return { schedule: schedule };
        });
    }

    private hostList = () => {
        if (this.state.RunOn !== 'hosts') { return null; }

        return (
            <Card.Card>
                <Card.Header>Hosts</Card.Header>
                <Card.Body>
                    <HostCheckList selectedHosts={this.state.schedule.Scope.HostIDs} onChange={this.changeHostIDs} />
                </Card.Body>
            </Card.Card>
        );
    }

    private groupList = () => {
        if (this.state.RunOn !== 'groups') { return null; }

        return (
            <Card.Card>
                <Card.Header>Groups</Card.Header>
                <Card.Body>
                    <GroupCheckList selectedGroups={this.state.schedule.Scope.GroupIDs} onChange={this.changeGroupIDs} />
                </Card.Body>
            </Card.Card>
        );
    }

    private formSave = () => {
        let promise: Promise<Schedule>;
        if (this.state.isNew) {
            promise = Schedule.New(this.state.schedule);
        } else {
            promise = this.state.schedule.Save();
        }

        promise.then(schedule => {
            Notification.success('Schedule Saved');
            Redirect.To('/schedules/schedule/' + schedule.ID);
        });
    }

    render(): JSX.Element {
        if (this.state.loading) { return (<PageLoading />); }

        const runOnChoices: RadioChoice[] = [
            {
                label: 'Individual Hosts',
                value: 'hosts'
            },
            {
                label: 'Groups',
                value: 'groups'
            }
        ];

        return (
        <Page title={ this.state.isNew ? 'New Schedule' : 'Edit Schedule' }>
            <Form showSaveButton onSubmit={this.formSave}>
                <Select
                    label="Script"
                    defaultValue={this.state.schedule.ScriptID}
                    onChange={this.changeScriptID}
                    required>
                    { this.state.scripts.map((script, idx) => {
                        return (<option value={script.ID} key={idx}>{script.Name}</option>);
                    })}
                </Select>
                <Select
                    label="Run Frequency"
                    defaultValue={this.state.patternTemplate}
                    onChange={this.changePatternTemplate}
                    required>
                        <option value="0 * * * *">Every Hour</option>
                        <option value="0 */4 * * *">Every 4 Hours</option>
                        <option value="0 0 * * *">Every Day at Midnight</option>
                        <option value="0 0 * * 1">Every Monday at Midnight</option>
                        <option value="custom">Custom</option>
                </Select>
                { this.enabledCheckbox() }
                { this.cronPatternInput() }
                <Radio
                    label="Run On"
                    onChange={this.changeRunOn}
                    choices={runOnChoices}
                    defaultValue={this.state.RunOn} />
                { this.hostList() }
                { this.groupList() }
            </Form>
        </Page>
        );
    }
}
