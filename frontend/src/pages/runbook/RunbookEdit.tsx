import * as React from 'react';
import { Runbook as RunbookAPI, RunbookType } from '../../types/Runbook';
import { Group as GroupAPI } from '../../types/Group';
import { useParams, useNavigate } from 'react-router-dom';
import { URLParams } from '../../services/Params';
import { Page } from '../../components/Page';
import { Form } from '../../components/Form';
import { Loading, PageLoading } from '../../components/Loading';
import { Input } from '../../components/input/Input';
import { Notification } from '../../components/Notification';
import { Card } from '../../components/Card';
import { GroupCheckList } from '../../components/CheckList';
import { LeftRightInput } from '../../components/LeftRightInput';
import { ScriptType } from '../../types/Script';
import { RunLevel } from '../../components/input/RunLevel';
import { ScriptRunLevel } from '../../types/gengo_enum';
import { Rand } from '../../services/Rand';

export const RunbookEdit: React.FC = () => {
    const { id } = useParams() as URLParams;
    const [IsLoading, SetIsLoading] = React.useState<boolean>(true);
    const [IsNew, SetIsNew] = React.useState<boolean>();
    const [Runbook, SetRunbook] = React.useState<RunbookType>();
    const [ScriptLoading, SetScriptLoading] = React.useState(true);
    const [PossibleScripts, SetPossibleScripts] = React.useState<ScriptType[]>([]);
    const [ScriptReloadID, SetScriptReloadID] = React.useState(Rand.ID());
    const navigate = useNavigate();

    React.useEffect(() => {
        loadRunbook();
    }, []);

    React.useEffect(() => {
        if (!Runbook) {
            return;
        }

        if (Runbook.GroupIDs.length == 0) {
            SetScriptLoading(false);
            SetPossibleScripts([]);
        } else {
            SetScriptLoading(true);
            Promise.all(Runbook.GroupIDs.map(groupId => {
                return GroupAPI.Scripts(groupId);
            })).then(results => {
                const possibleScripts: ScriptType[] = [];
                results.forEach(scripts => {
                    scripts.forEach(script => {
                        if (possibleScripts.findIndex(s => s.ID == script.ID) == -1) {
                            possibleScripts.push(script);
                        }
                    });
                });
                SetPossibleScripts(possibleScripts);
                SetScriptLoading(false);
            });
        }
    }, [ScriptReloadID]);

    const loadRunbook = () => {
        if (id == null) {
            SetIsNew(true);
            SetRunbook(RunbookAPI.Blank());
            SetIsLoading(false);
        } else {
            RunbookAPI.Get(id).then(runbook => {
                SetRunbook(runbook);
                SetIsNew(false);
                SetIsLoading(false);
            });
        }
    };

    const formSave = () => {
        let promise: Promise<RunbookType>;
        if (IsNew) {
            promise = RunbookAPI.New(Runbook);
        } else {
            promise = RunbookAPI.Save(Runbook);
        }

        return promise.then(runbook => {
            Notification.success('Group Saved');
            navigate('/runbooks/runbook/' + runbook.ID);
        });
    };

    const changeName = (Name: string) => {
        SetRunbook(runbook => {
            runbook.Name = Name;
            return {...runbook};
        });
    };

    const changeScriptFailureMode = (Mode: string) => {
        SetRunbook(runbook => {
            runbook.HaltOnFailure = Mode == 'halt';
            return {...runbook};
        });
    };

    const changeRunLevel = (RunLevel: ScriptRunLevel) => {
        SetRunbook(runbook => {
            runbook.RunLevel = RunLevel;
            return {...runbook};
        });
    };

    const changeGroupIDs = (GroupIDs: string[]) => {
        SetRunbook(runbook => {
            runbook.GroupIDs = GroupIDs;
            return {...runbook};
        });
        SetScriptReloadID(Rand.ID);
    };

    const changeScriptIDs = (ScriptIDs: string[]) => {
        SetRunbook(runbook => {
            runbook.ScriptIDs = ScriptIDs;
            return {...runbook};
        });
    };

    if (IsLoading) {
        return (<PageLoading />);
    }

    const breadcrumbs = [
        {
            title: 'Runbooks',
            href: '/runbooks',
        },
        {
            title: 'New Runbook'
        }
    ];
    if (!IsNew) {
        breadcrumbs[1] = {
            title: Runbook.Name,
            href: '/runbooks/runbook/' + Runbook.ID
        };
        breadcrumbs.push({
            title: 'Edit'
        });
    }

    const scriptFailureChoices = [
        {
            label: 'Stop Runbook',
            value: 'halt'
        },
        {
            label: 'Continue to Next Script',
            value: 'continue'
        }
    ];

    const scriptChoices = () => {
        return PossibleScripts.map(s => {
            return {
                label: s.Name,
                value: s.ID,
            };
        });
    };

    return (
        <Page title={breadcrumbs}>
            <Form showSaveButton onSubmit={formSave}>
                <Input.Text
                    label="Name"
                    type="text"
                    defaultValue={Runbook.Name}
                    onChange={changeName}
                    required />
                <Input.Radio
                    label="On Script Failure"
                    buttons
                    choices={scriptFailureChoices}
                    defaultValue={Runbook.HaltOnFailure ? 'halt' : 'continue'}
                    onChange={changeScriptFailureMode}
                    />
                <RunLevel
                    defaultValue={Runbook.RunLevel}
                    onChange={changeRunLevel}
                    />
                <Card.Card className="mt-3">
                    <Card.Header>Groups</Card.Header>
                    <Card.Body>
                        <GroupCheckList selectedGroups={Runbook.GroupIDs} onChange={changeGroupIDs} />
                    </Card.Body>
                </Card.Card>
                <Card.Card className="mt-3">
                    <Card.Header>Scripts</Card.Header>
                    <Card.Body>
                        { ScriptLoading ? (<Loading />) : (<LeftRightInput leftLabel="Availble Scripts" rightLabel="Selected Scripts" choices={scriptChoices()} selected={Runbook.ScriptIDs} onChange={changeScriptIDs} />) }
                    </Card.Body>
                </Card.Card>
            </Form>
        </Page>
    );
};
