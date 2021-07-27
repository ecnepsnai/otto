import * as React from 'react';
import { Icon } from './Icon';

export namespace Tabs {
    interface TabsProps {
        children?: React.ReactNode;
    }
    /**
     * An enclosed group of tabs. Children of this node must only be Tabs.Tab.
     * The first tab will be selected by default unless one of the children has the 'active' property.
     */
    export const Tabs: React.FC<TabsProps> = (props: TabsProps) => {
        const [SelectedIndex, setSelectedIndex] = React.useState(0);
        const [Ready, setReady] = React.useState(false);

        React.useEffect(() => {
            React.Children.forEach(props.children, (c, idx) => {
                const child = c as React.ReactElement<TabProps>;
                if (!child.props.title) {
                    throw 'Unrecognized child on <Tabs.Tab> element';
                }
                if (child.props.active) {
                    setSelectedIndex(idx);
                }
            });
            setReady(true);
        }, [props.children]);

        const tabClick = (index: number) => {
            return () => {
                setSelectedIndex(index);
            };
        };

        const content = () => {
            return React.Children.toArray(props.children)[SelectedIndex] as React.ReactElement<TabProps>;
        };

        if (!Ready) {
            return null;
        }

        return (<div>
            <ul className="nav nav-tabs" id="myTab" role="tablist">
                {React.Children.map(props.children, (c, idx) => {
                    const child = c as React.ReactElement<TabProps>;
                    const tabID = 'tab-' + idx;
                    const className = idx == SelectedIndex ? 'nav-link active' : 'nav-link';
                    const title = child.props.icon ? (<Icon.Label icon={child.props.icon} label={child.props.title} />) : (<span>{child.props.title}</span>);
                    return (<li className="nav-item" role="presentation">
                        <button
                            className={className}
                            id={tabID}
                            type="button"
                            role="tab"
                            aria-selected={idx == SelectedIndex}
                            onClick={tabClick(idx)}
                        >{title}</button>
                    </li>);
                })}
            </ul>
            <div className="mt-3">
                {content()}
            </div>
        </div>);
    };

    interface TabProps {
        icon?: JSX.Element;
        title: string;
        children?: React.ReactNode;
        active?: boolean;
        className?: string;
    }
    /**
     * A single tab. A title is required, an icon can be provided.
     */
    export const Tab: React.FC<TabProps> = (props: TabProps) => {
        return (<div className={props.className}>{props.children}</div>);
    };
}
