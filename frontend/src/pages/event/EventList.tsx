import * as React from 'react';
import { Event, EventType } from '../../types/Event';
import { PageLoading } from '../../components/Loading';
import { Page } from '../../components/Page';
import { Column, Table } from '../../components/Table';
import { Button } from '../../components/Button';
import { Style } from '../../components/Style';
import { Icon } from '../../components/Icon';
import { DateLabel } from '../../components/DateLabel';
import { DateSort } from '../../services/Sort';
import { GlobalModalFrame } from '../../components/Modal';
import { EventDialog } from './EventDialog';

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

    const showMoreClick = () => {
        SetShownEvents(events => {
            return AllEvents.slice(0, Math.min(events.length+20, AllEvents.length));
        });
    };

    const showMoreDiabled = () => {
        return ShownEvents.length >= AllEvents.length;
    };

    const viewClick = (event: EventType) => {
        return (mv: React.MouseEvent) => {
            mv.preventDefault();
            GlobalModalFrame.showModal(<EventDialog event={event} />);
        };
    };


    if (IsLoading) {
        return (<PageLoading />);
    }

    const tableCols: Column[] = [
        {
            title: 'Event',
            value: (v: EventType) => {
                return (<a href="#" onClick={viewClick(v)}>{v.Event}</a>);
            },
            sort: 'Event'
        },
        {
            title: 'Date',
            value: (v: EventType) => {
                return (<DateLabel date={v.Time} />);
            },
            sort: (asc: boolean, left: EventType, right: EventType) => {
                return DateSort(asc, left.Time, right.Time);
            }
        }
    ];

    return (
        <Page title="Events">
            <Table columns={tableCols} data={ShownEvents} defaultSort={{ ColumnIdx: 1, Ascending: false }} />
            <div className="mt-2">
                <Button color={Style.Palette.Primary} onClick={showMoreClick} disabled={showMoreDiabled()}><Icon.Label icon={<Icon.Plus />} label="Show More" /></Button>
                <span className="ms-1"><em>{ShownEvents.length} of {AllEvents.length}</em></span>
            </div>
        </Page>
    );
};
