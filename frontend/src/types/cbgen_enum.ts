// This file is was generated automatically by Codegen v1.11.0
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

export enum IPVersionOption { 
    /** IPv4 or IPv6 as chosen by the system automatically */
    Auto = 'auto',
    /** IPv4 only */
    IPv4 = 'ipv4',
    /** IPv6 only */
    IPv6 = 'ipv6',
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

export enum RequestResponseCode { 
    Output = 100,
    Keepalive = 101,
    Error = 400,
    Finished = 200,
}

export enum ScheduleResult { 
    /** All hosts executed the script successfully */
    Success = 0,
    /** Some hosts did not execute the script successfully */
    PartialSuccess = 1,
    /** No hosts executed the script successfully */
    Fail = 2,
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

