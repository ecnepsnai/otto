angular.module('otto').controller('scriptEdit', function($route, $q, $script, $location, title, notify) {
    var $ctrl = this;
    $ctrl.loaded = false;
    $ctrl.enabledGroups = [];

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

    function getScriptGroups() {
        $ctrl.loaded = false;
        if ($location.path() === '/scripts/script/') {
            return $q.resolve([]);
        } else {
            return $script.getGroups($route.current.params.id);
        }
    }

    $q.all({ script: getScript(), hosts: getScriptGroups() }).then(function(results) {
        $ctrl.script = results.script;
        results.hosts.forEach(function(host) {
            $ctrl.enabledGroups.push(host.ID);
        });

        if ($location.path() === '/scripts/script/') {
            $ctrl.isNew = true;
            $ctrl.title = 'New Script';
            title.set($ctrl.title);
        } else {
            $ctrl.isNew = false;
            $ctrl.title = 'Edit Script: ' + $ctrl.script.Name;
            title.set($ctrl.title);
        }

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
            $script.setGroups(script.ID, { Groups: $ctrl.enabledGroups }).then(function() {
                $ctrl.loading = false;
                $location.url('/scripts/script/' + script.ID);
                notify.success('Script Saved');
            });
        }, function() {
            $ctrl.loading = false;
        });
    };
});