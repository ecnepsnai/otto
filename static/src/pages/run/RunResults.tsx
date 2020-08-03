import * as React from 'react';
import { ScriptRun } from '../../types/Result';
import { EnvironmentVariableCard } from '../../components/EnvironmentVariableCard';
import { Card } from '../../components/Card';

export interface RunResultsProps {
    results: ScriptRun;
}
interface RunResultsState {}
export class RunResults extends React.Component<RunResultsProps, RunResultsState> {
    constructor(props: RunResultsProps) {
        super(props);
        this.state = { };
    }

    render(): JSX.Element {
        return (
            <Card.Body>
                <h4>Results</h4>
                <strong>Return Code</strong> {this.props.results.Result.Code}<br/>
                <strong>Duration</strong> {this.props.results.Duration}<br/>
                <EnvironmentVariableCard variables={this.props.results.Environment} />
                <h4>Output</h4>
                <Card.Card>
                    <Card.Header>Standard Out (stdout)</Card.Header>
                    <Card.Body>
                        <pre>{this.props.results.Result.Stdout}</pre>
                    </Card.Body>
                </Card.Card>
                <Card.Card>
                    <Card.Header>Standard Error (stderr)</Card.Header>
                    <Card.Body>
                        <pre>{this.props.results.Result.Stderr}</pre>
                    </Card.Body>
                </Card.Card>
            </Card.Body>
        );
    }
}
