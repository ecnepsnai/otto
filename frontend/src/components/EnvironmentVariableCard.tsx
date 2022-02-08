import * as React from 'react';
import { Card } from './Card';
import { ListGroup } from './ListGroup';
import { Nothing } from './Nothing';
import { Variable } from '../types/Variable';

interface EnvironmentVariableCardProps {
    variables: Variable[];
    className?: string;
}
export const EnvironmentVariableCard: React.FC<EnvironmentVariableCardProps> = (props: EnvironmentVariableCardProps) => {
    const list = () => {
        return (
            <ListGroup.List>
                {
                    (props.variables || []).map((variable, index) => {
                        const content = variable.Secret ? '******' : variable.Value;
                        return (
                            <ListGroup.TextItem title={variable.Key} key={index}>
                                <code>{content}</code>
                            </ListGroup.TextItem>
                        );
                    })
                }
            </ListGroup.List>
        );
    };

    const nothing = () => {
        return (
            <Card.Body>
                <Nothing />
            </Card.Body>
        );
    };

    const content = () => {
        if (Object.keys((props.variables || [])).length == 0) {
            return nothing();
        }
        return list();
    };

    return (
        <Card.Card className={props.className}>
            <Card.Header>Environment Variables</Card.Header>
            { content()}
        </Card.Card>
    );
};
