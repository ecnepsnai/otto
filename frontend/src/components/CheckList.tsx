import * as React from 'react';
import { Input } from './input/Input';
import { Loading } from './Loading';
import { Group, GroupType } from '../types/Group';
import { Script, ScriptType } from '../types/Script';
import { Host, HostType } from '../types/Host';
import { Nothing } from './Nothing';
import { Icon } from './Icon';
import { Button } from './Button';
import { Style } from './Style';

interface ItemType {
    key: string;
    value: string;
}

interface CheckListProps {
    selectedKeys: string[],
    keys: string[],
    labels: string[],
    onChange: (selected: string[]) => (void),
}
export const CheckList: React.FC<CheckListProps> = (props: CheckListProps) => {
    const initialChecked: { [id: string]: boolean } = {};
    (props.selectedKeys || []).forEach(key => {
        initialChecked[key] = true;
    });
    const [selected, setSelected] = React.useState<{ [id: string]: boolean }>(initialChecked);
    const [Query, setQuery] = React.useState<string>('');
    const [Items, setItems] = React.useState<ItemType[]>([]);
    const [ShowAll, SetShowAll] = React.useState(false);

    React.useEffect(() => {
        const keys: string[] = [];
        Object.keys(selected).forEach(key => {
            if (selected[key]) {
                keys.push(key);
            }
        });
        props.onChange(keys);
    }, [selected]);

    React.useEffect(() => {
        if (!Query) {
            setItems(props.keys.map((k, i) => {
                return {
                    key: k,
                    value: props.labels[i],
                };
            }));
            return;
        }

        const filteredItems: ItemType[] = [];
        props.keys.forEach((k, i) => {
            const value = props.labels[i].toLowerCase();
            if (value.includes(Query.toLowerCase())) {
                filteredItems.push({
                    key: k,
                    value: props.labels[i],
                });
            }
        });
        setItems(filteredItems);
    }, [Query]);

    const changeKey = (key: string) => {
        return (checked: boolean) => {
            setSelected(selected => {
                selected[key] = checked;
                return { ...selected };
            });
        };
    };

    const onSearch = (query: string) => {
        setQuery(query);
    };

    const showAllButton = () => {
        if (Items.length <= 10 || ShowAll) {
            return null;
        }

        return (<span>
            <Button onClick={() => {
                SetShowAll(true);
            }} color={Style.Palette.Primary} outline size={Style.Size.XS}>{'Show all (' + Items.length + ')'}</Button>
        </span>);
    };

    const checkList = () => {
        if (Items.length == 0) {
            return (<Nothing />);
        }

        return (<React.Fragment>
            {
                Items.map((item, idx) => {
                    if (idx > 9 && !ShowAll) {
                        return null;
                    }
                    return (
                        <Input.Checkbox
                            label={item.value}
                            onChange={changeKey(item.key)}
                            defaultValue={selected[item.key]}
                            key={idx}
                            thin />
                    );
                })
            }
            { showAllButton() }
        </React.Fragment>);
    };

    return (
        <div>
            <Input.Text type="search" placeholder="Filter" defaultValue={Query} onChange={onSearch} prepend={<Icon.MagnifyingGlass />}/>
            {checkList()}
        </div>
    );
};

interface GroupCheckListProps {
    selectedGroups: string[],
    onChange: (ids: string[]) => (void),
}
export const GroupCheckList: React.FC<GroupCheckListProps> = (props: GroupCheckListProps) => {
    const [loading, setLoading] = React.useState(true);
    const [groups, setGroups] = React.useState<GroupType[]>();

    const loadGroups = () => {
        Group.List().then(groups => {
            setGroups(groups);
            setLoading(false);
        });
    };

    React.useEffect(() => {
        loadGroups();
    }, []);

    if (loading) {
        return (<Loading />);
    }

    const keys = groups.map(group => {
        return group.ID;
    });
    const labels = groups.map(group => {
        return group.Name;
    });

    return (
        <CheckList
            selectedKeys={props.selectedGroups}
            keys={keys}
            labels={labels}
            onChange={props.onChange} />
    );
};

interface ScriptCheckListProps {
    selectedScripts: string[],
    onChange: (ids: string[]) => (void),
}
export const ScriptCheckList: React.FC<ScriptCheckListProps> = (props: ScriptCheckListProps) => {
    const [loading, setLoading] = React.useState(true);
    const [scripts, setScripts] = React.useState<ScriptType[]>();

    const loadScripts = () => {
        Script.List().then(scripts => {
            setScripts(scripts);
            setLoading(false);
        });
    };

    React.useEffect(() => {
        loadScripts();
    }, []);

    if (loading) {
        return (<Loading />);
    }

    const keys = scripts.map(script => {
        return script.ID;
    });
    const labels = scripts.map(script => {
        return script.Name;
    });

    return (
        <CheckList
            selectedKeys={props.selectedScripts}
            keys={keys}
            labels={labels}
            onChange={props.onChange} />
    );
};

interface HostCheckListProps {
    selectedHosts: string[],
    onChange: (ids: string[]) => (void),
}
export const HostCheckList: React.FC<HostCheckListProps> = (props: HostCheckListProps) => {
    const [loading, setLoading] = React.useState(true);
    const [hosts, setHosts] = React.useState<HostType[]>();

    const loadHosts = () => {
        Host.List().then(hosts => {
            setHosts(hosts);
            setLoading(false);
        });
    };

    React.useEffect(() => {
        loadHosts();
    }, []);

    if (loading) {
        return (<Loading />);
    }

    const keys = hosts.map(host => {
        return host.ID;
    });
    const labels = hosts.map(host => {
        return host.Name;
    });

    return (
        <CheckList
            selectedKeys={props.selectedHosts}
            keys={keys}
            labels={labels}
            onChange={props.onChange} />
    );
};
