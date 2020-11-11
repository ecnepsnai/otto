import * as React from 'react';
import { Link } from 'react-router-dom';
import { Group } from '../types/Group';
import { Card } from './Card';
import { Icon } from './Icon';
import { ListGroup } from './ListGroup';
import { Nothing } from './Nothing';

export interface GroupListCardProps {
    groups: Group[];
    className?: string;
}
interface GroupListCardState {}
export class GroupListCard extends React.Component<GroupListCardProps, GroupListCardState> {
    constructor(props: GroupListCardProps) {
        super(props);
        this.state = { };
    }

    private content = () => {
        if (!this.props.groups || this.props.groups.length == 0) { return (<Card.Body><Nothing /></Card.Body>); }

        return (<ListGroup.List>
            {
                this.props.groups.map((group, index) => {
                    return (
                    <ListGroup.Item key={index}>
                        <Icon.LayerGroup />
                        <Link to={'/groups/group/' + group.ID} className="ml-1">{ group.Name }</Link>
                    </ListGroup.Item>
                    );
                })
            }
        </ListGroup.List>);
    }

    render(): JSX.Element {
        return (
            <Card.Card className={this.props.className}>
                <Card.Header>Groups</Card.Header>
                { this.content() }
            </Card.Card>
        );
    }
}
