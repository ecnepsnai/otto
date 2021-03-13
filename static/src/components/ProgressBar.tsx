import * as React from 'react';
import { Button } from './Button';
import { Icon } from './Icon';
import { Style } from './Style';
import '../../css/progress.scss';

interface ProgressBarProps {
    percent?: number;
    intermediate?: boolean;
    cancelClick?: () => void;
}
export const ProgressBar: React.FC<ProgressBarProps> = (props: ProgressBarProps) => {
    const [cancelClicked, setCancelClicked] = React.useState(false);

    React.useEffect(() => {
        if (cancelClicked) {
            props.cancelClick();
        }
    }, [cancelClicked]);

    const cancelClick = () => {
        setCancelClicked(true);
    };

    let width = props.percent;
    let intermediate = props.intermediate;

    if (props.percent == undefined && props.intermediate != undefined) {
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

    let content = (<span>{props.percent}%</span>);
    if (intermediate) {
        content = null;
    }

    if (props.cancelClick) {
        return (
            <div className="progress-wrapper">
                <div className="progress">
                    <div className={className} role="progressbar" style={style} aria-valuenow={props.percent} aria-valuemin={0} aria-valuemax={100}>
                        {content}
                    </div>
                </div>
                <Button color={Style.Palette.Danger} outline disabled={cancelClicked} onClick={cancelClick} size={Style.Size.XS}>
                    <Icon.Label icon={<Icon.TimesCircle />} label="Cancel" />
                </Button>
            </div>
        );
    }

    return (
        <div className="progress">
            <div className={className} role="progressbar" style={style} aria-valuenow={props.percent} aria-valuemin={0} aria-valuemax={100}>
                {content}
            </div>
        </div>
    );
};
