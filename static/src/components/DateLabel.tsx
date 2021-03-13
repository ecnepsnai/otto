import * as React from 'react';
import { Popover } from './Popover';
import { DateFormatter } from '../services/Date';

interface DateLabelProps {
    date?: string;
}
export const DateLabel: React.FC<DateLabelProps> = (props: DateLabelProps) => {
    let d: DateFormatter;
    let never: boolean;
    if (props.date) {
        d = DateFormatter.fromDate(props.date);
        if (props.date.indexOf('0001-01-01') >= 0) {
            never = true;
        }
    } else {
        throw new Error('either date or timestamp must be specified for DateLabel');
    }
    if (never) {
        return <span>Never</span>;
    }

    return (
        <Popover content={d.formateDate()}>
            { d.timeFrom() }
        </Popover>
    );
};
