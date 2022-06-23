import * as React from 'react';
import { Icon } from './Icon';
import { useNavigate } from 'react-router-dom';
import '../../css/context-menu.scss';

export interface ContextMenuItem {
    title: string;
    icon?: JSX.Element;
    disabled?: boolean;
    href?: string;
    onClick?: () => void;
}

export interface ContextMenuProps {
    x: number;
    y: number;
    items: (ContextMenuItem | 'separator')[];
}

export const ContextMenu: React.FC<ContextMenuProps> = (props: ContextMenuProps) => {
    const navigate = useNavigate();

    const style: React.CSSProperties = {
        display: 'block',
        position: 'absolute',
        top: props.y,
    };

    // Determine if the mouse cursor/tap is close to the right hand edge
    // Average menu is 150px wide, assume 175px width
    if (document.body.offsetWidth - props.x > 175) {
        style.left = props.x;
    } else {
        style.left = props.x - 175;
    }

    const itemClick = (idx: number) => {
        return () => {
            const item = props.items[idx] as ContextMenuItem;

            GlobalContextMenuFrame.removeMenu();

            if (item.href) {
                navigate(item.href);
            }

            if (item.onClick) {
                item.onClick();
            }
        };
    };

    return (
        <ul className="dropdown-menu active" style={style}>
            {
                props.items.map((item, idx) => {
                    if (item === 'separator') {
                        return (<li key={idx}><hr className="dropdown-divider" /></li>);
                    }

                    let content = (<span>{item.title}</span>);
                    if (item.icon) {
                        content = (<Icon.Label icon={item.icon} label={item.title} />);
                    }

                    let className = 'dropdown-item';
                    if (item.disabled) {
                        className += ' disabled';
                    }

                    return (<li key={idx}><span className={className} aria-disabled={item.disabled} onClick={itemClick(idx)}>{content}</span></li>);
                })
            }
        </ul>
    );
};

interface GlobalContextMenuFrameState {
    menu?: JSX.Element;
}

export class GlobalContextMenuFrame extends React.Component<unknown, GlobalContextMenuFrameState> {
    constructor(props: unknown) {
        super(props);
        this.state = {};
        GlobalContextMenuFrame.instance = this;
    }

    private static instance: GlobalContextMenuFrame;

    public static showMenu(menu: JSX.Element): void {
        const menuBackdrop = document.createElement('div');
        menuBackdrop.id = 'menu-backdrop';
        menuBackdrop.onclick = (e: MouseEvent) => {
            e.preventDefault();
            this.removeMenu();
        };
        menuBackdrop.oncontextmenu = (e: MouseEvent) => {
            e.preventDefault();
            this.removeMenu();
        };
        document.body.appendChild(menuBackdrop);

        this.instance.setState(state => {
            if (state.menu != undefined) {
                throw new Error('Refusing to stack menus');
            }
            return { menu: menu };
        });
    }

    public static removeMenu(): void {
        try {
            document.querySelector('#menu-backdrop').remove();
        } catch (e) {
            //
        }
        this.instance.setState({ menu: undefined });
    }

    render(): JSX.Element {
        return (
            <div id="global-menu-frame">
                {
                    this.state.menu
                }
            </div>
        );
    }
}
