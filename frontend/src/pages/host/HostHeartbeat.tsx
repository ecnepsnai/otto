import * as React from 'react';
import { Host, HostType } from '../../types/Host';
import { ListGroup } from '../../components/ListGroup';
import { Dropdown, Menu } from '../../components/Menu';
import { Icon } from '../../components/Icon';
import { HeartbeatType } from '../../types/Heartbeat';
import { HeartbeatBadge } from '../../components/Badge';
import { Card } from '../../components/Card';
import { ClientVersion } from '../../components/ClientVersion';
import { DateLabel } from '../../components/DateLabel';

interface HostHeartbeatProps {
    host: HostType;
    defaultHeartbeat?: HeartbeatType;
    didUpdate?: (heartbeat: HeartbeatType) => void;
}
export const HostHeartbeat: React.FC<HostHeartbeatProps> = (props: HostHeartbeatProps) => {
    const [Heartbeat, setHeartbeat] = React.useState(props.defaultHeartbeat);
    const [IsLoading, setIsLoading] = React.useState(false);

    const triggerClick = () => {
        setIsLoading(true);
        Host.Heartbeat(props.host.ID).then(heartbeat => {
            setIsLoading(false);
            setHeartbeat(heartbeat);
            if (props.didUpdate) {
                props.didUpdate(heartbeat);
            }
        }, () => {
            setIsLoading(false);
        });
    };

    const lastReply = (): JSX.Element => {
        if (!Heartbeat) {
            return null;
        }

        return (<ListGroup.TextItem title="Last Heartbeat"><DateLabel date={Heartbeat.LastReply} /></ListGroup.TextItem>);
    };

    const clientVersion = (): JSX.Element => {
        if (!Heartbeat) {
            return null;
        }

        return (<ListGroup.TextItem title="Client Version"><ClientVersion heartbeat={Heartbeat} /></ListGroup.TextItem>);
    };

    const hostProperties = (): JSX.Element => {
        if (!Heartbeat || !Heartbeat.Properties) {
            return null;
        }

        return (
            <React.Fragment>
                {Object.keys(Heartbeat.Properties).map((key, idx) => {
                    return (
                        <ListGroup.TextItem title={key} key={idx}><code>{Heartbeat.Properties[key]}</code></ListGroup.TextItem>
                    );
                })}
            </React.Fragment>
        );
    };

    const menu = () => {
        if (IsLoading) {
            return (<Icon.Spinner pulse />);
        }

        return (<Dropdown label={<Icon.Bars />}>
            <Menu.Item icon={<Icon.QuestionCircle />} label="Check Now" onClick={triggerClick} />
        </Dropdown>);
    };

    return (
        <Card.Card className="mb-3">
            <Card.Header>Otto Client Information</Card.Header>
            <ListGroup.List>
                <ListGroup.TextItem title="Status">
                    <HeartbeatBadge heartbeat={Heartbeat} />
                    <span className="ms-2">
                        {menu()}
                    </span>
                </ListGroup.TextItem>
                {lastReply()}
                {clientVersion()}
                {hostProperties()}
            </ListGroup.List>
        </Card.Card>
    );
};
