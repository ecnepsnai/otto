import { Modal } from '../components/Modal';
import { Notification } from '../components/Notification';
import { API } from '../services/API';

export interface AttachmentType {
    ID?: string;
    Path?: string;
    Name?: string;
    MimeType?: string;
    UID?: number;
    GID?: number;
    Mode?: number;
    Size?: number;
}

export class Attachment {
    /**
     * Return a blank attachment
     */
    public static Blank(): AttachmentType {
        return {
            Path: '',
            Name: '',
            UID: 0,
            GID: 0,
            Mode: 644,
        };
    }

    /**
     * Create a new Attachment
     */
    public static async New(attachment: File, parameters: AttachmentType | NewAttachmentParameters): Promise<AttachmentType> {
        const data = await API.PUTFile('/api/attachments', attachment, {
            Path: parameters.Path,
            UID: parameters.UID.toString(),
            GID: parameters.GID.toString(),
            Mode: parameters.Mode.toString(),
        });
        return data as AttachmentType;
    }

    /**
     * Save this attachment
     */
    public static async Save(attachment: File, parameters: AttachmentType): Promise<AttachmentType> {
        const data = await API.POSTFile('/api/attachments/attachment/' + parameters.ID, attachment, {
            Path: parameters.Path,
            UID: parameters.UID.toString(),
            GID: parameters.GID.toString(),
            Mode: parameters.Mode.toString(),
        });
        return data as AttachmentType;
    }

    /**
     * Delete this attachment
     */
    public static async Delete(attachment: AttachmentType): Promise<unknown> {
        return await API.DELETE('/api/attachments/attachment/' + attachment.ID);
    }

    /**
     * Modify the attachment changing the properties specified
     * @param properties properties to change
     */
    public static async Update(attachment: AttachmentType, properties: { [key: string]: unknown }): Promise<AttachmentType> {
        const data = await API.PATCH('/api/attachments/attachment/' + attachment.ID, properties);
        return data as AttachmentType;
    }

    /**
     * Show a modal to delete this attachment
     */
    public static async DeleteModal(attachment: AttachmentType): Promise<boolean> {
        return new Promise(resolve => {
            Modal.delete('Delete Attachment?', 'Are you sure you want to delete this attachment? This can not be undone.').then(confirmed => {
                if (!confirmed) {
                    resolve(false);
                    return;
                }

                API.DELETE('/api/attachments/attachment/' + attachment.ID).then(() => {
                    Notification.success('Attachment Deleted');
                    resolve(true);
                });
            });
        });
    }

    /**
     * Get the specified attachment by its id
     */
    public static async Get(attachmentID: string): Promise<AttachmentType> {
        const data = await API.GET('/api/attachments/attachment/' + attachmentID);
        return data as AttachmentType;
    }

    /**
     * List all attachments
     */
    public static async List(): Promise<Attachment[]> {
        const data = await API.GET('/api/attachments');
        return data as AttachmentType[];
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
