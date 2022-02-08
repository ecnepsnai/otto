import * as React from 'react';
import { Link } from 'react-router-dom';
import { ScheduleType } from '../types/Schedule';
import { Card } from './Card';
import { Icon } from './Icon';
import { ListGroup } from './ListGroup';
import { Nothing } from './Nothing';

interface ScheduleListCardProps {
    schedules: ScheduleType[];
    className?: string;
}
export const ScheduleListCard: React.FC<ScheduleListCardProps> = (props: ScheduleListCardProps) => {
    const content = () => {
        if (!props.schedules || props.schedules.length == 0) {
            return (<Card.Body><Nothing /></Card.Body>);
        }

        return (<ListGroup.List>
            {
                props.schedules.map((schedule, index) => {
                    return (
                        <ListGroup.Item key={index}>
                            <Icon.Calendar />
                            <Link to={'/schedules/schedule/' + schedule.ID} className="ms-1">{schedule.Name}</Link>
                        </ListGroup.Item>
                    );
                })
            }
        </ListGroup.List>);
    };

    return (
        <Card.Card className={props.className}>
            <Card.Header>Schedules</Card.Header>
            { content()}
        </Card.Card>
    );
};
