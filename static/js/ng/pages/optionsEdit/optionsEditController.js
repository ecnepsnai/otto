angular.module('otto').controller('optionsEdit', function($scope, $api, $group, $q, notify, state, title) {
    var $ctrl = this;
    title.set('Options');
    $ctrl.state = state.current();
    $ctrl.urlPlaceholder = location.href.replace('/options/', '') + '/';

    $ctrl.loadData = () => {
        $ctrl.loading = true;
        return $q.all({
            groups: $group.list(),
            options: $api.get('/api/options'),
        }).then(function(results) {
            $ctrl.groups = results.groups;
            var options = results.options.data.data;
            $ctrl.originalConfig = angular.copy(options);
            $ctrl.options = options;
            $ctrl.loading = false;
        });
    };
    $ctrl.loadData();

    $scope.$watch('$ctrl.options.Register.Enabled', function(nv, ov) {
        if (nv === ov) {
            return;
        }

        if (nv && !$ctrl.options.Register.DefaultGroupID) {
            $ctrl.options.Register.DefaultGroupID = $ctrl.groups[0].ID;
        }
    });

    $ctrl.save = function(valid) {
        if (!valid) {
            return;
        }

        $ctrl.loading = true;
        $api.post('/api/options', $ctrl.options).then(() => {
            state.invalidate().then(() => {
                $ctrl.state = state.current();
                $ctrl.loading = false;
                notify.success('Options Updated');
            });
        }, () => {
            $ctrl.loading = false;
        });
    };
});
