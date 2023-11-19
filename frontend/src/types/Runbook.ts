import { API } from '../services/API';
import { Modal } from '../components/Modal';
import { Notification } from '../components/Notification';
import { ScriptRunLevel } from './cbgen_enum';

export interface RunbookType {
    ID?: string;
    Name?: string;
    GroupIDs?: string[];
    ScriptIDs?: string[];
    HaltOnFailure?: boolean;
    RunLevel?: ScriptRunLevel;
}

export class Runbook {
    /**
     * Return a blank runbook
     */
    public static Blank(): RunbookType {
        return {
            Name: '',
            ScriptIDs: [],
            GroupIDs: [],
            HaltOnFailure: true,
        };
    }

    /**
     * Create a new Runbook
     */
    public static async New(parameters: RunbookType): Promise<RunbookType> {
        const data = await API.PUT('/api/runbooks/runbook', parameters);
        return data as RunbookType;
    }

    /**
     * Save this runbook
     */
    public static async Save(runbook: RunbookType): Promise<RunbookType> {
        const data = await API.POST('/api/runbooks/runbook/' + runbook.ID, runbook);
        return data as RunbookType;
    }

    /**
     * Delete this runbook
     */
    public static async Delete(runbook: RunbookType): Promise<any> {
        return await API.DELETE('/api/runbooks/runbook/' + runbook.ID);
    }

    /**
     * Get the specified runbook by its id
     */
    public static async Get(id: string): Promise<RunbookType> {
        const data = await API.GET('/api/runbooks/runbook/' + id);
        return data as RunbookType;
    }

    /**
     * List all runbooks
     */
    public static async List(): Promise<RunbookType[]> {
        const data = await API.GET('/api/runbooks');
        return data as RunbookType[];
    }

    /**
     * Show a modal to delete this runbook
     */
    public static async DeleteModal(runbook: RunbookType): Promise<boolean> {
        return new Promise(resolve => {
            Modal.delete('Delete Runbook?', 'Are you sure you want to delete this runbook? This can not be undone.').then(confirmed => {
                if (!confirmed) {
                    resolve(false);
                    return;
                }

                API.DELETE('/api/runbooks/runbook/' + runbook.ID).then(() => {
                    Notification.success('Runbook Deleted');
                    resolve(true);
                });
            });
        });
    }
}
