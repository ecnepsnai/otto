import * as React from 'react';
import { User } from '../../types/User';
import { Loading } from '../../components/Loading';
import { Card } from '../../components/Card';
import { Icon } from '../../components/Icon';
import { CreateButton, Button } from '../../components/Button';
import { Table } from '../../components/Table';
import { EnabledBadge } from '../../components/Badge';
import { Style } from '../../components/Style';
import { Rand } from '../../services/Rand';
import { MenuItem } from '../../components/Menu';
import { StateManager } from '../../services/StateManager';
import { ModalButton, Modal, GlobalModalFrame } from '../../components/Modal';
import { Input, Checkbox, Form } from '../../components/Form';

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

export interface OptionsUsersProps {}
interface OptionsUsersState {
    loading: boolean;
    users?: User[];
}
export class OptionsUsers extends React.Component<OptionsUsersProps, OptionsUsersState> {
    constructor(props: OptionsUsersProps) {
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
            deleteMenuItem = (<MenuItem label="Delete" icon={<Icon.Delete />} onClick={this.deleteUserMenuClick(user)}/>);
        }

        return (
            <Table.Row key={Rand.ID()}>
                <td>{user.Username}</td>
                <td>{user.Email}</td>
                <td><EnabledBadge value={user.Enabled} /></td>
                <Table.Menu>
                    <MenuItem label="Edit" icon={<Icon.Edit />} onClick={this.editUserMenuClick(user)}/>
                    { deleteMenuItem }
                </Table.Menu>
            </Table.Row>
        );
    }

    private content = () => {
        if (this.state.loading) { return (<Loading />); }
        return (
            <div>
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
            </div>
        );
    }

    render(): JSX.Element {
        return (
            <Card.Card>
                <Card.Header>
                    <Icon.Label icon={<Icon.Users />} label="Users" />
                </Card.Header>
                <Card.Body>
                    { this.content() }
                </Card.Body>
            </Card.Card>
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

    render(): JSX.Element {
        const title = this.state.value.Username != '' ? 'Edit User' : 'New User';
        const buttons: ModalButton[] = [
            {
                label: 'Discard',
                color: Style.Palette.Secondary,
            },
            {
                label: 'Save',
                color: Style.Palette.Primary,
                onClick: () => {
                    this.props.onUpdate(this.state.value);
                }
            }
        ];

        return (
            <Modal title={title} static buttons={buttons}>
                <Form>
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
                </Form>
            </Modal>
        );
    }
}
