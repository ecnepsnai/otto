import * as React from 'react';
import { Event } from '../../types/Event';
import { Table } from '../../components/Table';
import { DateLabel } from '../../components/DateLabel';
import { Button } from '../../components/Button';
import { Style } from '../../components/Style';
import { Icon } from '../../components/Icon';
import { GlobalModalFrame } from '../../components/Modal';
import { EventDialog } from './EventDialog';

export interface EventListItemProps { event: Event }
export class EventListItem extends React.Component<EventListItemProps, {}> {
    private viewClick = () => {
        GlobalModalFrame.showModal(<EventDialog event={this.props.event} />);
    }

    render(): JSX.Element {
        return (
            <Table.Row>
                <td>{ this.props.event.Event }</td>
                <td><DateLabel date={this.props.event.Time} /></td>
                <td>
                    <Button color={Style.Palette.Secondary} outline size={Style.Size.XS} onClick={this.viewClick}>
                        <Icon.Label icon={<Icon.Eye />} label="View" />
                    </Button>
                </td>
            </Table.Row>
        );
    }
}
