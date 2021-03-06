import { Options } from './Options';
import { UserType } from './User';

export interface State {
    User: UserType;
    Runtime: Runtime;
    StartDate: string;
    Hostname: string;
    Warnings: ('default_user_password')[];
    Options: Options.OttoOptions;
    Enums: { [name: string]: { [key: string]: string; }[]; };
}

interface Runtime {
    ServerFQDN: string;
    Version: string;
    Config: string;
}
