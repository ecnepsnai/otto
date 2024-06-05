import * as React from 'react';
import { EnabledBadge } from '../../../components/Badge';
import { AddButton, Button, ConfirmButton } from '../../../components/Button';
import { Input } from '../../../components/input/Input';
import { Icon } from '../../../components/Icon';
import { PageLoading } from '../../../components/Loading';
import { GlobalModalFrame, Modal, ModalForm } from '../../../components/Modal';
import { Page } from '../../../components/Page';
import { Style } from '../../../components/Style';
import { Column, Table } from '../../../components/Table';
import { StateManager } from '../../../services/StateManager';
import { User, UserPermissions, UserType } from '../../../types/User';
import { ContextMenuItem } from '../../../components/ContextMenu';
import { ScriptRunLevel } from '../../../types/gengo_enum';
import { Permissions, UserAction } from '../../../services/Permissions';

export class UserManager {
    public static EditCurrentUser(): Promise<UserType> {
        return new Promise(resolve => {
            const editUser = (user: UserType) => {
                User.Save(user).then(resolve);
            };
            GlobalModalFrame.showModal(<OptionsUsersModal user={StateManager.Current().User} onUpdate={editUser} />);
        });
    }
}

export const SystemUsers: React.FC = () => {
    const [loading, setLoading] = React.useState<boolean>(true);
    const [users, setUsers] = React.useState<UserType[]>([]);

    React.useEffect(() => {
        loadUsers();
    }, []);

    const loadUsers = () => {
        User.List().then(users => {
            setLoading(false);
            setUsers(users);
        });
    };

    const newUser = (user: UserType) => {
        User.New({
            Username: user.Username,
            Password: user.Password,
            MustChangePassword: user.MustChangePassword,
        }).then(() => {
            loadUsers();
        });
    };

    const updateUser = (user: UserType) => {
        User.Save(user).then(() => {
            loadUsers();
        });
    };

    const newUserClick = () => {
        GlobalModalFrame.showModal(<OptionsUsersModal onUpdate={newUser} />);
    };

    const editUserMenuClick = (user: UserType) => {
        return () => {
            GlobalModalFrame.showModal(<OptionsUsersModal user={user} onUpdate={updateUser} />);
        };
    };
    const deleteUserMenuClick = (user: UserType) => {
        return () => {
            Modal.delete('Delete User?', 'Are you sure you want to delete this user? This can not be undone.').then(confirmed => {
                if (!confirmed) {
                    return;
                }

                User.Delete(user).then(() => {
                    loadUsers();
                });
            });
        };
    };

    if (loading) {
        return (<PageLoading />);
    }

    const toolbar = (
        <React.Fragment>
            <AddButton onClick={newUserClick} />
        </React.Fragment>
    );

    const tableCols: Column[] = [
        {
            title: 'Username',
            value: 'Username',
            sort: 'Username'
        },
        {
            title: 'Can Login',
            value: (v: UserType) => {
                return (<EnabledBadge value={v.CanLogIn} />);
            }
        },
    ];

    return (
        <Page title="Users" toolbar={toolbar}>
            <Table columns={tableCols} data={users} contextMenu={(a: UserType) => UserTableContextMenu(a, editUserMenuClick(a), deleteUserMenuClick(a))} defaultSort={{ ColumnIdx: 0, Ascending: true }} />
        </Page>
    );
};

const UserTableContextMenu = (user: UserType, didEditUser: () => void, didDeleteUser: () => void): (ContextMenuItem | 'separator')[] => {
    const contextMenu: (ContextMenuItem | 'separator')[] = [
        {
            title: 'Edit',
            icon: (<Icon.Edit />),
            onClick: () => {
                didEditUser();
            }
        },
    ];
    if (StateManager.Current().User.Username != user.Username) {
        contextMenu.push('separator');
        contextMenu.push({
            title: 'Delete',
            icon: (<Icon.Delete />),
            onClick: () => {
                didDeleteUser();
            }
        });
    }
    return contextMenu;
};

interface OptionsUsersModalProps {
    user?: UserType;
    onUpdate: (user: UserType) => (void);
}
export const OptionsUsersModal: React.FC<OptionsUsersModalProps> = (props: OptionsUsersModalProps) => {
    const [user, setUser] = React.useState<UserType>(props.user || User.Blank());
    const isNew = props.user == undefined;
    const [shouldShowPasswordField, setShouldShowPasswordField] = React.useState(isNew);

    const changeUsername = (Username: string) => {
        setUser(user => {
            user.Username = Username;
            return { ...user };
        });
    };

    const changePermissions = (Permissions: UserPermissions) => {
        setUser(user => {
            user.Permissions = Permissions;
            return { ...user };
        });
    };

    const changePassword = (Password: string) => {
        setUser(user => {
            user.Password = Password;
            return { ...user };
        });
    };

    const changeCanLogIn = (CanLogIn: boolean) => {
        setUser(user => {
            user.CanLogIn = CanLogIn;
            return { ...user };
        });
    };

    const changeMustChangePassword = (MustChangePassword: boolean) => {
        setUser(user => {
            user.MustChangePassword = MustChangePassword;
            return { ...user };
        });
    };

    const showPasswordField = () => {
        setShouldShowPasswordField(true);
    };

    const passwordField = () => {
        if (!shouldShowPasswordField) {
            return (
                <div className="mb-3">
                    <Button color={Style.Palette.Primary} outline onClick={showPasswordField}><Icon.Label icon={<Icon.Edit />} label="Change Password" /></Button>
                </div>
            );
        }

        return (
            <Input.Text
                type="password"
                label="Password"
                onChange={changePassword}
                required />
        );
    };

    const resetAPIKey = () => {
        if (isNew) {
            return null;
        }

        return (<UserAPIKeyEdit user={props.user} />);
    };

    const mustChangePasswordCheckbox = () => {
        if (props.user && StateManager.Current().User.Username == props.user.Username) {
            return null;
        }

        return (
            <Input.Checkbox
                label="Must Change Password"
                defaultValue={user.MustChangePassword}
                onChange={changeMustChangePassword}
                helpText="If checked this user must change their password the next time they log in" />
        );
    };

    const permissionEdit = () => {
        if (!Permissions.UserCan(UserAction.ModifyUsers)) {
            return null;
        }

        return (<UserPermissionsEdit permissions={user.Permissions} onUpdate={changePermissions} />);
    };

    const canLogInCheckbox = () => {
        if (props.user && StateManager.Current().User.Username == props.user.Username) {
            return null;
        }

        return (
            <Input.Checkbox
                label="User Can Log In"
                defaultValue={user.CanLogIn}
                onChange={changeCanLogIn}
                helpText="If unchecked this user cannot log in" />
        );
    };

    const onSubmit = (): Promise<void> => {
        return new Promise(resolve => {
            props.onUpdate(user);
            resolve();
        });
    };

    return (
        <ModalForm title={isNew ? 'New User' : 'Edit User'} onSubmit={onSubmit}>
            <Input.Text
                type="text"
                label="Username"
                defaultValue={user.Username}
                onChange={changeUsername}
                disabled={props.user != undefined}
                required />
            {passwordField()}
            {resetAPIKey()}
            {canLogInCheckbox()}
            {mustChangePasswordCheckbox()}
            {permissionEdit()}
        </ModalForm>
    );
};

interface UserAPIKeyEditProps {
    user: UserType;
}
const UserAPIKeyEdit: React.FC<UserAPIKeyEditProps> = (props: UserAPIKeyEditProps) => {
    const [loading, setLoading] = React.useState(false);
    const [newAPIKey, setNewAPIKey] = React.useState<string>();

    const resetAPIKey = () => {
        setLoading(true);
        User.ResetAPIKey(props.user).then(key => {
            setNewAPIKey(key);
            setLoading(false);
        });
    };

    if (newAPIKey) {
        return (<Input.Text
            type="text"
            label="API Key"
            helpText="This key is only shown here once and cannot be retrieved after closing this dialog."
            defaultValue={newAPIKey}
            onChange={() => { /* */ }}
            disabled />);
    }

    return (<ConfirmButton color={Style.Palette.Warning} size={Style.Size.S} outline onClick={resetAPIKey} disabled={loading}><Icon.Label icon={<Icon.Undo />} label="Reset API Key" /></ConfirmButton>);
};

interface UserPermissionsEditProps {
    permissions: UserPermissions;
    onUpdate: (value: UserPermissions) => void;
}
const UserPermissionsEdit: React.FC<UserPermissionsEditProps> = (props: UserPermissionsEditProps) => {
    const [Permissions, SetPermissions] = React.useState<UserPermissions>(props.permissions);

    React.useEffect(() => {
        if (!Permissions) {
            return;
        }

        props.onUpdate(Permissions);
    }, [Permissions]);

    const scriptRunLevelOptions = [
        {
            value: ScriptRunLevel.None,
            label: 'None'
        },
        {
            value: ScriptRunLevel.ReadOnly,
            label: 'Read Only'
        },
        {
            value: ScriptRunLevel.ReadWrite,
            label: 'Read Write'
        }
    ];

    const canActions = [
        {
            label: 'Can Modify Hosts',
            value: Permissions.CanModifyHosts,
            update: (v: boolean) => {
                SetPermissions(p => {
                    p.CanModifyHosts = v;
                    return { ...p };
                });
            }
        },
        {
            label: 'Can Modify Groups',
            value: Permissions.CanModifyGroups,
            update: (v: boolean) => {
                SetPermissions(p => {
                    p.CanModifyGroups = v;
                    return { ...p };
                });
            }
        },
        {
            label: 'Can Modify Scripts',
            value: Permissions.CanModifyScripts,
            update: (v: boolean) => {
                SetPermissions(p => {
                    p.CanModifyScripts = v;
                    return { ...p };
                });
            }
        },
        {
            label: 'Can Modify Schedules',
            value: Permissions.CanModifySchedules,
            update: (v: boolean) => {
                SetPermissions(p => {
                    p.CanModifySchedules = v;
                    return { ...p };
                });
            }
        },
        {
            label: 'Can Access Event Log',
            value: Permissions.CanAccessAuditLog,
            update: (v: boolean) => {
                SetPermissions(p => {
                    p.CanAccessAuditLog = v;
                    return { ...p };
                });
            }
        },
        {
            label: 'Can Modify Users',
            value: Permissions.CanModifyUsers,
            helpText: 'Users can always change their own password and API key',
            update: (v: boolean) => {
                SetPermissions(p => {
                    p.CanModifyUsers = v;
                    return { ...p };
                });
            }
        },
        {
            label: 'Can Modify Auto-Register Rules',
            value: Permissions.CanModifyAutoregister,
            update: (v: boolean) => {
                SetPermissions(p => {
                    p.CanModifyAutoregister = v;
                    return { ...p };
                });
            }
        },
        {
            label: 'Can Modify System Settings',
            value: Permissions.CanModifySystem,
            update: (v: boolean) => {
                SetPermissions(p => {
                    p.CanModifySystem = v;
                    return { ...p };
                });
            }
        },
    ];

    const onChangeScriptRunLevel = (level: ScriptRunLevel) => {
        SetPermissions(p => {
            p.ScriptRunLevel = level;
            return { ...p };
        });
    };

    return (<div className="mt-2">
        <h5>Permissions</h5>
        <Input.Radio label='Can run scripts' choices={scriptRunLevelOptions} buttons defaultValue={Permissions.ScriptRunLevel} onChange={onChangeScriptRunLevel} />
        <div className="checkboxes">
            {canActions.map((action, idx) => {
                return (<Input.Checkbox key={idx} label={action.label} helpText={action.helpText} defaultValue={action.value} onChange={action.update} />);
            })}
        </div>
    </div>);
};
