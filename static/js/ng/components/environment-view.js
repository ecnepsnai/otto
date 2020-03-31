angular.module('otto').component('environmentView', {
    bindings: {
        envs: '<'
    },
    controllerAs: '',
    template: '<nothing ng-if="$ctrl.envs.length === 0"></nothing><ul class="list-group list-group-flush"><li class="list-group-item" ng-repeat="environment in $ctrl.envs">{{ environment.Key }}: <code class="monospace">{{ environment.Value | truncate }}</code></li></ul>'
});