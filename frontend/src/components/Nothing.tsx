import * as React from 'react';
import { Icon } from './Icon';
import '../../css/nothing.scss';

export const Nothing: React.FC = () => {
    return (
        <span className="nothing">
            <Icon.StarOfLife />
            <span>None</span>
        </span>
    );
};
