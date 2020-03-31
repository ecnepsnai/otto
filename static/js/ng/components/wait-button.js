angular.module('otto').component('waitButton', {
    bindings: {
        title: '@',
        loading: '=',
        click: '&',
        icon: '@',
        buttonClass: '@',
        buttonDisabled: '<',
    },
    template: '<button type="button" class="btn" ng-class="$ctrl.buttonClassImpl()" ng-disabled="$ctrl.isDisabled()" ng-click="$ctrl.click()"><span ng-if="!$ctrl.loading"><i ng-class="$ctrl.iconClass" ng-if="$ctrl.showIcon"></i>{{:: $ctrl.title }}</span><span ng-if="$ctrl.loading"><i class="fas fa-spinner fa-pulse"></i> Loading...</span></button>',
    controllerAs: '',
    controller: function() {
        var $ctrl = this;

        $ctrl.buttonClassImpl = function() {
            var cls = {
                btn: true,
            };
            if ($ctrl.buttonClass) {
                cls[$ctrl.buttonClass] = true;
            } else {
                cls['btn-primary'] = true;
            }
            return cls;
        };

        $ctrl.$onInit = function() {
            if ($ctrl.icon) {
                $ctrl.iconClass = {
                    'pr-1': true,
                };
                $ctrl.iconClass[$ctrl.icon] = true;
                $ctrl.showIcon = true;
            }
        };

        $ctrl.isDisabled = function() {
            var disabled = false;
            if ($ctrl.loading) {
                disabled = true;
            }
            if ($ctrl.buttonDisabled) {
                disabled = true;
            }
            return disabled;
        };
    }
});
