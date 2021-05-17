import { API } from '../services/API';

export interface EventType {
    ID?: string;
    Event?: string;
    Time?: string;
    Details?: { [key: string]: string; };
}

export class Event {
    public static async List(count: number): Promise<EventType[]> {
        const data = await API.GET('/api/events?c=' + count);
        return data as EventType[];
    }
}
