angular.module('otto').controller('registerRuleEdit', function($scope) {
    var $ctrl = this;
    var $popupScope = $scope.$parent;
    var $popupData = $popupScope.popupData;
    $ctrl.groups = $popupData.groups;

    if ($popupData.rule) {
        $ctrl.rule = $popupData.rule;
        $ctrl.title = 'Edit Rule';
        if ($ctrl.rule.Uname) {
            $ctrl.ruleProperty = 'Uname';
            $ctrl.ruleValue = $ctrl.rule.Uname;
        } else if ($ctrl.rule.Hostname) {
            $ctrl.ruleProperty = 'Hostname';
            $ctrl.ruleValue = $ctrl.rule.Hostname;
        }
    } else {
        $ctrl.rule = {
            GroupID: $ctrl.groups[0].ID
        };
        $ctrl.ruleProperty = 'Uname';
        $ctrl.ruleValue = '';
        $ctrl.title = 'New Rule';
    }

    $ctrl.response = function(apply) {
        if (apply) {
            delete $ctrl.rule.Uname;
            delete $ctrl.rule.Hostname;
            if ($ctrl.ruleProperty === 'Uname') {
                $ctrl.rule.Uname = $ctrl.ruleValue;
            } else if ($ctrl.ruleProperty === 'Hostname') {
                $ctrl.rule.Hostname = $ctrl.ruleValue;
            }
            $popupScope.popupResolve($ctrl.rule);
        } else {
            $popupScope.popupResolve(false);
        }
        $popupScope.popupElement.modal('hide');
    };
});
