angular.module('otto').controller('scriptList', function($script, title) {
    var $ctrl = this;
    title.set('Scripts');

    $ctrl.loadData = () => {
        $ctrl.loading = true;
        $script.list().then(scripts => {
            $ctrl.loading = false;
            $ctrl.scripts = scripts;
        });
    };
    $ctrl.loadData();

    $ctrl.deleteScript = function(script) {
        $script.delete(script).then(() => {
            $ctrl.loadData();
        });
    };
});
