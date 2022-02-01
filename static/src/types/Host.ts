import { API } from '../services/API';
import { Modal } from '../components/Modal';
import { Notification } from '../components/Notification';
import { GroupType } from './Group';
import { Variable } from './Variable';
import { ScheduleType } from './Schedule';
import { HeartbeatType } from './Heartbeat';

export interface HostType {
    ID?: string;
    Name?: string;
    Address?: string;
    Port?: number;
    Trust?: TrustType;
    Enabled?: boolean;
    GroupIDs?: string[];
    Environment?: Variable[];
}

export interface TrustType {
    TrustedIdentity?: string;
    UntrustedIdentity?: string;
    LastTrustUpdate?: string;
}

export class Host {
    /**
     * Return a blank host
     */
    public static Blank(): HostType {
        return {
            Name: '',
            Address: '',
            Port: 12444,
            Enabled: true,
            GroupIDs: [],
            Environment: [],
        };
    }

    /**
     * Create a new Host
     */
    public static async New(parameters: HostType | NewHostParameters): Promise<HostType> {
        const data = await API.PUT('/api/hosts/host', parameters);
        return data as HostType;
    }

    /**
     * Save this host
     */
    public static async Save(host: HostType): Promise<HostType> {
        const data = await API.POST('/api/hosts/host/' + host.ID, host);
        return data as HostType;
    }

    /**
     * Delete this host
     */
    public static async Delete(host: HostType): Promise<unknown> {
        return await API.DELETE('/api/hosts/host/' + host.ID);
    }

    /**
     * Modify the host changing the properties specified
     * @param properties properties to change
     */
    public static async Update(host: HostType, properties: { [key: string]: unknown }): Promise<HostType> {
        const data = await API.PATCH('/api/hosts/host/' + host.ID, properties);
        return data as HostType;
    }

    /**
     * Show a modal to delete this host
     */
    public static async DeleteModal(host: HostType): Promise<boolean> {
        return new Promise(resolve => {
            Modal.delete('Delete Host?', 'Are you sure you want to delete this host? This can not be undone.').then(confirmed => {
                if (!confirmed) {
                    resolve(false);
                    return;
                }

                API.DELETE('/api/hosts/host/' + host.ID).then(() => {
                    Notification.success('Host Deleted');
                    resolve(true);
                });
            });
        });
    }

    /**
     * Get the specified host by its id
     */
    public static async Get(id: string): Promise<HostType> {
        const data = await API.GET('/api/hosts/host/' + id);
        return data as HostType;
    }

    /**
     * List all hosts
     */
    public static async List(): Promise<HostType[]> {
        const data = await API.GET('/api/hosts');
        return data as HostType[];
    }

    /**
     * List all scripts for a host
     */
    public static async Scripts(id: string): Promise<ScriptEnabledGroup[]> {
        const data = await API.GET('/api/hosts/host/' + id + '/scripts');
        return (data as ScriptEnabledGroup[]);
    }

    /**
     * List all groups for a host
     */
    public static async Groups(id: string): Promise<GroupType[]> {
        const data = await API.GET('/api/hosts/host/' + id + '/groups');
        return data as GroupType[];
    }

    /**
     * List all schedules for a host
     */
    public static async Schedules(id: string): Promise<ScheduleType[]> {
        const data = await API.GET('/api/hosts/host/' + id + '/schedules');
        return data as ScheduleType[];
    }

    /**
     * Rotate the PSK for a host, returning the new PSK
     */
    public static async RotatePSK(id: string): Promise<string> {
        const data = await API.POST('/api/hosts/host/' + id + '/psk', null);
        return data as string;
    }

    /**
     * Trigger a heartbeat for this host
     */
    public static async Heartbeat(id: string): Promise<HeartbeatType> {
        const data = await API.POST('/api/hosts/host/' + id + '/heartbeat', null);
        return data as HeartbeatType;
    }

    /**
     * Update the trust for this host
     */
     public static async UpdateTrust(id: string, action: ('permit'|'deny'), publicKey?: string): Promise<HeartbeatType> {
        const data = await API.POST('/api/hosts/host/' + id + '/trust', {
            Action: action,
            PublicKey: publicKey,
        });
        return data as HeartbeatType;
    }

    /**
     * Get the server ID for this host
     */
     public static async ServerID(id: string): Promise<string> {
        const data = await API.GET('/api/hosts/host/' + id + '/id');
        return data as string;
    }
}

export interface NewHostParameters {
    Name: string;
    Address: string;
    Port: number;
    PSK: string;
    GroupIDs: string[];
    Environment: Variable[];
}

export interface EditHostParameters {
    Name: string;
    Address: string;
    Port: number;
    PSK: string;
    LastPSKRotate: string;
    GroupIDs: string[];
    Enabled: boolean;
    Environment: Variable[];
}

export interface ScriptEnabledGroup {
    ScriptID: string;
    ScriptName: string;
    GroupID: string;
    GroupName: string;
}
