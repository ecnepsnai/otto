import * as React from 'react';
import { AddButton, Button } from '../../../components/Button';
import { Card } from '../../../components/Card';
import { Input } from '../../../components/input/Input';
import { Form } from '../../../components/Form';
import { Icon } from '../../../components/Icon';
import { PageLoading } from '../../../components/Loading';
import { Dropdown, Menu } from '../../../components/Menu';
import { GlobalModalFrame, ModalForm } from '../../../components/Modal';
import { Notification } from '../../../components/Notification';
import { Page } from '../../../components/Page';
import { Table } from '../../../components/Table';
import { Rand } from '../../../services/Rand';
import { StateManager } from '../../../services/StateManager';
import { Group, GroupType } from '../../../types/Group';
import { Options } from '../../../types/Options';
import { RegisterRuleType, RegisterRuleClauseType, RegisterRule } from '../../../types/RegisterRule';
import { Formatter } from '../../../services/Formatter';
import { Style } from '../../../components/Style';
import '../../../../css/registerrules.scss';

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
            <Input.Password
                label="Register PSK"
                helpText="Clients that wish to register with this server must specify this PSK to authenticate"
                defaultValue={options.PSK}
                onChange={changePSK}
                required />
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
                <td>{rule.Name}</td>
                <td>{ Formatter.ValueOrNothing(rule.Clauses.length) }</td>
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
                    <Table.Column>Name</Table.Column>
                    <Table.Column>Clauses</Table.Column>
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
        Name: '',
        Clauses: [{
            Property: 'hostname',
            Pattern: '',
        }],
        GroupID: props.groups[0].ID,
    });

    const changeName = (Name: string) => {
        setRule(rule => {
            rule.Name = Name;
            return {...rule};
        });
    };

    const changeClauses = (Clauses: RegisterRuleClauseType[]) => {
        setRule(rule => {
            rule.Clauses = Clauses;
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
    return (
        <ModalForm title={title} onSubmit={onSubmit}>
            <Input.Text
                label="Name"
                type="text"
                defaultValue={rule.Name}
                onChange={changeName}
                required />
            <RuleClauseListEdit defaultValue={rule.Clauses} onChange={changeClauses} />
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

interface RuleClauseListEditProps {
    defaultValue: RegisterRuleClauseType[];
    onChange: (clauses: RegisterRuleClauseType[]) => (void);
}
const RuleClauseListEdit: React.FC<RuleClauseListEditProps> = (props: RuleClauseListEditProps) => {
    const [clauses, setClauses] = React.useState(props.defaultValue);

    React.useEffect(() => {
        props.onChange(clauses);
    }, [clauses]);

    const changeClause = (idx: number) => {
        return (Clause: RegisterRuleClauseType) => {
            setClauses(clauses => {
                clauses[idx] = Clause;
                return [...clauses];
            });
        };
    };

    const removeButtonDisabled = () => {
        return clauses.length <= 1;
    };

    const addClauseClick = () => {
        setClauses(clauses => {
            return [...clauses, {
                Property: 'hostname',
                Pattern: ''
            }];
        });
    };

    const removeClauseClick = () => {
        setClauses(clauses => {
            clauses.splice(clauses.length-1, 1);
            return [...clauses];
        });
    };

    return (
        <React.Fragment>
            <strong>Clauses (All must match)</strong>
            {
                clauses.map((clause, idx) => {
                    return (<RuleClauseEdit key={idx} defaultValue={clause} onChange={changeClause(idx)} />);
                })
            }
            <Button color={Style.Palette.Secondary} size={Style.Size.XS} outline onClick={removeClauseClick} disabled={removeButtonDisabled()}>-</Button>
            <Button color={Style.Palette.Secondary} size={Style.Size.XS} outline onClick={addClauseClick}>+</Button>
        </React.Fragment>
    );
};

interface RuleClauseEditProps {
    defaultValue: RegisterRuleClauseType;
    onChange: (clauses: RegisterRuleClauseType) => (void);
}
const RuleClauseEdit: React.FC<RuleClauseEditProps> = (props: RuleClauseEditProps) => {
    const [clause, setClause] = React.useState(props.defaultValue);

    React.useEffect(() => {
        props.onChange(clause);
    }, [clause]);

    const changeProperty = (Property: string) => {
        setClause(clause => {
            clause.Property = Property;
            return {...clause};
        });
    };

    const changePattern = (Pattern: string) => {
        setClause(clause => {
            clause.Pattern = Pattern;
            return {...clause};
        });
    };

    const state = StateManager.Current();
    const properties = state.Enums['RegisterRuleProperty'];
    const radioChoices = properties.map(property => {
        return {
            value: property['value'],
            label: property['description'],
        };
    });

    return (
        <div className="horizontal-inputs">
            <Input.Select
                label="Property"
                defaultValue={clause.Property}
                onChange={changeProperty}>
                {
                    radioChoices.map((choice, idx) => {
                        return (<option key={idx} value={choice.value}>{choice.value}</option>);
                    })
                }
            </Input.Select>
            <Input.Text
                label="Regex Pattern"
                type="text"
                placeholder="Regular Expression"
                defaultValue={clause.Pattern}
                onChange={changePattern}
                required />
        </div>
    );
};
