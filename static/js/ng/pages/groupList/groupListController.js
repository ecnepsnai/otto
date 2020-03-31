angular.module('otto').controller('groupList', function($q, $group, title) {
    var $ctrl = this;
    title.set('Groups');

    $ctrl.loadData = function() {
        $ctrl.loading = true;
        return $q.all({
            groups: $group.list(),
            membership: $group.membership()
        });
    };

    $ctrl.loadAll = function() {
        $ctrl.loadData().then(function(results) {
            $ctrl.loading = false;
            $ctrl.groups = results.groups;
            for (i = 0; i < $ctrl.groups.length; i++) {
                $ctrl.groups[i].HostIDs = results.membership[$ctrl.groups[i].ID];
                $ctrl.groups[i].ScriptIDs = ($ctrl.groups[i].ScriptIDs || []);
            }
        });
    };
    $ctrl.loadAll();

    $ctrl.deleteHost = function(group) {
        $group.delete(group).then(function() {
            $ctrl.loadAll();
        });
    };
});
