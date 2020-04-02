angular.module('otto').factory('truncate', () => {
    return (input) => {
        if (input == undefined) {
            return;
        }

        if (input.length > 100) {
            return input.substring(0, 100) + '...';
        }
        return input;
    };
});