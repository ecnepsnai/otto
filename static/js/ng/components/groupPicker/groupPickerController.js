angular.module('otto').controller('groupPicker', function($scope, $group, popup) {
    var $ctrl = this;
    $ctrl.loading = true;

    $scope.$watch('$ctrl.model', function(model) {
        if (model === null || model === undefined) {
            return;
        }

        $ctrl.selectedGroups = {};
        model.forEach(function(groupID) {
            $ctrl.selectedGroups[groupID] = true;
        });
    });

    $group.list().then(groups => {
        $ctrl.groups = groups;
        $ctrl.buttonText = 'Select Groups';
        $ctrl.loading = false;
    });

    $ctrl.showPopup = () => {
        popup.new({
            template: '<group-picker-popup></group-picker-popup>',
            data: {
                selected: angular.copy($ctrl.selectedGroups),
                groups: $ctrl.groups
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
                notify.error('Invalid Selection', 'Requires at least ' + $ctrl.min + ' enabled group(s)');
                return;
            } else if (selected.count > $ctrl.max) {
                notify.error('Invalid Selection', 'Requires no more than ' + $ctrl.max + ' enabled group(s)');
                return;
            }

            $ctrl.model = selected;
        });
    };
});
