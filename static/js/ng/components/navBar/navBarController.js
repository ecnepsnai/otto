angular.module('otto').controller('navBar', function($api, $route, $location, $timeout, $user, notify, state, popup) {
    var $ctrl = this;
    $ctrl.state = state.current();

    $ctrl.downloadButtonClass = function(href) {
        var matches = $location.path().startsWith(href);
        return { 'btn-outline-dark-light': !matches, 'btn-light': matches };
    };

    $ctrl.navClass = function(tab) {
        var matches = $location.path().startsWith(tab);
        return { active: matches };
    };

    $ctrl.editUser = function() {
        popup.new({
            template: '<user-edit></user-edit>',
            data: {
                user: angular.copy($ctrl.state.User)
            }
        }).then(function(result) {
            if (!result) {
                return;
            }

            $user.update($ctrl.state.User.Username, result).then(function() {
                notify.success('Changed applied');
                state.invalidate().then(function() {
                    $ctrl.state = state.current();
                });
            });
        });
    };

    $ctrl.logout = function() {
        $api.post('/api/logout').then(function() {
            location.href = '/login?logout';
        }, function() {
            location.href = '/login?logout';
        });
    };

    $ctrl.items = [
        {
            link: '/hosts/',
            title: 'Hosts',
            icon: 'fas fa-desktop'
        },
        {
            link: '/groups/',
            title: 'Groups',
            icon: 'fas fa-layer-group'
        },
        {
            link: '/scripts/',
            title: 'Scripts',
            icon: 'fas fa-scroll'
        },
        {
            link: '/options/',
            title: 'Options',
            icon: 'fas fa-cog'
        }
    ];

    function doNavigate(href) {
        if ($location.path() === href) {
            $route.reload();
        } else {
            $location.url(href);
        }
    }

    $ctrl.navigate = function(href) {
        if (document.documentElement.clientWidth > 990) {
            doNavigate(href);
            return;
        }

        $ctrl.isLoading = true;
        $timeout(function() {
            $ctrl.isLoading = false;
            doNavigate(href);
        }, 300);
    };
});
