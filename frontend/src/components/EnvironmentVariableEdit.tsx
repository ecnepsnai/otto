import * as React from 'react';
import { Table } from './Table';
import { AddButton } from './Button';
import { Icon } from './Icon';
import { Modal, GlobalModalFrame, ModalForm } from './Modal';
import { Input } from './input/Input';
import { ValidationResult } from './Form';
import { Variable } from '../types/Variable';
import { ContextMenuItem } from './ContextMenu';

interface EnvironmentVariableEditProps {
    variables: Variable[];
    onChange: (variables: Variable[]) => (void);
}
export const EnvironmentVariableEdit: React.FC<EnvironmentVariableEditProps> = (props: EnvironmentVariableEditProps) => {
    const addNewVariable = (variable: Variable) => {
        const vars = props.variables;
        vars.push(variable);
        props.onChange(vars);
    };

    const createClick = () => {
        GlobalModalFrame.showModal(<EnvironmentVariableEditModal default={{}} onSave={addNewVariable} />);
    };

    const replaceVariable = (idx: number) => {
        return (variable: Variable) => {
            const vars = props.variables;
            vars[idx] = variable;
            props.onChange(vars);
        };
    };

    const editVar = (variable: Variable, idx: number): () => (void) => {
        return () => {
            GlobalModalFrame.showModal(<EnvironmentVariableEditModal default={variable} onSave={replaceVariable(idx)} />);
        };
    };

    const deleteVar = (idx: number): () => (void) => {
        return () => {
            Modal.delete('Delete Variable', 'Are you sure you want to delete this variable?').then(confirmed => {
                if (!confirmed) {
                    return;
                }
                const vars = props.variables;
                vars.splice(idx, 1);
                props.onChange(vars);
            });
        };
    };

    return (
        <React.Fragment>
            <AddButton onClick={createClick} />
            <Table.Table>
                <Table.Head>
                    <Table.Column>Key</Table.Column>
                    <Table.Column>Value</Table.Column>
                </Table.Head>
                <Table.Body>
                    {
                        (props.variables || []).map((variable, idx) => {
                            return (
                                <EnvironmentVariableEditListItem
                                    variable={variable}
                                    key={idx}
                                    requestEdit={editVar(variable, idx)}
                                    requestDelete={deleteVar(idx)} />
                            );
                        })
                    }
                </Table.Body>
            </Table.Table>
        </React.Fragment>
    );
};

interface EnvironmentVariableEditListItemProps {
    variable: Variable;
    requestEdit: () => (void);
    requestDelete: () => (void);
}
const EnvironmentVariableEditListItem: React.FC<EnvironmentVariableEditListItemProps> = (props: EnvironmentVariableEditListItemProps) => {
    const content = props.variable.Secret ? '******' : props.variable.Value;

    const contextMenu: (ContextMenuItem | 'separator')[] = [
        {
            title: 'Edit',
            icon: (<Icon.Edit />),
            onClick: () => {
                props.requestEdit();
            }
        },
        'separator',
        {
            title: 'Delete',
            icon: (<Icon.Delete />),
            onClick: () => {
                props.requestDelete();
            }
        },
    ];

    return (
        <Table.Row menu={contextMenu}>
            <td>
                {props.variable.Key}
            </td>
            <td>
                <code>{content}</code>
            </td>
        </Table.Row>
    );
};

interface EnvironmentVariableEditModalProps {
    default: Variable;
    onSave: (variable: Variable) => (void);
}
const EnvironmentVariableEditModal: React.FC<EnvironmentVariableEditModalProps> = (props: EnvironmentVariableEditModalProps) => {
    const [key, setKey] = React.useState(props.default.Key);
    const [value, setValue] = React.useState(props.default.Value);
    const [secret, setSecret] = React.useState(props.default.Secret);

    const changeKey = (key: string) => {
        setKey(key);
    };

    const changeValue = (value: string) => {
        setValue(value);
    };

    const changeSecret = (secret: boolean) => {
        setSecret(secret);
    };

    const onSave = (): Promise<void> => {
        return new Promise(resolve => {
            props.onSave({
                Key: key,
                Value: value,
                Secret: secret,
            });
            resolve();
        });
    };

    const validateKey = (value: string): Promise<ValidationResult> => {
        const reserved = [
            'OTTO_SERVER_VERSION',
            'OTTO_SERVER_URL',
            'OTTO_HOST_ADDRESS',
            'OTTO_HOST_PORT',
            'OTTO_HOST_PSK'
        ].includes(value);

        if (reserved) {
            return Promise.resolve({
                valid: false,
                invalidMessage: 'Key is reserved by the Otto system',
            });
        }
        return Promise.resolve({
            valid: true,
        });
    };

    const title = props.default.Key != '' ? 'Edit Variable' : 'New Variable';
    return (
        <ModalForm title={title} onSubmit={onSave}>
            <Input.Text
                label="Key"
                type="text"
                defaultValue={key}
                onChange={changeKey}
                fixedWidth
                validate={validateKey}
                required />
            <Input.Textarea
                label="Value"
                defaultValue={value}
                onChange={changeValue}
                fixedWidth />
            <Input.Checkbox
                label="Secret"
                defaultValue={secret}
                onChange={changeSecret}
                helpText="If checked then the value of this variable is obscured" />
        </ModalForm>
    );
};
