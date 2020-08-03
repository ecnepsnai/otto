import * as React from 'react';

export interface ProgressBarProps {
    percent?: number;
    intermediate?: boolean;
}
export class ProgressBar extends React.Component<ProgressBarProps, {}> {
    render(): JSX.Element {
        let width = this.props.percent;
        let intermediate = this.props.intermediate;

        if (this.props.percent == undefined && this.props.intermediate != undefined) {
            width = 100;
            intermediate = true;
        }

        const style = {
            width: width + '%',
        };

        let className = 'progress-bar';
        if (width >= 100 || intermediate) {
            className += ' progress-bar-striped progress-bar-animated';
        }

        let content = (<span>{this.props.percent}%</span>);
        if (intermediate) {
            content = null;
        }

        return (
        <div className="progress">
            <div className={className} role="progressbar" style={style} aria-valuenow={this.props.percent} aria-valuemin={0} aria-valuemax={100}>
                {content}
            </div>
        </div>
        );
    }
}
