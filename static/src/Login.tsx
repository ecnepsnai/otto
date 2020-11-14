import * as React from 'react';
import * as ReactDOM from 'react-dom';
import { Alert } from './components/Alert';
import { Style } from './components/Style';
import '../css/login.scss';

enum LoginError {
    Unauthorized = 1,
    LoggedOut = 2,
    IncorrectPassword = 3,
    LoginError = 4,
}

interface LoginProps {}
interface LoginState {
    username: string;
    password: string;
    loading?: boolean;
    error?: LoginError;
    redirect?: string;
}
class Login extends React.Component<LoginProps, LoginState> {
    constructor(props: LoginProps) {
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
            username: '',
            password: '',
            error: initialError,
            redirect: redirect,
        };
    }

    private loginError = () => {
        switch (this.state.error) {
            case LoginError.Unauthorized:
                return (<Alert color={Style.Palette.Danger}>You must be logged in to access that page</Alert>);
            case LoginError.LoggedOut:
                return (<Alert color={Style.Palette.Info}>You&apos;ve been logged out</Alert>);
            case LoginError.IncorrectPassword:
                return (<Alert color={Style.Palette.Danger}>Incorrect username or password</Alert>);
            case LoginError.LoginError:
                return (<Alert color={Style.Palette.Danger}>Internal Server Error</Alert>);
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

    private doLogin = (username: string, password: string) => {
        const credentials = {
            Username: username,
            Password: password,
        };
        return fetch('/api/login', { method: 'POST', body: JSON.stringify(credentials) }).then(response => {
            response.json().then(results => {
                if (results.code != 200) {
                    this.setState({
                        loading: false,
                        error: LoginError.IncorrectPassword,
                        username: '',
                        password: '',
                    });
                    return;
                }
                if (this.state.redirect) {
                    location.href = this.state.redirect;
                } else {
                    location.href = '/';
                }
            }, () => {
                this.setState({
                    loading: false,
                    error: LoginError.LoginError,
                    username: '',
                    password: '',
                });
                return;
            });
        }, () => {
            this.setState({
                loading: false,
                error: LoginError.LoginError,
                username: '',
                password: '',
            });
            return;
        }).catch(() => {
            this.setState({
                loading: false,
                error: LoginError.LoginError,
                username: '',
                password: '',
            });
            return;
        });
    }

    private formSubmit = () => {
        event.preventDefault();
        this.setState({ loading: true }, () => {
            this.doLogin(this.state.username, this.state.password);
        });
    }

    render(): JSX.Element {
        return (
        <div className="form-signin">
            <img className="mb-4" src="assets/img/logo_light.svg" alt="otto logo" width="72" height="72" />
            <form onSubmit={this.formSubmit}>
                { this.loginError() }
                <label htmlFor="username" className="visually-hidden">Username</label>
                <input type="text" value={this.state.username} onChange={this.changeUsername} className="form-control" placeholder="Username" required autoFocus disabled={this.state.loading}/>
                <label htmlFor="password" className="visually-hidden">Password</label>
                <input type="password" value={this.state.password} onChange={this.changePassword} className="form-control" placeholder="Password" required  disabled={this.state.loading}/>
                <div className="d-grid">
                    <button className="btn btn-lg login-button" id="login_button" type="submit" disabled={this.state.loading}>Sign in</button>
                </div>
            </form>
        </div>
        );
    }
}

ReactDOM.render(
    <Login />,
    document.getElementById('login')
);
