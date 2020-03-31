angular.module('otto').controller('userList', function($user, state) {
    var $ctrl = this;

    $ctrl.loadData = function() {
        $ctrl.loading = true;
        $user.list().then(users => {
            $ctrl.users = users;
            $ctrl.loading = false;
        });
    };
    $ctrl.loadData();

    $ctrl.canDeleteUser = function(user) {
        return user.Username !== state.current().User.Username;
    };
});
