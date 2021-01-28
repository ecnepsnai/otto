import * as React from 'react';
import { EnabledBadge } from '../../../components/Badge';
import { CreateButton, Button } from '../../../components/Button';
import { Input, Checkbox } from '../../../components/Form';
import { Icon } from '../../../components/Icon';
import { PageLoading } from '../../../components/Loading';
import { Menu } from '../../../components/Menu';
import { GlobalModalFrame, Modal, ModalForm } from '../../../components/Modal';
import { Page } from '../../../components/Page';
import { Style } from '../../../components/Style';
import { Table } from '../../../components/Table';
import { Rand } from '../../../services/Rand';
import { StateManager } from '../../../services/StateManager';
import { User } from '../../../types/User';

export class UserManager {
    public static EditCurrentUser(): Promise<User> {
        return new Promise(resolve => {
            const editUser = (user: User) => {
                user.Save().then(resolve);
            };
            GlobalModalFrame.showModal(<OptionsUsersModal user={StateManager.Current().User} onUpdate={editUser} />);
        });
    }
}

export interface SystemUsersProps {}
interface SystemUsersState {
    loading: boolean;
    users?: User[];
}
export class SystemUsers extends React.Component<SystemUsersProps, SystemUsersState> {
    constructor(props: SystemUsersProps) {
        super(props);
        this.state = {
            loading: true,
        };
    }

    private loadUsers = () => {
        this.setState({ loading: true });
        User.List().then(users => {
            this.setState({
                users: users,
                loading: false,
            });
        });
    }

    componentDidMount(): void {
        this.loadUsers();
    }

    private newUser = (user: User) => {
        User.New({
            Username: user.Username,
            Email: user.Email,
            Password: user.Password,
        }).then(() => {
            this.loadUsers();
        });
    }

    private updateUser = (user: User) => {
        user.Save().then(() => {
            this.loadUsers();
        });
    }

    private newUserClick = () => {
        GlobalModalFrame.showModal(<OptionsUsersModal onUpdate={this.newUser} />);
    }

    private editUserMenuClick = (user: User) => {
        return () => {
            GlobalModalFrame.showModal(<OptionsUsersModal user={user} onUpdate={this.updateUser} />);
        };
    }
    private deleteUserMenuClick = (user: User) => {
        return () => {
            Modal.delete('Delete User?', 'Are you sure you want to delete this user? This can not be undone.').then(confirmed => {
                if (!confirmed) {
                    return;
                }

                user.Delete().then(() => {
                    this.loadUsers();
                });
            });
        };
    }

    private userRow = (user: User) => {
        let deleteMenuItem: JSX.Element;
        if (StateManager.Current().User.Username != user.Username) {
            deleteMenuItem = (<React.Fragment>
                <Menu.Divider />
                <Menu.Item label="Delete" icon={<Icon.Delete />} onClick={this.deleteUserMenuClick(user)}/>
            </React.Fragment>);
        }

        return (
            <Table.Row key={Rand.ID()}>
                <td>{user.Username}</td>
                <td>{user.Email}</td>
                <td><EnabledBadge value={user.Enabled} /></td>
                <Table.Menu>
                    <Menu.Item label="Edit" icon={<Icon.Edit />} onClick={this.editUserMenuClick(user)}/>
                    { deleteMenuItem }
                </Table.Menu>
            </Table.Row>
        );
    }

    render(): JSX.Element {
        if (this.state.loading) { return (<PageLoading />); }

        return (
            <Page title="Users">
                <CreateButton onClick={this.newUserClick} />
                <Table.Table>
                    <Table.Head>
                        <Table.Column>Username</Table.Column>
                        <Table.Column>Email</Table.Column>
                        <Table.Column>Enabled</Table.Column>
                        <Table.MenuColumn />
                    </Table.Head>
                    <Table.Body>
                        { this.state.users.map(this.userRow) }
                    </Table.Body>
                </Table.Table>
            </Page>
        );
    }
}

interface OptionsUsersModalProps {
    user?: User;
    onUpdate: (user: User) => (void);
}
interface OptionsUsersModalState {
    value: User;
    showPasswordField: boolean;
    isNew: boolean;
}
class OptionsUsersModal extends React.Component<OptionsUsersModalProps, OptionsUsersModalState> {
    constructor(props: OptionsUsersModalProps) {
        super(props);
        const isNew = props.user == undefined;
        this.state = {
            value: props.user || User.Blank(),
            showPasswordField: isNew,
            isNew: isNew,
        };
    }

    private changeUsername = (Username: string) => {
        this.setState(state => {
            const user = state.value;
            user.Username = Username;
            return { value: user };
        });
    }

    private changeEmail = (Email: string) => {
        this.setState(state => {
            const user = state.value;
            user.Email = Email;
            return { value: user };
        });
    }

    private changePassword = (Password: string) => {
        this.setState(state => {
            const user = state.value;
            user.Password = Password;
            return { value: user };
        });
    }

    private changeEnabled = (Enabled: boolean) => {
        this.setState(state => {
            const user = state.value;
            user.Enabled = Enabled;
            return { value: user };
        });
    }

    private showPasswordField = () => {
        this.setState({ showPasswordField: true });
    }

    private passwordField = () => {
        if (!this.state.showPasswordField) {
            return (
                <div className="mb-3">
                    <Button color={Style.Palette.Primary} outline onClick={this.showPasswordField}><Icon.Label icon={<Icon.Edit />} label="Change Password" /></Button>
                </div>
            );
        }

        return (
            <Input
                type="password"
                label="Password"
                onChange={this.changePassword}
                required />
        );
    }

    private enabledCheckbox = () => {
        if (this.state.isNew || StateManager.Current().User.Username == this.props.user.Username) { return null; }

        return (
            <Checkbox
                label="Enabled"
                defaultValue={this.state.value.Enabled}
                onChange={this.changeEnabled} />
        );
    }

    private onSubmit = (): Promise<void> => {
        return new Promise(resolve => {
            this.props.onUpdate(this.state.value);
            resolve();
        });
    }

    render(): JSX.Element {
        const title = this.state.value.Username != '' ? 'Edit User' : 'New User';

        return (
            <ModalForm title={title} onSubmit={this.onSubmit}>
                <Input
                    type="text"
                    label="Username"
                    defaultValue={this.state.value.Username}
                    onChange={this.changeUsername}
                    disabled={this.props.user != undefined}
                    required />
                <Input
                    type="email"
                    label="Email"
                    defaultValue={this.state.value.Email}
                    onChange={this.changeEmail}
                    required />
                { this.passwordField() }
                { this.enabledCheckbox() }
            </ModalForm>
        );
    }
}
