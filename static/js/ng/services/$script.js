angular.module('otto').factory('$script', function($api, popup, notify, $q) {
    return {
        list: () => {
            return $api.get('/api/scripts').then(results => {
                return results.data.data;
            });
        },
        get: (ID) => {
            return $api.get('/api/scripts/script/' + ID).then(results => {
                return results.data.data;
            });
        },
        getGroups: (ID) => {
            return $api.get('/api/scripts/script/' + ID + '/groups').then(results => {
                return results.data.data;
            });
        },
        getHosts: (ID) => {
            return $api.get('/api/scripts/script/' + ID + '/hosts').then(results => {
                return results.data.data;
            });
        },
        setGroups: (ID, groups) => {
            return $api.post('/api/scripts/script/' + ID + '/groups', groups).then(results => {
                return results.data.data;
            });
        },
        toggle: (script) => {
            return $q((resolve) => {
                return $api.post('/api/scripts/script/' + script.ID + '/disable/').then(() => {
                    notify.success('Script Saved');
                    resolve();
                });
            });
        },
        delete: (script) => {
            return $q((resolve) => {
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
        new: (params) => {
            return $api.put('/api/scripts/script', params).then(results => {
                return results.data.data;
            });
        },
        update: (ID, params) => {
            return $api.post('/api/scripts/script/' + ID, params).then(results => {
                return results.data.data;
            });
        }
    };
});
