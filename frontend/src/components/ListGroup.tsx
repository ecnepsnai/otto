import * as React from 'react';
import '../../css/list-group.scss';

export namespace ListGroup {
    export const List: React.FC = (props: { children: React.ReactNode }) => {
        return (
            <ul className="list-group list-group-flush">{props.children}</ul>
        );
    };

    interface ItemProps {
        onClick?: () => (void);
        className?: string;
        children?: React.ReactNode;
    }
    export const Item: React.FC<ItemProps> = (props: ItemProps) => {
        const itemClicked = (event: React.MouseEvent<HTMLLIElement>) => {
            if (props.onClick) {
                props.onClick();
                event.preventDefault();
            }
        };

        let className = 'list-group-item';
        if (props.onClick) {
            className += ' list-group-item-hover';
        }
        if (props.className) {
            className += ' ' + props.className;
        }
        return (
            <li className={className} onClick={itemClicked}>{props.children}</li>
        );
    };

    interface TextItemProps {
        title: string;
        children?: React.ReactNode;
    }
    export const TextItem: React.FC<TextItemProps> = (props: TextItemProps) => {
        return (
            <li className="list-group-item"><strong>{props.title}</strong><span className="ms-1">{props.children}</span></li>
        );
    };
}
