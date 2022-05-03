import * as React from 'react';
import { Event, EventType } from '../../types/Event';
import { PageLoading } from '../../components/Loading';
import { Page } from '../../components/Page';
import { Table } from '../../components/Table';
import { EventListItem } from './EventListItem';
import { Button } from '../../components/Button';
import { Style } from '../../components/Style';
import { Icon } from '../../components/Icon';

export const EventList: React.FC = () => {
    const [IsLoading, SetIsLoading] = React.useState(true);
    const [ShownEvents, SetShownEvents] = React.useState<EventType[]>();
    const [AllEvents, SetAllEvents] = React.useState<EventType[]>();

    React.useEffect(() => {
        loadEvents();
    }, []);

    const loadEvents = () => {
        Event.List(0).then(events => {
            SetAllEvents(events);
            SetShownEvents(events.slice(0, Math.min(20, events.length)));
            SetIsLoading(false);
        });
    };

    const showMoreClient = () => {
        SetShownEvents(events => {
            return AllEvents.slice(0, Math.min(events.length+20, AllEvents.length));
        });
    };

    const showMoreDiabled = () => {
        return ShownEvents.length >= AllEvents.length;
    };


    if (IsLoading) {
        return (<PageLoading />);
    }
    return (
        <Page title="Event Log">
            <Table.Table>
                <Table.Head>
                    <Table.Column>Event</Table.Column>
                    <Table.Column>Date</Table.Column>
                </Table.Head>
                <Table.Body>
                    {
                        ShownEvents.map(event => {
                            return (<EventListItem event={event} key={event.ID}></EventListItem>);
                        })
                    }
                </Table.Body>
            </Table.Table>
            <div className="mt-2">
                <Button color={Style.Palette.Primary} onClick={showMoreClient} disabled={showMoreDiabled()}><Icon.Label icon={<Icon.Plus />} label="Show More" /></Button>
                <span className="ms-1"><em>{ShownEvents.length} of {AllEvents.length}</em></span>
            </div>
        </Page>
    );
};
