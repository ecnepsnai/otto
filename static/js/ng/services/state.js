angular.module('otto').factory('state', function($api, $q) {
    var currentState;

    var statePromise = () => {
        return $api.get('/api/state').then(response => {
            currentState = response.data.data;
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
            return $q((resolve) => {
                statePromise().then(() => {
                    window.postMessage('reload_state');
                    resolve();
                });
            });
        },
    };
});
