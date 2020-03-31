angular.module('otto').controller('userEdit', function ($user, $route, $q, $location, notify, state, popup) {
    var $ctrl = this;
    var username = $route.current.params.username;

    function getUser() {
        $ctrl.canDisableUser = true;

        if (username) {
            if (username === state.current().User.Username) {
                $ctrl.canDisableUser = false;
            }

            return $user.get(username);
        }
        return $q.resolve({
            Enabled: true,
        });
    }

    getUser().then(user => {
        $ctrl.loaded = true;
        $ctrl.isNew = username === undefined;
        if ($ctrl.isNew) {
            $ctrl.title = 'New User';
            $ctrl.showPassword = true;
            $ctrl.showResetPassword = false;
        } else {
            $ctrl.title = 'Edit User';
            $ctrl.showPassword = false;
            $ctrl.showResetPassword = true;
        }
        $ctrl.user = user;
    });

    $ctrl.resetPassword = function() {
        $ctrl.showPassword = true;
        $ctrl.showResetPassword = false;
    };

    $ctrl.save = function() {
        var savePromise;
        if ($ctrl.isNew) {
            savePromise = $user.create($ctrl.user);
        } else {
            savePromise = $user.update($ctrl.user.Username, $ctrl.user);
        }

        savePromise.then(function() {
            $location.url('/options/');
            notify.common.saved();
        });
    };
});