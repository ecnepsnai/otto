angular.module('otto').factory('$group', function($api, popup, notify, $q) {
    return {
        list: function() {
            return $api.get('/api/groups').then(results => {
                return results.data.data;
            });
        },
        membership: function() {
            return $api.get('/api/groups/membership').then(results => {
                return results.data.data;
            });
        },
        get: function(ID) {
            return $api.get('/api/groups/group/' + ID).then(results => {
                return results.data.data;
            });
        },
        getHosts: function(ID) {
            return $api.get('/api/groups/group/' + ID + '/hosts').then(results => {
                return results.data.data;
            });
        },
        getScripts: function(ID) {
            return $api.get('/api/groups/group/' + ID + '/scripts').then(results => {
                return results.data.data;
            });
        },
        toggle: function(group) {
            return $q(function(resolve) {
                return $api.post('/api/groups/group/' + group.ID + '/disable/').then(function() {
                    notify.success('Host Saved');
                    resolve();
                });
            });
        },
        delete: function(group) {
            return $q(function(resolve) {
                popup.confirm('Delete Host', 'Are you sure you want to delete the group "' + group.Name + '"?', ['Delete', 'Cancel']).then(result => {
                    if (result) {
                        $api.delete('/api/groups/group/' + group.ID).then(function() {
                            notify.success('Host Deleted');
                            resolve();
                        });
                    }
                });
            });
        },
        new: function(params) {
            return $api.put('/api/groups/group', params).then(results => {
                return results.data.data;
            });
        },
        update: function(ID, params) {
            return $api.post('/api/groups/group/' + ID, params).then(results => {
                return results.data.data;
            });
        },
    };
});
