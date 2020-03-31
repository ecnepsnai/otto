angular.module('otto').component('checkbox', {
    bindings: {
        label: '@',
        model: '=',
        disableIf: '<'
    },
    controllerAs: '',
    template: '<div><input type="checkbox" class="form-check-input" id="{{ $ctrl.id }}" ng-model="$ctrl.model" ng-disabled="{{ $ctrl.disableIf }}"><label class="form-check-label" for="{{ $ctrl.id }}">{{ $ctrl.label }}</label></div>',
    controller: function(rand) {
        var $ctrl = this;
        this.$onInit = function() {
            $ctrl.id = rand.ID();
        };
    }
});
