import { API } from "../services/API";
import { Modal } from "../components/Modal";
import { Notification } from "../components/Notification";

export interface RegisterRuleType {
    ID?: string;
    Property?: string;
    Pattern?: string;
    GroupID?: string;
}

export class RegisterRule {
    /**
     * Return a blank rule
     */
    public static Blank(): RegisterRuleType {
        return {
            Property: '',
            Pattern: '',
            GroupID: '',
        };
    }

    /**
     * Create a new RegisterRule
     */
    public static async New(parameters: RegisterRuleType|NewRegisterRuleParameters): Promise<RegisterRuleType> {
        const data = await API.PUT('/api/register/rules/rule', parameters);
        return data as RegisterRuleType;
    }

    /**
     * Save a rule
     */
    public static async Save(id: string, parameters: EditRegisterRuleParameters): Promise<RegisterRuleType> {
        const data = await API.POST('/api/register/rules/rule/' + id, parameters);
        return data as RegisterRuleType;
    }

    /**
     * Delete this rule
     */
    public static async Delete(rule: RegisterRuleType): Promise<any> {
        return await API.DELETE('/api/register/rules/rule/' + rule.ID);
    }

    /**
     * Modify the rule changing the properties specified
     * @param properties properties to change
     */
    public static async Update(rule: RegisterRuleType, properties: {[key:string]: any}): Promise<RegisterRuleType> {
        const data = await API.PATCH('/api/register/rules/rule/' + rule.ID, properties);
        return data as RegisterRuleType;
    }

    /**
     * Show a modal to delete this rule
     */
    public static async DeleteModal(rule: RegisterRuleType): Promise<boolean> {
        return new Promise(resolve => {
            Modal.delete('Delete RegisterRule?', 'Are you sure you want to delete this rule? This can not be undone.').then(confirmed => {
                if (!confirmed) {
                    resolve(false);
                    return;
                }

                API.DELETE('/api/register/rules/rule/' + rule.ID).then(() => {
                    Notification.success('RegisterRule Deleted');
                    resolve(true);
                });
            });
        });
    }

    /**
     * Get the specified rule by its id
     */
    public static async Get(id: string): Promise<RegisterRuleType> {
        const data = await API.GET('/api/register/rules/rule/' + id);
        return data as RegisterRuleType;
    }

    /**
     * List all rules
     */
    public static async List(): Promise<RegisterRuleType[]> {
        const data = await API.GET('/api/register/rules');
        return data as RegisterRuleType[];
    }
}

export interface NewRegisterRuleParameters {
    Property?: string;
    Pattern?: string;
    GroupID?: string;
}

export interface EditRegisterRuleParameters {
    Property?: string;
    Pattern?: string;
    GroupID?: string;
}
