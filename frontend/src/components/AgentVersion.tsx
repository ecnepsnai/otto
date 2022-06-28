import * as React from 'react';
import { StateManager } from '../services/StateManager';
import { HeartbeatType } from '../types/Heartbeat';
import { Icon } from './Icon';
import { Popover } from './Popover';
import { Style } from './Style';

interface AgentVersionProps {
    heartbeat: HeartbeatType;
}
export const AgentVersion: React.FC<AgentVersionProps> = (props: AgentVersionProps) => {
    let versionStr = '';
    if (props.heartbeat) {
        versionStr = props.heartbeat.Version;
    }

    let version = parseInt(versionStr.replace(/\./g, ''));
    if (isNaN(version)) {
        version = 0;
    }
    const agentVersion = versionStr;
    const agentVersionNumber = version;
    const serverVersionNumber = StateManager.VersionNumber();


    const isOutOfDate = () => {
        return agentVersionNumber < serverVersionNumber;
    };

    const versionString = agentVersionNumber == 0 ? 'Unknown' : agentVersion;

    if (isOutOfDate()) {
        return (<Popover content="A newer version of the Otto agent is available"><Icon.Label icon={<Icon.ExclamationTriangle color={Style.Palette.Warning} />} label={versionString} /></Popover>);
    }

    return (<span>{versionString}</span>);
};
