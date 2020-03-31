var otto = angular.module('otto', ['ngRoute', 'angularMoment', 'ngSanitize']);

otto.controller('otto', OttoController);

function OttoController($scope, $location, state) {
    state.start().then(function(state) {
        $scope.ready = true;
        $scope.state = state;
    });

    $scope.logout = function() {
        $api.post('/api/logout').then(function() {
            location.href = '/login?logout';
        }, function() {
            location.href = '/login?logout';
        });
    };

    $scope.navClass = function(tab) {
        var matches = $location.path().startsWith(tab);
        return { active: matches };
    };

    function doNavigate(href) {
        if ($location.path() === href) {
            $route.reload();
        } else {
            $location.url(href);
        }
    }

    $scope.navigate = function(href) {
        if (document.documentElement.clientWidth > 990) {
            doNavigate(href);
            return;
        }

        $scope.isLoading = true;
        $timeout(function() {
            $scope.isLoading = false;
            doNavigate(href);
        }, 300);
    };
}

otto.config(function($routeProvider, $locationProvider) {
    // Hosts
    $routeProvider.when('/hosts/', {
        template: '<host-list></host-list>'
    });
    $routeProvider.when('/hosts/host/', {
        template: '<host-edit></host-edit>'
    });
    $routeProvider.when('/hosts/host/:id/', {
        template: '<host-view></host-view>'
    });
    $routeProvider.when('/hosts/host/:id/edit/', {
        template: '<host-edit></host-edit>'
    });

    // Groups
    $routeProvider.when('/groups/', {
        template: '<group-list></group-list>'
    });
    $routeProvider.when('/groups/group/', {
        template: '<group-edit></group-edit>'
    });
    $routeProvider.when('/groups/group/:id/', {
        template: '<group-view></group-view>'
    });
    $routeProvider.when('/groups/group/:id/edit/', {
        template: '<group-edit></group-edit>'
    });

    // Scripts
    $routeProvider.when('/scripts/', {
        template: '<script-list></script-list>'
    });
    $routeProvider.when('/scripts/script/', {
        template: '<script-edit></script-edit>'
    });
    $routeProvider.when('/scripts/script/:id/', {
        template: '<script-view></script-view>'
    });
    $routeProvider.when('/scripts/script/:id/edit/', {
        template: '<script-edit></script-edit>'
    });
    $routeProvider.when('/scripts/script/:id/execute/', {
        template: '<script-execute></script-execute>'
    });

    // Options
    $routeProvider.when('/options/', {
        template: '<options-edit></options-edit>'
    });
    $routeProvider.when('/options/users/user/', {
        template: '<user-edit></user-edit>'
    });
    $routeProvider.when('/options/users/user/:username/', {
        template: '<user-edit></user-edit>'
    });

    $locationProvider.html5Mode(true);
});