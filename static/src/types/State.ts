import { Options } from './Options';
import { UserType } from './User';

export class State {
    public readonly User: UserType;
    public readonly Runtime: Runtime;
    public readonly StartDate: string;
    public readonly Hostname: string;
    public readonly Warnings: ('default_user_password')[];
    public readonly Options: Options.OttoOptions;
    public readonly Enums: { [name: string]: { [key: string] : string; }[]; };

    public constructor(json: any) {
        this.User = json.User as UserType;
        this.Runtime = json.Runtime as Runtime;
        this.StartDate = json.StartDate as string;
        this.Hostname = json.Hostname as string;
        this.Warnings = json.Warnings as ('default_user_password')[];
        this.Options = json.Options as Options.OttoOptions;
        this.Enums = json.Enums as { [name: string]: { [key: string] : string; }[]; };
    }

    /**
     * Get a integer representation of the current system version
     */
    public VersionNumber(): number {
        const version = parseInt(this.Runtime.Version.replace(/\./g, ''));
        if (isNaN(version)) {
            return 0;
        }

        return version;
    }
}

interface Runtime {
    ServerFQDN: string;
    Version: string;
    Config: string;
}
