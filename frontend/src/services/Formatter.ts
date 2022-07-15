export class Formatter {
    private static round(value: number, decimals: number): number {
        return Number(Math.round(parseFloat(value + 'e' + decimals)) + 'e-' + decimals);
    }

    /**
     * Return a formatted human readable strin representing the duration of the given amount of nanoseconds
     * @param input duration in nanoseconds
     * @returns A formatted string
     */
    public static DurationNS(input: number): string {
        return this.duration(input, 3600000000000, 60000000000, 1000000000);
    }

    /**
     * Return a formatted human readable strin representing the duration of the given amount of seconds
     * @param input duration in seconds
     * @returns A formatted string
     */
    public static DurationS(input: number): string {
        return this.duration(input, 3600, 60, 1);
    }

    private static duration(input: number, dHour: number, dMinute: number, dSecond: number): string {
        if (input == undefined) {
            console.warn('undefined input specified');
            return;
        }

        if (input < dSecond) {
            return 'Less than 1 second';
        }

        const hours = Math.floor(input / dHour);
        const minutes = Math.floor((input % dHour) / dMinute);
        const seconds = Math.floor(input % dMinute);
    
        let result = '';
        if (hours >= 1) {
            result += hours + ' hour' + (hours > 1 ? 's' : '') + ' ';
        }
        if (minutes >= 1) {
            result += minutes + ' minute' + (minutes > 1 ? 's' : '') + ' ';
        }
        if (seconds >= 1) {
            let s = seconds.toString();
            if (s.length > 2) {
                s = s.substring(0, 2);

                if (s.at(1) == '0') {
                    s = s.at(0);
                }
            }
            result += s + ' second' + (seconds > 1 ? 's' : '');
        }
        
        return result;
    }

    public static Bytes(input: number): string {
        const KB = 1024;
        const MB = 1024 * 1024;
        const GB = 1024 * 1024 * 1024;
        const TB = 1024 * 1024 * 1024 * 1024;

        if (input == undefined) {
            console.warn('undefined input specified');
            return;
        }

        if (input > TB) {
            return this.round(input / TB, 2) + ' TiB';
        } else if (input == TB) {
            return '1 TiB';
        } else if (input > GB) {
            return this.round(input / GB, 2) + ' GiB';
        } else if (input == GB) {
            return '1 GiB';
        } else if (input > MB) {
            return this.round(input / MB, 2) + ' MiB';
        } else if (input == MB) {
            return '1 MiB';
        } else if (input > KB) {
            return this.round(input / KB, 2) + ' KiB';
        } else if (input == KB) {
            return '1 KiB';
        }

        return input + ' B';
    }

    public static ValueOrNothing(v: number): string {
        if (!v || isNaN(v) || v == 0) {
            return '';
        }

        return v.toString();
    }
}