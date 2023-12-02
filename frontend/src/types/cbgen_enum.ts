// This file is was generated automatically by Codegen v1.12.2
// Do not make changes to this file as they will be lost

export enum AgentAction { 
    /** Ping the host */
    Ping = 'ping',
    /** Run the script on the host */
    RunScript = 'run_script',
    /** Reload the configuration of the agent */
    ReloadConfig = 'reload_config',
    /** Exit the agent on the host */
    ExitAgent = 'exit_agent',
    /** Reboot the host */
    Reboot = 'reboot',
    /** Power off the host */
    Shutdown = 'shutdown',
}

export function AgentActionAll() {
    return [ 
        AgentAction.Ping,
        AgentAction.RunScript,
        AgentAction.ReloadConfig,
        AgentAction.ExitAgent,
        AgentAction.Reboot,
        AgentAction.Shutdown,
    ];
}

export function AgentActionConfig() {
    return [
        {
            key: 'Ping',
            value: 'ping',
            description: 'Ping the host',
        },
        {
            key: 'RunScript',
            value: 'run_script',
            description: 'Run the script on the host',
        },
        {
            key: 'ReloadConfig',
            value: 'reload_config',
            description: 'Reload the configuration of the agent',
        },
        {
            key: 'ExitAgent',
            value: 'exit_agent',
            description: 'Exit the agent on the host',
        },
        {
            key: 'Reboot',
            value: 'reboot',
            description: 'Reboot the host',
        },
        {
            key: 'Shutdown',
            value: 'shutdown',
            description: 'Power off the host',
        },
    ];
}

export enum IPVersionOption { 
    /** IPv4 or IPv6 as chosen by the system automatically */
    Auto = 'auto',
    /** IPv4 only */
    IPv4 = 'ipv4',
    /** IPv6 only */
    IPv6 = 'ipv6',
}

export function IPVersionOptionAll() {
    return [ 
        IPVersionOption.Auto,
        IPVersionOption.IPv4,
        IPVersionOption.IPv6,
    ];
}

export function IPVersionOptionConfig() {
    return [
        {
            key: 'Auto',
            value: 'auto',
            description: 'IPv4 or IPv6 as chosen by the system automatically',
        },
        {
            key: 'IPv4',
            value: 'ipv4',
            description: 'IPv4 only',
        },
        {
            key: 'IPv6',
            value: 'ipv6',
            description: 'IPv6 only',
        },
    ];
}

export enum RegisterRuleProperty { 
    /** Hostname */
    Hostname = 'hostname',
    /** Kernel Name */
    KernelName = 'kernel_name',
    /** Kernel Version */
    KernelVersion = 'kernel_version',
    /** Distribution Name */
    DistributionName = 'distribution_name',
    /** Distribution Version */
    DistributionVersion = 'distribution_version',
}

export function RegisterRulePropertyAll() {
    return [ 
        RegisterRuleProperty.Hostname,
        RegisterRuleProperty.KernelName,
        RegisterRuleProperty.KernelVersion,
        RegisterRuleProperty.DistributionName,
        RegisterRuleProperty.DistributionVersion,
    ];
}

export function RegisterRulePropertyConfig() {
    return [
        {
            key: 'Hostname',
            value: 'hostname',
            description: 'Hostname',
        },
        {
            key: 'KernelName',
            value: 'kernel_name',
            description: 'Kernel Name',
        },
        {
            key: 'KernelVersion',
            value: 'kernel_version',
            description: 'Kernel Version',
        },
        {
            key: 'DistributionName',
            value: 'distribution_name',
            description: 'Distribution Name',
        },
        {
            key: 'DistributionVersion',
            value: 'distribution_version',
            description: 'Distribution Version',
        },
    ];
}

export enum RequestResponseCode { 
    Output = 100,
    Keepalive = 101,
    Error = 400,
    Finished = 200,
}

export function RequestResponseCodeAll() {
    return [ 
        RequestResponseCode.Output,
        RequestResponseCode.Keepalive,
        RequestResponseCode.Error,
        RequestResponseCode.Finished,
    ];
}

export function RequestResponseCodeConfig() {
    return [
        {
            key: 'Output',
            value: 100,
            
        },
        {
            key: 'Keepalive',
            value: 101,
            
        },
        {
            key: 'Error',
            value: 400,
            
        },
        {
            key: 'Finished',
            value: 200,
            
        },
    ];
}

export enum ScheduleResult { 
    /** All hosts executed the script successfully */
    Success = 0,
    /** Some hosts did not execute the script successfully */
    PartialSuccess = 1,
    /** No hosts executed the script successfully */
    Fail = 2,
}

export function ScheduleResultAll() {
    return [ 
        ScheduleResult.Success,
        ScheduleResult.PartialSuccess,
        ScheduleResult.Fail,
    ];
}

export function ScheduleResultConfig() {
    return [
        {
            key: 'Success',
            value: 0,
            description: 'All hosts executed the script successfully',
        },
        {
            key: 'PartialSuccess',
            value: 1,
            description: 'Some hosts did not execute the script successfully',
        },
        {
            key: 'Fail',
            value: 2,
            description: 'No hosts executed the script successfully',
        },
    ];
}

/** Permission level for users to run scripts */
export enum ScriptRunLevel { 
    /** No scripts can be executed */
    None = 0,
    /** Only scripts mark as read only can be executed */
    ReadOnly = 1,
    /** All scripts can be executed */
    ReadWrite = 2,
}

export function ScriptRunLevelAll() {
    return [ 
        ScriptRunLevel.None,
        ScriptRunLevel.ReadOnly,
        ScriptRunLevel.ReadWrite,
    ];
}

export function ScriptRunLevelConfig() {
    return [
        {
            key: 'None',
            value: 0,
            description: 'No scripts can be executed',
        },
        {
            key: 'ReadOnly',
            value: 1,
            description: 'Only scripts mark as read only can be executed',
        },
        {
            key: 'ReadWrite',
            value: 2,
            description: 'All scripts can be executed',
        },
    ];
}

