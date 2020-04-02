angular.module('otto').controller('scriptPicker', function($scope, $script, popup) {
    var $ctrl = this;
    $ctrl.loading = true;

    $scope.$watch('$ctrl.model', (model) => {
        if (model === null || model === undefined) {
            return;
        }

        $ctrl.selectedScripts = {};
        model.forEach((scriptID) => {
            $ctrl.selectedScripts[scriptID] = true;
        });
    });

    $script.list().then(scripts => {
        $ctrl.scripts = scripts;
        $ctrl.buttonText = 'Select Scripts';
        $ctrl.loading = false;
    });

    $ctrl.showPopup = () => {
        popup.new({
            template: '<script-picker-popup></script-picker-popup>',
            data: {
                selected: angular.copy($ctrl.selectedScripts),
                scripts: $ctrl.scripts
            }
        }).then((result) => {
            if (result === false || result === undefined) {
                return;
            }

            var selected = [];
            Object.keys(result).forEach((key) => {
                if (result[key]) {
                    selected.push(key);
                }
            });

            if (selected.count < $ctrl.min) {
                notify.error('Invalid Selection', 'Requires at least ' + $ctrl.min + ' enabled script(s)');
                return;
            } else if (selected.count > $ctrl.max) {
                notify.error('Invalid Selection', 'Requires no more than ' + $ctrl.max + ' enabled script(s)');
                return;
            }

            $ctrl.model = selected;
        });
    };
});
