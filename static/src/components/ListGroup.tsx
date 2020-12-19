import * as React from 'react';
import '../../css/list-group.scss';

export namespace ListGroup {
    export class List extends React.Component<{}, {}> {
        render(): JSX.Element {
            return (
                <ul className="list-group list-group-flush">{ this.props.children }</ul>
            );
        }
    }

    export interface ItemProps {
        onClick?: () => (void);
        className?: string;
    }
    export class Item extends React.Component<ItemProps, {}> {
        private itemClicked = (event: React.MouseEvent<HTMLLIElement>) => {
            if (this.props.onClick) {
                this.props.onClick();
                event.preventDefault();
            }
        }
        render(): JSX.Element {
            let className = 'list-group-item';
            if (this.props.onClick) {
                className += ' list-group-item-hover';
            }
            if (this.props.className) {
                className += ' ' + this.props.className;
            }
            return (
                <li className={className} onClick={this.itemClicked}>{ this.props.children }</li>
            );
        }
    }

    export interface TextItemProps { title: string; }
    export class TextItem extends React.Component<TextItemProps, {}> {
        render(): JSX.Element {
            return (
                <li className="list-group-item"><strong>{ this.props.title }</strong><span className="ms-1">{ this.props.children }</span></li>
            );
        }
    }
}
