import * as React from 'react';
import { Popover } from './Popover';
import { DateFormatter } from '../services/Date';

export interface DateLabelProps {
    date?: string;
    timestamp?: number;
}

interface DateLabelState {
    d: DateFormatter;
    never?: boolean;
}

export class DateLabel extends React.Component<DateLabelProps, DateLabelState> {
    constructor(props: DateLabelProps) {
        super(props);
        let d: DateFormatter;
        let never: boolean;
        if (props.date) {
            d = DateFormatter.fromDate(props.date);
            if (props.date.indexOf('0001-01-01') >= 0) {
                never = true;
            }
        } else if (props.timestamp >= 0) {
            d = DateFormatter.fromUNIXTimestamp(props.timestamp);
            if (props.timestamp == 0) {
                never = true;
            }
        } else {
            throw new Error("either date or timestamp must be specified for DateLabel");
        }
        this.state = { d: d, never: never };
    }
    render(): JSX.Element {
        if (this.state.never) {
            return <span>Never</span>;
        }

        return (
            <Popover content={this.state.d.formateDate()}>
                { this.state.d.timeFrom() }
            </Popover>
        );
    }
}