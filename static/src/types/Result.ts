import { API } from "../services/API";

export interface ScriptRun {
    ScriptID: string;
    Duration: number;
    Environment: {[id: string]: string};
    Result: ScriptResultDetails;
}

export interface ScriptResultDetails {
    Success: boolean;
    ExecError: string;
    Code: number;
    Stdout: string;
    Stderr: string;
}

export class ScriptRequest {
    public static async Run(scriptID: string, hostID: string): Promise<ScriptRun> {
        const results = await API.PUT('/api/request', {
            HostID: hostID,
            Action: 'run_script',
            ScriptID: scriptID,
        });
        return results as ScriptRun;
    }
}
