import moment = require('moment');

export class DateFormatter {
    private m: moment.Moment;

    private constructor(m: moment.Moment) {
        this.m = m;
    }

    public static fromDate(date: string): DateFormatter {
        return new DateFormatter(moment(date));
    }

    public static fromUNIXTimestamp(timestamp: number): DateFormatter {
        return new DateFormatter(moment.unix(timestamp));
    }

    public timeFrom(): string {
        return this.m.fromNow();
    }

    public formateDate(): string {
        return this.m.format('dddd, MMMM Do YYYY, h:mm:ss a');
    }
}