angular.module('otto').factory('state', function($api, $q) {
    var currentState;
    var watchCallbacks = [];

    var statePromise = () => {
        return $api.get('/api/state').then(response => {
            currentState = response.data.data;

            watchCallbacks.forEach((cb) => {
                try {
                    cb(currentState);
                } catch (e) {
                    // Don't worry about it
                }
            });

            return response.data.data;
        });
    };

    return {
        current: () => {
            return currentState;
        },
        start: statePromise,
        invalidate: () => {
            currentState = undefined;
            return $q(function(resolve) {
                statePromise().then(() => {
                    window.postMessage('reload_state');
                    resolve();
                });
            });
        },
        watch: (cb) => {
            watchCallbacks.push(cb);
        },
    };
});
