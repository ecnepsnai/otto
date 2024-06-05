import { API } from '../services/API';
import { Modal } from '../components/Modal';
import { Notification } from '../components/Notification';
import { GroupType } from './Group';
import { Variable } from './Variable';
import { ScheduleType } from './Schedule';
import { AttachmentType } from './Attachment';
import { ScriptRunLevel } from './gengo_enum';

export interface ScriptType {
    ID?: string;
    Name?: string;
    Executable?: string;
    Script?: string;
    Environment?: Variable[];
    RunAs?: RunAs;
    WorkingDirectory?: string;
    AfterExecution?: string;
    AttachmentIDs?: string[];
    RunLevel: ScriptRunLevel;
}

export class Script {
    /**
     * Return a blank script
     */
    public static Blank(): ScriptType {
        return {
            Name: '',
            Executable: '/bin/bash',
            Script: '',
            RunAs: {
                UID: 0,
                GID: 0,
                Inherit: true,
            },
            WorkingDirectory: '',
            AfterExecution: '',
            Environment: [],
            AttachmentIDs: [],
            RunLevel: ScriptRunLevel.ReadOnly,
        };
    }

    /**
     * Create a new Script
     */
    public static async New(parameters: ScriptType | NewScriptParameters): Promise<ScriptType> {
        const data = await API.PUT('/api/scripts/script', parameters);
        return data as ScriptType;
    }

    /**
     * Save this script
     */
    public static async Save(script: ScriptType): Promise<ScriptType> {
        const data = await API.POST('/api/scripts/script/' + script.ID, script);
        return data as ScriptType;
    }

    /**
     * Delete this script
     */
    public static async Delete(script: ScriptType): Promise<unknown> {
        return await API.DELETE('/api/scripts/script/' + script.ID);
    }

    /**
     * Modify the script changing the properties specified
     * @param properties properties to change
     */
    public static async Update(script: ScriptType, properties: { [key: string]: unknown }): Promise<ScriptType> {
        const data = await API.PATCH('/api/scripts/script/' + script.ID, properties);
        return data as ScriptType;
    }

    /**
     * Show a modal to delete this script
     */
    public static async DeleteModal(script: ScriptType): Promise<boolean> {
        return new Promise(resolve => {
            Modal.delete('Delete Script?', 'Are you sure you want to delete this script? This can not be undone.').then(confirmed => {
                if (!confirmed) {
                    resolve(false);
                    return;
                }

                API.DELETE('/api/scripts/script/' + script.ID).then(() => {
                    Notification.success('Script Deleted');
                    resolve(true);
                });
            });
        });
    }

    /**
     * Get the specified script by its id
     */
    public static async Get(id: string): Promise<ScriptType> {
        const data = await API.GET('/api/scripts/script/' + id);
        return data as ScriptType;
    }

    /**
     * List all scripts
     */
    public static async List(): Promise<ScriptType[]> {
        const data = await API.GET('/api/scripts');
        return data as ScriptType[];
    }

    /**
     * List all hosts for a script
     */
    public static async Hosts(scriptID: string): Promise<ScriptEnabledHost[]> {
        const data = await API.GET('/api/scripts/script/' + scriptID + '/hosts');
        return data as ScriptEnabledHost[];
    }

    /**
     * List all group for a script
     */
    public static async Groups(scriptID: string): Promise<GroupType[]> {
        const data = await API.GET('/api/scripts/script/' + scriptID + '/groups');
        return data as GroupType[];
    }

    /**
     * Set the groups for this script
     * @param groupIDs array of group IDs
     */
    public static async SetGroups(scriptID: string, groupIDs: string[]): Promise<ScriptEnabledHost[]> {
        const data = await API.POST('/api/scripts/script/' + scriptID + '/groups', { Groups: groupIDs });
        return data as ScriptEnabledHost[];
    }

    /**
     * List all schedules for a script
     */
    public static async Schedules(id: string): Promise<ScheduleType[]> {
        const data = await API.GET('/api/scripts/script/' + id + '/schedules');
        return data as ScheduleType[];
    }

    /**
     * List all attachments for a script
     */
    public static async Attachments(id: string): Promise<AttachmentType[]> {
        const data = await API.GET('/api/scripts/script/' + id + '/attachments');
        return data as AttachmentType[];
    }

    /**
     * Cancel a running script
     */
    public static async Cancel(hostID: string, scriptID: string): Promise<boolean> {
        const data = await API.POST('/api/action/cancel', {
            HostID: hostID,
            ScriptID: scriptID,
        });
        return data as boolean;
    }
}

export interface RunAs {
    UID: number;
    GID: number;
    Inherit: boolean;
}

export interface ScriptEnabledHost {
    ScriptID: string;
    ScriptName: string;
    GroupID: string;
    GroupName: string;
    HostID: string;
    HostName: string;
}

export interface NewScriptParameters {
    Name: string;
    Executable: string;
    Script: string;
    Environment: Variable[];
    RunAs: RunAs;
    WorkingDirectory: string;
    AfterExecution: string;
    AttachmentIDs: string[];
}

export interface EditScriptParameters {
    Name: string;
    Enabled: boolean;
    Executable: string;
    Script: string;
    Environment: Variable[];
    RunAs: RunAs;
    WorkingDirectory: string;
    AfterExecution: string;
    AttachmentIDs: string[];
}
