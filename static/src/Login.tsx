import * as React from 'react';
import * as ReactDOM from 'react-dom';
import { Alert } from './components/Alert';
import '../css/login.scss';

enum LoginFlowStage {
    Login = 1,
    ChangePassword = 2
}

enum LoginError {
    Unauthorized = 1,
    LoggedOut = 2,
    IncorrectPassword = 3,
    LoginError = 4,
}

enum LoginStatus {
    Error = -1,
    Success = 0,
    MustChangePassword = 1,
}

interface LoginFormProps {
    doLogin: (username: string, password: string) => Promise<unknown>;
    loading?: boolean;
    error?: LoginError;
}
interface LoginFormState {
    username: string;
    password: string;
}
class LoginForm extends React.Component<LoginFormProps, LoginFormState> {
    constructor(props: LoginFormProps) {
        super(props);
        this.state = {
            username: '',
            password: '',
        };
    }

    private loginError = () => {
        switch (this.props.error) {
        case LoginError.Unauthorized:
            return (<Alert.Danger>You must be logged in to access that page</Alert.Danger>);
        case LoginError.LoggedOut:
            return (<Alert.Info>You&apos;ve been logged out</Alert.Info>);
        case LoginError.IncorrectPassword:
            return (<Alert.Danger>Incorrect username or password</Alert.Danger>);
        case LoginError.LoginError:
            return (<Alert.Danger>Internal Server Error</Alert.Danger>);
        }
    }

    private changeUsername = (event: React.FormEvent<HTMLInputElement>) => {
        const target = event.target as HTMLInputElement;
        this.setState({ username: target.value });
    }

    private changePassword = (event: React.FormEvent<HTMLInputElement>) => {
        const target = event.target as HTMLInputElement;
        this.setState({ password: target.value });
    }

    private loginFormSubmit = (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault();
        this.props.doLogin(this.state.username, this.state.password).then(() => {
            this.setState({ password: '' });
        });
    }

    render() {
        return (
            <form onSubmit={this.loginFormSubmit}>
                { this.loginError() }
                <label htmlFor="username" className="visually-hidden">Username</label>
                <input type="text" value={this.state.username} onChange={this.changeUsername} className="form-control input-first" placeholder="Username" required autoFocus disabled={this.props.loading}/>
                <label htmlFor="password" className="visually-hidden">Password</label>
                <input type="password" value={this.state.password} onChange={this.changePassword} className="form-control input-second" placeholder="Password" required disabled={this.props.loading}/>
                <div className="d-grid">
                    <button className="btn btn-lg login-button" id="login_button" type="submit" disabled={this.props.loading}>Sign in</button>
                </div>
            </form>
        );
    }
}

interface ChangePasswordFormProps {
    doChangePassword: (password: string) => void;
    loading?: boolean;
    error?: LoginError;
}
interface ChangePasswordFormState {
    password1: string;
    password2: string;
}
class ChangePasswordForm extends React.Component<ChangePasswordFormProps, ChangePasswordFormState> {
    constructor(props: ChangePasswordFormProps) {
        super(props);
        this.state = {
            password1: '',
            password2: '',
        };
    }

    private changePasswordError = () => {
        switch (this.props.error) {
        case LoginError.LoginError:
            return (<Alert.Danger>Internal Server Error</Alert.Danger>);
        }

        return (<Alert.Warning>You Must Change Your Password</Alert.Warning>);
    }

    private changePassword1 = (event: React.FormEvent<HTMLInputElement>) => {
        const target = event.target as HTMLInputElement;
        this.setState({ password1: target.value });
    }

    private changePassword2 = (event: React.FormEvent<HTMLInputElement>) => {
        const target = event.target as HTMLInputElement;
        this.setState({ password2: target.value });
    }

    private changePasswordFormSubmit = (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault();
        this.props.doChangePassword(this.state.password1);
    }

    render() {
        return (
            <form onSubmit={this.changePasswordFormSubmit}>
                { this.changePasswordError() }
                <label htmlFor="password" className="visually-hidden">New Password</label>
                <input type="password" value={this.state.password1} onChange={this.changePassword1} className="form-control input-first" placeholder="New Password" required autoFocus disabled={this.props.loading}/>
                <label htmlFor="password" className="visually-hidden">Confirm New Password</label>
                <input type="password" value={this.state.password2} onChange={this.changePassword2} className="form-control input-second" placeholder="Confirm New Password" required disabled={this.props.loading}/>
                <div className="d-grid">
                    <button className="btn btn-lg login-button" id="login_button" type="submit" disabled={this.props.loading}>Change Password</button>
                </div>
            </form>
        );
    }
}

interface LoginState {
    stage: LoginFlowStage;
    loading?: boolean;
    error?: LoginError;
    redirect?: string;
}
class Login extends React.Component<unknown, LoginState> {
    constructor(props: unknown) {
        super(props);

        const urlParams = new URLSearchParams(window.location.search);
        let initialError: LoginError;
        if (urlParams.has('unauthorized')) {
            initialError = LoginError.Unauthorized;
        } else if (urlParams.has('logged_out')) {
            initialError = LoginError.LoggedOut;
        }
        let redirect: string;
        if (urlParams.has('redirect')) {
            redirect = urlParams.get('redirect');
            if (!redirect.startsWith('/')) {
                redirect = '/' + redirect;
            }
        }

        this.state = {
            stage: LoginFlowStage.Login,
            error: initialError,
            redirect: redirect,
        };
    }

    private doLogin = (username: string, password: string): Promise<void> => {
        this.setState({ loading: true });
        const credentials = {
            Username: username,
            Password: password,
        };
        return fetch('/api/login', { method: 'POST', body: JSON.stringify(credentials) }).then(response => {
            return response.json().then(results => {
                console.log(results);

                if (results.code != 200) {
                    this.setState({
                        loading: false,
                        error: LoginError.IncorrectPassword,
                    });
                    return;
                }

                const status = results.data as LoginStatus;
                if (status === LoginStatus.Success) {
                    this.finishLogin();
                } else if (status === LoginStatus.MustChangePassword) {
                    this.setState({
                        loading: false,
                        stage: LoginFlowStage.ChangePassword
                    });
                }
                return;
            }, () => {
                this.setState({
                    loading: false,
                    error: LoginError.LoginError,
                });
                return;
            });
        }, () => {
            this.setState({
                loading: false,
                error: LoginError.LoginError,
            });
            return;
        }).catch(() => {
            this.setState({
                loading: false,
                error: LoginError.LoginError,
            });
            return;
        });
    }

    private doChangePassword = (password: string): Promise<void> => {
        const request = {
            Password: password,
        };
        return fetch('/api/users/reset_password', { method: 'POST', body: JSON.stringify(request) }).then(response => {
            return response.json().then(results => {
                console.log(results);

                if (results.code != 200) {
                    this.setState({
                        loading: false,
                        error: LoginError.LoginError,
                    });
                    return;
                }

                this.finishLogin();
                return;
            }, () => {
                this.setState({
                    loading: false,
                    error: LoginError.LoginError,
                });
                return;
            });
        }, () => {
            this.setState({
                loading: false,
                error: LoginError.LoginError,
            });
            return;
        }).catch(() => {
            this.setState({
                loading: false,
                error: LoginError.LoginError,
            });
            return;
        });
    }

    private finishLogin = () => {
        if (this.state.redirect) {
            location.href = this.state.redirect;
        } else {
            location.href = '/';
        }
    }

    private content = () => {
        switch (this.state.stage) {
        case LoginFlowStage.Login:
            return (<LoginForm doLogin={this.doLogin} loading={this.state.loading} error={this.state.error}/>);
        case LoginFlowStage.ChangePassword:
            return (<ChangePasswordForm doChangePassword={this.doChangePassword} loading={this.state.loading} error={this.state.error}/>);
        }

        return null;
    }

    render(): JSX.Element {
        return (
            <div className="form-signin">
                <img className="mb-4" src="assets/img/logo_light.svg" alt="otto logo" width="72" height="72" />
                { this.content() }
            </div>
        );
    }
}

ReactDOM.render(
    <Login />,
    document.getElementById('login')
);
