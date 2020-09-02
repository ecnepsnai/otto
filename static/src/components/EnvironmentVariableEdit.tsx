import * as React from 'react';
import { Table } from './Table';
import { CreateButton } from './Button';
import { Icon } from './Icon';
import { Style } from './Style';
import { Dropdown, MenuItem } from './Menu';
import { Modal, ModalButton, GlobalModalFrame } from './Modal';
import { Form, Input, Textarea, Checkbox } from './Form';
import { Variable } from '../types/Variable';

export interface EnvironmentVariableEditProps {
    variables: Variable[];
    onChange: (variables: Variable[]) => (void);
}
export class EnvironmentVariableEdit extends React.Component<EnvironmentVariableEditProps, {}> {
    private addNewVariable = (variable: Variable) => {
        const vars = this.props.variables;
        vars.push(variable);
        this.props.onChange(vars);
    }

    private createClick = () => {
        GlobalModalFrame.showModal(<EnvironmentVariableEditModal default={{}} onSave={this.addNewVariable} />);
    }

    private replaceVariable = (idx: number) => {
        return (varible: Variable) => {
            const vars = this.props.variables;
            vars[idx] = varible;
            this.props.onChange(vars);
        };
    }

    private editVar = (variable: Variable, idx: number): () => (void) => {
        return () => {
            GlobalModalFrame.showModal(<EnvironmentVariableEditModal default={variable} onSave={this.replaceVariable(idx)} />);
        };
    }

    private deleteVar = (idx: number): () => (void) => {
        return () => {
            Modal.delete('Delete Variable', 'Are you sure you want to delete this variable?').then(confirmed => {
                if (!confirmed) { return; }
                const vars = this.props.variables;
                vars.splice(idx, 1);
                this.props.onChange(vars);
            });
        };
    }

    render(): JSX.Element {
        return (
            <React.Fragment>
                <CreateButton onClick={this.createClick} />
                <Table.Table>
                    <Table.Head>
                        <Table.Column>Key</Table.Column>
                        <Table.Column>Value</Table.Column>
                        <Table.MenuColumn />
                    </Table.Head>
                    <Table.Body>
                        {
                            this.props.variables.map((variable, idx) => {
                                return (
                                    <EnvironmentVariableEditListItem
                                        variable={variable}
                                        key={idx}
                                        requestEdit={this.editVar(variable, idx)}
                                        requestDelete={this.deleteVar(idx)}
                                        />
                                );
                            })
                        }
                    </Table.Body>
                </Table.Table>
            </React.Fragment>
        );
    }
}

interface EnvironmentVariableEditListItemProps {
    variable: Variable;
    requestEdit: () => (void);
    requestDelete: () => (void);
}
class EnvironmentVariableEditListItem extends React.Component<EnvironmentVariableEditListItemProps, {}> {
    render(): JSX.Element {
        const dropdownLabel = <Icon.Bars />;
        const buttonProps = {
            color: Style.Palette.Secondary,
            outline: true,
            size: Style.Size.XS,
        };

        const content = this.props.variable.Secret ? '******' : this.props.variable.Value;

        return (
            <Table.Row>
                <td>
                    { this.props.variable.Key }
                </td>
                <td>
                    <code>{content}</code>
                </td>
                <td>
                    <Dropdown label={dropdownLabel} button={buttonProps}>
                        <MenuItem label="Edit" icon={<Icon.Edit />} onClick={this.props.requestEdit}/>
                        <MenuItem label="Delete" icon={<Icon.Delete />} onClick={this.props.requestDelete}/>
                    </Dropdown>
                </td>
            </Table.Row>
        );
    }
}

interface EnvironmentVariableEditModalProps {
    default: Variable;
    onSave: (variable: Variable) => (void);
}
interface EnvironmentVariableEditModalState {
    key: string;
    value: string;
    secret: boolean;
}
class EnvironmentVariableEditModal extends React.Component<EnvironmentVariableEditModalProps, EnvironmentVariableEditModalState> {
    constructor(props: EnvironmentVariableEditModalProps) {
        super(props);
        this.state = {
            key: props.default.Key,
            value: props.default.Value,
            secret: props.default.Secret,
        };
    }

    private changeKey = (key: string) => {
        this.setState({ key: key });
    }

    private changeValue = (value: string) => {
        this.setState({ value: value });
    }

    private changeSecret = (secret: boolean) => {
        this.setState({ secret: secret });
    }

    render(): JSX.Element {
        const title = this.props.default.Key != '' ? 'Edit Variable' : 'New Variable';
        const buttons: ModalButton[] = [
            {
                label: 'Discard',
                color: Style.Palette.Secondary,
            },
            {
                label: 'Add',
                onClick: () => {
                    this.props.onSave({
                        Key: this.state.key,
                        Value: this.state.value,
                        Secret: this.state.secret,
                    });
                },
            }
        ];
        return (
            <Modal title={title} buttons={buttons} static>
                <Form>
                    <Input
                        label="Key"
                        type="text"
                        defaultValue={this.state.key}
                        onChange={this.changeKey} />
                    <Textarea
                        label="Value"
                        defaultValue={this.state.value}
                        onChange={this.changeValue}
                        fixedWidth />
                    <Checkbox
                        label="Secret"
                        defaultValue={this.state.secret}
                        onChange={this.changeSecret}
                        helpText="If checked then the value of this variable is obscured" />
                </Form>
            </Modal>
        );
    }
}