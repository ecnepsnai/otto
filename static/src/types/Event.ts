import { API } from "../services/API";

export class Event {
    ID: string;
    Event: string;
    Time: string;
    Details: {[key: string]: string;};

    constructor(json: any) {
        this.ID = json.ID as string;
        this.Event = json.Event as string;
        this.Time = json.Time as string;
        this.Details = json.Details as {[key: string]: string;};
    }

    public static async List(count: number): Promise<Event[]> {
        const objects = await API.GET('/api/events?c=' + count);
        return (objects as any[]).map(o => new Event(o));
    }
}
