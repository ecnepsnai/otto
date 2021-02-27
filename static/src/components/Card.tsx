import * as React from 'react';
import { Style } from './Style';
import '../../css/card.scss';

export namespace Card {
    interface CardProps {
        color?: Style.Palette;
        className?: string;
        onClick?: () => (void);
        children?: React.ReactNode;
    }
    export const Card: React.FC<CardProps> = (props: CardProps) => {
        let className = 'card';
        if (props.className) {
            className += ' ' + props.className;
        }
        if (props.color) {
            className += ' card-' + props.color.toString();
        }

        return (
            <div className={className} onClick={props.onClick}>
                { props.children }
            </div>
        );
    };

    interface HeaderProps {
        className?: string;
        children?: React.ReactNode;
    }
    export const Header: React.FC<HeaderProps> = (props: HeaderProps) => {
        let className = 'card-header';
        if (props.className) {
            className += ' ' + props.className;
        }
        return (
            <div className={className}>
                { props.children }
            </div>
        );
    };

    interface BodyProps {
        className?: string;
        children?: React.ReactNode;
    }
    export const Body: React.FC<BodyProps> = (props: BodyProps) => {
        let className = 'card-body';
        if (props.className) {
            className += ' ' + props.className;
        }
        return (
            <div className={className}>
                { props.children }
            </div>
        );
    };

    interface ImageProps {
        src: string;
        children?: React.ReactNode;
    }
    export const Image: React.FC<ImageProps> = (props: ImageProps) => {
        return (
            <img className="card-img-top" src={props.src} />
        );
    };

    interface FooterProps {
        children?: React.ReactNode;
    }
    export const Footer: React.FC<FooterProps> = (props: FooterProps) => {
        return (
            <div className="card-footer">
                { props.children }
            </div>
        );
    };
}
