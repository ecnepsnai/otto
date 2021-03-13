import * as React from 'react';
import { EventType } from '../../types/Event';
import { Table } from '../../components/Table';
import { DateLabel } from '../../components/DateLabel';
import { Button } from '../../components/Button';
import { Style } from '../../components/Style';
import { Icon } from '../../components/Icon';
import { GlobalModalFrame } from '../../components/Modal';
import { EventDialog } from './EventDialog';

interface EventListItemProps { event: EventType }
export const EventListItem: React.FC<EventListItemProps> = (props: EventListItemProps) => {
    const viewClick = () => {
        GlobalModalFrame.showModal(<EventDialog event={props.event} />);
    };

    return (
        <Table.Row>
            <td>{ props.event.Event }</td>
            <td><DateLabel date={props.event.Time} /></td>
            <td>
                <Button color={Style.Palette.Secondary} outline size={Style.Size.XS} onClick={viewClick}>
                    <Icon.Label icon={<Icon.Eye />} label="View" />
                </Button>
            </td>
        </Table.Row>
    );
};
