import * as React from 'react';
import { Schedule, ScheduleType } from '../../types/Schedule';
import { PageLoading } from '../../components/Loading';
import { Page } from '../../components/Page';
import { Buttons, CreateButton } from '../../components/Button';
import { Table } from '../../components/Table';
import { ScheduleListItem } from './ScheduleListItem';
import { Script, ScriptType } from '../../types/Script';

export const ScheduleList: React.FC = () => {
    const [loading, setLoading] = React.useState(true);
    const [schedules, setSchedules] = React.useState<ScheduleType[]>();
    const [scripts, setScripts] = React.useState<Map<string, ScriptType>>();

    React.useEffect(() => {
        loadData();
    }, []);

    const loadSchedules = () => {
        return Schedule.List().then(schedules => {
            setSchedules(schedules);
        });
    };

    const loadScripts = () => {
        return Script.List().then(scripts => {
            const m = new Map<string, ScriptType>();
            scripts.forEach(script => {
                m.set(script.ID, script);
            });
            setScripts(m);
        });
    };

    const loadData = () => {
        Promise.all([loadSchedules(), loadScripts()]).then(() => {
            setLoading(false);
        });
    };

    if (loading) {
        return (<PageLoading />);
    }

    return (
        <Page title="Schedules">
            <Buttons>
                <CreateButton to="/schedules/schedule/" />
            </Buttons>
            <Table.Table>
                <Table.Head>
                    <Table.Column>Name</Table.Column>
                    <Table.Column>Script</Table.Column>
                    <Table.Column>Frequency</Table.Column>
                    <Table.Column>Scope</Table.Column>
                    <Table.Column>Last Run</Table.Column>
                    <Table.Column>Enabled</Table.Column>
                    <Table.MenuColumn />
                </Table.Head>
                <Table.Body>
                    {
                        schedules.map(schedule => {
                            const script = scripts.get(schedule.ScriptID);
                            return <ScheduleListItem schedule={schedule} script={script} key={schedule.ID} onReload={loadData}></ScheduleListItem>;
                        })
                    }
                </Table.Body>
            </Table.Table>
        </Page>
    );
};
