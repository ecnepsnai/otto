angular.module('otto').controller('hostList', function($q, $heartbeat, $host, title) {
    var $ctrl = this;
    title.set('Hosts');

    $ctrl.loadData = () => {
        $ctrl.loading = true;
        return $q.all({
            hosts: $host.list(),
            heartbeats: $heartbeat.list().then(hosts => {
                var results = {};
                hosts.forEach(function(hb) {
                    results[hb.Address] = hb;
                });
                return results;
            })
        });
    };

    $ctrl.loadAll = () => {
        $ctrl.loadData().then(function(results) {
            $ctrl.loading = false;
            $ctrl.hosts = results.hosts;
            $ctrl.heartbeats = results.heartbeats;
        });
    };
    $ctrl.loadAll();

    $ctrl.deleteHost = function(host) {
        $host.delete(host).then(() => {
            $ctrl.loadAll();
        });
    };
});
