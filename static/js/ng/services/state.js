angular.module('otto').factory('state', function($api, $q) {
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
            return $q(function(resolve) {
                statePromise().then(function() {
                    window.postMessage('reload_state');
                });
            });
        },
    };
});
