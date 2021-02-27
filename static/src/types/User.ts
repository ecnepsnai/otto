import { API } from "../services/API";

export interface UserType {
    Username?: string;
    Email?: string;
    Password?: string;
    CanLogIn?: boolean;
    MustChangePassword?: boolean;
}

export class User {
    public static Blank(): UserType {
        return {
            Username: '',
            Email: '',
            CanLogIn: true,
            MustChangePassword: false,
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
}

export interface NewUserParameters {
    Username: string;
    Email: string;
    Password: string;
    MustChangePassword: boolean;
}

export interface EditUserParameters {
    Email: string;
    Password?: string;
    CanLogIn: boolean;
    MustChangePassword: boolean;
}
