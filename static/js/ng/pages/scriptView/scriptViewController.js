angular.module('otto').controller('scriptView', function($q, $script, $route, $location, title) {
    var $ctrl = this;
    var id = $route.current.params.id;

    $ctrl.getData = () => {
        return $q.all({
            script: $script.get(id),
            hosts: $script.getHosts(id)
        });
    };

    $ctrl.deleteScript = () => {
        $script.delete($ctrl.script).then(() => {
            $location.url('/scripts/');
        });
    };

    title.set('View Script');
    $ctrl.loaded = false;
    $ctrl.getData().then(function(result) {
        $ctrl.script = result.script;
        $ctrl.hosts = result.hosts;
        title.set('View Script: ' + $ctrl.script.Name);

        var keys = Object.keys($ctrl.script.Environment).sort();
        var environmentListSorted = [];
        keys.forEach(function(key) {
            environmentListSorted.push({
                Key: key,
                Value: $ctrl.script.Environment[key],
            });
        });
        $ctrl.environmentListSorted = environmentListSorted;
        $ctrl.loaded = true;
    });
});
