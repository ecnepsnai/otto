import { API } from '../services/API';
import { Modal } from '../components/Modal';
import { Notification } from '../components/Notification';
import { GroupType } from './Group';
import { HostType } from './Host';
import { ScriptType } from './Script';

export interface ScheduleType {
    ID?: string;
    Name?: string;
    ScriptID?: string;
    Scope?: ScheduleScope;
    Pattern?: string;
    Enabled?: boolean;
    LastRunTime?: string;
}

export class Schedule {
    /**
     * Return a blank schedule
     */
    public static Blank(): ScheduleType {
        return {
            Name: '',
            ScriptID: '',
            Scope: {
                HostIDs: [],
                GroupIDs: [],
            },
            Pattern: '',
        };
    }

    /**
     * Create a new Schedule
     */
    public static async New(parameters: ScheduleType|NewScheduleParameters): Promise<ScheduleType> {
        const data = await API.PUT('/api/schedules/schedule', parameters);
        return data as ScheduleType;
    }

    /**
     * Save this schedule
     */
    public static async Save(schedule: ScheduleType): Promise<ScheduleType> {
        const data = await API.POST('/api/schedules/schedule/' + schedule.ID, schedule);
        return data as ScheduleType;
    }

    /**
     * Delete this schedule
     */
    public static async Delete(schedule: ScheduleType): Promise<any> {
        return await API.DELETE('/api/schedules/schedule/' + schedule.ID);
    }

    /**
     * Get the specified schedule by its id
     */
    public static async Get(id: string): Promise<ScheduleType> {
        const data = await API.GET('/api/schedules/schedule/' + id);
        return data as ScheduleType;
    }

    /**
     * Get all reports for a schedule
     */
    public static async Reports(id: string): Promise<ScheduleReport[]> {
        const data = await API.GET('/api/schedules/schedule/' + id + '/reports');
        return (data as any[]).map(obj => {
            return obj as ScheduleReport;
        });
    }

    /**
     * Get all groups for a schedule
     */
    public static async Groups(id: string): Promise<GroupType[]> {
        const data = await API.GET('/api/schedules/schedule/' + id + '/groups');
        return data as GroupType[];
    }

    /**
     * Get all hosts for a schedule
     */
    public static async Hosts(id: string): Promise<HostType[]> {
        const data = await API.GET('/api/schedules/schedule/' + id + '/hosts');
        return data as HostType[];
    }

    /**
     * Get all script for a schedule
     */
    public static async Script(id: string): Promise<ScriptType> {
        const data = await API.GET('/api/schedules/schedule/' + id + '/script');
        return data as ScriptType;
    }

    /**
     * List all schedules
     */
    public static async List(): Promise<ScheduleType[]> {
        const data = await API.GET('/api/schedules');
        return data as ScheduleType[];
    }

    /**
     * Show a modal to delete this schedule
     */
    public static async DeleteModal(schedule: ScheduleType): Promise<boolean> {
        return new Promise(resolve => {
            Modal.delete('Delete Schedule?', 'Are you sure you want to delete this schedule? This can not be undone.').then(confirmed => {
                if (!confirmed) {
                    resolve(false);
                    return;
                }

                API.DELETE('/api/schedules/schedule/' + schedule.ID).then(() => {
                    Notification.success('Schedule Deleted');
                    resolve(true);
                });
            });
        });
    }
}

export interface ScheduleScope {
    HostIDs: string[];
    GroupIDs: string[];
}

export interface NewScheduleParameters {
    Name: string;
    ScriptID: string;
    Scope: ScheduleScope;
    Pattern: string;
}

export interface EditScheduleParameters {
    Name: string;
    Scope: ScheduleScope;
    Pattern: string;
    Enabled: boolean;
}

export interface ScheduleReport {
    ID: string;
    ScheduleID: string;
    HostIDs: string[];
    Time: ScheduleReportTime;
    Result: number;
    HostResult: {[HostID: string]: number};
}

export interface ScheduleReportTime {
    Start: string;
    Finished: string;
    ElapsedSeconds: number;
}
