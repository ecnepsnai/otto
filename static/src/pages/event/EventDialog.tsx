import * as React from 'react';
import { Card } from '../../components/Card';
import { DateLabel } from '../../components/DateLabel';
import { Icon } from '../../components/Icon';
import { ListGroup } from '../../components/ListGroup';
import { Modal } from '../../components/Modal';
import { EventType } from '../../types/Event';
import { Group } from '../../types/Group';
import { Host } from '../../types/Host';
import { Schedule } from '../../types/Schedule';
import { Script } from '../../types/Script';

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

interface EventDetailProps {
    prop: string;
    value: string;
}
export const EventDetail: React.FC<EventDetailProps> = (props: EventDetailProps) => {
    const [loading, setLoading] = React.useState<boolean>(false);
    const [itemLabel, setItemLabel] = React.useState<string>();
    const [itemValue, setItemValue] = React.useState<string>();

    React.useEffect(() => {
        switch (props.prop) {
        case 'group_id':
            setLoading(true);
            Group.Get(props.value).then(group => {
                setLoading(false);
                setItemLabel('Group');
                setItemValue(group.Name);
            }, () => {
                setLoading(false);
            });
            break;
        case 'host_id':
            setLoading(true);
            Host.Get(props.value).then(host => {
                setLoading(false);
                setItemLabel('Host');
                setItemValue(host.Name);
            }, () => {
                setLoading(false);
            });
            break;
        case 'schedule_id':
            setLoading(true);
            Schedule.Get(props.value).then(schedule => {
                setLoading(false);
                setItemLabel('Schedule');
                setItemValue(schedule.Name);
            }, () => {
                setLoading(false);
            });
            break;
        case 'script_id':
            setLoading(true);
            Script.Get(props.value).then(script => {
                setLoading(false);
                setItemLabel('Script');
                setItemValue(script.Name);
            }, () => {
                setLoading(false);
            });
            break;
        }
    }, []);

    const linkButton = (): JSX.Element => {
        let link: string;

        switch (props.prop) {
        case 'group_id':
            link = '/groups/group/' + props.value;
            break;
        case 'host_id':
            link = '/hosts/host/' + props.value;
            break;
        case 'schedule_id':
            link = '/schedules/schedule/' + props.value;
            break;
        case 'script_id':
            link = '/scripts/script/' + props.value;
            break;
        }

        if (!link) {
            return null;
        }

        return (<span className="ms-1"><a href={link} rel="noreferrer" target="_blank"><Icon.ExternalLinkAlt /></a></span>);
    };

    if (itemLabel && itemValue) {
        return (<ListGroup.TextItem title={itemLabel}>
            <span>{itemValue}</span>
            { linkButton() }
        </ListGroup.TextItem>);
    }

    const spinner = loading ? (<Icon.Spinner pulse />) : null;
    return (<ListGroup.TextItem title={props.prop}>
        <code>{props.value}</code>
        { spinner }
        { linkButton() }
    </ListGroup.TextItem>);
};
