angular.module('otto').component('registerRuleList', {
    templateUrl: '/ottodev/static/html/registerRuleList.html',
    bindings: {
        rules: '=',
        groups: '<'
    },
    controller: 'registerRuleList',
    controllerAs: ''
});
