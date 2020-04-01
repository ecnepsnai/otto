angular.module('otto').factory('state', function($api) {
    var currentState;

    var statePromise = function() {
        return $api.get('/api/state').then(response => {
            currentState = response.data.data;
            return response.data.data;
        });
    };

    return {
        current: function() {
            return currentState;
        },
        start: statePromise,
        invalidate: function() {
            currentState = undefined;
            return statePromise();
        },
    };
});