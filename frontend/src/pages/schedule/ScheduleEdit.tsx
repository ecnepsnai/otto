import * as React from 'react';
import { Schedule, ScheduleType } from '../../types/Schedule';
import { Link, useParams } from 'react-router-dom';
import { URLParams } from '../../services/Params';
import { PageLoading } from '../../components/Loading';
import { Page } from '../../components/Page';
import { Input } from '../../components/input/Input';
import { Form } from '../../components/Form';
import { Notification } from '../../components/Notification';
import { Redirect } from '../../components/Redirect';
import { Script, ScriptType } from '../../types/Script';
import { GroupCheckList, HostCheckList } from '../../components/CheckList';
import { Card } from '../../components/Card';
import { Alert } from '../../components/Alert';
import { Icon } from '../../components/Icon';
import { Checkbox } from '../../components/input/Checkbox';
import { RadioChoice } from '../../components/input/Radio';

export const ScheduleEdit: React.FC = () => {
    const { id } = useParams() as URLParams;
    const [noData, setNoData] = React.useState<boolean>();
    const [loading, setLoading] = React.useState(true);
    const [schedule, setSchedule] = React.useState<ScheduleType>();
    const [isNew, setIsNew] = React.useState<boolean>();
    const [runOn, setRunOn] = React.useState<'groups' | 'hosts'>();
    const [patternTemplate, setPatternTemplate] = React.useState<string>();
    const [scripts, setScripts] = React.useState<ScriptType[]>();

    React.useEffect(() => {
        loadData();
    }, []);

    const loadSchedule = () => {
        if (id == null) {
            return Promise.resolve(Schedule.Blank());
        } else {
            return Schedule.Get(id);
        }
    };

    const loadData = () => {
        Promise.all([loadSchedule(), Script.List()]).then(results => {
            const isNew = results[0].ID == undefined;
            const schedule = results[0];
            const scripts = results[1];
            let runOn: ('groups' | 'hosts') = 'groups';
            let patternTemplate = '';

            if (isNew && (!scripts || scripts.length === 0)) {
                setIsNew(isNew);
                setNoData(true);
                return;
            }

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

                if (schedule.Scope.HostIDs && schedule.Scope.HostIDs.length > 0) {
                    runOn = 'hosts';
                }
            }

            setSchedule(schedule);
            setIsNew(isNew);
            setScripts(scripts);
            setRunOn(runOn);
            setPatternTemplate(patternTemplate);
            setLoading(false);
        });
    };

    const changeName = (Name: string) => {
        setSchedule(schedule => {
            schedule.Name = Name;
            return { ...schedule };
        });
    };

    const changeScriptID = (ScriptID: string) => {
        setSchedule(schedule => {
            schedule.ScriptID = ScriptID;
            return { ...schedule };
        });
    };

    const changePatternTemplate = (pattern: string) => {
        setSchedule(schedule => {
            if (pattern !== 'custom') {
                schedule.Pattern = pattern;
            } else {
                schedule.Pattern = '';
            }
            setPatternTemplate(pattern);
            return { ...schedule };
        });
    };

    const changeEnabled = (Enabled: boolean) => {
        setSchedule(schedule => {
            schedule.Enabled = Enabled;
            return { ...schedule };
        });
    };

    const enabledCheckbox = () => {
        if (isNew) {
            return null;
        }

        return (
            <Checkbox
                label="Enabled"
                defaultValue={schedule.Enabled}
                onChange={changeEnabled} />
        );
    };

    const cronPatternInput = () => {
        if (patternTemplate !== 'custom') {
            return null;
        }

        return (
            <Input.Text
                label="Frequency Expression"
                type="text"
                helpText="Cron expression"
                defaultValue={schedule.Pattern}
                onChange={changePattern}
                required />
        );
    };

    const changePattern = (Pattern: string) => {
        setSchedule(schedule => {
            schedule.Pattern = Pattern;
            return { ...schedule };
        });
    };

    const changeRunOn = (RunOn: string) => {
        if (runOn === RunOn) {
            return;
        }

        setSchedule(schedule => {
            schedule.Scope.HostIDs = [];
            schedule.Scope.GroupIDs = [];
            setRunOn(RunOn as 'groups' | 'hosts');
            return { ...schedule };
        });
    };

    const changeHostIDs = (HostIDs: string[]) => {
        setSchedule(schedule => {
            schedule.Scope.HostIDs = HostIDs;
            return { ...schedule };
        });
    };

    const changeGroupIDs = (GroupIDs: string[]) => {
        setSchedule(schedule => {
            schedule.Scope.GroupIDs = GroupIDs;
            return { ...schedule };
        });
    };

    const hostList = () => {
        if (runOn !== 'hosts') {
            return null;
        }

        return (
            <Card.Card>
                <Card.Header>Hosts</Card.Header>
                <Card.Body>
                    <HostCheckList selectedHosts={schedule.Scope.HostIDs} onChange={changeHostIDs} />
                </Card.Body>
            </Card.Card>
        );
    };

    const groupList = () => {
        if (runOn !== 'groups') {
            return null;
        }

        return (
            <Card.Card>
                <Card.Header>Groups</Card.Header>
                <Card.Body>
                    <GroupCheckList selectedGroups={schedule.Scope.GroupIDs} onChange={changeGroupIDs} />
                </Card.Body>
            </Card.Card>
        );
    };

    const formSave = () => {
        let promise: Promise<ScheduleType>;
        if (isNew) {
            promise = Schedule.New(schedule);
        } else {
            promise = Schedule.Save(schedule);
        }

        return promise.then(schedule => {
            Notification.success('Schedule Saved');
            Redirect.To('/schedules/schedule/' + schedule.ID);
        });
    };

    if (isNew && noData) {
        return (<Page title="New Schedule">
            <Alert.Danger>
                <p>At least one script is required before you can create a schedule</p>
                <Link to="/schedules"><Icon.Label icon={<Icon.ArrowLeft />} label="Go Back" /></Link>
            </Alert.Danger>
        </Page>);
    }

    if (loading) {
        return (<PageLoading />);
    }

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

    const breadcrumbs = [
        {
            title: 'Schedules',
            href: '/schedules',
        },
        {
            title: 'New Schedule'
        }
    ];
    if (!isNew) {
        breadcrumbs[1] = {
            title: schedule.Name,
            href: '/schedules/schedule/' + schedule.ID
        };
        breadcrumbs.push({
            title: 'Edit'
        });
    }

    return (
        <Page title={breadcrumbs}>
            <Form showSaveButton onSubmit={formSave}>
                <Input.Text
                    label="Name"
                    type="text"
                    defaultValue={schedule.Name}
                    onChange={changeName}
                    required />
                {enabledCheckbox()}
                <Input.Select
                    label="Script"
                    defaultValue={schedule.ScriptID}
                    onChange={changeScriptID}
                    required>
                    {scripts.map((script, idx) => {
                        return (<option value={script.ID} key={idx}>{script.Name}</option>);
                    })}
                </Input.Select>
                <Input.Select
                    label="Run Frequency"
                    helpText="All schedules trigger in UTC/GMT time"
                    defaultValue={patternTemplate}
                    onChange={changePatternTemplate}
                    required>
                    <option value="0 * * * *">Every Hour</option>
                    <option value="0 */4 * * *">Every 4 Hours</option>
                    <option value="0 0 * * *">Every Day at Midnight</option>
                    <option value="0 0 * * 1">Every Monday at Midnight</option>
                    <option value="custom">Custom</option>
                </Input.Select>
                {cronPatternInput()}
                <Input.Radio
                    label="Run On"
                    onChange={changeRunOn}
                    choices={runOnChoices}
                    defaultValue={runOn} />
                {hostList()}
                {groupList()}
            </Form>
        </Page>
    );
};
