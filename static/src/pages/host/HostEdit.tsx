import * as React from 'react';
import { Host, HostType } from '../../types/Host';
import { match } from 'react-router-dom';
import { URLParams } from '../../services/Params';
import { PageLoading } from '../../components/Loading';
import { Page } from '../../components/Page';
import { Input } from '../../components/input/Input';
import { Form } from '../../components/Form';
import { EnvironmentVariableEdit } from '../../components/EnvironmentVariableEdit';
import { GroupCheckList } from '../../components/CheckList';
import { Card } from '../../components/Card';
import { Notification } from '../../components/Notification';
import { Redirect } from '../../components/Redirect';
import { Variable } from '../../types/Variable';
import { RandomPSK } from '../../components/RandomPSK';

interface HostEditProps {
    match: match;
}
export const HostEdit: React.FC<HostEditProps> = (props: HostEditProps) => {
    const [loading, setLoading] = React.useState(true);
    const [host, setHost] = React.useState<HostType>();
    const [isNew, setIsNew] = React.useState<boolean>();
    const [useHostName, setUseHostname] = React.useState<boolean>();

    React.useEffect(() => {
        loadHost();
    }, []);

    const loadHost = () => {
        const id = (props.match.params as URLParams).id;
        if (id == null) {
            setIsNew(true);
            setHost(Host.Blank());
            setUseHostname(true);
            setLoading(false);
        } else {
            Host.Get(id).then(host => {
                setIsNew(false);
                setHost(host);
                setUseHostname(host.Name == host.Address);
                setLoading(false);
            });
        }
    };

    const changeName = (Name: string) => {
        setHost(host => {
            host.Name = Name;
            if (useHostName) {
                host.Address = Name;
            }
            return {...host};
        });
    };

    const changeAddress = (Address: string) => {
        setHost(host => {
            host.Address = Address;
            return {...host};
        });
    };

    const changePort = (Port: number) => {
        setHost(host => {
            host.Port = Port;
            return {...host};
        });
    };

    const changePSK = (PSK: string) => {
        setHost(host => {
            host.PSK = PSK;
            return {...host};
        });
    };

    const enabledCheckbox = () => {
        if (isNew) {
            return null;
        }

        return (
            <Input.Checkbox
                label="Enabled"
                helpText=""
                defaultValue={host.Enabled}
                onChange={changeEnabled} />
        );
    };

    const changeEnabled = (Enabled: boolean) => {
        setHost(host => {
            host.Enabled = Enabled;
            return {...host};
        });
    };

    const changeEnvironment = (Environment: Variable[]) => {
        setHost(host => {
            host.Environment = Environment;
            return {...host};
        });
    };

    const changeGroupIDs = (GroupIDs: string[]) => {
        setHost(host => {
            host.GroupIDs = GroupIDs;
            return {...host};
        });
    };

    const formSave = () => {
        let promise: Promise<HostType>;
        if (isNew) {
            promise = Host.New(host);
        } else {
            promise = Host.Save(host);
        }

        return promise.then(host => {
            Notification.success('Host Saved');
            Redirect.To('/hosts/host/' + host.ID);
        });
    };

    const changeUseHostName = (useHostName: boolean) => {
        setUseHostname(useHostName);
    };

    const addressInput = () => {
        if (useHostName) {
            return null;
        }

        return (
            <Input.Text
                label="Address"
                type="text"
                defaultValue={host.Address}
                onChange={changeAddress}
                required />
        );
    };

    if (loading) {
        return (<PageLoading />);
    }

    return (
        <Page title={ isNew ? 'New Host' : 'Edit Host' }>
            <Form showSaveButton onSubmit={formSave}>
                <Input.Text
                    label="Name"
                    type="text"
                    defaultValue={host.Name}
                    onChange={changeName}
                    required />
                <Input.Checkbox label="Connect to host using this name" defaultValue={useHostName} onChange={changeUseHostName} />
                { addressInput() }
                <Input.Number
                    label="Port"
                    defaultValue={host.Port}
                    onChange={changePort}
                    required />
                <Input.Text
                    label="Pre-Shared Key"
                    type="password"
                    defaultValue={host.PSK}
                    onChange={changePSK}
                    required />
                <RandomPSK newPSK={changePSK} />
                { enabledCheckbox() }
                <Card.Card className="mt-3">
                    <Card.Header>Environment Variables</Card.Header>
                    <Card.Body>
                        <EnvironmentVariableEdit
                            variables={host.Environment}
                            onChange={changeEnvironment} />
                    </Card.Body>
                </Card.Card>
                <Card.Card className="mt-3">
                    <Card.Header>Groups</Card.Header>
                    <Card.Body>
                        <GroupCheckList selectedGroups={host.GroupIDs} onChange={changeGroupIDs}/>
                    </Card.Body>
                </Card.Card>
            </Form>
        </Page>
    );
};
