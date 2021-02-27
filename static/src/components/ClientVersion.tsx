import * as React from 'react';
import { StateManager } from '../services/StateManager';
import { HeartbeatType } from '../types/Heartbeat';
import { Icon } from './Icon';
import { Popover } from './Popover';
import { Style } from './Style';

interface ClientVersionProps {
    heartbeat: HeartbeatType;
}
export const ClientVersion: React.FC<ClientVersionProps> = (props: ClientVersionProps) => {
    let versionStr = '';
    if (props.heartbeat) {
        versionStr = props.heartbeat.Version;
    }

    let version = parseInt(versionStr.replace(/\./g, ''));
    if (isNaN(version)) {
        version = 0;
    }
    const clientVersion = versionStr;
    const clientVersionNumber = version;
    const serverVersionNumber = StateManager.Current().VersionNumber();


    const isOutOfDate = () => {
        return clientVersionNumber < serverVersionNumber;
    };

    const versionString = clientVersionNumber == 0 ? 'Unknown' : clientVersion;

    if (isOutOfDate()) {
        return (<Popover content="A newer version of the Otto client is available"><Icon.Label icon={<Icon.ExclamationTriangle color={Style.Palette.Warning} />} label={versionString} /></Popover>);
    }

    return (<span>{versionString}</span>);
};
