import * as React from 'react';
import { Runbook as RunbookAPI, RunbookType } from '../../types/Runbook';
import { useParams, useNavigate } from 'react-router-dom';
import { URLParams } from '../../services/Params';
import { Page } from '../../components/Page';
import { Form } from '../../components/Form';
import { PageLoading } from '../../components/Loading';
import { Input } from '../../components/input/Input';
import { Notification } from '../../components/Notification';
import { Card } from '../../components/Card';
import { GroupCheckList, ScriptCheckList } from '../../components/CheckList';
import { Rand } from '../../services/Rand';

export const RunbookEdit: React.FC = () => {
    const { id } = useParams() as URLParams;
    const [IsLoading, SetIsLoading] = React.useState<boolean>(true);
    const [IsNew, SetIsNew] = React.useState<boolean>();
    const [Runbook, SetRunbook] = React.useState<RunbookType>();
    const [ScriptListID, SetScriptListID] = React.useState(Rand.ID());
    const navigate = useNavigate();

    React.useEffect(() => {
        loadRunbook();
    }, []);

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

    const changeGroupIDs = (GroupIDs: string[]) => {
        SetRunbook(runbook => {
            runbook.GroupIDs = GroupIDs;
            SetScriptListID(Rand.ID());
            return {...runbook};
        });
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

    return (
        <Page title={breadcrumbs}>
            <Form showSaveButton onSubmit={formSave}>
                <Input.Text
                    label="Name"
                    type="text"
                    defaultValue={Runbook.Name}
                    onChange={changeName}
                    required />
                <Card.Card className="mt-3">
                    <Card.Header>Groups</Card.Header>
                    <Card.Body>
                        <GroupCheckList selectedGroups={Runbook.GroupIDs} onChange={changeGroupIDs} />
                    </Card.Body>
                </Card.Card>
                <Card.Card className="mt-3">
                    <Card.Header>Scripts (Available to Groups)</Card.Header>
                    <Card.Body>
                        <ScriptCheckList selectedScripts={Runbook.ScriptIDs} groupIds={Runbook.GroupIDs} onChange={changeScriptIDs} key={ScriptListID} />
                    </Card.Body>
                </Card.Card>
            </Form>
        </Page>
    );
};
