angular.module('otto').controller('hostEdit', function($route, $host, $location, title, notify) {
    var $ctrl = this;
    $ctrl.loading = false;

    if ($location.path() === '/hosts/host/') {
        $ctrl.isNew = true;
        $ctrl.title = 'New Host';
        title.set($ctrl.title);
        $ctrl.useNameAsAddress = true;
        $ctrl.host = {
            Port: 12444,
            Environment: {},
            GroupIDs: [],
        };
        $ctrl.loaded = true;
    } else {
        $ctrl.title = 'Edit Host';
        var id = $route.current.params.id;
        $host.get(id).then(host => {
            $ctrl.host = host;
            title.set('Edit Host: ' + host.Name);
            $ctrl.loaded = true;
            $ctrl.useNameAsAddress = host.Name === host.Address;
        });
        $ctrl.isNew = false;
    }

    $ctrl.save = function(isValid) {
        if (!isValid) {
            return;
        }

        if ($ctrl.useNameAsAddress) {
            $ctrl.host.Address = $ctrl.host.Name;
        }

        var savePromise;
        if ($ctrl.isNew) {
            savePromise = $host.new($ctrl.host);
        } else {
            savePromise = $host.update($ctrl.host.ID, $ctrl.host);
        }
        $ctrl.loading = true;
        savePromise.then(host => {
            $ctrl.loading = false;
            $location.url('/hosts/host/' + host.ID);
            notify.success('Host Saved');
        }, function() {
            $ctrl.loading = false;
        });
    };
});