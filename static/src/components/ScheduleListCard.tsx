import * as React from 'react';
import { Link } from 'react-router-dom';
import { Schedule } from '../types/Schedule';
import { Card } from './Card';
import { Icon } from './Icon';
import { ListGroup } from './ListGroup';
import { Nothing } from './Nothing';

export interface ScheduleListCardProps {
    schedules: Schedule[];
    className?: string;
}
interface ScheduleListCardState {}
export class ScheduleListCard extends React.Component<ScheduleListCardProps, ScheduleListCardState> {
    constructor(props: ScheduleListCardProps) {
        super(props);
        this.state = { };
    }

    private content = () => {
        if (!this.props.schedules || this.props.schedules.length == 0) { return (<Card.Body><Nothing /></Card.Body>); }

        return (<ListGroup.List>
            {
                this.props.schedules.map((schedule, index) => {
                    return (
                    <ListGroup.Item key={index}>
                        <Icon.Calendar />
                        <Link to={'/schedules/schedule/' + schedule.ID} className="ml-1">{ schedule.Name }</Link>
                    </ListGroup.Item>
                    );
                })
            }
        </ListGroup.List>);
    }

    render(): JSX.Element {
        return (
            <Card.Card className={this.props.className}>
                <Card.Header>Schedules</Card.Header>
                { this.content() }
            </Card.Card>
        );
    }
}
