import * as React from 'react';
import { Style } from './Style';
import '../../css/card.scss';

export namespace Card {
    export interface CardProps {
        color?: Style.Palette;
        className?: string;
    }

    export class Card extends React.Component<CardProps, {}> {
        render(): JSX.Element {
            let className = 'card';
            if (this.props.className) {
                className += ' ' + this.props.className;
            }
            if (this.props.color) {
                className += ' card-' + this.props.color.toString();
            }

            return (
                <div className={className}>
                    { this.props.children }
                </div>
            );
        }
    }

    export interface HeaderProps {
        className?: string;
    }
    export class Header extends React.Component<HeaderProps, {}> {
        render(): JSX.Element {
            let className = 'card-header';
            if (this.props.className) {
                className += ' ' + this.props.className;
            }
            return (
                <div className={className}>
                    { this.props.children }
                </div>
            );
        }
    }

    export interface BodyProps {
        className?: string;
    }
    export class Body extends React.Component<BodyProps, {}> {
        render(): JSX.Element {
            let className = 'card-body';
            if (this.props.className) {
                className += ' ' + this.props.className;
            }
            return (
                <div className={className}>
                    { this.props.children }
                </div>
            );
        }
    }

    export interface ImageProps {
        src: string;
    }

    export class Image extends React.Component<ImageProps, {}> {
        render(): JSX.Element {
            return (
                <img className="card-img-top" src={this.props.src} />
            );
        }
    }

    export class Footer extends React.Component<{}, {}> {
        render(): JSX.Element {
            return (
                <div className="card-footer">
                    { this.props.children }
                </div>
            );
        }
    }
}
