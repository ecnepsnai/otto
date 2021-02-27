import * as React from 'react';
import { Link } from 'react-router-dom';
import { GroupType } from '../types/Group';
import { Card } from './Card';
import { Icon } from './Icon';
import { ListGroup } from './ListGroup';
import { Nothing } from './Nothing';

export interface GroupListCardProps {
    groups: GroupType[];
    className?: string;
}
export const GroupListCard: React.FC<GroupListCardProps> = (props: GroupListCardProps) => {
    const content = () => {
        if (!props.groups || props.groups.length == 0) {
            return (<Card.Body><Nothing /></Card.Body>);
        }

        return (<ListGroup.List>
            {
                props.groups.map((group, index) => {
                    return (
                        <ListGroup.Item key={index}>
                            <Icon.LayerGroup />
                            <Link to={'/groups/group/' + group.ID} className="ms-1">{ group.Name }</Link>
                        </ListGroup.Item>
                    );
                })
            }
        </ListGroup.List>);
    };

    return (
        <Card.Card className={props.className}>
            <Card.Header>Groups</Card.Header>
            { content() }
        </Card.Card>
    );
};
