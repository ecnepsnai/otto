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

interface SystemRegisterState {
    loading: boolean;
    rules?: RegisterRuleType[];
    groups?: GroupType[];
    options?: Options.Register;
}
export class SystemRegister extends React.Component<unknown, SystemRegisterState> {
    constructor(props: unknown) {
        super(props);
        this.state = {
            loading: true
        };
    }

    private loadRules = () => {
        return RegisterRule.List().then(rules => {
            this.setState({rules: rules});
        });
    }

    private loadGroups = () => {
        return Group.List().then(groups => {
            this.setState({groups: groups});
        });
    }

    private loadOptions = () => {
        return Options.Options.Get().then(o => {
            this.setState({options: o.Register});
        });
    }

    componentDidMount(): void {
        Promise.all([this.loadRules(), this.loadGroups(), this.loadOptions()]).then(() => {
            this.setState({loading: false});
        });
    }

    private onSubmit = () => {
        return Options.Options.Get().then(options => {
            options.Register = this.state.options;
            Options.Options.Save(options).then(() => {
                Notification.success('Changes Saved');
            });
        });
    }

    private changeEnabled = (Enabled: boolean) => {
        this.setState(state => {
            const options = state.options;
            options.Enabled = Enabled;
            return { options: options };
        });
    }

    private changePSK = (PSK: string) => {
        this.setState(state => {
            const options = state.options;
            options.PSK = PSK;
            return { options: options };
        });
    }

    private addRule = (rule: RegisterRuleType) => {
        RegisterRule.New(rule).then(() => {
            Notification.success('Rule Added');
            this.loadRules().then(() => {
                this.setState({ loading: false });
            });
        });
    }

    private modifyRule = (id: string, rule: RegisterRuleType) => {
        RegisterRule.Save(id, rule).then(() => {
            Notification.success('Rule Modified');
            this.loadRules().then(() => {
                this.setState({ loading: false });
            });
        });
    }

    private deleteRule = (rule: RegisterRuleType) => {
        RegisterRule.DeleteModal(rule).then(confirmed => {
            if (!confirmed) {
                return;
            }

            this.loadRules().then(() => {
                this.setState({ loading: false });
            });
        });
    }

    private changeDefaultGroupID = (DefaultGroupID: string) => {
        this.setState(state => {
            const options = state.options;
            options.DefaultGroupID = DefaultGroupID;
            return { options: options };
        });
    }

    private enabledContent = () => {
        if (!this.state.options.Enabled) {
            return null;
        }

        return (<React.Fragment>
            <Input.Text
                type="password"
                label="Register PSK"
                helpText="Clients that wish to register with this server must specify this PSK to authenticate"
                defaultValue={this.state.options.PSK}
                onChange={this.changePSK}
                required />
            <RandomPSK newPSK={this.changePSK} />
            <Card.Card className="mb-2">
                <Card.Header>
                    Rules
                </Card.Header>
                <Card.Body>
                    <RegisterRules rules={this.state.rules} onAdd={this.addRule} onChange={this.modifyRule} onDelete={this.deleteRule} groups={this.state.groups}/>
                </Card.Body>
            </Card.Card>
            <Input.Select
                label="Default Group"
                helpText="If none of the above rules match the client will be added to this group"
                defaultValue={this.state.options.DefaultGroupID}
                onChange={this.changeDefaultGroupID}>
                { this.state.groups.map((group, idx) => {
                    return ( <option key={idx} value={group.ID}>{group.Name}</option> );
                }) }
            </Input.Select>
        </React.Fragment>);
    }

    render(): JSX.Element {
        if (this.state.loading) {
            return (<PageLoading />);
        }

        return (
            <Page title="Host Registration">
                <Form showSaveButton={true} onSubmit={this.onSubmit}>
                    <Input.Checkbox
                        label="Allow Hosts to Register Themselves"
                        helpText="If checked hosts can automatically register themselves with this Otto server"
                        defaultValue={this.state.options.Enabled}
                        onChange={this.changeEnabled} />
                    { this.enabledContent() }
                </Form>
            </Page>
        );
    }
}

interface RegisterRulesProps {
    rules: RegisterRuleType[];
    onAdd: (rule: RegisterRuleType) => (void);
    onChange: (id: string, rule: RegisterRuleType) => (void);
    onDelete: (rule: RegisterRuleType) => (void);
    groups: GroupType[];
}
class RegisterRules extends React.Component<RegisterRulesProps, unknown> {
    private createNew = () => {
        GlobalModalFrame.showModal(<RuleModal onSave={this.props.onAdd} groups={this.props.groups}/>);
    }

    private modifyRule = (rule: RegisterRuleType) => {
        return (params: RegisterRuleType) => {
            this.props.onChange(rule.ID, params);
        };
    }

    private deleteRuleMenuClick = (rule: RegisterRuleType) => {
        return () => {
            this.props.onDelete(rule);
        };
    }

    private editRuleMenuClick = (rule: RegisterRuleType) => {
        return () => {
            GlobalModalFrame.showModal(<RuleModal defaultValue={rule} onSave={this.modifyRule(rule)} groups={this.props.groups}/>);
        };
    }

    private ruleRow = (rule: RegisterRuleType) => {
        let groupName = '';
        this.props.groups.forEach(group => {
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
                        <Menu.Item label="Edit" icon={<Icon.Edit />} onClick={this.editRuleMenuClick(rule)}/>
                        <Menu.Divider />
                        <Menu.Item label="Delete" icon={<Icon.Delete />} onClick={this.deleteRuleMenuClick(rule)}/>
                    </Dropdown>
                </td>
            </Table.Row>
        );
    }

    render(): JSX.Element {
        return (
            <div>
                <AddButton onClick={this.createNew} />
                <Table.Table>
                    <Table.Head>
                        <Table.Column>Property</Table.Column>
                        <Table.Column>Matches</Table.Column>
                        <Table.Column>Add To Group</Table.Column>
                        <Table.MenuColumn />
                    </Table.Head>
                    <Table.Body>
                        { this.props.rules.map(rule => {
                            return this.ruleRow(rule);
                        }) }
                    </Table.Body>
                </Table.Table>
            </div>
        );
    }
}

interface RuleModalProps {
    defaultValue?: RegisterRuleType;
    onSave: (rule: RegisterRuleType) => (void);
    groups: GroupType[];
}
interface RuleModalState {
    value: RegisterRuleType;
}
class RuleModal extends React.Component<RuleModalProps, RuleModalState> {
    constructor(props: RuleModalProps) {
        super(props);
        this.state = {
            value: props.defaultValue || {
                Property: 'hostname',
                Pattern: '',
                GroupID: props.groups[0].ID,
            },
        };
    }

    private changePropType = (propType: string) => {
        this.setState(state => {
            const rule = state.value;
            rule.Property = propType;
            return { value: rule };
        });
    }

    private changePattern = (Pattern: string) => {
        this.setState(state => {
            const rule = state.value;
            rule.Pattern = Pattern;
            return { value: rule };
        });
    }

    private changeGroupID = (GroupID: string) => {
        this.setState(state => {
            const rule = state.value;
            rule.GroupID = GroupID;
            return { value: rule };
        });
    }

    private onSubmit = (): Promise<void> => {
        return new Promise(resolve => {
            this.props.onSave(this.state.value);
            resolve();
        });
    }

    render(): JSX.Element {
        const title = this.props.defaultValue ? 'Edit Rule' : 'New Rule';

        const state = StateManager.Current();
        const properties = state.Enums['RegisterRuleProperty'];
        const radioChoices = properties.map(property => {
            return {
                value: property['value'],
                label: property['description'],
            };
        });

        return (
            <ModalForm title={title} onSubmit={this.onSubmit}>
                <Input.Radio
                    label="Property"
                    choices={radioChoices}
                    defaultValue={this.state.value.Property}
                    onChange={this.changePropType} />
                <Input.Text
                    label="Regex Pattern"
                    type="text"
                    placeholder="Regular Expression"
                    defaultValue={this.state.value.Pattern}
                    onChange={this.changePattern}
                    required />
                <Input.Select
                    label="Add To Group"
                    defaultValue={this.state.value.GroupID}
                    onChange={this.changeGroupID}>
                    { this.props.groups.map((group, idx) => {
                        return ( <option key={idx} value={group.ID}>{group.Name}</option> );
                    }) }
                </Input.Select>
            </ModalForm>
        );
    }
}

