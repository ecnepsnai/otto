import * as React from 'react';
import { Link } from 'react-router-dom';
import { RunModal } from '../pages/run/RunModal';
import { Rand } from '../services/Rand';
import { ScriptType } from '../types/Script';
import { SmallPlayButton } from './Button';
import { Card } from './Card';
import { Icon } from './Icon';
import { ListGroup } from './ListGroup';
import { GlobalModalFrame } from './Modal';
import { Nothing } from './Nothing';
import { Permissions } from '../services/Permissions';

interface ScriptListCardProps {
    scripts: ScriptType[];
    hostIDs: string[];
    className?: string;
}
export const ScriptListCard: React.FC<ScriptListCardProps> = (props: ScriptListCardProps) => {
    const runScriptClick = (scriptID: string, hostIDs: string[]) => {
        return () => {
            GlobalModalFrame.showModal(<RunModal scriptID={scriptID} hostIDs={hostIDs} key={Rand.ID()} />);
        };
    };

    const content = () => {
        if (!props.scripts || props.scripts.length == 0) {
            return (<Card.Body><Nothing /></Card.Body>);
        }

        return (<ListGroup.List>{props.scripts.map((script, index) => {
            return (
                <ListGroup.Item key={index}>
                    <div className="d-flex justify-content-between">
                        <div>
                            <Icon.Scroll />
                            <Link to={'/scripts/script/' + script.ID} className="ms-1">{script.Name}</Link>
                        </div>
                        <div>
                            { Permissions.UserCanRunScript(script.RunLevel) ? (<SmallPlayButton onClick={runScriptClick(script.ID, props.hostIDs)} />) : null }
                        </div>
                    </div>
                </ListGroup.Item>
            );
        })} </ListGroup.List>);
    };

    return (
        <Card.Card className={props.className}>
            <Card.Header>Scripts</Card.Header>
            { content()}
        </Card.Card>
    );
};
