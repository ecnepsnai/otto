import * as React from 'react';
import { Button } from './Button';
import { Icon } from './Icon';
import { Style } from './Style';
import '../../css/progress.scss';

export interface ProgressBarProps {
    percent?: number;
    intermediate?: boolean;
    cancelClick?: () => void;
}
interface ProgressBarState {
    cancelClicked?: boolean;
}
export class ProgressBar extends React.Component<ProgressBarProps, ProgressBarState> {
    constructor(props: ProgressBarProps) {
        super(props);
        this.state = {};
    }

    private cancelClick = () => {
        this.setState({ cancelClicked: true }, () => {
            this.props.cancelClick();
        });
    }

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

        if (this.props.cancelClick) {
            return (
                <div className="progress-wrapper">
                    <div className="progress">
                        <div className={className} role="progressbar" style={style} aria-valuenow={this.props.percent} aria-valuemin={0} aria-valuemax={100}>
                            {content}
                        </div>
                    </div>
                    <Button color={Style.Palette.Danger} outline disabled={this.state.cancelClicked} onClick={this.cancelClick} size={Style.Size.XS}>
                        <Icon.Label icon={<Icon.TimesCircle />} label="Cancel" />
                    </Button>
                </div>
            );
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
