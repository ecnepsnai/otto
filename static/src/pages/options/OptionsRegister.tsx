import * as React from 'react';
import { Card } from '../../components/Card';
import { Icon } from '../../components/Icon';
import { Checkbox, Input, Select, RadioChoice, Radio } from '../../components/Form';
import { Options } from '../../types/Options';
import { CreateButton } from '../../components/Button';
import { Modal, GlobalModalFrame, ModalForm } from '../../components/Modal';
import { Loading } from '../../components/Loading';
import { Group } from '../../types/Group';
import { Style } from '../../components/Style';
import { Table } from '../../components/Table';
import { Rand } from '../../services/Rand';
import { Dropdown, MenuItem } from '../../components/Menu';

export interface OptionsRegisterProps {
    defaultValue: Options.Register;
    onUpdate: (value: Options.Register) => (void);
}
interface OptionsRegisterState {
    loading: boolean;
    value: Options.Register;
    groups?: Group[];
}
export class OptionsRegister extends React.Component<OptionsRegisterProps, OptionsRegisterState> {
    constructor(props: OptionsRegisterProps) {
        super(props);
        this.state = {
            loading: true,
            value: props.defaultValue,
        };
    }

    private loadGroups = () => {
        Group.List().then(groups => {
            this.setState(state => {
                const value = state.value;
                if (!state.value.DefaultGroupID) {
                    value.DefaultGroupID = groups[0].ID;
                }
                return { value: value, loading: false, groups: groups };
            });
        });
    }

    componentDidMount(): void {
        this.loadGroups();
    }

    private changeEnabled = (Enabled: boolean) => {
        this.setState(state => {
            const options = state.value;
            options.Enabled = Enabled;
            return {
                value: options,
            };
        }, () => {
            this.props.onUpdate(this.state.value);
        });
    }

    private changePSK = (PSK: string) => {
        this.setState(state => {
            const options = state.value;
            options.PSK = PSK;
            return {
                value: options,
            };
        }, () => {
            this.props.onUpdate(this.state.value);
        });
    }

    private changeRules = (Rules: Options.RegisterRule[]) => {
        this.setState(state => {
            const options = state.value;
            options.Rules = Rules;
            return {
                value: options,
            };
        }, () => {
            this.props.onUpdate(this.state.value);
        });
    }

    private changeDefaultGroupID = (DefaultGroupID: string) => {
        this.setState(state => {
            const options = state.value;
            options.DefaultGroupID = DefaultGroupID;
            return {
                value: options,
            };
        }, () => {
            this.props.onUpdate(this.state.value);
        });
    }

    private enabledContent = () => {
        if (!this.state.value.Enabled) { return null; }

        return (
            <React.Fragment>
                <Input
                    type="password"
                    label="Register PSK"
                    helpText="Clients that wish to register with this server must specify this PSK to authenticate"
                    defaultValue={this.state.value.PSK}
                    onChange={this.changePSK}
                    required />
                <label className="form-label">Rules</label>
                <RegisterRules defaultValue={this.state.value.Rules} onChange={this.changeRules} groups={this.state.groups}/>
                <Select
                    label="Default Group"
                    helpText="If none of the above rules match the client will be added to this group"
                    defaultValue={this.state.value.DefaultGroupID}
                    onChange={this.changeDefaultGroupID}>
                        { this.state.groups.map((group, idx) => {
                            return ( <option key={idx} value={group.ID}>{group.Name}</option> );
                        }) }
                    </Select>
            </React.Fragment>
        );
    }

    private content = () => {
        if (this.state.loading) { return (<Loading />); }

        return (
            <React.Fragment>
                <Checkbox
                    label="Allow Hosts to Register Themselves"
                    helpText="If checked hosts can automatically register themselves with this Otto server"
                    defaultValue={this.state.value.Enabled}
                    onChange={this.changeEnabled} />
                { this.enabledContent() }
            </React.Fragment>
        );
    }

    render(): JSX.Element {
        return (
            <Card.Card>
                <Card.Header>
                    <Icon.Label icon={<Icon.Magic />} label="Register" />
                </Card.Header>
                <Card.Body>
                    { this.content() }
                </Card.Body>
            </Card.Card>
        );
    }
}

interface RegisterRulesProps {
    defaultValue: Options.RegisterRule[];
    onChange: (rules: Options.RegisterRule[]) => (void);
    groups: Group[];
}
interface RegisterRulesState {
    value: Options.RegisterRule[];
}
class RegisterRules extends React.Component<RegisterRulesProps, RegisterRulesState> {
    constructor(props: RegisterRulesProps) {
        super(props);
        this.state = {
            value: props.defaultValue,
        };
    }

    private addRule = (rule: Options.RegisterRule) => {
        this.setState(state => {
            const rules = state.value;
            rules.push(rule);
            return { value: rules };
        }, () => {
            this.props.onChange(this.state.value);
        });
    }

    private updateRule = (idx: number, ) => {
        return (rule: Options.RegisterRule) => {
            this.setState(state => {
                const rules = state.value;
                rules[idx] = rule;
                return { value: rules };
            }, () => {
                this.props.onChange(this.state.value);
            });
        };
    }

    private deleteRule = (idx: number) => {
        this.setState(state => {
            const rules = state.value;
            rules.splice(idx, 1);
            return { value: rules };
        }, () => {
            this.props.onChange(this.state.value);
        });
    }

    private createNew = () => {
        GlobalModalFrame.showModal(<RuleModal onSave={this.addRule} groups={this.props.groups}/>);
    }

    private deleteRuleMenuClick = (ruleIdx: number) => {
        return () => {
            Modal.delete('Delete Rule', 'Are you sure you want to delete this rule?').then(confirmed => {
                if (!confirmed) {
                    return;
                }

                this.deleteRule(ruleIdx);
            });
        };
    }

    private editRuleMenuClick = (ruleIdx: number) => {
        return () => {
            GlobalModalFrame.showModal(<RuleModal defaultValue={this.state.value[ruleIdx]} onSave={this.updateRule(ruleIdx)} groups={this.props.groups}/>);
        };
    }

    private ruleRow = (rule: Options.RegisterRule, idx: number) => {
        const prop = rule.Uname ? 'Uname' : 'Hostname';
        const propValue = rule.Uname || rule.Hostname;
        let groupName = '';
        this.props.groups.forEach(group => {
            if (group.ID === rule.GroupID) {
                groupName = group.Name;
            }
        });

        const dropdownLabel = <Icon.Bars />;
        const buttonProps = {
            color: Style.Palette.Secondary,
            outline: true,
            size: Style.Size.XS,
        };

        return (
            <Table.Row key={Rand.ID()}>
                <td>{prop}</td>
                <td>{propValue}</td>
                <td>{groupName}</td>
                <td>
                    <Dropdown label={dropdownLabel} button={buttonProps}>
                        <MenuItem label="Edit" icon={<Icon.Edit />} onClick={this.editRuleMenuClick(idx)}/>
                        <MenuItem label="Delete" icon={<Icon.Delete />} onClick={this.deleteRuleMenuClick(idx)}/>
                    </Dropdown>
                </td>
            </Table.Row>
        );
    }

    render(): JSX.Element {
        return (
            <div>
                <CreateButton onClick={this.createNew} />
                <Table.Table>
                    <Table.Head>
                        <Table.Column>Property</Table.Column>
                        <Table.Column>Matches</Table.Column>
                        <Table.Column>Add To Group</Table.Column>
                        <Table.MenuColumn />
                    </Table.Head>
                    <Table.Body>
                        { this.state.value.map((rule, idx) => { return this.ruleRow(rule, idx); }) }
                    </Table.Body>
                </Table.Table>
            </div>
        );
    }
}

interface RuleModalProps {
    defaultValue?: Options.RegisterRule;
    onSave: (rule: Options.RegisterRule) => (void);
    groups: Group[];
}
interface RuleModalState {
    propType: string;
    value?: Options.RegisterRule;
}
class RuleModal extends React.Component<RuleModalProps, RuleModalState> {
    constructor(props: RuleModalProps) {
        super(props);
        let propType = 'uname';
        if (props.defaultValue && props.defaultValue.Hostname) {
            propType = 'hostname';
        }
        this.state = {
            propType: propType,
            value: props.defaultValue || {
                GroupID: props.groups[0].ID,
            },
        };
    }

    private changePropType = (propType: string) => {
        this.setState(state => {
            const rule = state.value;
            if (propType === 'uname') {
                rule.Hostname = undefined;
            } else if (propType === 'hostname') {
                rule.Uname = undefined;
            }
            return { propType: propType, value: rule };
        });
    }

    private changeUname = (Uname: string) => {
        this.setState(state => {
            const rule = state.value;
            rule.Uname = Uname;
            return { value: rule };
        });
    }

    private unameInput = () => {
        if (this.state.propType !== 'uname') { return null; }

        return (
            <Input
                label="Uname (Regex)"
                type="text"
                placeholder="Regular Expression"
                defaultValue={this.state.value.Uname}
                onChange={this.changeUname}
                required />
        );
    }

    private changeHostname = (Hostname: string) => {
        this.setState(state => {
            const rule = state.value;
            rule.Hostname = Hostname;
            return { value: rule };
        });
    }

    private hostnameInput = () => {
        if (this.state.propType !== 'hostname') { return null; }

        return (
            <Input
                label="Hostname (Regex)"
                type="text"
                placeholder="Regular Expression"
                defaultValue={this.state.value.Hostname}
                onChange={this.changeHostname}
                required />
        );
    }

    private changeGroupID = (GroupID: string) => {
        this.setState(state => {
            const rule = state.value;
            rule.GroupID = GroupID;
            return { value: rule };
        });
    }

    private onSubmit = () => {
        return new Promise(resolve => {
            this.props.onSave(this.state.value);
            resolve();
        });
    }

    render(): JSX.Element {
        const title = this.props.defaultValue ? 'Edit Rule' : 'New Rule';

        const radioChoices: RadioChoice[] = [
            {
                value: 'uname',
                label: 'Uname',
            },
            {
                value: 'hostname',
                label: 'Hostname',
            }
        ];

        return (
            <ModalForm title={title} onSubmit={this.onSubmit}>
                <Radio
                    label="Property"
                    choices={radioChoices}
                    defaultValue={this.state.propType}
                    onChange={this.changePropType} />
                { this.unameInput() }
                { this.hostnameInput() }
                <Select
                    label="Add To Group"
                    defaultValue={this.state.value.GroupID}
                    onChange={this.changeGroupID}>
                        { this.props.groups.map((group, idx) => {
                            return ( <option key={idx} value={group.ID}>{group.Name}</option> );
                        }) }
                </Select>
            </ModalForm>
        );
    }
}
