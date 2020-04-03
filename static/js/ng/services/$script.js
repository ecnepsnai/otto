angular.module('otto').factory('$script', function($api, popup, notify, $q) {
    return {
        list: () => {
            return $api.get('/api/scripts').then(results => {
                return results.data.data;
            });
        },
        get: function(ID) {
            return $api.get('/api/scripts/script/' + ID).then(results => {
                return results.data.data;
            });
        },
        getGroups: function(ID) {
            return $api.get('/api/scripts/script/' + ID + '/groups').then(results => {
                return results.data.data;
            });
        },
        getHosts: function(ID) {
            return $api.get('/api/scripts/script/' + ID + '/hosts').then(results => {
                return results.data.data;
            });
        },
        setGroups: function(ID, groups) {
            return $api.post('/api/scripts/script/' + ID + '/groups', groups).then(results => {
                return results.data.data;
            });
        },
        toggle: function(script) {
            return $q(function(resolve) {
                return $api.post('/api/scripts/script/' + script.ID + '/disable/').then(() => {
                    notify.success('Script Saved');
                    resolve();
                });
            });
        },
        delete: function(script) {
            return $q(function(resolve) {
                popup.confirm('Delete Script', 'Are you sure you want to delete the script "' + script.Name + '"?', ['Delete', 'Cancel']).then(result => {
                    if (result) {
                        $api.delete('/api/scripts/script/' + script.ID).then(() => {
                            notify.success('Script Deleted');
                            resolve();
                        });
                    }
                });
            });
        },
        new: function(params) {
            return $api.put('/api/scripts/script', params).then(results => {
                return results.data.data;
            });
        },
        update: function(ID, params) {
            return $api.post('/api/scripts/script/' + ID, params).then(results => {
                return results.data.data;
            });
        }
    };
});
