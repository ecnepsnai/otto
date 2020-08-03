import * as React from 'react';
import { Bootstrap } from '../services/Bootstrap';
import { Rand } from '../services/Rand';
import '../../css/popover.scss';

export interface PopoverProps {
    content: string;
}

export class Popover extends React.Component<PopoverProps, {}> {
    private id: string = Rand.ID()
    componentDidMount(): void {
        const element = document.getElementById(this.id);
        Bootstrap.Popover(element);
    }
    render(): JSX.Element {
        return (
            <span className="popover-hover" data-toggle="popover" data-trigger="hover" data-placement="top" id={this.id} data-content={this.props.content}>
                { this.props.children }
            </span>
        );
    }
}