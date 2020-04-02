angular.module('otto').controller('userList', function($user, popup, state) {
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
        }).then((result) => {
            if (!result) {
                return;
            }

            $user.new(result).then(() => {
                notify.success('Changed applied');
                $ctrl.loadData();
            });
        });
    };

    $ctrl.editUser = (user) => {
        popup.new({
            template: '<user-edit></user-edit>',
            data: {
                user: angular.copy(user)
            }
        }).then((result) => {
            if (!result) {
                return;
            }

            $user.update(user.Username, result).then(() => {
                notify.success('Changed applied');
                $ctrl.loadData();
            });
        });
    };

    $ctrl.deleteUser = (user) => {
        $user.delete(user).then(() => {
            $ctrl.loadData();
        });
    };

    $ctrl.canDeleteUser = (user) => {
        return user.Username !== state.current().User.Username;
    };
});
