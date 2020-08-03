import * as React from 'react';
import { Card } from './Card';
import { Table } from './Table';
import { CreateButton } from './Button';
import { Icon } from './Icon';
import { Style } from './Style';
import { Dropdown, MenuItem } from './Menu';
import { Modal, ModalButton } from './Modal';
import { Form, Input, Textarea } from './Form';

export interface EnvironmentVariableEditProps {
    variables: {[id: string]: string};
    onChange: (variables: {[id: string]: string}) => (void);
}
interface EnvironmentVariableEditState {
    editKey?: string;
}
export class EnvironmentVariableEdit extends React.Component<EnvironmentVariableEditProps, EnvironmentVariableEditState> {
    constructor(props: EnvironmentVariableEditProps) {
        super(props);
        this.state = { };
    }

    private createClick = () => {
        this.setState({ editKey: '' });
    }

    private editVar = (key: string): () => (void) => {
        return () => {
            this.setState({ editKey: key });
        };
    }

    private deleteVar = (key: string): () => (void) => {
        return () => {
            Modal.delete('Delete Variable', 'Are you sure you want to delete this variable?').then(confirmed => {
                if (!confirmed) { return; }

                const vars = this.props.variables;
                delete vars[key];
                this.props.onChange(vars);
            });
        };
    }

    private editDialogDismissed = (key: string, value: string) => {
        this.setState({ editKey: undefined });
        if (key == undefined) {
            return;
        }

        const vars = this.props.variables;
        vars[key] = value;
        this.props.onChange(vars);
    }

    private editDialog = () => {
        if (this.state.editKey == undefined) { return null; }
        const value = this.state.editKey != '' ? this.props.variables[this.state.editKey] : undefined;
        return (
            <EnvironmentVariableEditModal
                defaultKey={this.state.editKey}
                defaultValue={value}
                onDismiss={this.editDialogDismissed} />
        );
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
                            Object.keys(this.props.variables).map((key, idx) => {
                                return (
                                    <EnvironmentVariableEditListItem
                                        envKey={key}
                                        value={this.props.variables[key]}
                                        key={idx}
                                        requestEdit={this.editVar(key)}
                                        requestDelete={this.deleteVar(key)}
                                        />
                                );
                            })
                        }
                    </Table.Body>
                </Table.Table>
                { this.editDialog() }
            </React.Fragment>
        );
    }
}

interface EnvironmentVariableEditListItemProps {
    envKey: string;
    value: string;
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
        return (
            <Table.Row>
                <td>
                    { this.props.envKey }
                </td>
                <td>
                    <code>{ this.props.value }</code>
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
    defaultKey: string;
    defaultValue: string;
    onDismiss: (key: string, value: string) => (void);
}
interface EnvironmentVariableEditModalState {
    key: string;
    value: string;
}
class EnvironmentVariableEditModal extends React.Component<EnvironmentVariableEditModalProps, EnvironmentVariableEditModalState> {
    constructor(props: EnvironmentVariableEditModalProps) {
        super(props);
        this.state = {
            key: props.defaultKey,
            value: props.defaultValue,
        };
    }

    private changeKey = (key: string) => {
        this.setState({ key: key });
    }

    private changeValue = (value: string) => {
        this.setState({ value: value });
    }

    render(): JSX.Element {
        const title = this.props.defaultKey != '' ? 'Edit Variable' : 'New Variable';
        const buttons: ModalButton[] = [
            {
                label: 'Discard',
                color: Style.Palette.Secondary,
                onClick: () => {
                    this.props.onDismiss(undefined, undefined);
                },
            },
            {
                label: 'Add',
                onClick: () => {
                    this.props.onDismiss(this.state.key, this.state.value);
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
                </Form>
            </Modal>
        );
    }
}