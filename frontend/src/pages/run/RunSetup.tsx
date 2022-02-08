import * as React from 'react';
import { Script } from '../../types/Script';
import { Loading } from '../../components/Loading';
import { Input } from '../../components/input/Input';
import { Form } from '../../components/Form';
import { Card } from '../../components/Card';
import { Alert } from '../../components/Alert';

interface SGroup {
    ID: string;
    Name: string;
}

interface SHost {
    ID: string;
    Name: string;
}

interface RunSetupProps {
    scriptID: string;
    onSelectedHosts: (hostIDs: string[]) => (void);
}
export const RunSetup: React.FC<RunSetupProps> = (props: RunSetupProps) => {
    const [loading, setLoading] = React.useState<boolean>(true);
    const [groups, setGroups] = React.useState<SGroup[]>();
    const [hosts, setHosts] = React.useState<SHost[]>();
    const [selectedGroups, setSelectedGroups] = React.useState<{ [id: string]: number }>();
    const [selectedHosts, setSelectedHosts] = React.useState<{ [id: string]: number }>();
    const [groupMembers, setGroupMembers] = React.useState<{ [id: string]: string[] }>();

    const loadData = () => {
        Script.Get(props.scriptID).then(script => {
            Script.Hosts(script.ID).then(hosts => {
                const groupMap: { [id: string]: string } = {};
                const groupMembership: { [id: string]: string[] } = {};
                const selectedGroups: { [id: string]: number } = {};
                const selectedHosts: { [id: string]: number } = {};

                const shosts: SHost[] = [];
                hosts.forEach(host => {
                    groupMap[host.GroupID] = host.GroupName;
                    shosts.push({ ID: host.HostID, Name: host.HostName });
                    if (groupMembership[host.GroupID] == undefined) {
                        groupMembership[host.GroupID] = [];
                    }
                    groupMembership[host.GroupID].push(host.HostID);
                    selectedHosts[host.HostID] = 0;
                    selectedGroups[host.GroupID] = 0;
                });

                const groups: SGroup[] = [];
                Object.keys(groupMap).forEach(id => {
                    groups.push({ ID: id, Name: groupMap[id] });
                });

                setGroups(groups);
                setHosts(shosts);
                setSelectedGroups(selectedGroups);
                setSelectedHosts(selectedHosts);
                setGroupMembers(groupMembership);
                setLoading(false);
            });
        });
    };

    React.useEffect(() => {
        loadData();
    }, []);

    React.useEffect(() => {
        if (selectedHosts == undefined || selectedGroups == undefined) {
            return;
        }

        props.onSelectedHosts(selectedHostIDs());
    }, [selectedHosts, selectedGroups]);

    const selectGroup = (groupID: string) => {
        return (checked: boolean) => {
            setSelectedGroups(selectedGroups => {
                if (checked) {
                    selectedGroups[groupID] = 1;
                } else {
                    selectedGroups[groupID] = 0;
                }
                return { ...selectedGroups };
            });
            setSelectedHosts(selectedHosts => {
                groupMembers[groupID].forEach(host => {
                    if (checked) {
                        selectedHosts[host] = 1;
                    } else {
                        selectedHosts[host] = 0;
                    }
                });
                return { ...selectedHosts };
            });
        };
    };

    const selectHost = (hostID: string) => {
        return (checked: boolean) => {
            setSelectedHosts(selectedHosts => {
                const selected: { [id: string]: number } = selectedHosts;
                if (checked) {
                    selected[hostID] = 1;
                } else {
                    selected[hostID] = 0;
                }
                return { ...selectedHosts };
            });
        };
    };

    const selectedHostIDs = () => {
        const hostIDs: string[] = [];
        Object.keys(selectedHosts).forEach(hostID => {
            if (selectedHosts[hostID] > 0) {
                hostIDs.push(hostID);
            }
        });
        return hostIDs;
    };

    if (loading) {
        return (<Loading />);
    }

    if ((hosts || []).length == 0) {
        return (<Alert.Warning>
            <h5>No Hosts Found</h5>
            <p>There are no hosts in any of the groups associated with this script.</p>
        </Alert.Warning>);
    }

    return (
        <Form>
            <Card.Card>
                <Card.Header>Groups</Card.Header>
                <Card.Body>
                    {
                        groups.map((group, idx) => {
                            return (
                                <Input.Checkbox label={group.Name} onChange={selectGroup(group.ID)} defaultValue={selectedGroups[group.ID] > 0} key={idx} />
                            );
                        })
                    }
                </Card.Body>
            </Card.Card>
            <Card.Card className="mt-3">
                <Card.Header>Hosts</Card.Header>
                <Card.Body>
                    {
                        hosts.map((host, idx) => {
                            return (
                                <Input.Checkbox label={host.Name} onChange={selectHost(host.ID)} defaultValue={selectedHosts[host.ID] > 0} key={idx} />
                            );
                        })
                    }
                </Card.Body>
            </Card.Card>
        </Form>
    );
};
