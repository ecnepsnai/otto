import { Modal } from "../components/Modal";
import { Notification } from "../components/Notification";
import { API } from "../services/API";

export class Attachment {
    ID: string;
    Path: string;
    Name: string;
    MimeType: string;
    UID: number;
    GID: number;
    Mode: number;
    Size: number;

    constructor(json: any) {
        this.ID = json.ID as string;
        this.Path = json.Path as string;
        this.Name = json.Name as string;
        this.MimeType = json.MimeType as string;
        this.UID = json.UID as number;
        this.GID = json.GID as number;
        this.Mode = json.Mode as number;
        this.Size = json.Size as number;
    }

    /**
     * Return a blank attachment
     */
    public static Blank(): Attachment {
        return new Attachment({
            Path: '',
            Name: '',
            UID: 0,
            GID: 0,
            Mode: 644,
        });
    }

    /**
     * Create a new Attachment
     */
    public static async New(attachment: File, parameters: NewAttachmentParameters): Promise<Attachment> {
        const data = await API.PUTFile('/api/attachments', attachment, {
            Path: parameters.Path,
            UID: parameters.UID.toString(),
            GID: parameters.GID.toString(),
            Mode: parameters.Mode.toString(),
        });
        return new Attachment(data);
    }

    /**
     * Save this attachment
     */
    public async Save(): Promise<Attachment> {
        const data = await API.POST('/api/attachments/attachment/' + this.ID, this as EditAttachmentParameters);
        return new Attachment(data);
    }

    /**
     * Delete this attachment
     */
    public async Delete(): Promise<any> {
        return await API.DELETE('/api/attachments/attachment/' + this.ID);
    }

    /**
     * Modify the attachment changing the properties specified
     * @param properties properties to change
     */
    public async Update(properties: {[key:string]: any}): Promise<Attachment> {
        const data = await API.PATCH('/api/attachments/attachment/' + this.ID, properties);
        return new Attachment(data);
    }

    /**
     * Show a modal to delete this attachment
     */
    public async DeleteModal(): Promise<boolean> {
        return new Promise(resolve => {
            Modal.delete('Delete Attachment?', 'Are you sure you want to delete this attachment? This can not be undone.').then(confirmed => {
                if (!confirmed) {
                    resolve(false);
                    return;
                }

                API.DELETE('/api/attachments/attachment/' + this.ID).then(() => {
                    Notification.success('Attachment Deleted');
                    resolve(true);
                });
            });
        });
    }

    /**
     * Get the specified attachment by its id
     */
    public static async Get(attachmentID: string): Promise<Attachment> {
        const data = await API.GET('/api/attachments/attachment/' + attachmentID);
        return new Attachment(data);
    }

    /**
     * List all attachments
     */
    public static async List(): Promise<Attachment[]> {
        const data = await API.GET('/api/attachments');
        return (data as any[]).map(obj => {
            return new Attachment(obj);
        });
    }
}

export interface NewAttachmentParameters {
    Path: string;
    UID: number;
    GID: number;
    Mode: number;
}

export interface EditAttachmentParameters {
    Path: string;
    UID: number;
    GID: number;
    Mode: number;
}
