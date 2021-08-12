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
            return { ...host };
        });
    };

    const changeAddress = (Address: string) => {
        setHost(host => {
            host.Address = Address;
            return { ...host };
        });
    };

    const changePort = (Port: number) => {
        setHost(host => {
            host.Port = Port;
            return { ...host };
        });
    };

    const changePSK = (PSK: string) => {
        setHost(host => {
            host.PSK = PSK;
            return { ...host };
        });
    };

    const enabledCheckbox = () => {
        if (isNew) {
            return null;
        }

        return (
            <Input.Checkbox
                label="Include Host in Scripts"
                helpText="If unchecked this host will not be included in scripts or schedules that target this host"
                defaultValue={host.Enabled}
                onChange={changeEnabled} />
        );
    };

    const changeEnabled = (Enabled: boolean) => {
        setHost(host => {
            host.Enabled = Enabled;
            return { ...host };
        });
    };

    const changeEnvironment = (Environment: Variable[]) => {
        setHost(host => {
            host.Environment = Environment;
            return { ...host };
        });
    };

    const changeGroupIDs = (GroupIDs: string[]) => {
        setHost(host => {
            host.GroupIDs = GroupIDs;
            return { ...host };
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
        changeAddress(useHostName ? host.Name : '');
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

    const breadcrumbs = [
        {
            title: 'Hosts',
            href: '/hosts',
        },
        {
            title: 'New Host'
        }
    ];
    if (!isNew) {
        breadcrumbs[1] = {
            title: host.Name,
            href: '/hosts/host/' + host.ID
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
                    defaultValue={host.Name}
                    onChange={changeName}
                    required />
                <Input.Checkbox label="Connect to host using this name" defaultValue={useHostName} onChange={changeUseHostName} />
                {addressInput()}
                <Input.Number
                    label="Port"
                    defaultValue={host.Port}
                    onChange={changePort}
                    required />
                <Input.Password
                    label="Pre-Shared Key"
                    defaultValue={host.PSK}
                    onChange={changePSK}
                    required />
                {enabledCheckbox()}
                <Card.Card className="mt-3">
                    <Card.Header>Environment Variables</Card.Header>
                    <Card.Body>
                        <EnvironmentVariableEdit
                            variables={host.Environment || []}
                            onChange={changeEnvironment} />
                    </Card.Body>
                </Card.Card>
                <Card.Card className="mt-3">
                    <Card.Header>Groups</Card.Header>
                    <Card.Body>
                        <GroupCheckList selectedGroups={host.GroupIDs} onChange={changeGroupIDs} />
                    </Card.Body>
                </Card.Card>
            </Form>
        </Page>
    );
};
