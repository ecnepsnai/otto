angular.module('otto').component('confirmPopup', {
    bindings: {},
    controller: function($scope) {
        var $ctrl = this;
        this.$onInit = function () {
            var $popupScope = $scope.$parent;
            $ctrl.title = $popupScope.popupData.title;
            $ctrl.body = $popupScope.popupData.body;
            $ctrl.choices = $popupScope.popupData.choices;
            $ctrl.response = function(response) {
                $popupScope.popupResolve(response);
                $popupScope.popupElement.modal('hide');
            };
        };
    },
    controllerAs: '',
    template: '<div class="modal-header"><h5 class="modal-title">{{:: $ctrl.title }}</h5></div>' +
        '<div class="modal-body"><p>{{:: $ctrl.body }}</p></div>' +
        '<div class="modal-footer">' +
        '<button type="button" class="btn btn-danger" ng-click="$ctrl.response(true)">{{:: $ctrl.choices[0] }}</button>' +
        '<button type="button" class="btn btn-secondary" ng-click="$ctrl.response(false)">{{:: $ctrl.choices[1] }}</button>' +
        '</div>'
});