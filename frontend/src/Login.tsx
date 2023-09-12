import * as React from 'react';
import { createRoot } from 'react-dom/client';
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
const LoginForm: React.FC<LoginFormProps> = (props: LoginFormProps) => {
    const [username, setUsername] = React.useState('');
    const [password, setPassword] = React.useState('');

    const loginError = () => {
        switch (props.error) {
            case LoginError.Unauthorized:
                return (<Alert.Danger>You must be logged in to access that page</Alert.Danger>);
            case LoginError.LoggedOut:
                return (<Alert.Info>You&apos;ve been logged out</Alert.Info>);
            case LoginError.IncorrectPassword:
                return (<Alert.Danger>Incorrect username or password</Alert.Danger>);
            case LoginError.LoginError:
                return (<Alert.Danger>Internal Server Error</Alert.Danger>);
        }
    };

    const changeUsername = (event: React.FormEvent<HTMLInputElement>) => {
        const target = event.target as HTMLInputElement;
        setUsername(target.value);
    };

    const changePassword = (event: React.FormEvent<HTMLInputElement>) => {
        const target = event.target as HTMLInputElement;
        setPassword(target.value);
    };

    const loginFormSubmit = (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault();
        props.doLogin(username, password).then(() => {
            setPassword('');
        });
    };

    return (
        <form onSubmit={loginFormSubmit}>
            { loginError()}
            <label htmlFor="username" className="visually-hidden">Username</label>
            <input type="text" id="username" value={username} onChange={changeUsername} className="form-control input-first" placeholder="Username" required autoFocus disabled={props.loading} />
            <label htmlFor="password" className="visually-hidden">Password</label>
            <input type="password" id="password" value={password} onChange={changePassword} className="form-control input-second" placeholder="Password" required disabled={props.loading} />
            <div className="d-grid">
                <button className="btn btn-lg login-button" id="login_button" type="submit" disabled={props.loading}>Sign in</button>
            </div>
        </form>
    );
};

interface ChangePasswordFormProps {
    doChangePassword: (password: string) => void;
    loading?: boolean;
    error?: LoginError;
}
const ChangePasswordForm: React.FC<ChangePasswordFormProps> = (props: ChangePasswordFormProps) => {
    const [password1, setPassword1] = React.useState('');
    const [password2, setPassword2] = React.useState('');

    const changePasswordError = () => {
        switch (props.error) {
            case LoginError.LoginError:
                return (<Alert.Danger>Internal Server Error</Alert.Danger>);
        }

        return (<Alert.Warning>You Must Change Your Password</Alert.Warning>);
    };

    const changePassword1 = (event: React.FormEvent<HTMLInputElement>) => {
        const target = event.target as HTMLInputElement;
        setPassword1(target.value);
    };

    const changePassword2 = (event: React.FormEvent<HTMLInputElement>) => {
        const target = event.target as HTMLInputElement;
        setPassword2(target.value);
    };

    const changePasswordFormSubmit = (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault();
        props.doChangePassword(password1);
    };

    return (
        <form onSubmit={changePasswordFormSubmit}>
            { changePasswordError()}
            <label htmlFor="password1" className="visually-hidden">New Password</label>
            <input type="password" id="password1" value={password1} onChange={changePassword1} className="form-control input-first" placeholder="New Password" required autoFocus disabled={props.loading} />
            <label htmlFor="password2" className="visually-hidden">Confirm New Password</label>
            <input type="password" id="password2" value={password2} onChange={changePassword2} className="form-control input-second" placeholder="Confirm New Password" required disabled={props.loading} />
            <div className="d-grid">
                <button className="btn btn-lg login-button" id="change_password_button" type="submit" disabled={props.loading}>Change Password</button>
            </div>
        </form>
    );
};

export const Login: React.FC = () => {
    const [stage, setStage] = React.useState<LoginFlowStage>(LoginFlowStage.Login);
    const [loading, setLoading] = React.useState<boolean>(false);
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
    const [error, setError] = React.useState<LoginError>(initialError);

    const doLogin = async (username: string, password: string): Promise<void> => {
        setLoading(true);
        const credentials = {
            Username: username,
            Password: password,
        };

        try {
            const response = await fetch('/api/login', { method: 'POST', body: JSON.stringify(credentials) });
            const results = await response.json();
            if (results.code != 200) {
                setLoading(false);
                setError(LoginError.IncorrectPassword);
                return;
            }

            const status = results.data as LoginStatus;
            if (status === LoginStatus.Success) {
                finishLogin();
            } else if (status === LoginStatus.MustChangePassword) {
                setLoading(false);
                setStage(LoginFlowStage.ChangePassword);
            }
        } catch (err) {
            console.error('Login error', err);
            setLoading(false);
            setError(LoginError.LoginError);
        }
    };

    const doChangePassword = async (password: string): Promise<void> => {
        const request = {
            Password: password,
        };

        try {
            const response = await fetch('/api/users/reset_password', { method: 'POST', body: JSON.stringify(request) });
            const results = await response.json();
            if (results.code != 200) {
                setLoading(false);
                setError(LoginError.LoginError);
                return;
            }

            finishLogin();
        } catch (err) {
            console.error('Error changing password', err);
            setLoading(false);
            setError(LoginError.LoginError);
        }
    };

    const finishLogin = () => {
        if (redirect) {
            location.href = redirect;
        } else {
            location.href = '/';
        }
    };

    const content = () => {
        switch (stage) {
            case LoginFlowStage.Login:
                return (<LoginForm doLogin={doLogin} loading={loading} error={error} />);
            case LoginFlowStage.ChangePassword:
                return (<ChangePasswordForm doChangePassword={doChangePassword} loading={loading} error={error} />);
        }

        return null;
    };

    return (
        <div className="form-signin">
            <img className="mb-4" src="assets/img/logo_light.svg" alt="otto logo" width="72" height="72" />
            { content()}
        </div>
    );
};

const container = document.getElementById('login');
const root = createRoot(container);
root.render(<Login />);
