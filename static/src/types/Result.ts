import { Variable } from './Variable';

export interface ScriptRun {
    ScriptID?: string;
    Duration?: number;
    Environment?: Variable[];
    Result?: ScriptResultDetails;
    RunError?: string;
}

export interface ScriptResultDetails {
    success?: boolean;
    exec_error?: string;
    code?: number;
    stdout?: string;
    stderr?: string;
}
