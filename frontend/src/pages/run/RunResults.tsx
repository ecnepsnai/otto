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
export const RunResults: React.FC<RunResultsProps> = (props: RunResultsProps) => {
    const error = () => {
        const errorMessage = props.results.RunError || props.results.Result.exec_error || 'Unknown Error';

        return (
            <Card.Body>
                <h4>Error Running Script</h4>
                <strong>Details</strong>
                <Pre>{errorMessage}</Pre>
            </Card.Body>
        );
    };

    if (props.results.RunError || props.results.Result.exec_error) {
        return error();
    }

    let returnCodeIcon = (<Icon.CheckCircle color={Style.Palette.Success} />);
    if (props.results.Result.code !== 0) {
        returnCodeIcon = (<Icon.ExclamationCircle color={Style.Palette.Danger} />);
    }

    return (
        <Card.Body>
            <Card.Card>
                <Card.Header>Details</Card.Header>
                <ListGroup.List>
                    <ListGroup.TextItem title="Return Code">{props.results.Result.code} {returnCodeIcon}</ListGroup.TextItem>
                    <ListGroup.TextItem title="Duration">{Formatter.DurationNS(props.results.Duration)}</ListGroup.TextItem>
                </ListGroup.List>
            </Card.Card>
            <EnvironmentVariableCard variables={props.results.Environment} />
            <RunOutput stdout={props.results.Output.Stdout} stderr={props.results.Output.Stderr} />
        </Card.Body>
    );
};

interface RunOutputProps {
    stdout: string;
    stderr: string;
}
export const RunOutput: React.FC<RunOutputProps> = (props: RunOutputProps) => {
    const content = () => {
        if (!props.stdout && !props.stderr) {
            return (<Card.Body><em className="text-muted">Script produced no output</em></Card.Body>);
        }

        let stdout: JSX.Element;
        if (props.stdout) {
            stdout = (<ListGroup.TextItem title="stdout"><Pre>{props.stdout}</Pre></ListGroup.TextItem>);
        }
        let stderr: JSX.Element;
        if (props.stderr) {
            stderr = (<ListGroup.TextItem title="stderr"><Pre>{props.stderr}</Pre></ListGroup.TextItem>);
        }

        return (<ListGroup.List>
            {stdout}
            {stderr}
        </ListGroup.List>);
    };

    return (
        <Card.Card>
            <Card.Header>Output</Card.Header>
            {content()}
        </Card.Card>
    );
};
