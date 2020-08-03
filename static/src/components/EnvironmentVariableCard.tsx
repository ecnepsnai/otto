import * as React from 'react';
import { Card } from './Card';
import { ListGroup } from './ListGroup';
import { Nothing } from './Nothing';

export interface EnvironmentVariableCardProps {
    variables: {[id: string]: string};
}
interface vbtype {
    key: string;
    value: string;
}
export class EnvironmentVariableCard extends React.Component<EnvironmentVariableCardProps, {}> {
    private list = () => {
        const variables: vbtype[] = [];
        Object.keys(this.props.variables).forEach(key => {
            const value = this.props.variables[key];
            variables.push({ key: key, value: value });
        });

        return (
        <ListGroup.List>
            {
                variables.map((variable, index) => {
                    return (
                    <ListGroup.TextItem title={variable.key} key={index}>
                        <code>{ variable.value }</code>
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
        <Card.Card>
            <Card.Header>Environment Variables</Card.Header>
            { this.content() }
        </Card.Card>
        );
    }
}
