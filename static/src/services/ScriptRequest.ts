import { ScriptRun } from '../types/Result';

enum RequestResponseCode {
    Output = 100,
    Keepalive = 101,
    Error = 400,
    Finished = 200,
}

interface RequestResponse {
    Code: RequestResponseCode;
    Error?: string;
    Stdout?: string;
    Stderr?: string;
    Result?: ScriptRun;
}

export class ScriptRequest {
    private scriptID: string;
    private hostID: string;
    private socket: WebSocket;

    constructor(scriptID: string, hostID: string) {
        this.scriptID = scriptID;
        this.hostID = hostID;
    }

    public Stream(onOutput: (stdout: string, stderr: string) => (void)): Promise<ScriptRun> {
        return new Promise((resolve, reject) => {
            let protocol = 'wss:';
            if (window.location.protocol === 'http:') {
                protocol = 'ws:';
            }
            const url = protocol + '//' + window.location.host + '/api/action/async';
            this.socket = new WebSocket(url);

            this.socket.addEventListener('open', () => {
                this.socket.send(JSON.stringify({
                    HostID: this.hostID,
                    Action: 'run_script',
                    ScriptID: this.scriptID,
                }));
            });

            let result: ScriptRun;

            this.socket.addEventListener('message', message => {
                const response = JSON.parse(message.data) as RequestResponse;
                if (!response) {
                    return;
                }

                switch (response.Code) {
                    case RequestResponseCode.Keepalive:
                        break;
                    case RequestResponseCode.Error:
                        reject(response.Error);
                        break;
                    case RequestResponseCode.Output:
                        onOutput(response.Stdout, response.Stderr);
                        break;
                    case RequestResponseCode.Finished:
                        result = response.Result;
                        break;
                    default:
                        console.error('Unknown response from server', response);
                        break;
                }
            });

            this.socket.addEventListener('error', error => {
                console.error('ws error', error);
                reject(error);
                return;
            });

            this.socket.addEventListener('close', () => {
                if (result) {
                    resolve(result);
                }
            });
        });
    }

    public Cancel(): void {
        if (!this.socket) {
            return;
        }

        this.socket.send(JSON.stringify({ Cancel: true }));
    }
}
