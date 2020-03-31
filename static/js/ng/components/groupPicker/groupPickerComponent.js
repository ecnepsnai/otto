angular.module('otto').component('groupPicker', {
    templateUrl: '/ottodev/static/html/groupPicker.html',
    bindings: {
        model: '=',
        max: '<',
        min: '<'
    },
    controller: 'groupPicker',
    controllerAs: ''
});
