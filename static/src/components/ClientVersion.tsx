import * as React from 'react';
import { StateManager } from '../services/StateManager';
import { Heartbeat } from '../types/Heartbeat';
import { Icon } from './Icon';
import { Popover } from './Popover';
import { Style } from './Style';

export interface ClientVersionProps {
    heartbeat: Heartbeat;
}
interface ClientVersionState {
    clientVersion: string;
    clientVersionNumber: number;
    serverVersionNumber: number;
}
export class ClientVersion extends React.Component<ClientVersionProps, ClientVersionState> {
    constructor(props: ClientVersionProps) {
        super(props);

        let versionStr = '';
        if (props.heartbeat) {
            versionStr = props.heartbeat.Version;
        }

        let version = parseInt(versionStr.replace(/\./g, ''));
        if (isNaN(version)) {
            version = 0;
        }

        this.state = {
            clientVersion: versionStr,
            clientVersionNumber: version,
            serverVersionNumber: StateManager.Current().VersionNumber(),
        };
    }

    private isOutOfDate = () => {
        return this.state.clientVersionNumber < this.state.serverVersionNumber;
    }

    render(): JSX.Element {
        const versionString = this.state.clientVersionNumber == 0 ? 'Unknown' : this.state.clientVersion;

        if (this.isOutOfDate()) {
            return (<Popover content="A newer version of the Otto client is available"><Icon.Label icon={<Icon.ExclamationTriangle color={Style.Palette.Warning} />} label={versionString} /></Popover>);
        }

        return (<span>{versionString}</span>);
    }
}
