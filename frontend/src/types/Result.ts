import { Variable } from './Variable';

export interface ScriptRun {
    ScriptID?: string;
    Duration?: number;
    Environment?: Variable[];
    Result?: ScriptResultDetails;
    Output?: ScriptOutput;
    RunError?: string;
}

export interface ScriptResultDetails {
    success?: boolean;
    exec_error?: string;
    code?: number;
    stdout_len?: string;
    stderr_len?: string;
    duration?: number;
}

export interface ScriptOutput {
    Stdout?: string;
    Stderr?: string;
}
