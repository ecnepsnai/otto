import { API } from "../services/API";

export class User {
    public Username: string;
    public Email: string;
    public Enabled: boolean;
    public Password?: string;

    public constructor(json: any) {
        this.Username = json.Username as string;
        this.Email = json.Email as string;
        this.Enabled = json.Enabled as boolean;
    }

    public static Blank(): User {
        return new User({
            Username: '',
            Email: '',
            Enabled: true
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
}

export interface EditUserParameters {
    Email: string;
    Enabled: boolean;
    Password?: string;
}
