angular.module('otto').controller('userList', function($user, popup, state, notify) {
    var $ctrl = this;

    $ctrl.loadData = () => {
        $ctrl.loading = true;
        $user.list().then(users => {
            $ctrl.users = users;
            $ctrl.loading = false;
        });
    };
    $ctrl.loadData();

    $ctrl.newUser = () => {
        popup.new({
            template: '<user-edit></user-edit>',
            data: {},
        }).then(function(result) {
            if (!result) {
                return;
            }

            $user.new(result).then(() => {
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

            $user.update(user.Username, result).then(() => {
                notify.success('Changed applied');
                $ctrl.loadData();
            });
        });
    };

    $ctrl.deleteUser = function(user) {
        $user.delete(user).then(() => {
            $ctrl.loadData();
        });
    };

    $ctrl.canDeleteUser = function(user) {
        return user.Username !== state.current().User.Username;
    };
});
