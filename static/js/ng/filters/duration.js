angular.module('otto').filter('duration', () => {
    function round(value, decimals) {
        return Number(Math.round(value+'e'+decimals)+'e-'+decimals);
    }

    return function(input, uppercase) {
        if (input != undefined) {
            if (input > 600000000000) {
                nHours = round(input/600000000000, 2);
                return nHours + ' hours';
            } else if (input > 100000000000) {
                nMinutes = round(input/100000000000, 0);
                return nMinutes + ' minutes';
            } else if (input > 1000000000) {
                nSeconds = round(input/1000000000, 0);
                return nSeconds + ' seconds';
            } else {
                return 'Less than 1 second';
            }
        }
    };
});