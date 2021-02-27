import * as React from 'react';
import { Event, EventType } from '../../types/Event';
import { PageLoading } from '../../components/Loading';
import { Page } from '../../components/Page';
import { Table } from '../../components/Table';
import { EventListItem } from './EventListItem';

export const EventList: React.FC = () => {
    const [loading, setLoading] = React.useState(true);
    const [events, setEvents] = React.useState<EventType[]>();

    React.useEffect(() => {
        loadData();
    }, []);

    const loadEvents = () => {
        return Event.List(20).then(events => {
            setEvents(events);
        });
    };

    const loadData = () => {
        Promise.all([loadEvents()]).then(() => {
            setLoading(false);
        });
    };

    if (loading) {
        return ( <PageLoading /> );
    }
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
                        events.map(event => {
                            return <EventListItem event={event} key={event.ID}></EventListItem>;
                        })
                    }
                </Table.Body>
            </Table.Table>
        </Page>
    );
};
