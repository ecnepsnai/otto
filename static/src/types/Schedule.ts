import { API } from "../services/API";
import { Modal } from "../components/Modal";
import { Notification } from "../components/Notification";
import { Group } from "./Group";
import { Host } from "./Host";
import { Script } from "./Script";

export class Schedule {
    ID: string;
    Name: string;
    ScriptID: string;
    Scope: ScheduleScope;
    Pattern: string;
    Enabled: boolean;
    LastRunTime: string;

    public constructor(json: any) {
        this.ID = json.ID as string;
        this.Name = json.Name as string;
        this.ScriptID = json.ScriptID as string;
        this.Scope = json.Scope as ScheduleScope;
        this.Pattern = json.Pattern as string;
        this.Enabled = json.Enabled as boolean;
        this.LastRunTime = json.LastRunTime as string;
    }

    /**
     * Return a blank schedule
     */
    public static Blank(): Schedule {
        return new Schedule({
            Name: '',
            ScriptID: '',
            Scope: {
                HostIDs: [],
                GroupIDs: [],
            },
            Pattern: '',
        });
    }

    /**
     * Create a new Schedule
     */
    public static async New(parameters: NewScheduleParameters): Promise<Schedule> {
        const data = await API.PUT('/api/schedules/schedule', parameters);
        return new Schedule(data);
    }

    /**
     * Save this schedule
     */
    public async Save(): Promise<Schedule> {
        const data = await API.POST('/api/schedules/schedule/' + this.ID, this as EditScheduleParameters);
        return new Schedule(data);
    }

    /**
     * Delete this schedule
     */
    public async Delete(): Promise<any> {
        return await API.DELETE('/api/schedules/schedule/' + this.ID);
    }

    /**
     * Get the specified schedule by its id
     */
    public static async Get(id: string): Promise<Schedule> {
        const data = await API.GET('/api/schedules/schedule/' + id);
        return new Schedule(data);
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
     * Get all reports for this schedule
     */
    public async Reports(): Promise<ScheduleReport[]> {
        return Schedule.Reports(this.ID);
    }

    /**
     * Get all groups for a schedule
     */
    public static async Groups(id: string): Promise<Group[]> {
        const data = await API.GET('/api/schedules/schedule/' + id + '/groups');
        return (data as any[]).map(obj => {
            return new Group(obj);
        });
    }

    /**
     * Get all groups for this schedule
     */
    public async Groups(): Promise<Group[]> {
        return Schedule.Groups(this.ID);
    }

    /**
     * Get all hosts for a schedule
     */
    public static async Hosts(id: string): Promise<Host[]> {
        const data = await API.GET('/api/schedules/schedule/' + id + '/hosts');
        return (data as any[]).map(obj => {
            return new Host(obj);
        });
    }

    /**
     * Get all hosts for this schedule
     */
    public async Hosts(): Promise<Host[]> {
        return Schedule.Hosts(this.ID);
    }

    /**
     * Get all script for a schedule
     */
    public static async Script(id: string): Promise<Script> {
        const data = await API.GET('/api/schedules/schedule/' + id + '/script');
        return data as Script;
    }

    /**
     * Get all script for this schedule
     */
    public async Script(): Promise<Script> {
        return Schedule.Script(this.ID);
    }

    /**
     * List all schedules
     */
    public static async List(): Promise<Schedule[]> {
        const data = await API.GET('/api/schedules');
        return (data as any[]).map(obj => {
            return new Schedule(obj);
        });
    }

    /**
     * Show a modal to delete this schedule
     */
    public async DeleteModal(): Promise<boolean> {
        return new Promise(resolve => {
            Modal.delete('Delete Schedule?', 'Are you sure you want to delete this schedule? This can not be undone.').then(confirmed => {
                if (!confirmed) {
                    resolve(false);
                    return;
                }

                API.DELETE('/api/schedules/schedule/' + this.ID).then(() => {
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
}

export interface ScheduleReportTime {
    Start: string;
    Finished: string;
    ElapsedSeconds: number;
}
