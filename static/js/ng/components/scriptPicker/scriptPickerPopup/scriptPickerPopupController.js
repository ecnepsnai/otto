angular.module('otto').controller('scriptPickerPopup', function($scope) {
    var $ctrl = this;
    var $popupScope = $scope.$parent;
    var $popupData = $popupScope.popupData;

    $ctrl.selectedScripts = $popupData.selected;
    $ctrl.scripts = $popupData.scripts;

    $ctrl.response = (apply) => {
        if (apply) {
            $popupScope.popupResolve($ctrl.selectedScripts);
        } else {
            $popupScope.popupResolve(false);
        }
        $popupScope.popupElement.modal('hide');
    };
});
