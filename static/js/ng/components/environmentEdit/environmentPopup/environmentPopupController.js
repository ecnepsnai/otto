angular.module('otto').controller('environmentPopup', function($scope) {
    var $ctrl = this;

    var popupScope = $scope.$parent;
    var popupData = popupScope.popupData;

    if (popupData) {
        $ctrl.title = 'Edit Environment Variable';
        $ctrl.isNew = false;
        $ctrl.environment = angular.copy(popupData.environment);
    } else {
        $ctrl.title = 'New Environment Variable';
        $ctrl.isNew = true;
        $ctrl.environment = {
            Key: '',
            Value: ''
        };
    }

    $ctrl.dismiss = function(apply) {
        if (apply) {
            popupScope.popupResolve($ctrl.environment);
        }
        popupScope.popupElement.modal('hide');
    };
});
