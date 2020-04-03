angular.module('otto').controller('hostView', function($q, $host, $group, $location, $route, title) {
    var $ctrl = this;
    var id = $route.current.params.id;

    $ctrl.getData = () => {
        return $q.all({
            host: $host.get(id),
            scripts: $host.getScripts(id),
            groups: $group.list(),
        });
    };

    $ctrl.deleteHost = () => {
        $host.delete($ctrl.host).then(() => {
            $location.url('/hosts/');
        });
    };

    title.set('View Host');
    $ctrl.loaded = false;
    $ctrl.getData().then(function(result) {
        $ctrl.host = result.host;
        $ctrl.scripts = result.scripts;
        var groupMap = {};
        result.groups.forEach(function(group) {
            groupMap[group.ID] = group;
        });
        $ctrl.groups = [];
        ($ctrl.host.GroupIDs || []).forEach(function(groupID) {
            $ctrl.groups.push(groupMap[groupID]);
        });

        title.set('View Host: ' + $ctrl.host.Name);

        var keys = Object.keys(($ctrl.host.Environment || {})).sort();
        var environmentListSorted = [];
        keys.forEach(function(key) {
            environmentListSorted.push({
                Key: key,
                Value: $ctrl.host.Environment[key],
            });
        });
        $ctrl.environmentListSorted = environmentListSorted;
        $ctrl.loaded = true;
    });
});
