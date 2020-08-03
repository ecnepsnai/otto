import * as React from 'react';
import { Icon } from './Icon';
import '../../css/nothing.scss';

export class Nothing extends React.Component<{}, {}> {
    render(): JSX.Element {
        return (
            <div className="nothing">
                <Icon.Asterisk />
                <span className="">Nothing here...</span>
            </div>
        );
    }
}
