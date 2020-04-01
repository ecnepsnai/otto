angular.module('otto').controller('userEdit', function ($scope, state) {
    var $ctrl = this;
    var $popupScope = $scope.$parent;
    var $popupData = $popupScope.popupData;
    $ctrl.user = {
        Enabled: true
    };
    $ctrl.title = 'New User';
    $ctrl.canDisableUser = true;
    $ctrl.showPasswordField = true;
    $ctrl.isNew = true;
    if ($popupData.user) {
        $ctrl.user = $popupData.user;
        $ctrl.title = 'Edit User';
        $ctrl.canDisableUser = state.current().User.Username !== $ctrl.user.Username;
        $ctrl.showPasswordField = false;
        $ctrl.isNew = false;
    }

    $ctrl.resetPassword = function() {
        $ctrl.showPasswordField = true;
    };

    $ctrl.response = function(apply) {
        if (apply) {
            $popupScope.popupResolve($ctrl.user);
        } else {
            $popupScope.popupResolve(false);
        }
        $popupScope.popupElement.modal('hide');
    };
});