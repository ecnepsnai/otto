angular.module('otto').controller('groupEdit', function($route, $group, $q, $location, title, notify) {
    var $ctrl = this;
    $ctrl.loaded = false;

    $ctrl.getData = () => {
        if ($location.path() === '/groups/group/') {
            return $q.all({
                group: $q.resolve({
                    Environment: {},
                    ScriptIDs: [],
                }),
                hosts: $q.resolve([]),
            });
        } else {
            var id = $route.current.params.id;
            return $q.all({
                group: $group.get(id),
                hosts: $group.getHosts(id),
            });
        }
    };

    $ctrl.getData().then((results) => {
        $ctrl.group = results.group;
        $ctrl.selectedHosts = results.hosts.map((host) => {
            return host.ID;
        });
        $ctrl.loaded = true;

        if ($location.path() === '/groups/group/') {
            $ctrl.isNew = true;
            $ctrl.title = 'New Group';
            title.set($ctrl.title);
        } else {
            $ctrl.isNew = false;
            $ctrl.title = 'Edit Group: ' + $ctrl.group.Name;
            title.set($ctrl.title);
        }
    });

    $ctrl.save = (isValid) => {
        if (!isValid) {
            return;
        }

        var savePromise;
        if ($ctrl.isNew) {
            savePromise = $group.new($ctrl.group);
        } else {
            savePromise = $group.update($ctrl.group.ID, $ctrl.group);
        }
        $ctrl.loading = true;
        savePromise.then(group => {
            $group.setHosts(group.ID, { Hosts: $ctrl.selectedHosts }).then(() => {
                $ctrl.loading = false;
                $location.url('/groups/group/' + group.ID);
                notify.success('Group Saved');
            });
        }, () => {
            $ctrl.loading = false;
        });
    };
});