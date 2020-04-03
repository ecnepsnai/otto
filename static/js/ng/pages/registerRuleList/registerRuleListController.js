angular.module('otto').controller('registerRuleList', function(popup) {
    var $ctrl = this;
    $ctrl.groupNameMap = {};

    this.$onInit = function () {
        $ctrl.groups.forEach(function(group) {
            $ctrl.groupNameMap[group.ID] = group.Name;
        });
        $ctrl.loadRules();
    };

    $ctrl.loadRules = () => {
        if (!$ctrl.rules) {
            $ctrl.rules = [];
        }
        $ctrl.rules.forEach(function(rule) {
            if (rule.Uname) {
                rule.Property = 'Uname';
                rule.Matches = rule.Uname;
            } else if (rule.Hostname) {
                rule.Property = 'Hostname';
                rule.Matches = rule.Hostname;
            } else {
                rule.Property = 'Unknown';
            }
            rule.GroupName = $ctrl.groupNameMap[rule.GroupID];
        });
    };

    $ctrl.newRule = () => {
        popup.new({
            template: '<register-rule-edit></register-rule-edit>',
            data: {
                groups: $ctrl.groups,
            },
        }).then(function(result) {
            if (result) {
                $ctrl.rules.push(result);
                $ctrl.loadRules();
            }
        });
    };

    $ctrl.editRule = function(index) {
        var rule = angular.copy($ctrl.rules[index]);
        popup.new({
            template: '<register-rule-edit></register-rule-edit>',
            data: {
                rule: rule,
                groups: $ctrl.groups,
            },
        }).then(function(result) {
            if (result) {
                $ctrl.rules[index] = rule;
                $ctrl.loadRules();
            }
        });
    };

    $ctrl.deleteRule = function(index) {
        popup.confirm('Delete Register Rule', 'Are you sure you want to delete this register rule?', ['Delete', 'Cancel']).then(function(result) {
            if (result) {
                $ctrl.rules.splice(index, 1);
                $ctrl.loadRules();
            }
        });
    };

});
