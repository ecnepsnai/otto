import * as React from 'react';
import { ButtonLink } from '../../components/Button';
import { Card } from '../../components/Card';
import { DateLabel } from '../../components/DateLabel';
import { Icon } from '../../components/Icon';
import { ListGroup } from '../../components/ListGroup';
import { Modal } from '../../components/Modal';
import { Style } from '../../components/Style';
import { Event } from '../../types/Event';

export interface EventDialogProps { event: Event; }
interface EventDialogState {}
export class EventDialog extends React.Component<EventDialogProps, EventDialogState> {
    constructor(props: EventDialogProps) {
        super(props);
        this.state = { };
    }

    render(): JSX.Element {

        const details = Object.keys(this.props.event.Details).sort().map(key => {
            return {
                key: key,
                value: this.props.event.Details[key],
            };
        });

        return (
            <Modal title="View Event">
                <Card.Card>
                    <Card.Header>Event Metadata</Card.Header>
                    <ListGroup.List>
                        <ListGroup.TextItem title="Event">{ this.props.event.Event }</ListGroup.TextItem>
                        <ListGroup.TextItem title="Time"><DateLabel date={this.props.event.Time} /></ListGroup.TextItem>
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
    }
}

interface EventDetailProps { prop: string; value: string; }
interface EventDetailState {}
class EventDetail extends React.Component<EventDetailProps, EventDetailState> {
    constructor(props: EventDetailProps) {
        super(props);
        this.state = { };
    }

    private linkButton = (): JSX.Element => {
        switch (this.props.prop) {
            case 'group_id':
                return (<ButtonLink to={'/groups/group/' + this.props.value} outline color={Style.Palette.Primary} size={Style.Size.XS}><Icon.Label icon={<Icon.Eye />} label="View" /></ButtonLink>);
            case 'host_id':
                return (<ButtonLink to={'/hosts/host/' + this.props.value} outline color={Style.Palette.Primary} size={Style.Size.XS}><Icon.Label icon={<Icon.Eye />} label="View" /></ButtonLink>);
            case 'schedule_id':
                return (<ButtonLink to={'/schedules/schedule/' + this.props.value} outline color={Style.Palette.Primary} size={Style.Size.XS}><Icon.Label icon={<Icon.Eye />} label="View" /></ButtonLink>);
            case 'script_id':
                return (<ButtonLink to={'/scripts/script/' + this.props.value} outline color={Style.Palette.Primary} size={Style.Size.XS}><Icon.Label icon={<Icon.Eye />} label="View" /></ButtonLink>);
        }

        return null;
    }

    render(): JSX.Element {
        return (
        <ListGroup.TextItem title={this.props.prop}>
            <code>{this.props.value}</code>
            { this.linkButton() }
        </ListGroup.TextItem>
        );
    }
}
