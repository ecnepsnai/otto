import * as React from 'react';
import { Popover as BSPopover } from 'bootstrap';
import '../../css/popover.scss';

interface PopoverProps {
    content: string;
    children?: React.ReactNode;
}
export const Popover: React.FC<PopoverProps> = (props: PopoverProps) => {
    const spanRef: React.RefObject<HTMLSpanElement> = React.createRef();

    React.useEffect(() => {
        // Popovers must be initialized as they are opt-in
        new BSPopover(spanRef.current, {
            trigger: 'hover',
            placement: 'top',
            content: props.content,
        });
    }, []);

    return (<span ref={spanRef} className="popover-hover" data-toggle="popover">{props.children}</span>);
};
