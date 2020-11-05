import { API } from "../services/API";
import { Modal } from "../components/Modal";
import { Notification } from "../components/Notification";
import { Group } from "./Group";
import { Variable } from "./Variable";
import { Schedule } from "./Schedule";
import { Attachment } from "./Attachment";

export class Script {
    ID: string;
    Name: string;
    Enabled: boolean;
    Executable: string;
    Script: string;
    Environment: Variable[];
    UID: number;
    GID: number;
    WorkingDirectory: string;
    AfterExecution: string;
    AttachmentIDs: string[];

    constructor(json: any) {
        this.ID = json.ID as string;
        this.Name = json.Name as string;
        this.Enabled = json.Enabled as boolean;
        this.Executable = json.Executable as string;
        this.Script = json.Script as string;
        this.Environment = (json.Environment || []) as Variable[];
        this.UID = json.UID as number;
        this.GID = json.GID as number;
        this.WorkingDirectory = json.WorkingDirectory as string;
        this.AfterExecution = json.AfterExecution as string;
        this.AttachmentIDs = (json.AttachmentIDs || []) as string[];
    }

    /**
     * Return a blank script
     */
    public static Blank(): Script {
        return new Script({
            Name: '',
            Address: '',
            Port: 12444,
            PSK: '',
            Enabled: '',
            GroupIDs: [],
            Environment: [],
            AttachmentIDs: [],
        });
    }

    /**
     * Create a new Script
     */
    public static async New(parameters: NewScriptParameters): Promise<Script> {
        const data = await API.PUT('/api/scripts/script', parameters);
        return new Script(data);
    }

    /**
     * Save this script
     */
    public async Save(): Promise<Script> {
        const data = await API.POST('/api/scripts/script/' + this.ID, this as EditScriptParameters);
        return new Script(data);
    }

    /**
     * Delete this script
     */
    public async Delete(): Promise<any> {
        return await API.DELETE('/api/scripts/script/' + this.ID);
    }

    /**
     * Modify the script changing the properties specified
     * @param properties properties to change
     */
    public async Update(properties: {[key:string]: any}): Promise<Script> {
        const data = await API.PATCH('/api/scripts/script/' + this.ID, properties);
        return new Script(data);
    }

    /**
     * Show a modal to delete this script
     */
    public async DeleteModal(): Promise<boolean> {
        return new Promise(resolve => {
            Modal.delete('Delete Script?', 'Are you sure you want to delete this script? This can not be undone.').then(confirmed => {
                if (!confirmed) {
                    resolve(false);
                    return;
                }

                API.DELETE('/api/scripts/script/' + this.ID).then(() => {
                    Notification.success('Script Deleted');
                    resolve(true);
                });
            });
        });
    }

    /**
     * Get the specified script by its id
     */
    public static async Get(id: string): Promise<Script> {
        const data = await API.GET('/api/scripts/script/' + id);
        return new Script(data);
    }

    /**
     * List all scripts
     */
    public static async List(): Promise<Script[]> {
        const data = await API.GET('/api/scripts');
        return (data as any[]).map(obj => {
            return new Script(obj);
        });
    }

    /**
     * List all hosts for this script
     */
    public async Hosts(): Promise<ScriptEnabledHost[]> {
        return Script.Hosts(this.ID);
    }

    /**
     * List all group for this script
     */
    public async Groups(): Promise<Group[]> {
        return Script.Groups(this.ID);
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
    public static async Groups(scriptID: string): Promise<Group[]> {
        const data = await API.GET('/api/scripts/script/' + scriptID + '/groups');
        return (data as any[]).map(obj => {
            return new Group(obj);
        });
    }

    /**
     * Set the groups for this script
     * @param groupIDs array of group IDs
     */
    public async SetGroups(groupIDs: string[]): Promise<ScriptEnabledHost[]> {
        const data = await API.POST('/api/scripts/script/' + this.ID + '/groups', { Groups: groupIDs });
        return data as ScriptEnabledHost[];
    }

    /**
     * List all schedules for this script
     */
    public async Schedules(): Promise<Schedule[]> {
        return Script.Schedules(this.ID);
    }

    /**
     * List all schedules for a script
     */
    public static async Schedules(id: string): Promise<Schedule[]> {
        const data = await API.GET('/api/scripts/script/' + id + '/schedules');
        return (data as any[]).map(obj => {
            return new Schedule(obj);
        });
    }

    /**
     * List all attachments for this script
     */
    public async Attachments(): Promise<Attachment[]> {
        return Script.Attachments(this.ID);
    }

    /**
     * List all attachments for a script
     */
    public static async Attachments(id: string): Promise<Attachment[]> {
        const data = await API.GET('/api/scripts/script/' + id + '/attachments');
        return (data as any[]).map(obj => {
            return new Attachment(obj);
        });
    }
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
    UID: number;
    GID: number;
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
    UID: number;
    GID: number;
    WorkingDirectory: string;
    AfterExecution: string;
    AttachmentIDs: string[];
}
