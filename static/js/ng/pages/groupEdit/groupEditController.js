angular.module('otto').controller('groupEdit', function($route, $group, $location, title, notify) {
    var $ctrl = this;
    $ctrl.loading = false;

    if ($location.path() === '/groups/group/') {
        $ctrl.isNew = true;
        $ctrl.title = 'New Group';
        title.set($ctrl.title);
        $ctrl.group = {
            Environment: {},
            ScriptIDs: [],
        };
        $ctrl.loaded = true;
    } else {
        $ctrl.title = 'Edit Group';
        title.set($ctrl.title);
        var id = $route.current.params.id;
        $group.get(id).then(group => {
            $ctrl.group = group;
            $ctrl.loaded = true;
        });
        $ctrl.isNew = false;
    }

    $ctrl.save = function(isValid) {
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
            $ctrl.loading = false;
            $location.url('/groups/group/' + group.ID);
            notify.success('Group Saved');
        }, function() {
            $ctrl.loading = false;
        });
    };
});