import * as React from 'react';
import { Loading } from '../../components/Loading';
import { Card } from '../../components/Card';
import { Host, HostType } from '../../types/Host';
import { ProgressBar } from '../../components/ProgressBar';
import { ScriptRun } from '../../types/Result';
import { RunOutput, RunResults } from './RunResults';
import { Style } from '../../components/Style';
import { ScriptRequest } from '../../services/ScriptRequest';

interface RunScriptProps {
    hostID: string;
    scriptID: string;
    onFinished: (results?: ScriptRun) => (void);
}
export const RunScript: React.FC<RunScriptProps> = (props: RunScriptProps) => {
    const [loadingHost, setLoadingHost] = React.useState<boolean>(true);
    const [runningScript, setRunningScript] = React.useState<boolean>(false);
    const [host, setHost] = React.useState<HostType>();
    const [results, setResults] = React.useState<ScriptRun>();
    const [stdout, setStdout] = React.useState<string>();
    const [stderr, setStderr] = React.useState<string>();
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    const [scriptConnection, setScriptConnection] = React.useState<ScriptRequest>(new ScriptRequest(props.scriptID, props.hostID));

    React.useEffect(() => {
        loadHost();
    }, []);

    React.useEffect(() => {
        if (runningScript) {
            startScript();
        }
    }, [runningScript]);

    React.useEffect(() => {
        if (results) {
            props.onFinished(results);
        }
    }, [results]);

    const loadHost = () => {
        Host.Get(props.hostID).then(host => {
            setHost(host);
            setLoadingHost(false);
            setRunningScript(true);
        });
    };

    const startScript = () => {
        scriptConnection.Stream((stdout: string, stderr: string) => {
            setStdout(stdout);
            setStderr(stderr);
        }).then(results => {
            setResults(results);
            setRunningScript(false);
            setStdout(undefined);
            setStderr(undefined);
        }, error => {
            setRunningScript(false);
            setResults({
                Result: {
                    success: false,
                },
                RunError: error,
            });
            props.onFinished();
        });
    };

    const cancelClick = () => {
        scriptConnection.Cancel();
    };

    const content = () => {
        if (!runningScript && results != undefined) {
            return (<RunResults results={results} />);
        }

        if (stdout || stderr) {
            return (
                <Card.Body>
                    <ProgressBar intermediate cancelClick={cancelClick} />
                    <RunOutput stdout={stdout} stderr={stderr} />
                </Card.Body>
            );
        }

        return (
            <Card.Body>
                <ProgressBar intermediate />
            </Card.Body>
        );
    };

    if (loadingHost) {
        return (<Loading />);
    }

    let color: Style.Palette;
    if (results && results.Result && !results.Result.success) {
        color = Style.Palette.Danger;
    }

    return (
        <Card.Card color={color}>
            <Card.Header>{host.Name}</Card.Header>
            {content()}
        </Card.Card>
    );
};
