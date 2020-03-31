angular.module('otto').controller('hostView', function($q, $host, $group, $location, $route, title) {
    var $ctrl = this;
    var id = $route.current.params.id;

    $ctrl.getData = function() {
        return $q.all({
            host: $host.get(id),
            scripts: $host.getScripts(id),
            groups: $group.list(),
        });
    };

    $ctrl.deleteHost = function() {
        $host.delete($ctrl.host).then(function() {
            $location.url('/hosts/');
        });
    };

    title.set('View Host');
    $ctrl.loaded = false;
    $ctrl.getData().then(function(result) {
        $ctrl.host = result.host;
        $ctrl.scripts = result.scripts;
        var groupMap = {};
        result.groups.forEach((group) => {
            groupMap[group.ID] = group;
        });
        $ctrl.groups = [];
        $ctrl.host.GroupIDs.forEach((groupID) => {
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
