import * as React from 'react';
import { Runbook, RunbookType } from '../../types/Runbook';
import { Link, useParams, useNavigate } from 'react-router-dom';
import { URLParams } from '../../services/Params';
import { PageLoading } from '../../components/Loading';
import { Page } from '../../components/Page';
import { Input } from '../../components/input/Input';
import { Form } from '../../components/Form';
import { Notification } from '../../components/Notification';
import { Script, ScriptType } from '../../types/Script';
import { GroupCheckList, HostCheckList } from '../../components/CheckList';
import { Card } from '../../components/Card';
import { Alert } from '../../components/Alert';
import { Icon } from '../../components/Icon';
import { Checkbox } from '../../components/input/Checkbox';
import { RadioChoice } from '../../components/input/Radio';

export const RunbookEdit: React.FC = () => {
    const { id } = useParams() as URLParams;
    const [noData, setNoData] = React.useState<boolean>();
    const [loading, setLoading] = React.useState(true);
    const [runbook, setRunbook] = React.useState<RunbookType>();
    const [isNew, setIsNew] = React.useState<boolean>();
    const [runOn, setRunOn] = React.useState<'groups' | 'hosts'>();
    const [patternTemplate, setPatternTemplate] = React.useState<string>();
    const [scripts, setScripts] = React.useState<ScriptType[]>();
    const navigate = useNavigate();

    React.useEffect(() => {
        loadData();
    }, []);

    const loadRunbook = () => {
        if (id == null) {
            return Promise.resolve(Runbook.Blank());
        } else {
            return Runbook.Get(id);
        }
    };

    const loadData = () => {
        Promise.all([loadRunbook(), Script.List()]).then(results => {
            const isNew = results[0].ID == undefined;
            const runbook = results[0];
            const scripts = results[1];
            let runOn: ('groups' | 'hosts') = 'groups';
            let patternTemplate = '';

            if (isNew && (!scripts || scripts.length === 0)) {
                setIsNew(isNew);
                setNoData(true);
                return;
            }

            if (isNew) {
                runbook.ScriptID = scripts[0].ID;
                runbook.Pattern = '0 * * * *';
                patternTemplate = '0 * * * *';
            } else {
                switch (runbook.Pattern) {
                    case '0 * * * *':
                    case '0 */4 * * *':
                    case '0 0 * * *':
                    case '0 0 * * 1':
                        patternTemplate = runbook.Pattern;
                        break;
                    default:
                        patternTemplate = 'custom';
                        break;
                }

                if (runbook.Scope.HostIDs && runbook.Scope.HostIDs.length > 0) {
                    runOn = 'hosts';
                }
            }

            setRunbook(runbook);
            setIsNew(isNew);
            setScripts(scripts);
            setRunOn(runOn);
            setPatternTemplate(patternTemplate);
            setLoading(false);
        });
    };

    const changeName = (Name: string) => {
        setRunbook(runbook => {
            runbook.Name = Name;
            return { ...runbook };
        });
    };

    const changeScriptID = (ScriptID: string) => {
        setRunbook(runbook => {
            runbook.ScriptID = ScriptID;
            return { ...runbook };
        });
    };

    const changePatternTemplate = (pattern: string) => {
        setRunbook(runbook => {
            if (pattern !== 'custom') {
                runbook.Pattern = pattern;
            } else {
                runbook.Pattern = '';
            }
            setPatternTemplate(pattern);
            return { ...runbook };
        });
    };

    const changeEnabled = (Enabled: boolean) => {
        setRunbook(runbook => {
            runbook.Enabled = Enabled;
            return { ...runbook };
        });
    };

    const enabledCheckbox = () => {
        if (isNew) {
            return null;
        }

        return (
            <Checkbox
                label="Enabled"
                defaultValue={runbook.Enabled}
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
                defaultValue={runbook.Pattern}
                onChange={changePattern}
                required />
        );
    };

    const changePattern = (Pattern: string) => {
        setRunbook(runbook => {
            runbook.Pattern = Pattern;
            return { ...runbook };
        });
    };

    const changeRunOn = (RunOn: string) => {
        if (runOn === RunOn) {
            return;
        }

        setRunbook(runbook => {
            runbook.Scope.HostIDs = [];
            runbook.Scope.GroupIDs = [];
            setRunOn(RunOn as 'groups' | 'hosts');
            return { ...runbook };
        });
    };

    const changeHostIDs = (HostIDs: string[]) => {
        setRunbook(runbook => {
            runbook.Scope.HostIDs = HostIDs;
            return { ...runbook };
        });
    };

    const changeGroupIDs = (GroupIDs: string[]) => {
        setRunbook(runbook => {
            runbook.Scope.GroupIDs = GroupIDs;
            return { ...runbook };
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
                    <HostCheckList selectedHosts={runbook.Scope.HostIDs} onChange={changeHostIDs} />
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
                    <GroupCheckList selectedGroups={runbook.Scope.GroupIDs} onChange={changeGroupIDs} />
                </Card.Body>
            </Card.Card>
        );
    };

    const formSave = () => {
        let promise: Promise<RunbookType>;
        if (isNew) {
            promise = Runbook.New(runbook);
        } else {
            promise = Runbook.Save(runbook);
        }

        return promise.then(runbook => {
            Notification.success('Runbook Saved');
            navigate('/runbooks/runbook/' + runbook.ID);
        });
    };

    if (isNew && noData) {
        return (<Page title="New Runbook">
            <Alert.Danger>
                <p>At least one script is required before you can create a runbook</p>
                <Link to="/runbooks"><Icon.Label icon={<Icon.ArrowLeft />} label="Go Back" /></Link>
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
            title: 'Runbooks',
            href: '/runbooks',
        },
        {
            title: 'New Runbook'
        }
    ];
    if (!isNew) {
        breadcrumbs[1] = {
            title: runbook.Name,
            href: '/runbooks/runbook/' + runbook.ID
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
                    defaultValue={runbook.Name}
                    onChange={changeName}
                    required />
                {enabledCheckbox()}
                <Input.Select
                    label="Script"
                    defaultValue={runbook.ScriptID}
                    onChange={changeScriptID}
                    required>
                    {scripts.map((script, idx) => {
                        return (<option value={script.ID} key={idx}>{script.Name}</option>);
                    })}
                </Input.Select>
                <Input.Select
                    label="Run Frequency"
                    helpText="All runbooks trigger in UTC/GMT time"
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
