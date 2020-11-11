import * as React from 'react';
import { Card } from './Card';
import { ListGroup } from './ListGroup';
import { Nothing } from './Nothing';
import { Variable } from '../types/Variable';

export interface EnvironmentVariableCardProps {
    variables: Variable[];
    className?: string;
}
export class EnvironmentVariableCard extends React.Component<EnvironmentVariableCardProps, {}> {
    private list = () => {
        return (
        <ListGroup.List>
            {
                this.props.variables.map((variable, index) => {
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
    }
    private nothing = () => {
        return (
        <Card.Body>
            <Nothing />
        </Card.Body>
        );
    }
    private content = () => {
        if (Object.keys(this.props.variables).length == 0) {
            return this.nothing();
        }
        return this.list();
    }

    render(): JSX.Element {
        return (
        <Card.Card className={this.props.className}>
            <Card.Header>Environment Variables</Card.Header>
            { this.content() }
        </Card.Card>
        );
    }
}
