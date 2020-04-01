angular.module('otto').factory('$user', function($api, popup, notify, $q) {
    return {
        list: function() {
            return $api.get('/api/users').then(results => {
                return results.data.data;
            });
        },
        get: function(username) {
            return $api.get('/api/users/user/' + username).then(results => {
                return results.data.data;
            });
        },
        toggle: function(user) {
            return $q(function(resolve) {
                return $api.post('/api/users/user/' + user.Username + '/disable/').then(function() {
                    notify.success('User Saved');
                    resolve();
                });
            });
        },
        delete: function(user) {
            return $q(function(resolve) {
                popup.confirm('Delete User', 'Are you sure you want to delete the user "' + user.Username + '"?', ['Delete', 'Cancel']).then(result => {
                    if (result) {
                        $api.delete('/api/users/user/' + user.Username).then(function() {
                            notify.success('User Deleted');
                            resolve();
                        });
                    }
                });
            });
        },
        new: function(params) {
            return $api.put('/api/users/user', params).then(results => {
                return results.data.data;
            });
        },
        update: function(username, params) {
            return $api.post('/api/users/user/' + username, params).then(results => {
                return results.data.data;
            });
        },
    };
});
