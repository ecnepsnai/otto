angular.module('otto').controller('hostPickerPopup', function($scope) {
    var $ctrl = this;
    var $popupScope = $scope.$parent;
    var $popupData = $popupScope.popupData;

    $ctrl.selectedHosts = $popupData.selected;
    $ctrl.hosts = $popupData.hosts;

    $ctrl.response = (apply) => {
        if (apply) {
            $popupScope.popupResolve($ctrl.selectedHosts);
        } else {
            $popupScope.popupResolve(false);
        }
        $popupScope.popupElement.modal('hide');
    };
});
