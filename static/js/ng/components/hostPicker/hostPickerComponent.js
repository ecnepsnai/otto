angular.module('otto').component('hostPicker', {
    templateUrl: '/ottodev/static/html/hostPicker.html',
    bindings: {
        model: '=',
        max: '<',
        min: '<'
    },
    controller: 'hostPicker',
    controllerAs: ''
});
