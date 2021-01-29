import { API } from "../services/API";

export class User {
    public Username: string;
    public Email: string;
    public Password?: string;
    public CanLogIn: boolean;
    public MustChangePassword: boolean;

    public constructor(json: any) {
        this.Username = json.Username as string;
        this.Email = json.Email as string;
        this.CanLogIn = json.CanLogIn as boolean;
        this.MustChangePassword = json.MustChangePassword as boolean;
    }

    public static Blank(): User {
        return new User({
            Username: '',
            Email: '',
            CanLogIn: true,
            MustChangePassword: false,
        });
    }

    public static async List(): Promise<User[]> {
        const result = await API.GET('/api/users');
        return (result as any[]).map(user => {
            return new User(user);
        });
    }

    public static async New(parameters: NewUserParameters): Promise<User> {
        const data = await API.PUT('/api/users/user', parameters);
        return new User(data);
    }

    public async Save(): Promise<User> {
        const data = await API.POST('/api/users/user/' + this.Username, this as EditUserParameters);
        return new User(data);
    }

    public async Delete(): Promise<any> {
        return await API.DELETE('/api/users/user/' + this.Username);
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
