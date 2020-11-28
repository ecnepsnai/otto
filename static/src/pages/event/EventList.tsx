import * as React from 'react';
import { Event } from '../../types/Event';
import { PageLoading } from '../../components/Loading';
import { Page } from '../../components/Page';
import { Table } from '../../components/Table';
import { EventListItem } from './EventListItem';

export interface EventListProps {}
interface EventListState {
    loading: boolean;
    events?: Event[];
}
export class EventList extends React.Component<EventListProps, EventListState> {
    constructor(props: EventListProps) {
        super(props);
        this.state = {
            loading: true,
        };
    }

    componentDidMount(): void {
        this.loadData();
    }

    private loadEvents = () => {
        return Event.List(20).then(events => {
            this.setState({
                events: events,
            });
        });
    }

    private loadData = () => {
        Promise.all([this.loadEvents()]).then(() => {
            this.setState({ loading: false });
        });
    }

    render(): JSX.Element {
        if (this.state.loading) { return ( <PageLoading /> ); }

        return (
            <Page title="Audit Log">
                <Table.Table>
                    <Table.Head>
                        <Table.Column>Event</Table.Column>
                        <Table.Column>Date</Table.Column>
                        <Table.MenuColumn />
                    </Table.Head>
                    <Table.Body>
                        {
                            this.state.events.map(event => {
                                return <EventListItem event={event} key={event.ID}></EventListItem>;
                            })
                        }
                    </Table.Body>
                </Table.Table>
            </Page>
        );
    }
}
