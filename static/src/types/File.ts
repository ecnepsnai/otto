import { Modal } from "../components/Modal";
import { Notification } from "../components/Notification";
import { API } from "../services/API";

export class File {
    ID: string;
    Path: string;
    UID: number;
    GID: number;
    Mode: number;

    constructor(json: any) {
        this.ID = json.ID as string;
        this.Path = json.Path as string;
        this.UID = json.UID as number;
        this.GID = json.GID as number;
        this.Mode = json.Mode as number;
    }

    /**
     * Return a blank script
     */
    public static Blank(): File {
        return new File({
            Path: '',
            UID: 0,
            GID: 0,
            Mode: 644,
        });
    }

    /**
     * Create a new File
     */
    public static async New(file: Blob, parameters: NewFileParameters): Promise<File> {
        const data = await API.POSTFile('/api/files', file, {
            Path: parameters.Path,
            UID: parameters.UID.toString(),
            GID: parameters.GID.toString(),
            Mode: parameters.Mode.toString(),
        });
        return new File(data);
    }

    /**
     * Save this script
     */
    public async Save(): Promise<File> {
        const data = await API.POST('/api/files/file/' + this.ID, this as EditFileParameters);
        return new File(data);
    }

    /**
     * Delete this script
     */
    public async Delete(): Promise<any> {
        return await API.DELETE('/api/files/file/' + this.ID);
    }

    /**
     * Modify the script changing the properties specified
     * @param properties properties to change
     */
    public async Update(properties: {[key:string]: any}): Promise<File> {
        const data = await API.PATCH('/api/files/file/' + this.ID, properties);
        return new File(data);
    }

    /**
     * Show a modal to delete this file
     */
    public async DeleteModal(): Promise<boolean> {
        return new Promise(resolve => {
            Modal.delete('Delete File?', 'Are you sure you want to delete this file? This can not be undone.').then(confirmed => {
                if (!confirmed) {
                    resolve(false);
                    return;
                }

                API.DELETE('/api/files/file/' + this.ID).then(() => {
                    Notification.success('File Deleted');
                    resolve(true);
                });
            });
        });
    }

    /**
     * Get the specified file by its id
     */
    public static async Get(fileID: string): Promise<File> {
        const data = await API.GET('/api/files/file/' + fileID);
        return new File(data);
    }

    /**
     * List all files
     */
    public static async List(): Promise<File[]> {
        const data = await API.GET('/api/files');
        return (data as any[]).map(obj => {
            return new File(obj);
        });
    }
}

export interface NewFileParameters {
    Path: string;
    UID: number;
    GID: number;
    Mode: number;
}

export interface EditFileParameters {
    Path: string;
    UID: number;
    GID: number;
    Mode: number;
}
