angular.module('otto').controller('scriptEdit', function($route, $q, $script, $location, title, notify) {
    var $ctrl = this;
    $ctrl.loaded = false;
    $ctrl.enabledHosts = [];

    if ($location.path() === '/scripts/script/') {
        $ctrl.isNew = true;
        $ctrl.title = 'New Script';
        title.set($ctrl.title);
    } else {
        $ctrl.isNew = false;
        $ctrl.title = 'Edit Script';
        title.set($ctrl.title);
    }

    function getScript() {
        $ctrl.loaded = false;
        if ($location.path() === '/scripts/script/') {
            return $q.resolve({
                UID: 0,
                GID: 0,
                Environment: {},
                Executable: '/bin/bash',
                AfterExecution: ''
            });
        } else {
            return $script.get($route.current.params.id);
        }
    }

    function getScriptHosts() {
        $ctrl.loaded = false;
        if ($location.path() === '/scripts/script/') {
            return $q.resolve([]);
        } else {
            return $script.getHosts($route.current.params.id);
        }
    }

    $q.all({ script: getScript(), hosts: getScriptHosts() }).then(function(results) {
        $ctrl.script = results.script;
        results.hosts.forEach(function(host) {
            $ctrl.enabledHosts.push(host.ID);
        });
        $ctrl.loaded = true;
    });

    $ctrl.save = function(isValid) {
        if (!isValid) {
            return;
        }

        var savePromise;
        if ($ctrl.isNew) {
            savePromise = $script.new($ctrl.script);
        } else {
            savePromise = $script.update($ctrl.script.ID, $ctrl.script);
        }

        $ctrl.loading = true;
        savePromise.then(function(script) {
            $ctrl.loading = false;

            var id = script.ID;
            $script.setHosts(ID, { Hosts: $ctrl.enabledHosts }).then(function() {
                $location.url('/scripts/script/' + script.ID);
                notify.success('Script Saved');
            });
        }, function() {
            $ctrl.loading = false;
        });
    };
});