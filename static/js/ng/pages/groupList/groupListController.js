angular.module('otto').controller('groupList', function($q, $group, title, notify) {
    var $ctrl = this;
    title.set('Groups');

    $ctrl.loadData = () => {
        $ctrl.loading = true;
        return $q.all({
            groups: $group.list(),
            membership: $group.membership()
        });
    };

    $ctrl.loadAll = () => {
        $ctrl.loadData().then((results) => {
            $ctrl.loading = false;
            $ctrl.groups = results.groups;
            for (i = 0; i < $ctrl.groups.length; i++) {
                $ctrl.groups[i].HostIDs = results.membership[$ctrl.groups[i].ID];
                $ctrl.groups[i].ScriptIDs = ($ctrl.groups[i].ScriptIDs || []);
            }
        });
    };
    $ctrl.loadAll();

    $ctrl.deleteHost = (group) => {
        if (group.HostIDs.length > 0) {
            notify.error('Group must have no host members before it can be deleted', 'Unable to Delete Group');
            return;
        }

        $group.delete(group).then(() => {
            $ctrl.loadAll();
        });
    };
});
