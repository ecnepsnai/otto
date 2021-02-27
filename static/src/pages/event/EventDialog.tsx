import * as React from 'react';
import { ButtonLink } from '../../components/Button';
import { Card } from '../../components/Card';
import { DateLabel } from '../../components/DateLabel';
import { Icon } from '../../components/Icon';
import { ListGroup } from '../../components/ListGroup';
import { Modal } from '../../components/Modal';
import { Style } from '../../components/Style';
import { EventType } from '../../types/Event';

interface EventDialogProps { event: EventType; }
export const EventDialog: React.FC<EventDialogProps> = (props: EventDialogProps) => {
    const details = Object.keys(props.event.Details).sort().map(key => {
        return {
            key: key,
            value: props.event.Details[key],
        };
    });

    return (
        <Modal title="View Event">
            <Card.Card>
                <Card.Header>Event Metadata</Card.Header>
                <ListGroup.List>
                    <ListGroup.TextItem title="Event">{ props.event.Event }</ListGroup.TextItem>
                    <ListGroup.TextItem title="Time"><DateLabel date={props.event.Time} /></ListGroup.TextItem>
                </ListGroup.List>
            </Card.Card>
            <Card.Card className="mt-3">
                <Card.Header>Event Details</Card.Header>
                <ListGroup.List>
                    { details.map((o, idx) => {
                        return (<EventDetail prop={o.key} value={o.value} key={idx} />);
                    })}
                </ListGroup.List>
            </Card.Card>
        </Modal>
    );
};

interface EventDetailProps { prop: string; value: string; }
export const EventDetail: React.FC<EventDetailProps> = (props: EventDetailProps) => {
    const linkButton = (): JSX.Element => {
        switch (props.prop) {
        case 'group_id':
            return (<ButtonLink to={'/groups/group/' + props.value} outline color={Style.Palette.Primary} size={Style.Size.XS}><Icon.Label icon={<Icon.Eye />} label="View" /></ButtonLink>);
        case 'host_id':
            return (<ButtonLink to={'/hosts/host/' + props.value} outline color={Style.Palette.Primary} size={Style.Size.XS}><Icon.Label icon={<Icon.Eye />} label="View" /></ButtonLink>);
        case 'schedule_id':
            return (<ButtonLink to={'/schedules/schedule/' + props.value} outline color={Style.Palette.Primary} size={Style.Size.XS}><Icon.Label icon={<Icon.Eye />} label="View" /></ButtonLink>);
        case 'script_id':
            return (<ButtonLink to={'/scripts/script/' + props.value} outline color={Style.Palette.Primary} size={Style.Size.XS}><Icon.Label icon={<Icon.Eye />} label="View" /></ButtonLink>);
        }

        return null;
    };

    return (
        <ListGroup.TextItem title={props.prop}>
            <code>{props.value}</code>
            { linkButton() }
        </ListGroup.TextItem>
    );
};
