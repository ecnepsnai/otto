import * as React from 'react';
import { StateManager } from './services/StateManager';
import { BrowserRouter as Router, Switch, Route } from 'react-router-dom';
import { GlobalModalFrame } from './components/Modal';
import { GlobalNotificationFrame } from './components/Notification';
import { GlobalRedirectFrame } from './components/Redirect';
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
import { OptionsMain } from './pages/options/OptionsMain';
import '../css/main.scss';
import { Loading } from './components/Loading';

export interface AppProps {}
interface AppState { loading: boolean; }
export class App extends React.Component<AppProps, AppState> {
    constructor(props: AppProps) {
        super(props);
        this.state = { loading: true };
    }
    componentDidMount(): void {
        StateManager.Refresh().then(() => {
            console.log('App bootstrapped');
            this.setState({ loading: false });
        });
    }

    render(): JSX.Element {
        if (this.state.loading) {
            return (<div className="mt-3 ml-3"><Loading /></div>);
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
                    <Route path="/options" component={OptionsMain} />
                </Switch>
                <GlobalModalFrame />
                <GlobalNotificationFrame />
                <GlobalRedirectFrame />
            </Router>
        );
    }
}
