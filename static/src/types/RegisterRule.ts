import { API } from "../services/API";
import { Modal } from "../components/Modal";
import { Notification } from "../components/Notification";

export class RegisterRule {
    ID: string;
    Property: string;
    Pattern: string;
    GroupID: string;

    constructor(json: any) {
        this.ID = json.ID as string;
        this.Property = json.Property as string;
        this.Pattern = json.Pattern as string;
        this.GroupID = json.GroupID as string;
    }

    /**
     * Return a blank rule
     */
    public static Blank(): RegisterRule {
        return new RegisterRule({
            Property: '',
            Pattern: '',
            GroupID: '',
        });
    }

    /**
     * Create a new RegisterRule
     */
    public static async New(parameters: NewRegisterRuleParameters): Promise<RegisterRule> {
        const data = await API.PUT('/api/register/rules/rule', parameters);
        return new RegisterRule(data);
    }

    /**
     * Save this rule
     */
    public async Save(): Promise<RegisterRule> {
        return RegisterRule.Save(this.ID, this as EditRegisterRuleParameters);
    }

    /**
     * Save a rule
     */
    public static async Save(id: string, parameters: EditRegisterRuleParameters): Promise<RegisterRule> {
        const data = await API.POST('/api/register/rules/rule/' + id, parameters);
        return new RegisterRule(data);
    }

    /**
     * Delete this rule
     */
    public async Delete(): Promise<any> {
        return await API.DELETE('/api/register/rules/rule/' + this.ID);
    }

    /**
     * Modify the rule changing the properties specified
     * @param properties properties to change
     */
    public async Update(properties: {[key:string]: any}): Promise<RegisterRule> {
        const data = await API.PATCH('/api/register/rules/rule/' + this.ID, properties);
        return new RegisterRule(data);
    }

    /**
     * Show a modal to delete this rule
     */
    public async DeleteModal(): Promise<boolean> {
        return new Promise(resolve => {
            Modal.delete('Delete RegisterRule?', 'Are you sure you want to delete this rule? This can not be undone.').then(confirmed => {
                if (!confirmed) {
                    resolve(false);
                    return;
                }

                API.DELETE('/api/register/rules/rule/' + this.ID).then(() => {
                    Notification.success('RegisterRule Deleted');
                    resolve(true);
                });
            });
        });
    }

    /**
     * Get the specified rule by its id
     */
    public static async Get(id: string): Promise<RegisterRule> {
        const data = await API.GET('/api/register/rules/rule/' + id);
        return new RegisterRule(data);
    }

    /**
     * List all rules
     */
    public static async List(): Promise<RegisterRule[]> {
        const data = await API.GET('/api/register/rules');
        return (data as any[]).map(obj => {
            return new RegisterRule(obj);
        });
    }
}

export interface NewRegisterRuleParameters {
    Property: string;
    Pattern: string;
    GroupID: string;
}

export interface EditRegisterRuleParameters {
    Property: string;
    Pattern: string;
    GroupID: string;
}
