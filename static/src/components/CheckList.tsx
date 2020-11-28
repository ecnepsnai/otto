import * as React from 'react';
import { Checkbox } from './Form';
import { Loading } from './Loading';
import { Group } from '../types/Group';
import { Script } from '../types/Script';
import { Host } from '../types/Host';
import { Nothing } from './Nothing';

interface CheckListProps {
    selectedKeys: string[],
    keys: string[],
    labels: string[],
    onChange: (selected: string[]) => (void),
}
interface CheckListState {
    selected: {[id: string]: boolean},
}
class CheckList extends React.Component<CheckListProps, CheckListState> {
    constructor(props: CheckListProps) {
        super(props);
        const checked: {[id: string]: boolean} = {};
        (this.props.selectedKeys || []).forEach(key => {
            checked[key] = true;
        });
        this.state = {
            selected: checked,
        };
    }

    private changeKey = (key: string) => {
        return (checked: boolean) => {
            this.setState(state => {
                state.selected[key] = checked;
            }, () => {
                const keys: string[] = [];
                const selected = this.state.selected;
                Object.keys(selected).forEach(key => {
                    if (selected[key]) {
                        keys.push(key);
                    }
                });
                this.props.onChange(keys);
            });
        };
    };

    render(): JSX.Element {
        if (!this.props.keys || this.props.keys.length == 0) { return (<div><Nothing /></div>); }

        return (
            <div>
                {
                    this.props.keys.map((key, idx) => {
                        return (
                            <Checkbox
                                label={this.props.labels[idx]}
                                onChange={this.changeKey(key)}
                                defaultValue={this.state.selected[key]}
                                key={idx} />
                        );
                    })
                }
            </div>
        );
    }
}

export interface GroupCheckListProps {
    selectedGroups: string[],
    onChange: (ids: string[]) => (void),
}
interface GroupCheckListState {
    loading: boolean;
    groups?: Group[];
}
export class GroupCheckList extends React.Component<GroupCheckListProps, GroupCheckListState> {
    constructor(props: GroupCheckListProps) {
        super(props);
        this.state = {
            loading: true,
        };
    }

    componentDidMount(): void {
        Group.List().then(groups => {
            this.setState({ loading: false, groups: groups });
        });
    }

    render(): JSX.Element {
        if (this.state.loading) { return (<Loading />); }

        const keys = this.state.groups.map(group => { return group.ID; });
        const labels = this.state.groups.map(group => { return group.Name; });

        return (
            <CheckList
                selectedKeys={this.props.selectedGroups}
                keys={keys}
                labels={labels}
                onChange={this.props.onChange}/>
        );
    }
}

export interface ScriptCheckListProps {
    selectedScripts: string[],
    onChange: (ids: string[]) => (void),
}
interface ScriptCheckListState {
    loading: boolean;
    scripts?: Script[];
}
export class ScriptCheckList extends React.Component<ScriptCheckListProps, ScriptCheckListState> {
    constructor(props: ScriptCheckListProps) {
        super(props);
        this.state = {
            loading: true,
        };
    }

    componentDidMount(): void {
        Script.List().then(scripts => {
            this.setState({ loading: false, scripts: scripts });
        });
    }

    render(): JSX.Element {
        if (this.state.loading) { return (<Loading />); }

        const keys = this.state.scripts.map(script => { return script.ID; });
        const labels = this.state.scripts.map(script => { return script.Name; });

        return (
            <CheckList
                selectedKeys={this.props.selectedScripts}
                keys={keys}
                labels={labels}
                onChange={this.props.onChange}/>
        );
    }
}

export interface HostCheckListProps {
    selectedHosts: string[],
    onChange: (ids: string[]) => (void),
}
interface HostCheckListState {
    loading: boolean;
    hosts?: Host[];
}
export class HostCheckList extends React.Component<HostCheckListProps, HostCheckListState> {
    constructor(props: HostCheckListProps) {
        super(props);
        this.state = {
            loading: true,
        };
    }

    componentDidMount(): void {
        Host.List().then(hosts => {
            this.setState({ loading: false, hosts: hosts });
        });
    }

    render(): JSX.Element {
        if (this.state.loading) { return (<Loading />); }

        const keys = this.state.hosts.map(host => { return host.ID; });
        const labels = this.state.hosts.map(host => { return host.Name; });

        return (
            <CheckList
                selectedKeys={this.props.selectedHosts}
                keys={keys}
                labels={labels}
                onChange={this.props.onChange}/>
        );
    }
}
