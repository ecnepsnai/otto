angular.module('otto').filter('truncate', (truncate) => {
    return (input, uppercase) => {
        return truncate(input);
    };
});