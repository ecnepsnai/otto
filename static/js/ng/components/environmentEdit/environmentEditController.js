angular.module('otto').controller('environmentEdit', function($scope, popup, notify, truncate) {
    var $ctrl = this;

    $scope.$watch('$ctrl.environment', (environment) => {
        if (!environment) {
            $ctrl.environmentListSorted = [];
            return;
        }
        var keys = Object.keys(environment).sort();
        var environmentListSorted = [];
        keys.forEach((key) => {
            environmentListSorted.push({
                Key: key,
                Value: environment[key],
            });
        });
        $ctrl.environmentListSorted = environmentListSorted;
    }, true);

    $ctrl.newEnvironment = () => {
        popup.new({
            template: '<environment-popup></environment-popup>',
        }).then((result) => {
            if (!result) {
                return;
            }

            if (!result.Key || result.Key === '') {
                notify.error('Invalid Environment Variable');
                return;
            }

            if ($ctrl.environment[result.Key] !== undefined) {
                notify.error('Duplicate Environment Variable');
                return;
            }

            $ctrl.environment[result.Key] = result.Value;
        });
    };

    $ctrl.editEnvironment = (environment) => {
        popup.new({
            template: '<environment-popup></environment-popup>',
            data: {
                environment: environment
            }
        }).then((result) => {
            if (!result) {
                return;
            }

            if (!result.Key || result.Key === '') {
                notify.error('Invalid Environment Variable');
                return;
            }

            $ctrl.environment[result.Key] = result.Value;
        });
    };

    $ctrl.deleteEnvironment = (environment) => {
        popup.confirm(
            'Delete Environment',
            'Are you sure you want to delete this environment?',
            ['Delete', 'Cancel']
        ).then((result) => {
            if (result) {
                delete $ctrl.environment[environment.Key];
            }
        });
    };
});
