import * as React from 'react';
import '../../css/layout.scss';

export namespace Layout {
    export interface LayoutProps {
        className?: string;
    }

    export class Container extends React.Component<LayoutProps, {}> {
        render(): JSX.Element {
            let className = 'container';
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

    export class Row extends React.Component<LayoutProps, {}> {
        render(): JSX.Element {
            let className = 'row';
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

    export class Column extends React.Component<LayoutProps, {}> {
        render(): JSX.Element {
            let className = 'col';
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
}
