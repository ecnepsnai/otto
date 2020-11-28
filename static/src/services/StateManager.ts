import { State } from '../types/State';
import { API } from './API';

export class StateManager {
    private static current: State;
    public static Refresh(): Promise<State> {
        return API.GET('/api/state').then(data => {
            this.current = new State(data);
            console.log('State loaded', this.current);
            return this.current;
        });
    }
    public static Current(): State {
        if (this.current == undefined) {
            throw new Error('State was not ready yet!');
        }

        return this.current;
    }
}