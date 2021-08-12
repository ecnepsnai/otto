import * as React from 'react';
import { EnabledBadge } from '../../../components/Badge';
import { AddButton, Button, ConfirmButton } from '../../../components/Button';
import { Input } from '../../../components/input/Input';
import { Icon } from '../../../components/Icon';
import { PageLoading } from '../../../components/Loading';
import { GlobalModalFrame, Modal, ModalForm } from '../../../components/Modal';
import { Page } from '../../../components/Page';
import { Style } from '../../../components/Style';
import { Table } from '../../../components/Table';
import { Rand } from '../../../services/Rand';
import { StateManager } from '../../../services/StateManager';
import { User, UserType } from '../../../types/User';
import { ContextMenuItem } from '../../../components/ContextMenu';

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
            Email: user.Email,
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

    const userRow = (user: UserType) => {
        const contextMenu: (ContextMenuItem | 'separator')[] = [
            {
                title: 'Edit',
                icon: (<Icon.Edit />),
                onClick: () => {
                    editUserMenuClick(user);
                }
            },
        ];
        if (StateManager.Current().User.Username != user.Username) {
            contextMenu.push('separator');
            contextMenu.push({
                title: 'Delete',
                icon: (<Icon.Delete />),
                onClick: () => {
                    deleteUserMenuClick(user);
                }
            });
        }

        return (
            <Table.Row key={Rand.ID()} menu={contextMenu}>
                <td>{user.Username}</td>
                <td>{user.Email}</td>
                <td><EnabledBadge value={user.CanLogIn} trueText="Yes" falseText="No" /></td>
            </Table.Row>
        );
    };

    if (loading) {
        return (<PageLoading />);
    }

    const toolbar = (
        <React.Fragment>
            <AddButton onClick={newUserClick} />
        </React.Fragment>
    );

    return (
        <Page title="Users" toolbar={toolbar}>
            <Table.Table>
                <Table.Head>
                    <Table.Column>Username</Table.Column>
                    <Table.Column>Email</Table.Column>
                    <Table.Column>Can Login</Table.Column>
                </Table.Head>
                <Table.Body>
                    {users.map(userRow)}
                </Table.Body>
            </Table.Table>
        </Page>
    );
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

    const changeEmail = (Email: string) => {
        setUser(user => {
            user.Email = Email;
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

    const title = user.Username != '' ? 'Edit User' : 'New User';

    return (
        <ModalForm title={title} onSubmit={onSubmit}>
            <Input.Text
                type="text"
                label="Username"
                defaultValue={user.Username}
                onChange={changeUsername}
                disabled={props.user != undefined}
                required />
            <Input.Text
                type="email"
                label="Email"
                defaultValue={user.Email}
                onChange={changeEmail}
                required />
            {passwordField()}
            {resetAPIKey()}
            {canLogInCheckbox()}
            {mustChangePasswordCheckbox()}
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

