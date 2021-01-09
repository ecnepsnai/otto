import * as React from 'react';
import { Loading } from '../../components/Loading';
import { Card } from '../../components/Card';
import { Host } from '../../types/Host';
import { ProgressBar } from '../../components/ProgressBar';
import { ScriptRequest, ScriptRun } from '../../types/Result';
import { RunOutput, RunResults } from './RunResults';
import { Style } from '../../components/Style';

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
    stdout?: string;
    stderr?: string;
}
export class RunScript extends React.Component<RunScriptProps, RunScriptState> {
    private scriptConnection: ScriptRequest;

    constructor(props: RunScriptProps) {
        super(props);
        this.state = {
            loadingHost: true,
            runningScript: false,
        };
        this.scriptConnection = new ScriptRequest(props.scriptID, props.hostID);
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
        this.scriptConnection.Stream((stdout: string, stderr: string) => {
            this.setState({
                stdout: stdout,
                stderr: stderr,
            });
        }).then(results => {
            this.props.onFinished(results);
            this.setState({
                runningScript: false,
                results: results,
                stdout: undefined,
                stderr: undefined,
            });
        }, error => {
            this.setState({
                runningScript: false,
                results: {
                    Result: {
                        Success: false,
                    },
                    RunError: error,
                },
            });
            this.props.onFinished();
        });
    }

    private cancelClick = () => {
        this.scriptConnection.Cancel();
    }

    private content = () => {
        if (!this.state.runningScript) {
            return ( <RunResults results={this.state.results} /> );
        }

        if (this.state.stdout || this.state.stderr) {
            return (
                <Card.Body>
                    <ProgressBar intermediate cancelClick={this.cancelClick} />
                    <RunOutput stdout={this.state.stdout } stderr={this.state.stderr}/>
                </Card.Body>
            );
        }

        return (
            <Card.Body>
                <ProgressBar intermediate />
            </Card.Body>
        );
    };

    render(): JSX.Element {
        if (this.state.loadingHost) { return (<Loading />); }

        let color: Style.Palette;
        if (this.state.results && this.state.results.Result && !this.state.results.Result.Success) {
            color = Style.Palette.Danger;
        }

        return (
            <Card.Card color={color}>
                <Card.Header>{this.state.host.Name}</Card.Header>
                { this.content() }
            </Card.Card>
        );
    }
}
