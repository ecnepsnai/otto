import { State } from '../types/State';
import { API } from './API';

export class StateManager {
    private static current: State;
    private static refreshTimeout: NodeJS.Timeout;

    /**
     * Setup the state manager. Will start the state refresh timer and populate the cache
     * @returns A populated state object
     */
    public static async Initialize(): Promise<State> {
        if (!StateManager.refreshTimeout) {
            StateManager.refreshTimeout = setInterval(() => {
                StateManager.Refresh();
            }, 10000);
        }

        return StateManager.Refresh();
    }

    /**
     * Refresh the current state object
     * @returns A populated state object
     */
    public static async Refresh(): Promise<State> {
        try {
            const data = await API.GET('/api/state');
            this.current = data as State;
        } catch (err) {
            location.href = '/login';
        }

        return this.current;
    }

    /**
     * Get the current state object
     * @returns A state object.
     * @throws If the state cache has not been loaded yet
     */
    public static Current(): State {
        if (this.current == undefined) {
            throw new Error('State was not ready yet!');
        }

        return this.current;
    }

    /**
     * Get a numerical representation of the version number for the Otto service
     * @returns A number representing the version of the Otto service
     */
    public static VersionNumber(): number {
        const version = parseInt(this.current.Runtime.Version.replace(/\./g, ''));
        if (isNaN(version)) {
            return 0;
        }

        return version;
    }
}