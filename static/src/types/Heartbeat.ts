import { API } from "../services/API";

export class Heartbeat {
    Address?: string;
    IsReachable: boolean;
    LastReply?: string;
    LastAttempt?: string;
    Version?: string;
    Properties?: { [key: string]: string};

    constructor(json: any) {
        this.Address = json.Address as string;
        this.IsReachable = json.IsReachable as boolean;
        this.LastReply = json.LastReply as string;
        this.LastAttempt = json.LastAttempt as string;
        this.Version = json.Version as string;
        this.Properties = json.Properties as { [key: string]: string};
    }

    public static List(): Promise<Map<string, Heartbeat>> {
        const map = new Map<string, Heartbeat>();
        return API.GET('/api/heartbeat').then(data => {
            (data as any[]).forEach(obj => {
                const heartbeat = new Heartbeat(obj);
                map.set(heartbeat.Address, heartbeat);
            });
            return map;
        });
    }
}