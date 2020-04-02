angular.module('otto').factory('$host', function($api, popup, notify, $q) {
    return {
        list: () => {
            return $api.get('/api/hosts').then(results => {
                return results.data.data;
            });
        },
        get: (ID) => {
            return $api.get('/api/hosts/host/' + ID).then(results => {
                return results.data.data;
            });
        },
        getScripts: (ID) => {
            return $api.get('/api/hosts/host/' + ID + '/scripts').then(results => {
                return results.data.data;
            });
        },
        toggle: (host) => {
            return $q((resolve) => {
                return $api.post('/api/hosts/host/' + host.ID + '/disable/').then(() => {
                    notify.success('Host Saved');
                    resolve();
                });
            });
        },
        delete: (host) => {
            return $q((resolve) => {
                popup.confirm('Delete Host', 'Are you sure you want to delete the host "' + host.Name + '"?', ['Delete', 'Cancel']).then(result => {
                    if (result) {
                        $api.delete('/api/hosts/host/' + host.ID).then(() => {
                            notify.success('Host Deleted');
                            resolve();
                        });
                    }
                });
            });
        },
        new: (params) => {
            return $api.put('/api/hosts/host', params).then(results => {
                return results.data.data;
            });
        },
        update: (ID, params) => {
            return $api.post('/api/hosts/host/' + ID, params).then(results => {
                return results.data.data;
            });
        },
    };
});
