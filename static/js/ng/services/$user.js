angular.module('otto').factory('$user', function($api, popup, notify, $q) {
    return {
        list: function() {
            return $api.get('/api/users').then(results => {
                return results.data.data;
            });
        },
        get: function(ID) {
            return $api.get('/api/users/user/' + ID).then(results => {
                return results.data.data;
            });
        },
        toggle: function(user) {
            return $q(function(resolve) {
                return $api.post('/api/users/user/' + user.ID + '/disable/').then(function() {
                    notify.success('User Saved');
                    resolve();
                });
            });
        },
        delete: function(user) {
            return $q(function(resolve) {
                popup.confirm('Delete User', 'Are you sure you want to delete the user "' + user.Name + '"?', ['Delete', 'Cancel']).then(result => {
                    if (result) {
                        $api.delete('/api/users/user/' + user.ID).then(function() {
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
        update: function(ID, params) {
            return $api.post('/api/users/user/' + ID, params).then(results => {
                return results.data.data;
            });
        },
    };
});
