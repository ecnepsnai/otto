import { API } from "../services/API";

export namespace Options {
    export interface OttoOptions {
        General: General;
        Network: Network;
        Register: Register;
    }

    export interface General {
        ServerURL: string;
        GlobalEnvironment: {[id: string]: string};
    }

    export interface Network {
        ForceIPVersion: string;
        Timeout: number;
        HeartbeatFrequency: number;
    }

    export interface Register {
        Enabled: boolean;
        PSK: string;
        Rules: RegisterRule[];
        DefaultGroupID: string;
    }

    export interface RegisterRule {
        Uname?: string;
        Hostname?: string;
        GroupID?: string;
    }

    export class Options {
        public static async Get(): Promise<OttoOptions> {
            const results = await API.GET('/api/options');
            return results as OttoOptions;
        }

        public static async Save(options: OttoOptions): Promise<OttoOptions> {
            const results = await API.POST('/api/options', options);
            return results as OttoOptions;
        }
    }
}
