import * as React from 'react';
import { Script } from '../../types/Script';
import { Loading } from '../../components/Loading';
import { Form, Checkbox } from '../../components/Form';
import { Card } from '../../components/Card';
import { Rand } from '../../services/Rand';

interface SGroup {
    ID: string;
    Name: string;
}

interface SHost {
    ID: string;
    Name: string;
}

export interface RunSetupProps {
    scriptID: string;
    onSelectedHosts: (hostIDs: string[]) => (void);
}
interface RunSetupState {
    loading: boolean;
    script?: Script;
    groups?: SGroup[];
    hosts?: SHost[];
    selectedGroups?: {[id: string]: number};
    selectedHosts?: {[id: string]: number};
    groupMembers?: {[id: string]: string[]};
}
export class RunSetup extends React.Component<RunSetupProps, RunSetupState> {
    constructor(props: RunSetupProps) {
        super(props);
        this.state = { loading: true };
    }

    private loadData = () => {
        Script.Get(this.props.scriptID).then(script => {
            script.Hosts().then(hosts => {
                const groupMap: {[id: string]: string} = {};
                const groupMembership: {[id: string]: string[]} = {};
                const selectedGroups: {[id: string]: number} = {};
                const selectedHosts: {[id: string]: number} = {};

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

                this.setState({
                    loading: false,
                    script: script,
                    groups: groups,
                    hosts: shosts,
                    selectedGroups: selectedGroups,
                    selectedHosts: selectedHosts,
                    groupMembers: groupMembership,
                });
            });
        });
    }

    componentDidMount(): void {
        this.loadData();
    }

    private selectGroup = (groupID: string) => {
        return (checked: boolean) => {
            this.setState(state => {
                const selectedGroups = state.selectedGroups;
                if (checked) {
                    selectedGroups[groupID]++;
                } else {
                    selectedGroups[groupID]--;
                }
                const selectedHosts = state.selectedHosts;
                state.groupMembers[groupID].forEach(host => {
                    if (checked) {
                        selectedHosts[host]++;
                    } else {
                        selectedHosts[host]--;
                    }
                });
                return {
                    selectedGroups: selectedGroups,
                    selectedHosts: selectedHosts,
                };
            }, () => {
                this.props.onSelectedHosts(this.selectedHostIDs());
            });
        };
    }

    private selectHost = (hostID: string) => {
        return (checked: boolean) => {
            this.setState(state => {
                const selected: {[id: string]: number} = state.selectedHosts;
                if (checked) {
                    selected[hostID]++;
                } else {
                    selected[hostID]--;
                }
                return {
                    selectedHosts: selected,
                };
            }, () => {
                this.props.onSelectedHosts(this.selectedHostIDs());
            });
        };
    }

    private selectedHostIDs = () => {
        const hostIDs: string[] = [];
        Object.keys(this.state.selectedHosts).forEach(hostID => {
            if (this.state.selectedHosts[hostID] > 0) {
                hostIDs.push(hostID);
            }
        });
        return hostIDs;
    }

    render(): JSX.Element {
        if (this.state.loading) { return (<Loading />); }

        return (
            <Form>
                <Card.Card>
                    <Card.Header>Groups</Card.Header>
                    <Card.Body>
                        {
                            this.state.groups.map(group => {
                                return (
                                    <Checkbox label={group.Name} onChange={this.selectGroup(group.ID)} defaultValue={this.state.selectedGroups[group.ID]>0} key={Rand.ID()}/>
                                );
                            })
                        }
                    </Card.Body>
                </Card.Card>
                <Card.Card className="mt-3">
                    <Card.Header>Hosts</Card.Header>
                    <Card.Body>
                        {
                            this.state.hosts.map(host => {
                                return (
                                    <Checkbox label={host.Name} onChange={this.selectHost(host.ID)} defaultValue={this.state.selectedHosts[host.ID]>0} key={Rand.ID()}/>
                                );
                            })
                        }
                    </Card.Body>
                </Card.Card>
            </Form>
        );
    }
}
