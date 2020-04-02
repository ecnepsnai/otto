angular.module('otto').factory('$group', function($api, popup, notify, $q) {
    return {
        list: () => {
            return $api.get('/api/groups').then(results => {
                return results.data.data;
            });
        },
        membership: () => {
            return $api.get('/api/groups/membership').then(results => {
                return results.data.data;
            });
        },
        get: (ID) => {
            return $api.get('/api/groups/group/' + ID).then(results => {
                return results.data.data;
            });
        },
        getHosts: (ID) => {
            return $api.get('/api/groups/group/' + ID + '/hosts').then(results => {
                return results.data.data;
            });
        },
        setHosts: (ID, hosts) => {
            return $api.post('/api/groups/group/' + ID + '/hosts', hosts).then(results => {
                return results.data.data;
            });
        },
        getScripts: (ID) => {
            return $api.get('/api/groups/group/' + ID + '/scripts').then(results => {
                return results.data.data;
            });
        },
        toggle: (group) => {
            return $q((resolve) => {
                return $api.post('/api/groups/group/' + group.ID + '/disable/').then(() => {
                    notify.success('Host Saved');
                    resolve();
                });
            });
        },
        delete: (group) => {
            return $q((resolve) => {
                popup.confirm('Delete Host', 'Are you sure you want to delete the group "' + group.Name + '"?', ['Delete', 'Cancel']).then(result => {
                    if (result) {
                        $api.delete('/api/groups/group/' + group.ID).then(() => {
                            notify.success('Host Deleted');
                            resolve();
                        });
                    }
                });
            });
        },
        new: (params) => {
            return $api.put('/api/groups/group', params).then(results => {
                return results.data.data;
            });
        },
        update: (ID, params) => {
            return $api.post('/api/groups/group/' + ID, params).then(results => {
                return results.data.data;
            });
        },
    };
});
