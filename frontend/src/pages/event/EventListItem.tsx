import * as React from 'react';
import { EventType } from '../../types/Event';
import { Table } from '../../components/Table';
import { DateLabel } from '../../components/DateLabel';
import { GlobalModalFrame } from '../../components/Modal';
import { EventDialog } from './EventDialog';

interface EventListItemProps { event: EventType }
export const EventListItem: React.FC<EventListItemProps> = (props: EventListItemProps) => {
    const viewClick = (event: React.MouseEvent) => {
        event.preventDefault();
        GlobalModalFrame.showModal(<EventDialog event={props.event} />);
    };

    return (
        <Table.Row>
            <td><a href="#" onClick={viewClick}>{props.event.Event}</a></td>
            <td><DateLabel date={props.event.Time} /></td>
        </Table.Row>
    );
};
