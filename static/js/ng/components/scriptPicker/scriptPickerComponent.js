angular.module('otto').component('scriptPicker', {
    templateUrl: '/ottodev/static/html/scriptPicker.html',
    bindings: {
        model: '=',
        max: '<',
        min: '<'
    },
    controller: 'scriptPicker',
    controllerAs: ''
});
