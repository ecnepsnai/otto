angular.module('otto').factory('$user', function($api, popup, notify, $q) {
    return {
        list: () => {
            return $api.get('/api/users').then(results => {
                return results.data.data;
            });
        },
        get: (username) => {
            return $api.get('/api/users/user/' + username).then(results => {
                return results.data.data;
            });
        },
        toggle: (user) => {
            return $q((resolve) => {
                return $api.post('/api/users/user/' + user.Username + '/disable/').then(() => {
                    notify.success('User Saved');
                    resolve();
                });
            });
        },
        delete: (user) => {
            return $q((resolve) => {
                popup.confirm('Delete User', 'Are you sure you want to delete the user "' + user.Username + '"?', ['Delete', 'Cancel']).then(result => {
                    if (result) {
                        $api.delete('/api/users/user/' + user.Username).then(() => {
                            notify.success('User Deleted');
                            resolve();
                        });
                    }
                });
            });
        },
        new: (params) => {
            return $api.put('/api/users/user', params).then(results => {
                return results.data.data;
            });
        },
        update: (username, params) => {
            return $api.post('/api/users/user/' + username, params).then(results => {
                return results.data.data;
            });
        },
    };
});
