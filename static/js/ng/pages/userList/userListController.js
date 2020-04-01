angular.module('otto').controller('userList', function($user, popup, state) {
    var $ctrl = this;

    $ctrl.loadData = function() {
        $ctrl.loading = true;
        $user.list().then(users => {
            $ctrl.users = users;
            $ctrl.loading = false;
        });
    };
    $ctrl.loadData();

    $ctrl.newUser = function() {
        popup.new({
            template: '<user-edit></user-edit>',
            data: {},
        }).then(function(result) {
            if (!result) {
                return;
            }

            $user.new(result).then(function() {
                notify.success('Changed applied');
                $ctrl.loadData();
            });
        });
    };

    $ctrl.editUser = function(user) {
        popup.new({
            template: '<user-edit></user-edit>',
            data: {
                user: angular.copy(user)
            }
        }).then(function(result) {
            if (!result) {
                return;
            }

            $user.update(user.Username, result).then(function() {
                notify.success('Changed applied');
                $ctrl.loadData();
            });
        });
    };

    $ctrl.deleteUser = function(user) {
        $user.delete(user).then(function() {
            $ctrl.loadData();
        });
    };

    $ctrl.canDeleteUser = function(user) {
        return user.Username !== state.current().User.Username;
    };
});
