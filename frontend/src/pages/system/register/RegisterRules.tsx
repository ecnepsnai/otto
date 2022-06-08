import * as React from 'react';
import { AddButton, Button } from '../../../components/Button';
import { ContextMenuItem } from '../../../components/ContextMenu';
import { Input } from '../../../components/input/Input';
import { Icon } from '../../../components/Icon';
import { GlobalModalFrame, ModalForm } from '../../../components/Modal';
import { Column, Table } from '../../../components/Table';
import { StateManager } from '../../../services/StateManager';
import { GroupType } from '../../../types/Group';
import { RegisterRuleType, RegisterRuleClauseType } from '../../../types/RegisterRule';
import { Formatter } from '../../../services/Formatter';
import { Style } from '../../../components/Style';
import '../../../../css/registerrules.scss';
import { DefaultSort } from '../../../services/Sort';

interface RegisterRulesProps {
    rules: RegisterRuleType[];
    onAdd: (rule: RegisterRuleType) => (void);
    onChange: (id: string, rule: RegisterRuleType) => (void);
    onDelete: (rule: RegisterRuleType) => (void);
    groups: GroupType[];
}
export const RegisterRules: React.FC<RegisterRulesProps> = (props: RegisterRulesProps) => {
    const groupMap = new Map(props.rules.map(r => [r.GroupID, r]));

    const createNew = () => {
        GlobalModalFrame.showModal(<RuleModal onSave={props.onAdd} groups={props.groups} />);
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
            GlobalModalFrame.showModal(<RuleModal defaultValue={rule} onSave={modifyRule(rule)} groups={props.groups} />);
        };
    };

    const tableCols: Column[] = [
        {
            title: 'Name',
            value: 'Name',
            sort: 'Name'
        },
        {
            title: 'Clauses',
            value: (v: RegisterRuleType) => {
                return (<span>{Formatter.ValueOrNothing((v.Clauses || []).length)}</span>);
            },
            sort: (asc: boolean, left: RegisterRuleType, right: RegisterRuleType) => {
                return DefaultSort(asc, (left.Clauses || []).length, (right.Clauses || []).length);
            }
        },
        {
            title: 'Add To Group',
            value: (v: RegisterRuleType) => {
                return (<span>{groupMap.get(v.GroupID).Name}</span>);
            },
            sort: (asc: boolean, left: RegisterRuleType, right: RegisterRuleType) => {
                return DefaultSort(asc, groupMap.get(left.GroupID).Name, groupMap.get(right.GroupID).Name);
            }
        },
    ];

    return (
        <div>
            <AddButton onClick={createNew} />
            <Table columns={tableCols} data={props.rules} contextMenu={(a: RegisterRuleType) => RuleTableContextMenu(a, editRuleMenuClick(a), deleteRuleMenuClick(a))} defaultSort={{ ColumnIdx: 0, Ascending: true }} />
        </div>
    );
};

const RuleTableContextMenu = (rule: RegisterRuleType, onEdit: () => void, onDelete: () => void): (ContextMenuItem | 'separator')[] => {
    return [
        {
            title: 'Edit',
            icon: (<Icon.Edit />),
            onClick: () => {
                onEdit();
            }
        },
        'separator',
        {
            title: 'Delete',
            icon: (<Icon.Delete />),
            onClick: () => {
                onDelete();
            }
        },
    ];
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
            return { ...rule };
        });
    };

    const changeClauses = (Clauses: RegisterRuleClauseType[]) => {
        setRule(rule => {
            rule.Clauses = Clauses;
            return { ...rule };
        });
    };

    const changeGroupID = (GroupID: string) => {
        setRule(rule => {
            rule.GroupID = GroupID;
            return { ...rule };
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
                {props.groups.map((group, idx) => {
                    return (<option key={idx} value={group.ID}>{group.Name}</option>);
                })}
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
            clauses.splice(clauses.length - 1, 1);
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
            return { ...clause };
        });
    };

    const changePattern = (Pattern: string) => {
        setClause(clause => {
            clause.Pattern = Pattern;
            return { ...clause };
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
