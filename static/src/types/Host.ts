import { API } from "../services/API";
import { Modal } from "../components/Modal";
import { Notification } from "../components/Notification";
import { Group } from "./Group";

export class Host {
    ID: string;
    Name: string;
    Address: string;
    Port: number;
    PSK: string;
    Enabled: boolean;
    GroupIDs: string[];
    Environment: {[id: string]: string};

    constructor(json: any) {
        this.ID = json.ID as string;
        this.Name = json.Name as string;
        this.Address = json.Address as string;
        this.Port = json.Port as number;
        this.PSK = json.PSK as string;
        this.Enabled = json.Enabled as boolean;
        this.GroupIDs = (json.GroupIDs || []) as string[];
        this.Environment = (json.Environment || {}) as {[id: string]: string};
    }

    /**
     * Return a blank host
     */
    public static Blank(): Host {
        return new Host({
            Name: '',
            Address: '',
            Port: 12444,
            PSK: '',
            Enabled: '',
            GroupIDs: [],
            Environment: {},
        });
    }

    /**
     * Create a new Host
     */
    public static async New(parameters: NewHostParameters): Promise<Host> {
        const data = await API.PUT('/api/hosts/host', parameters);
        return new Host(data);
    }

    /**
     * Save this host
     */
    public async Save(): Promise<Host> {
        const data = await API.POST('/api/hosts/host/' + this.ID, this as EditHostParameters);
        return new Host(data);
    }

    /**
     * Delete this host
     */
    public async Delete(): Promise<any> {
        return await API.DELETE('/api/hosts/host/' + this.ID);
    }

    /**
     * Modify the host changing the properties specified
     * @param properties properties to change
     */
    public async Update(properties: {[key:string]: any}): Promise<Host> {
        const data = await API.PATCH('/api/hosts/host/' + this.ID, properties);
        return new Host(data);
    }

    /**
     * Show a modal to delete this host
     */
    public async DeleteModal(): Promise<boolean> {
        return new Promise(resolve => {
            Modal.delete('Delete Host?', 'Are you sure you want to delete this host? This can not be undone.').then(confirmed => {
                if (!confirmed) {
                    resolve(false);
                    return;
                }

                API.DELETE('/api/hosts/host/' + this.ID).then(() => {
                    Notification.success('Host Deleted', this.Name);
                    resolve(true);
                });
            });
        });
    }

    /**
     * Get the specified host by its id
     */
    public static async Get(id: string): Promise<Host> {
        const data = await API.GET('/api/hosts/host/' + id);
        return new Host(data);
    }

    /**
     * List all hosts
     */
    public static async List(): Promise<Host[]> {
        const data = await API.GET('/api/hosts');
        return (data as any[]).map(obj => {
            return new Host(obj);
        });
    }

    /**
     * List all scripts for this host
     */
    public async Scripts(): Promise<ScriptEnabledGroup[]> {
        return Host.Scripts(this.ID);
    }

    /**
     * List all scripts for a host
     */
    public static async Scripts(id: string): Promise<ScriptEnabledGroup[]> {
        const data = await API.GET('/api/hosts/host/' + id + '/scripts');
        return (data as ScriptEnabledGroup[]);
    }

    /**
     * List all groups for this host
     */
    public async Groups(): Promise<Group[]> {
        return Host.Groups(this.ID);
    }

    /**
     * List all groups for a host
     */
    public static async Groups(id: string): Promise<Group[]> {
        const data = await API.GET('/api/hosts/host/' + id + '/groups');
        return (data as any[]).map(obj => {
            return new Group(obj);
        });
    }
}

export interface NewHostParameters {
    Name: string;
    Address: string;
    Port: number;
    PSK: string;
    GroupIDs: string[];
    Environment: {[id: string]: string};
}

export interface EditHostParameters {
    Name: string;
    Address: string;
    Port: number;
    PSK: string;
    GroupIDs: string[];
    Enabled: boolean;
    Environment: {[id: string]: string};
}

export interface ScriptEnabledGroup {
    ScriptID: string;
    ScriptName: string;
    GroupID: string;
    GroupName: string;
}
