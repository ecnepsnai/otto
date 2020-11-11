import * as React from 'react';
import { Link } from 'react-router-dom';
import { RunModal } from '../pages/run/RunModal';
import { Rand } from '../services/Rand';
import { ScriptEnabledGroup } from '../types/Host';
import { Script } from '../types/Script';
import { SmallPlayButton } from './Button';
import { Card } from './Card';
import { Icon } from './Icon';
import { ListGroup } from './ListGroup';
import { GlobalModalFrame } from './Modal';
import { Nothing } from './Nothing';

export interface ScriptListCardProps {
    scripts: ScriptEnabledGroup[] | Script[];
    hostIDs: string[];
    className?: string;
}
interface ScriptListCardState {}
interface CommonScriptType {
    ID: string;
    Name: string;
}
export class ScriptListCard extends React.Component<ScriptListCardProps, ScriptListCardState> {
    constructor(props: ScriptListCardProps) {
        super(props);
        this.state = { };
    }

    private runScriptClick = (scriptID: string, hostIDs: string[]) => {
        return () => {
            GlobalModalFrame.showModal(<RunModal scriptID={scriptID} hostIDs={hostIDs} key={Rand.ID()}/>);
        };
    }

    private content = () => {
        if (!this.props.scripts || this.props.scripts.length == 0) { return (<Card.Body><Nothing /></Card.Body>); }

        const scripts: CommonScriptType[] = [];
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        this.props.scripts.forEach((s: any) => {
            scripts.push({
                Name: s.ScriptName || s.Name || '',
                ID: s.ScriptID || s.ID || '',
            });
        });

        return (<ListGroup.List>{ scripts.map((script, index) => { return (
            <ListGroup.Item key={index}>
                <div className="d-flex justify-content-between">
                    <div>
                        <Icon.Scroll />
                        <Link to={'/scripts/script/' + script.ID} className="ml-1">{ script.Name }</Link>
                    </div>
                    <div>
                        <SmallPlayButton onClick={this.runScriptClick(script.ID, this.props.hostIDs)} />
                    </div>
                </div>
            </ListGroup.Item>
        );})} </ListGroup.List>);
    }

    render(): JSX.Element {


        return (
            <Card.Card className={this.props.className}>
                <Card.Header>Scripts</Card.Header>
                { this.content() }
            </Card.Card>
        );
    }
}
