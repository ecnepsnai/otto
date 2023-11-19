import * as React from 'react';
import { StateManager } from './services/StateManager';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { GlobalModalFrame } from './components/Modal';
import { GlobalNotificationFrame } from './components/Notification';
import { GlobalContextMenuFrame } from './components/ContextMenu';
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
import { RunbookEdit } from './pages/runbook/RunbookEdit';
import { RunbookView } from './pages/runbook/RunbookView';
import { RunbookList } from './pages/runbook/RunbookList';
import { EventList } from './pages/event/EventList';
import { SystemOptions } from './pages/system/options/SystemOptions';
import { SystemUsers } from './pages/system/users/SystemUsers';
import { SystemRegister } from './pages/system/register/SystemRegister';
import { ErrorBoundary } from './components/ErrorBoundary';
import '../css/main.scss';

export const App: React.FC = () => {
    const [loading, setLoading] = React.useState(true);

    React.useEffect(() => {
        StateManager.Initialize().then(() => {
            setLoading(false);
        });
    }, []);

    if (loading) {
        return (
            <div className="container mt-2">
                <Loading />
            </div>
        );
    }

    return (<ErrorBoundary>
        <BrowserRouter>
            <Nav />
            <Routes>
                <Route path="/hosts/host/:id/edit" element={<HostEdit />} />
                <Route path="/hosts/host/:id" element={<HostView />} />
                <Route path="/hosts/host" element={<HostEdit />} />
                <Route path="/hosts" element={<HostList />} />
                <Route path="/groups/group/:id/edit" element={<GroupEdit />} />
                <Route path="/groups/group/:id" element={<GroupView />} />
                <Route path="/groups/group" element={<GroupEdit />} />
                <Route path="/groups" element={<GroupList />} />
                <Route path="/scripts/script/:id/edit" element={<ScriptEdit />} />
                <Route path="/scripts/script/:id" element={<ScriptView />} />
                <Route path="/scripts/script" element={<ScriptEdit />} />
                <Route path="/scripts" element={<ScriptList />} />
                <Route path="/schedules/schedule/:id/edit" element={<ScheduleEdit />} />
                <Route path="/schedules/schedule/:id" element={<ScheduleView />} />
                <Route path="/schedules/schedule" element={<ScheduleEdit />} />
                <Route path="/schedules" element={<ScheduleList />} />
                <Route path="/runbooks/runbook/:id/edit" element={<RunbookEdit />} />
                <Route path="/runbooks/runbook/:id" element={<RunbookView />} />
                <Route path="/runbooks/runbook" element={<RunbookEdit />} />
                <Route path="/runbooks" element={<RunbookList />} />
                <Route path="/system/options" element={<SystemOptions />} />
                <Route path="/system/users" element={<SystemUsers />} />
                <Route path="/system/register" element={<SystemRegister />} />
                <Route path="/events" element={<EventList />} />
            </Routes>
            <GlobalModalFrame />
            <GlobalNotificationFrame />
            <GlobalContextMenuFrame />
        </BrowserRouter>
    </ErrorBoundary>);
};
