import * as React from 'react';
import '../../css/pre.scss';

interface PreProps {
    children?: React.ReactNode;
}
export const Pre: React.FC<PreProps> = (props: PreProps) => {
    return (
        <div className="pre-wrapper"><pre>{props.children}</pre></div>
    );
};
