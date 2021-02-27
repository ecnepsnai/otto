import * as React from 'react';

interface SchedulePatternProps { pattern: string; }
export class SchedulePattern extends React.Component<SchedulePatternProps, unknown> {
    render(): JSX.Element {
        let value: string;

        switch (this.props.pattern) {
        case '0 * * * *':
            value = 'Every Hour';
            break;
        case '0 */4 * * *':
            value = 'Every 4 Hours';
            break;
        case '0 0 * * *':
            value = 'Every Day at Midnight';
            break;
        case '0 0 * * 1':
            value = 'Every Monday at Midnight';
            break;
        default:
            value = 'Custom';
            break;
        }

        return (
            <span>{value}</span>
        );
    }
}
