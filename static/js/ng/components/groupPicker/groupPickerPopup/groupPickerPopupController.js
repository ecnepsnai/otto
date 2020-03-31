angular.module('otto').controller('groupPickerPopup', function($scope) {
    var $ctrl = this;
    var $popupScope = $scope.$parent;
    var $popupData = $popupScope.popupData;

    $ctrl.selectedGroups = $popupData.selected;
    $ctrl.groups = $popupData.groups;

    $ctrl.response = function(apply) {
        if (apply) {
            $popupScope.popupResolve($ctrl.selectedGroups);
        } else {
            $popupScope.popupResolve(false);
        }
        $popupScope.popupElement.modal('hide');
    };
});
