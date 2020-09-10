import { API } from "../services/API";
import { Modal } from "../components/Modal";
import { Notification } from "../components/Notification";
import { Host } from "./Host";
import { Script } from "./Script";
import { Variable } from "./Variable";
import { Schedule } from "./Schedule";

export class Group {
    ID: string;
    Name: string;
    ScriptIDs: string[];
    Environment: Variable[];

    constructor(json: any) {
        this.ID = json.ID as string;
        this.Name = json.Name as string;
        this.ScriptIDs = (json.ScriptIDs || []) as string[];
        this.Environment = (json.Environment || []) as Variable[];
    }


    /**
     * Return a blank group
     */
    public static Blank(): Group {
        return new Group({
            Name: '',
            ScriptIDs: [],
            Environment: [],
        });
    }

    /**
     * Create a new Group
     */
    public static async New(parameters: NewGroupParameters): Promise<Group> {
        const data = await API.PUT('/api/groups/group', parameters);
        return new Group(data);
    }

    /**
     * Save this group
     */
    public async Save(): Promise<Group> {
        const data = await API.POST('/api/groups/group/' + this.ID, this as EditGroupParameters);
        return new Group(data);
    }

    /**
     * Delete this group
     */
    public async Delete(): Promise<any> {
        return await API.DELETE('/api/groups/group/' + this.ID);
    }

    /**
     * Modify the group changing the properties specified
     * @param properties properties to change
     */
    public async Update(properties: {[key:string]: any}): Promise<Group> {
        const data = await API.PATCH('/api/groups/group/' + this.ID, properties);
        return new Group(data);
    }

    /**
     * Show a modal to delete this group
     */
    public async DeleteModal(): Promise<boolean> {
        return new Promise(resolve => {
            Modal.delete('Delete Group?', 'Are you sure you want to delete this group? This can not be undone.').then(confirmed => {
                if (!confirmed) {
                    resolve(false);
                    return;
                }

                API.DELETE('/api/groups/group/' + this.ID).then(() => {
                    Notification.success('Group Deleted', this.Name);
                    resolve(true);
                });
            });
        });
    }

    /**
     * Get the specified group by its id
     */
    public static async Get(id: string): Promise<Group> {
        const data = await API.GET('/api/groups/group/' + id);
        return new Group(data);
    }

    /**
     * List all groups
     */
    public static async List(): Promise<Group[]> {
        const data = await API.GET('/api/groups');
        return (data as any[]).map(obj => {
            return new Group(obj);
        });
    }

    /**
     * Return a mapping of group ID to an array of host IDs
     */
    public static async Membership(): Promise<{[id: string]: string[]}> {
        const data = await API.GET('/api/groups/membership');
        return (data as {[id: string]: string[]});
    }

    /**
     * List all hosts for this group
     */
    public async Hosts(): Promise<Host[]> {
        return Group.Hosts(this.ID);
    }

    /**
     * List all hosts for a group
     */
    public static async Hosts(groupID: string): Promise<Host[]> {
        const data = await API.GET('/api/groups/group/' + groupID + '/hosts');
        return (data as any[]).map(obj => {
            return new Host(obj);
        });
    }

    /**
     * List all scripts for this group
     */
    public async Scripts(): Promise<Script[]> {
        return Group.Scripts(this.ID);
    }

    /**
     * List all scripts for a group
     */
    public static async Scripts(groupID: string): Promise<Script[]> {
        const data = await API.GET('/api/groups/group/' + groupID + '/scripts');
        return (data as any[]).map(obj => {
            return new Script(obj);
        });
    }

    /**
     * Set the hosts for this group
     * @param hostIDs array of host IDs
     */
    public async SetHosts(hostIDs: string[]): Promise<Host[]> {
        const data = await API.POST('/api/groups/group/' + this.ID + '/hosts', { Hosts: hostIDs });
        return (data as any[]).map(obj => {
            return new Host(obj);
        });
    }

    /**
     * List all schedules for this group
     */
    public async Schedules(): Promise<Schedule[]> {
        return Group.Schedules(this.ID);
    }

    /**
     * List all schedules for a group
     */
    public static async Schedules(groupID: string): Promise<Schedule[]> {
        const data = await API.GET('/api/groups/group/' + groupID + '/schedules');
        return (data as any[]).map(obj => {
            return new Schedule(obj);
        });
    }
}

export interface NewGroupParameters {
    Name: string;
    ScriptIDs: string[];
    Environment: Variable[];
}

export interface EditGroupParameters {
    Name: string;
    ScriptIDs: string[];
    Environment: Variable[];
}
