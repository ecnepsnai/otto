import * as React from 'react';
import '../../css/layout.scss';
import { Style } from './Style';

export namespace Layout {
    interface LayoutProps {
        className?: string;
        size?: Style.Size;
        children?: React.ReactNode;
    }

    export const Container: React.FC<LayoutProps> = (props: LayoutProps) => {
        let className = 'container';
        if (props.className) {
            className += ' ' + props.className;
        }
        return (
            <div className={className}>
                { props.children}
            </div>
        );
    };

    export const Row: React.FC<LayoutProps> = (props: LayoutProps) => {
        let className = 'row';
        if (props.className) {
            className += ' ' + props.className;
        }
        return (
            <div className={className}>
                { props.children}
            </div>
        );
    };

    export const Column: React.FC<LayoutProps> = (props: LayoutProps) => {
        const className = props.className || 'col-md';
        return (
            <div className={className}>
                { props.children}
            </div>
        );
    };
}
