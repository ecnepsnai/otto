angular.module('otto').controller('groupView', function($q, $group, $location, $route, title) {
    var $ctrl = this;
    var id = $route.current.params.id;

    $ctrl.getData = function() {
        return $q.all({
            group: $group.get(id),
            hosts: $group.getHosts(id),
            scripts: $group.getScripts(id),
        });
    };

    $ctrl.deleteGroup = function() {
        $group.delete($ctrl.group).then(function() {
            $location.url('/groups/');
        });
    };

    title.set('View Group');
    $ctrl.loaded = false;
    $ctrl.getData().then(function(results) {
        $ctrl.group = results.group;
        $ctrl.hosts = results.hosts;
        $ctrl.scripts = results.scripts;
        title.set('View Group: ' + $ctrl.group.Name);

        var keys = Object.keys($ctrl.group.Environment).sort();
        var environmentListSorted = [];
        keys.forEach(function(key) {
            environmentListSorted.push({
                Key: key,
                Value: $ctrl.group.Environment[key],
            });
        });
        $ctrl.environmentListSorted = environmentListSorted;
        $ctrl.loaded = true;
    });
});
