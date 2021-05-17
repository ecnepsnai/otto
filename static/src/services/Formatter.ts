export class Formatter {
    private static round(value: number, decimals: number): number {
        return Number(Math.round(parseFloat(value + 'e' + decimals)) + 'e-' + decimals);
    }

    public static Duration(input: number): string {
        if (input == undefined) {
            console.warn('undefined input specified');
            return;
        }
        if (input < 10000000000) {
            return 'Less than 1 minute';
        }

        if (input > 600000000000) {
            const nHours = this.round(input / 600000000000, 2);
            return nHours + ' hours';
        }

        const nSeconds = this.round(input / 10000000000, 2);
        return nSeconds + ' minutes';
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