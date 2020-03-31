angular.module('otto').controller('scriptList', function($script) {
    var $ctrl = this;

    $ctrl.loadData = function() {
        $ctrl.loading = true;
        $script.list().then(scripts => {
            $ctrl.loading = false;
            $ctrl.scripts = scripts;
        });
    };
    $ctrl.loadData();

    $ctrl.deleteScript = function(script) {
        $script.delete(script).then(function() {
            $ctrl.loadData();
        });
    };
});
