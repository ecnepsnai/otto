angular.module('otto').factory('$heartbeat', function($api) {
    return {
        list: () => {
            return $api.get('/api/heartbeat').then(results => {
                return results.data.data;
            });
        },
    };
});
