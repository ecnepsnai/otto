import * as React from 'react';
import { Column, Table } from './Table';
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

    const replaceVariable = (variable: Variable) => {
        const vars = props.variables;
        let idx = -1;
        for (let i = 0; i < vars.length; i++) {
            if (vars[i].Key === variable.Key) {
                idx = i;
                break;
            }
        }
        vars[idx] = variable;
        props.onChange(vars);
    };

    const didEditVariable = (variable: Variable): () => (void) => {
        return () => {
            GlobalModalFrame.showModal(<EnvironmentVariableEditModal default={variable} onSave={replaceVariable} />);
        };
    };

    const didDeleteVariable = (variable: Variable): () => (void) => {
        return () => {
            Modal.delete('Delete Variable', 'Are you sure you want to delete this variable?').then(confirmed => {
                if (!confirmed) {
                    return;
                }
                const vars = props.variables;
                let idx = -1;
                for (let i = 0; i < vars.length; i++) {
                    if (vars[i].Key === variable.Key) {
                        idx = i;
                        break;
                    }
                }

                vars.splice(idx, 1);
                props.onChange(vars);
            });
        };
    };

    const tableCols: Column[] = [
        {
            title: 'Key',
            value: (v: Variable) => {
                return (<span>{v.Key}</span>);
            },
            sort: 'Key'
        },
        {
            title: 'Value',
            value: (v: Variable) => {
                if (v.Secret) {
                    return (<span>******</span>);
                }
                return (<code>{v.Value}</code>);
            },
            sort: 'Value'
        }
    ];

    return (
        <React.Fragment>
            <AddButton onClick={createClick} />
            <Table columns={tableCols} data={props.variables} contextMenu={(a: Variable) => VariableTableContextMenu(a, didEditVariable(a), didDeleteVariable(a))} defaultSort={{ ColumnIdx: 0, Ascending: true }} />
        </React.Fragment>
    );
};

const VariableTableContextMenu = (variable: Variable, didEditVariable: () => void, didDeleteVariable: () => void): (ContextMenuItem | 'separator')[] => {
    return [
        {
            title: 'Edit',
            icon: (<Icon.Edit />),
            onClick: () => {
                didEditVariable();
            }
        },
        'separator',
        {
            title: 'Delete',
            icon: (<Icon.Delete />),
            onClick: () => {
                didDeleteVariable();
            }
        },
    ];
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
