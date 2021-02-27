import * as React from 'react';
import { Schedule, ScheduleType } from '../../types/Schedule';
import { PageLoading } from '../../components/Loading';
import { Page } from '../../components/Page';
import { Buttons, CreateButton } from '../../components/Button';
import { Table } from '../../components/Table';
import { ScheduleListItem } from './ScheduleListItem';
import { Script } from '../../types/Script';

interface ScheduleListState {
    loading: boolean;
    schedules?: ScheduleType[];
    scripts?: Map<string, Script>;
}
export class ScheduleList extends React.Component<unknown, ScheduleListState> {
    constructor(props: unknown) {
        super(props);
        this.state = {
            loading: true,
        };
    }

    componentDidMount(): void {
        this.loadData();
    }

    private loadSchedules = () => {
        return Schedule.List().then(schedules => {
            this.setState({
                schedules: schedules,
            });
        });
    }

    private loadScripts = () => {
        return Script.List().then(scripts => {
            const m = new Map<string, Script>();
            scripts.forEach(script => {
                m.set(script.ID, script);
            });
            this.setState({
                scripts: m,
            });
        });
    }

    private loadData = () => {
        Promise.all([this.loadSchedules(), this.loadScripts()]).then(() => {
            this.setState({ loading: false });
        });
    }

    render(): JSX.Element {
        if (this.state.loading) {
            return ( <PageLoading /> );
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
                            this.state.schedules.map(schedule => {
                                const script = this.state.scripts.get(schedule.ScriptID);
                                return <ScheduleListItem schedule={schedule} script={script} key={schedule.ID} onReload={this.loadData}></ScheduleListItem>;
                            })
                        }
                    </Table.Body>
                </Table.Table>
            </Page>
        );
    }
}
