import * as React from 'react';

interface SchedulePatternProps { pattern: string; }
export const SchedulePattern: React.FC<SchedulePatternProps> = (props: SchedulePatternProps) => {
    let value: string;

    switch (props.pattern) {
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
};
