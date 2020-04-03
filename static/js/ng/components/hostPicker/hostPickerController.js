angular.module('otto').controller('hostPicker', function($scope, $host, popup) {
    var $ctrl = this;
    $ctrl.loading = true;

    $scope.$watch('$ctrl.model', function(model) {
        if (model === null || model === undefined) {
            return;
        }

        $ctrl.selectedHosts = {};
        model.forEach(function(hostID) {
            $ctrl.selectedHosts[hostID] = true;
        });
    });

    $host.list().then(hosts => {
        $ctrl.hosts = hosts;
        $ctrl.buttonText = 'Select Hosts';
        $ctrl.loading = false;
    });

    $ctrl.showPopup = () => {
        popup.new({
            template: '<host-picker-popup></host-picker-popup>',
            data: {
                selected: angular.copy($ctrl.selectedHosts),
                hosts: $ctrl.hosts
            }
        }).then(function(result) {
            if (result === false || result === undefined) {
                return;
            }

            var selected = [];
            Object.keys(result).forEach(function(key) {
                if (result[key]) {
                    selected.push(key);
                }
            });

            if (selected.count < $ctrl.min) {
                notify.error('Invalid Selection', 'Requires at least ' + $ctrl.min + ' enabled host(s)');
                return;
            } else if (selected.count > $ctrl.max) {
                notify.error('Invalid Selection', 'Requires no more than ' + $ctrl.max + ' enabled host(s)');
                return;
            }

            $ctrl.model = selected;
        });
    };
});
