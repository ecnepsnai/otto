import * as React from 'react';
import { Popover as BSPopover } from 'bootstrap';
import '../../css/popover.scss';

export interface PopoverProps {
    content: string;
}

export class Popover extends React.Component<PopoverProps, {}> {
    private spanRef: React.RefObject<HTMLSpanElement>;

    constructor(props: PopoverProps) {
        super(props);
        this.spanRef = React.createRef();
    }

    componentDidMount(): void {
        // Popovers must be initialized as they are opt-in
        new BSPopover(this.spanRef.current, {
            trigger: 'hover',
            placement: 'top',
            content: this.props.content,
        });
    }

    render(): JSX.Element {
        return (
            <span ref={this.spanRef} className="popover-hover" data-toggle="popover">
                { this.props.children }
            </span>
        );
    }
}