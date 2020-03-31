angular.module('otto').factory('$heartbeat', function($api) {
    return {
        list: function() {
            return $api.get('/api/heartbeat').then(results => {
                return results.data.data;
            });
        },
    };
});
