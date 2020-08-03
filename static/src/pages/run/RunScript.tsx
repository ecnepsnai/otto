import * as React from 'react';
import { Loading } from '../../components/Loading';
import { Card } from '../../components/Card';
import { Host } from '../../types/Host';
import { ProgressBar } from '../../components/ProgressBar';
import { ScriptRequest, ScriptRun } from '../../types/Result';
import { RunResults } from './RunResults';

export interface RunScriptProps {
    hostID: string;
    scriptID: string;
    onFinished: (results?: ScriptRun) => (void);
}
interface RunScriptState {
    loadingHost: boolean;
    runningScript: boolean;
    host?: Host;
    results?: ScriptRun;
}
export class RunScript extends React.Component<RunScriptProps, RunScriptState> {
    constructor(props: RunScriptProps) {
        super(props);
        this.state = {
            loadingHost: true,
            runningScript: false,
        };
    }

    private loadHost = () => {
        Host.Get(this.props.hostID).then(host => {
            this.setState({
                loadingHost: false,
                runningScript: true,
                host: host,
            }, () => {
                this.startScript();
            });
        });
    }

    componentDidMount(): void {
        this.loadHost();
    }

    private startScript = () => {
        ScriptRequest.Run(this.props.scriptID, this.props.hostID).then(results => {
            this.props.onFinished(results);
            this.setState({
                runningScript: false,
                results: results,
            });
        }, () => {
            this.props.onFinished();
        });
    }

    private content = () => {
        if (!this.state.runningScript) {
            return ( <RunResults results={this.state.results} /> );
        }

        return (
            <Card.Body>
                <ProgressBar intermediate />
            </Card.Body>
        );
    };

    render(): JSX.Element {
        if (this.state.loadingHost) { return (<Loading />); }
        return (
            <Card.Card>
                <Card.Header>{this.state.host.Name}</Card.Header>
                { this.content() }
            </Card.Card>
        );
    }
}
