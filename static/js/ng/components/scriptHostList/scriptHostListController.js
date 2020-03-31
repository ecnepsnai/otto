angular.module('otto').controller('scriptHostList', function($scope) {
    var $ctrl = this;

    this.$onInit = function() {
        $ctrl.groupMap = {};
        $ctrl.hostsByGroup = {};
        $ctrl.hosts.forEach(host => {
            $ctrl.groupMap[host.GroupID] = {
                GroupIP: host.GroupID,
                GroupName: host.GroupName
            };

            var hosts = [];
            if ($ctrl.hostsByGroup[host.GroupID]) {
                hosts = $ctrl.hostsByGroup[host.GroupID];
            }
            hosts.push(host);
            $ctrl.hostsByGroup[host.GroupID] = hosts;
        });
        $ctrl.results = [];
        Object.keys($ctrl.hostsByGroup).forEach(groupID => {
            $ctrl.results.push({
                GroupID: groupID,
                GroupName: $ctrl.groupMap[groupID].GroupName,
                Hosts: $ctrl.hostsByGroup[groupID]
            });
        });
    };
});
