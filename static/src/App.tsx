import * as React from 'react';
import { StateManager } from './services/StateManager';
import { BrowserRouter as Router, Switch, Route } from 'react-router-dom';
import { GlobalModalFrame } from './components/Modal';
import { GlobalNotificationFrame } from './components/Notification';
import { GlobalRedirectFrame } from './components/Redirect';
import { Loading } from './components/Loading';
import { Nav } from './components/Nav';
import { HostEdit } from './pages/host/HostEdit';
import { HostList } from './pages/host/HostList';
import { HostView } from './pages/host/HostView';
import { GroupEdit } from './pages/group/GroupEdit';
import { GroupList } from './pages/group/GroupList';
import { GroupView } from './pages/group/GroupView';
import { ScriptEdit } from './pages/script/ScriptEdit';
import { ScriptList } from './pages/script/ScriptList';
import { ScriptView } from './pages/script/ScriptView';
import { ScheduleEdit } from './pages/schedule/ScheduleEdit';
import { ScheduleView } from './pages/schedule/ScheduleView';
import { ScheduleList } from './pages/schedule/ScheduleList';
import { EventList } from './pages/event/EventList';
import { SystemOptions } from './pages/system/options/SystemOptions';
import '../css/main.scss';
import { SystemUsers } from './pages/system/users/SystemUsers';
import { SystemRegister } from './pages/system/register/SystemRegister';

export interface AppProps {}
interface AppState {
    loading?: boolean;
    errorDetails?: string;
}
export class App extends React.Component<AppProps, AppState> {
    constructor(props: AppProps) {
        super(props);
        this.state = { loading: true };
    }

    componentDidMount(): void {
        StateManager.Refresh().then(() => {
            this.setState({ loading: false });
        });
    }

    static getDerivedStateFromError(error: Error): AppState {
        const errData = JSON.stringify(error, Object.getOwnPropertyNames(error), 4);
        console.log(errData);
        return { errorDetails: errData };
    }

    componentDidCatch(error: Error, errorInfo: React.ErrorInfo): void {
        console.error('An error occurred: ' + error, errorInfo);
    }

    render(): JSX.Element {
        if (this.state.errorDetails) {
            return (
                <div className="container">
                    <div className="card mt-3">
                        <div className="card-header">
                            An Error Occurred
                        </div>
                        <div className="card-body">
                            <p>An unrecoverable error occurred while attempting to render this page. Please report this as an issue on <a href="https://github.com/ecnepsnai/otto/issues/new/choose" target="_blank" rel="noreferrer">Github</a> and include the following information:</p>
                            <pre>{ this.state.errorDetails }</pre>
                        </div>
                    </div>
                </div>
            );
        }

        if (this.state.loading) {
            return (<div className="mt-3 ms-3"><Loading /></div>);
        }

        return (
            <Router>
                <Nav />
                <Switch>
                    <Route path="/hosts/host/:id/edit" component={HostEdit} />
                    <Route path="/hosts/host/:id" component={HostView} />
                    <Route path="/hosts/host" component={HostEdit} />
                    <Route path="/hosts" component={HostList} />
                    <Route path="/groups/group/:id/edit" component={GroupEdit} />
                    <Route path="/groups/group/:id" component={GroupView} />
                    <Route path="/groups/group" component={GroupEdit} />
                    <Route path="/groups" component={GroupList} />
                    <Route path="/scripts/script/:id/edit" component={ScriptEdit} />
                    <Route path="/scripts/script/:id" component={ScriptView} />
                    <Route path="/scripts/script" component={ScriptEdit} />
                    <Route path="/scripts" component={ScriptList} />
                    <Route path="/schedules/schedule/:id/edit" component={ScheduleEdit} />
                    <Route path="/schedules/schedule/:id" component={ScheduleView} />
                    <Route path="/schedules/schedule" component={ScheduleEdit} />
                    <Route path="/schedules" component={ScheduleList} />
                    <Route path="/system/options" component={SystemOptions} />
                    <Route path="/system/users" component={SystemUsers} />
                    <Route path="/system/register" component={SystemRegister} />
                    <Route path="/events" component={EventList} />
                </Switch>
                <GlobalModalFrame />
                <GlobalNotificationFrame />
                <GlobalRedirectFrame />
            </Router>
        );
    }
}
