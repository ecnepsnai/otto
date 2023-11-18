import { API } from '../services/API';
import { ScriptRunLevel } from './cbgen_enum';

export interface UserType {
    Username?: string;
    Password?: string;
    CanLogIn?: boolean;
    MustChangePassword?: boolean;
    Permissions?: UserPermissions;
}

export interface UserPermissions {
    ScriptRunLevel?: ScriptRunLevel;
    CanModifyHosts?: boolean;
    CanModifyGroups?: boolean;
    CanModifyScripts?: boolean;
    CanModifySchedules?: boolean;
    CanModifyRunbooks?: boolean;
    CanAccessAuditLog?: boolean;
    CanModifyUsers?: boolean;
    CanModifyAutoregister?: boolean;
    CanModifySystem?: boolean;
}

export class User {
    public static Blank(): UserType {
        return {
            Username: '',
            CanLogIn: true,
            MustChangePassword: false,
            Permissions: {
                ScriptRunLevel: ScriptRunLevel.ReadOnly,
            }
        };
    }

    public static async List(): Promise<UserType[]> {
        const data = await API.GET('/api/users');
        return data as UserType[];
    }

    public static async New(parameters: NewUserParameters): Promise<UserType> {
        const data = await API.PUT('/api/users/user', parameters);
        return data as UserType;
    }

    public static async Save(user: UserType): Promise<UserType> {
        const data = await API.POST('/api/users/user/' + user.Username, user);
        return data as UserType;
    }

    public static async Delete(user: UserType): Promise<unknown> {
        return await API.DELETE('/api/users/user/' + user.Username);
    }

    public static async ResetAPIKey(user: UserType): Promise<string> {
        const data = await API.POST('/api/users/user/' + user.Username + '/apikey', user);
        return data as string;
    }
}

export interface NewUserParameters {
    Username: string;
    Password: string;
    MustChangePassword: boolean;
    Permissions?: UserPermissions;
}

export interface EditUserParameters {
    Password?: string;
    CanLogIn: boolean;
    MustChangePassword: boolean;
    Permissions?: UserPermissions;
}
