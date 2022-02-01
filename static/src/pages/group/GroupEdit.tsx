import * as React from 'react';
import { Group, GroupType } from '../../types/Group';
import { useParams } from 'react-router-dom';
import { URLParams } from '../../services/Params';
import { PageLoading } from '../../components/Loading';
import { Page } from '../../components/Page';
import { Input } from '../../components/input/Input';
import { Form } from '../../components/Form';
import { EnvironmentVariableEdit } from '../../components/EnvironmentVariableEdit';
import { ScriptCheckList, HostCheckList } from '../../components/CheckList';
import { Card } from '../../components/Card';
import { Notification } from '../../components/Notification';
import { Redirect } from '../../components/Redirect';
import { Variable } from '../../types/Variable';

export const GroupEdit: React.FC = () => {
    const { id } = useParams() as URLParams;
    const [loading, setLoading] = React.useState(true);
    const [group, setGroup] = React.useState<GroupType>();
    const [isNew, setIsNew] = React.useState<boolean>();
    const [hostIDs, setHostIDs] = React.useState<string[]>();

    React.useEffect(() => {
        loadGroup();
    }, []);

    const loadGroup = () => {
        if (id == null) {
            setIsNew(true);
            setGroup(Group.Blank());
            setLoading(false);
            setHostIDs([]);
        } else {
            Group.Get(id).then(group => {
                Group.Hosts(group.ID).then(hostIDs => {
                    setIsNew(false);
                    setGroup(group);
                    setHostIDs(hostIDs.map(host => host.ID));
                    setLoading(false);
                });
            });
        }
    };

    const changeName = (Name: string) => {
        setGroup(group => {
            group.Name = Name;
            return { ...group };
        });
    };

    const changeEnvironment = (Environment: Variable[]) => {
        setGroup(group => {
            group.Environment = Environment;
            return { ...group };
        });
    };

    const changeScriptIDs = (ScriptIDs: string[]) => {
        setGroup(group => {
            group.ScriptIDs = ScriptIDs;
            return { ...group };
        });
    };

    const changeHostIDs = (HostIDs: string[]) => {
        setHostIDs(HostIDs);
    };

    const formSave = () => {
        let promise: Promise<GroupType>;
        if (isNew) {
            promise = Group.New(group);
        } else {
            promise = Group.Save(group);
        }

        return promise.then(group => {
            Group.SetHosts(group.ID, hostIDs).then(() => {
                Notification.success('Group Saved');
                Redirect.To('/groups/group/' + group.ID);
            });
        });
    };

    if (loading) {
        return (<PageLoading />);
    }

    const breadcrumbs = [
        {
            title: 'Groups',
            href: '/groups',
        },
        {
            title: 'New Group'
        }
    ];
    if (!isNew) {
        breadcrumbs[1] = {
            title: group.Name,
            href: '/groups/group/' + group.ID
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
                    defaultValue={group.Name}
                    onChange={changeName}
                    required />
                <Card.Card className="mt-3">
                    <Card.Header>Environment Variables</Card.Header>
                    <Card.Body>
                        <EnvironmentVariableEdit
                            variables={group.Environment || []}
                            onChange={changeEnvironment} />
                    </Card.Body>
                </Card.Card>
                <Card.Card className="mt-3">
                    <Card.Header>Scripts</Card.Header>
                    <Card.Body>
                        <ScriptCheckList selectedScripts={group.ScriptIDs} onChange={changeScriptIDs} />
                    </Card.Body>
                </Card.Card>
                <Card.Card className="mt-3">
                    <Card.Header>Hosts</Card.Header>
                    <Card.Body>
                        <HostCheckList selectedHosts={hostIDs} onChange={changeHostIDs} />
                    </Card.Body>
                </Card.Card>
            </Form>
        </Page>
    );
};
