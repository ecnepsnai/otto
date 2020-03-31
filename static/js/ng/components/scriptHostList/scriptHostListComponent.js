angular.module('otto').component('scriptHostList', {
    templateUrl: '/ottodev/static/html/scriptHostList.html',
    bindings: {
        script: '<',
        hosts: '<'
    },
    controller: 'scriptHostList',
    controllerAs: ''
});
