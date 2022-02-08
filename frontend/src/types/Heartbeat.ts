import { API } from '../services/API';

export interface HeartbeatType {
    Address?: string;
    IsReachable: boolean;
    LastReply?: string;
    LastAttempt?: string;
    Version?: string;
    Properties?: { [key: string]: string };
}

export class Heartbeat {
    public static async List(): Promise<Map<string, HeartbeatType>> {
        const map = new Map<string, HeartbeatType>();
        const data = await API.GET('/api/heartbeat');
        (data as HeartbeatType[]).forEach(heartbeat => {
            map.set(heartbeat.Address, heartbeat);
        });
        return map;
    }
}