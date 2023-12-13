import { ScriptRunLevel } from '../types/gengo_enum';
import { StateManager } from './StateManager';

/**
 * Possible actions a user can take
 */
export enum UserAction {
    ModifyHosts,
    ModifyGroups,
    ModifyScripts,
    ModifySchedules,
    ModifyRunbooks,
    AccessAuditLog,
    ModifyUsers,
    ModifyAutoregister,
    ModifySystem,
}

/**
 * Permissions manager
 */
export class Permissions {
    /**
     * Can the current user perform the given action
     * @param action the action to take
     * @returns If the user can take this action or not
     */
    public static UserCan(action: UserAction): boolean {
        const permissions = StateManager.Current().User.Permissions;

        switch (action) {
            case UserAction.ModifyHosts:
                return permissions.CanModifyHosts;
            case UserAction.ModifyGroups:
                return permissions.CanModifyGroups;
            case UserAction.ModifyScripts:
                return permissions.CanModifyScripts;
            case UserAction.ModifySchedules:
                return permissions.CanModifySchedules;
            case UserAction.ModifyRunbooks:
                return permissions.CanModifyRunbooks;
            case UserAction.AccessAuditLog:
                return permissions.CanAccessAuditLog;
            case UserAction.ModifyUsers:
                return permissions.CanModifyUsers;
            case UserAction.ModifyAutoregister:
                return permissions.CanModifyAutoregister;
            case UserAction.ModifySystem:
                return permissions.CanModifySystem;
        }

        return false;
    }

    /**
     * Can the user run a script of this level
     * @param scriptLevel the script run level
     * @returns If the user can run the script
     */
    public static UserCanRunScript(scriptLevel: ScriptRunLevel): boolean {
        return StateManager.Current().User.Permissions.ScriptRunLevel >= scriptLevel;
    }
}