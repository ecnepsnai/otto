export class DateFormatter {
    private date: Date;

    private constructor(d: Date) {
        this.date = d;
    }

    public static fromDate(date: string): DateFormatter {
        return new DateFormatter(new Date(date));
    }

    public static fromUNIXTimestamp(timestamp: number): DateFormatter {
        return new DateFormatter(new Date(timestamp * 1000));
    }

    public timeFrom(): string {
        const secondsSinceNow = Date.now() - this.date.valueOf();
        if (secondsSinceNow === 0) {
            return 'just now';
        }

        const seconds = 1000; // we're working in ms
        const secondsPerMinute = seconds * 60;
        const secondsPerHour = secondsPerMinute * 60;
        const secondsPerDay = secondsPerHour * 24;
        const secondsPerWeek = secondsPerDay * 7;
        const secondsPerMonth = secondsPerDay * 30;
        const secondsPerYear = secondsPerMonth * 12;

        if (secondsSinceNow > secondsPerYear) {
            const years = Math.round(secondsSinceNow / secondsPerYear);
            return years + ' year' + (years === 1 ? '' : 's') + ' ago';
        } else if (secondsSinceNow > secondsPerMonth) {
            const months = Math.round(secondsSinceNow / secondsPerMonth);
            return months + ' month' + (months === 1 ? '' : 's') + ' ago';
        } else if (secondsSinceNow > secondsPerWeek) {
            const weeks = Math.round(secondsSinceNow / secondsPerWeek);
            return weeks + ' week' + (weeks === 1 ? '' : 's') + ' ago';
        } else if (secondsSinceNow > secondsPerDay) {
            const days = Math.round(secondsSinceNow / secondsPerDay);
            return days + ' day' + (days === 1 ? '' : 's') + ' ago';
        } else if (secondsSinceNow > secondsPerHour) {
            const hours = Math.round(secondsSinceNow / secondsPerHour);
            return hours + ' hour' + (hours === 1 ? '' : 's') + ' ago';
        } else if (secondsSinceNow > secondsPerMinute) {
            const minutes = Math.round(secondsSinceNow / secondsPerMinute);
            return minutes + ' minute' + (minutes === 1 ? '' : 's') + ' ago';
        }

        return 'just now';
    }

    public formateDate(): string {
        return this.date.toLocaleDateString('en-US', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' });
    }
}