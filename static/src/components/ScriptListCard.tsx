import * as React from 'react';
import { Link } from 'react-router-dom';
import { RunModal } from '../pages/run/RunModal';
import { Rand } from '../services/Rand';
import { ScriptEnabledGroup } from '../types/Host';
import { ScriptType } from '../types/Script';
import { SmallPlayButton } from './Button';
import { Card } from './Card';
import { Icon } from './Icon';
import { ListGroup } from './ListGroup';
import { GlobalModalFrame } from './Modal';
import { Nothing } from './Nothing';

interface ScriptListCardProps {
    scripts: ScriptEnabledGroup[] | ScriptType[];
    hostIDs: string[];
    className?: string;
}
interface CommonScriptType {
    ID: string;
    Name: string;
}
export const ScriptListCard: React.FC<ScriptListCardProps> = (props: ScriptListCardProps) => {
    const runScriptClick = (scriptID: string, hostIDs: string[]) => {
        return () => {
            GlobalModalFrame.showModal(<RunModal scriptID={scriptID} hostIDs={hostIDs} key={Rand.ID()}/>);
        };
    };

    const content = () => {
        if (!props.scripts || props.scripts.length == 0) {
            return (<Card.Body><Nothing /></Card.Body>);
        }

        const scripts: CommonScriptType[] = [];
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        props.scripts.forEach((s: any) => {
            scripts.push({
                Name: s.ScriptName || s.Name || '',
                ID: s.ScriptID || s.ID || '',
            });
        });

        return (<ListGroup.List>{ scripts.map((script, index) => {
            return (
                <ListGroup.Item key={index}>
                    <div className="d-flex justify-content-between">
                        <div>
                            <Icon.Scroll />
                            <Link to={'/scripts/script/' + script.ID} className="ms-1">{ script.Name }</Link>
                        </div>
                        <div>
                            <SmallPlayButton onClick={runScriptClick(script.ID, props.hostIDs)} />
                        </div>
                    </div>
                </ListGroup.Item>
            );
        })} </ListGroup.List>);
    };

    return (
        <Card.Card className={props.className}>
            <Card.Header>Scripts</Card.Header>
            { content() }
        </Card.Card>
    );
};
