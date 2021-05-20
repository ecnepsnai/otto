import { API } from '../services/API';
import { StateManager } from '../services/StateManager';
import { Variable } from './Variable';

export namespace Options {
    export interface OttoOptions {
        General: General;
        Authentication: Authentication;
        Network: Network;
        Register: Register;
        Security: Security;
    }

    export interface General {
        ServerURL: string;
        GlobalEnvironment: Variable[];
    }

    export interface Authentication {
        MaxAgeMinutes: number;
        SecureOnly: boolean;
    }

    export interface Network {
        ForceIPVersion: string;
        Timeout: number;
        HeartbeatFrequency: number;
    }

    export interface Register {
        Enabled: boolean;
        Key: string;
        DefaultGroupID: string;
        RunScriptsOnRegister: boolean;
    }

    export interface Security {
        IncludePSKEnv: boolean;
    }

    export class Options {
        public static async Get(): Promise<OttoOptions> {
            const results = await API.GET('/api/options');
            return results as OttoOptions;
        }

        public static async Save(options: OttoOptions): Promise<OttoOptions> {
            const results = await API.POST('/api/options', options);
            await StateManager.Refresh();
            return results as OttoOptions;
        }
    }
}
