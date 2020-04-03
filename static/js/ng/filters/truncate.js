angular.module('otto').filter('truncate', function(truncate) {
    return function(input, uppercase) {
        return truncate(input);
    };
});