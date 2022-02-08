import { API } from '../services/API';
import { Modal } from '../components/Modal';
import { Notification } from '../components/Notification';
import { HostType } from './Host';
import { ScriptType } from './Script';
import { Variable } from './Variable';
import { ScheduleType } from './Schedule';

export interface GroupType {
    ID?: string;
    Name?: string;
    ScriptIDs?: string[];
    Environment?: Variable[];
}

export class Group {
    /**
     * Return a blank group
     */
    public static Blank(): GroupType {
        return {
            Name: '',
            ScriptIDs: [],
            Environment: [],
        };
    }

    /**
     * Create a new Group
     */
    public static async New(parameters: GroupType | NewGroupParameters): Promise<GroupType> {
        const data = await API.PUT('/api/groups/group', parameters);
        return data as GroupType;
    }

    /**
     * Save this group
     */
    public static async Save(group: GroupType): Promise<GroupType> {
        const data = await API.POST('/api/groups/group/' + group.ID, group);
        return data as GroupType;
    }

    /**
     * Delete this group
     */
    public static async Delete(group: GroupType): Promise<any> {
        return await API.DELETE('/api/groups/group/' + group.ID);
    }

    /**
     * Modify the group changing the properties specified
     * @param properties properties to change
     */
    public static async Update(group: GroupType, properties: { [key: string]: any }): Promise<GroupType> {
        const data = await API.PATCH('/api/groups/group/' + group.ID, properties);
        return data as GroupType;
    }

    /**
     * Show a modal to delete this group
     */
    public static async DeleteModal(group: GroupType): Promise<boolean> {
        return new Promise(resolve => {
            Modal.delete('Delete Group?', 'Are you sure you want to delete this group? This can not be undone.').then(confirmed => {
                if (!confirmed) {
                    resolve(false);
                    return;
                }

                API.DELETE('/api/groups/group/' + group.ID).then(() => {
                    Notification.success('Group Deleted');
                    resolve(true);
                });
            });
        });
    }

    /**
     * Get the specified group by its id
     */
    public static async Get(id: string): Promise<GroupType> {
        const data = await API.GET('/api/groups/group/' + id);
        return data as GroupType;
    }

    /**
     * List all groups
     */
    public static async List(): Promise<GroupType[]> {
        const data = await API.GET('/api/groups');
        return data as GroupType[];
    }

    /**
     * Return a mapping of group ID to an array of host IDs
     */
    public static async Membership(): Promise<{ [id: string]: string[] }> {
        const data = await API.GET('/api/groups/membership');
        return (data as { [id: string]: string[] });
    }

    /**
     * List all hosts for a group
     */
    public static async Hosts(groupID: string): Promise<HostType[]> {
        const data = await API.GET('/api/groups/group/' + groupID + '/hosts');
        return data as HostType[];
    }

    /**
     * List all scripts for a group
     */
    public static async Scripts(groupID: string): Promise<ScriptType[]> {
        const data = await API.GET('/api/groups/group/' + groupID + '/scripts');
        return data as ScriptType[];
    }

    /**
     * Set the hosts for this group
     * @param hostIDs array of host IDs
     */
    public static async SetHosts(groupID: string, hostIDs: string[]): Promise<HostType[]> {
        const data = await API.POST('/api/groups/group/' + groupID + '/hosts', { Hosts: hostIDs });
        return data as HostType[];
    }

    /**
     * List all schedules for a group
     */
    public static async Schedules(groupID: string): Promise<ScheduleType[]> {
        const data = await API.GET('/api/groups/group/' + groupID + '/schedules');
        return data as ScheduleType[];
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
