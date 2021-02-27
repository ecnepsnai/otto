import * as React from 'react';
import { ScriptRun } from '../../types/Result';
import { EnvironmentVariableCard } from '../../components/EnvironmentVariableCard';
import { Card } from '../../components/Card';
import { Pre } from '../../components/Pre';
import { ListGroup } from '../../components/ListGroup';
import { Icon } from '../../components/Icon';
import { Style } from '../../components/Style';
import { Formatter } from '../../services/Formatter';

interface RunResultsProps {
    results: ScriptRun;
}
export class RunResults extends React.Component<RunResultsProps, unknown> {
    private error = () => {
        const errorMessage = this.props.results.RunError || this.props.results.Result.ExecError || 'Unknown Error';

        return (
            <Card.Body>
                <h4>Error Running Script</h4>
                <strong>Details</strong>
                <Pre>{errorMessage}</Pre>
            </Card.Body>
        );
    }

    render(): JSX.Element {
        if (this.props.results.RunError || this.props.results.Result.ExecError) {
            return this.error();
        }

        let returnCodeIcon = (<Icon.CheckCircle color={Style.Palette.Success} />);
        if (this.props.results.Result.Code !== 0) {
            returnCodeIcon = (<Icon.ExclamationCircle color={Style.Palette.Danger} />);
        }

        return (
            <Card.Body>
                <Card.Card>
                    <Card.Header>Details</Card.Header>
                    <ListGroup.List>
                        <ListGroup.TextItem title="Return Code">{this.props.results.Result.Code} {returnCodeIcon}</ListGroup.TextItem>
                        <ListGroup.TextItem title="Duration">{Formatter.Duration(this.props.results.Duration)}</ListGroup.TextItem>
                    </ListGroup.List>
                </Card.Card>
                <EnvironmentVariableCard variables={this.props.results.Environment} />
                <RunOutput stdout={this.props.results.Result.Stdout} stderr={this.props.results.Result.Stderr} />
            </Card.Body>
        );
    }
}

interface RunOutputProps {
    stdout: string;
    stderr: string;
}
export class RunOutput extends React.Component<RunOutputProps, unknown> {
    private content = () => {
        if (!this.props.stdout && !this.props.stderr) {
            return (<Card.Body><em className="text-muted">Script produced no output</em></Card.Body>);
        }

        let stdout: JSX.Element;
        if (this.props.stdout) {
            stdout = (<ListGroup.TextItem title="stdout"><Pre>{this.props.stdout}</Pre></ListGroup.TextItem>);
        }
        let stderr: JSX.Element;
        if (this.props.stderr) {
            stderr = (<ListGroup.TextItem title="stderr"><Pre>{this.props.stderr}</Pre></ListGroup.TextItem>);
        }

        return (<ListGroup.List>
            {stdout}
            {stderr}
        </ListGroup.List>);
    }

    render(): JSX.Element {
        return (
            <Card.Card>
                <Card.Header>Output</Card.Header>
                { this.content() }
            </Card.Card>
        );
    }
}
