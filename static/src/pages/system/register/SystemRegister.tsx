import * as React from 'react';
import { AddButton } from '../../../components/Button';
import { Card } from '../../../components/Card';
import { Input } from '../../../components/input/Input';
import { Form } from '../../../components/Form';
import { Icon } from '../../../components/Icon';
import { PageLoading } from '../../../components/Loading';
import { Dropdown, Menu } from '../../../components/Menu';
import { GlobalModalFrame, ModalForm } from '../../../components/Modal';
import { Notification } from '../../../components/Notification';
import { Page } from '../../../components/Page';
import { RandomPSK } from '../../../components/RandomPSK';
import { Table } from '../../../components/Table';
import { Rand } from '../../../services/Rand';
import { StateManager } from '../../../services/StateManager';
import { Group, GroupType } from '../../../types/Group';
import { Options } from '../../../types/Options';
import { RegisterRuleType, RegisterRule } from '../../../types/RegisterRule';

export const SystemRegister: React.FC = () => {
    const [loading, setLoading] = React.useState<boolean>(true);
    const [rules, setRules] = React.useState<RegisterRuleType[]>();
    const [groups, setGroups] = React.useState<GroupType[]>();
    const [options, setOptions] = React.useState<Options.Register>();

    React.useEffect(() => {
        loadData();
    }, []);

    const loadRules = () => {
        return RegisterRule.List().then(rules => {
            setRules(rules);
        });
    };

    const loadGroups = () => {
        return Group.List().then(groups => {
            setGroups(groups);
        });
    };

    const loadOptions = () => {
        return Options.Options.Get().then(o => {
            setOptions(o.Register);
        });
    };

    const loadData = () => {
        Promise.all([loadRules(), loadGroups(), loadOptions()]).then(() => {
            setLoading(false);
        });
    };

    const onSubmit = () => {
        return Options.Options.Get().then(o => {
            o.Register = options;
            Options.Options.Save(o).then(() => {
                Notification.success('Changes Saved');
            });
        });
    };

    const changeEnabled = (Enabled: boolean) => {
        setOptions(options => {
            options.Enabled = Enabled;
            return {...options};
        });
    };

    const changePSK = (PSK: string) => {
        setOptions(options => {
            options.PSK = PSK;
            return {...options};
        });
    };

    const addRule = (rule: RegisterRuleType) => {
        RegisterRule.New(rule).then(() => {
            Notification.success('Rule Added');
            loadRules().then(() => {
                setLoading(false);
            });
        });
    };

    const modifyRule = (id: string, rule: RegisterRuleType) => {
        RegisterRule.Save(id, rule).then(() => {
            Notification.success('Rule Modified');
            loadRules().then(() => {
                setLoading(false);
            });
        });
    };

    const deleteRule = (rule: RegisterRuleType) => {
        RegisterRule.DeleteModal(rule).then(confirmed => {
            if (!confirmed) {
                return;
            }

            loadRules().then(() => {
                setLoading(false);
            });
        });
    };

    const changeDefaultGroupID = (DefaultGroupID: string) => {
        setOptions(options => {
            options.DefaultGroupID = DefaultGroupID;
            return {...options};
        });
    };

    const enabledContent = () => {
        if (!options.Enabled) {
            return null;
        }

        return (<React.Fragment>
            <Input.Text
                type="password"
                label="Register PSK"
                helpText="Clients that wish to register with this server must specify this PSK to authenticate"
                defaultValue={options.PSK}
                onChange={changePSK}
                required />
            <RandomPSK newPSK={changePSK} />
            <Card.Card className="mb-2">
                <Card.Header>
                    Rules
                </Card.Header>
                <Card.Body>
                    <RegisterRules rules={rules} onAdd={addRule} onChange={modifyRule} onDelete={deleteRule} groups={groups}/>
                </Card.Body>
            </Card.Card>
            <Input.Select
                label="Default Group"
                helpText="If none of the above rules match the client will be added to this group"
                defaultValue={options.DefaultGroupID}
                onChange={changeDefaultGroupID}>
                { groups.map((group, idx) => {
                    return ( <option key={idx} value={group.ID}>{group.Name}</option> );
                }) }
            </Input.Select>
        </React.Fragment>);
    };

    if (loading) {
        return (<PageLoading />);
    }

    return (
        <Page title="Host Registration">
            <Form showSaveButton={true} onSubmit={onSubmit}>
                <Input.Checkbox
                    label="Allow Hosts to Register Themselves"
                    helpText="If checked hosts can automatically register themselves with this Otto server"
                    defaultValue={options.Enabled}
                    onChange={changeEnabled} />
                { enabledContent() }
            </Form>
        </Page>
    );
};

interface RegisterRulesProps {
    rules: RegisterRuleType[];
    onAdd: (rule: RegisterRuleType) => (void);
    onChange: (id: string, rule: RegisterRuleType) => (void);
    onDelete: (rule: RegisterRuleType) => (void);
    groups: GroupType[];
}
export const RegisterRules: React.FC<RegisterRulesProps> = (props: RegisterRulesProps) => {
    const createNew = () => {
        GlobalModalFrame.showModal(<RuleModal onSave={props.onAdd} groups={props.groups}/>);
    };

    const modifyRule = (rule: RegisterRuleType) => {
        return (params: RegisterRuleType) => {
            props.onChange(rule.ID, params);
        };
    };

    const deleteRuleMenuClick = (rule: RegisterRuleType) => {
        return () => {
            props.onDelete(rule);
        };
    };

    const editRuleMenuClick = (rule: RegisterRuleType) => {
        return () => {
            GlobalModalFrame.showModal(<RuleModal defaultValue={rule} onSave={modifyRule(rule)} groups={props.groups}/>);
        };
    };

    const ruleRow = (rule: RegisterRuleType) => {
        let groupName = '';
        props.groups.forEach(group => {
            if (group.ID === rule.GroupID) {
                groupName = group.Name;
            }
        });

        return (
            <Table.Row key={Rand.ID()}>
                <td>{rule.Property}</td>
                <td>{rule.Pattern}</td>
                <td>{groupName}</td>
                <td>
                    <Dropdown label={<Icon.Bars />}>
                        <Menu.Item label="Edit" icon={<Icon.Edit />} onClick={editRuleMenuClick(rule)}/>
                        <Menu.Divider />
                        <Menu.Item label="Delete" icon={<Icon.Delete />} onClick={deleteRuleMenuClick(rule)}/>
                    </Dropdown>
                </td>
            </Table.Row>
        );
    };

    return (
        <div>
            <AddButton onClick={createNew} />
            <Table.Table>
                <Table.Head>
                    <Table.Column>Property</Table.Column>
                    <Table.Column>Matches</Table.Column>
                    <Table.Column>Add To Group</Table.Column>
                    <Table.MenuColumn />
                </Table.Head>
                <Table.Body>
                    { props.rules.map(rule => {
                        return ruleRow(rule);
                    }) }
                </Table.Body>
            </Table.Table>
        </div>
    );
};

interface RuleModalProps {
    defaultValue?: RegisterRuleType;
    onSave: (rule: RegisterRuleType) => (void);
    groups: GroupType[];
}
const RuleModal: React.FC<RuleModalProps> = (props: RuleModalProps) => {
    const [rule, setRule] = React.useState<RegisterRuleType>(props.defaultValue || {
        Property: 'hostname',
        Pattern: '',
        GroupID: props.groups[0].ID,
    });

    const changeProperty = (Property: string) => {
        setRule(rule => {
            rule.Property = Property;
            return {...rule};
        });
    };

    const changePattern = (Pattern: string) => {
        setRule(rule => {
            rule.Pattern = Pattern;
            return {...rule};
        });
    };

    const changeGroupID = (GroupID: string) => {
        setRule(rule => {
            rule.GroupID = GroupID;
            return {...rule};
        });
    };

    const onSubmit = (): Promise<void> => {
        return new Promise(resolve => {
            props.onSave(rule);
            resolve();
        });
    };

    const title = props.defaultValue ? 'Edit Rule' : 'New Rule';

    const state = StateManager.Current();
    const properties = state.Enums['RegisterRuleProperty'];
    const radioChoices = properties.map(property => {
        return {
            value: property['value'],
            label: property['description'],
        };
    });

    return (
        <ModalForm title={title} onSubmit={onSubmit}>
            <Input.Radio
                label="Property"
                choices={radioChoices}
                defaultValue={rule.Property}
                onChange={changeProperty} />
            <Input.Text
                label="Regex Pattern"
                type="text"
                placeholder="Regular Expression"
                defaultValue={rule.Pattern}
                onChange={changePattern}
                required />
            <Input.Select
                label="Add To Group"
                defaultValue={rule.GroupID}
                onChange={changeGroupID}>
                { props.groups.map((group, idx) => {
                    return ( <option key={idx} value={group.ID}>{group.Name}</option> );
                }) }
            </Input.Select>
        </ModalForm>
    );
};
