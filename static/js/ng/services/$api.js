angular.module('otto').factory('$api', function($http, notify, $q) {
    function dealWithError(error) {
        // Session no longer valid
        if (error.status === 403) {
            location.href = '/login?unauthorized';
            return;
        }

        var message = 'Internal Server Error';
        if (error && error.data) {
            if (error.data.message) {
                message = error.data.message;
            } else if (error.data.error && error.data.error.message) {
                message = error.data.error.message;
            }
        }
        console.error(error);
        notify.error(message, 'Error Processing Request');
    }

    return {
        get: function(url) {
            return $q(function(resolve, reject) {
                $http.get(url).then(results => {
                    resolve(results);
                }, function(error) {
                    dealWithError(error);
                    reject(error);
                }).catch(function(exception) {
                    dealWithError(exception);
                    reject('Internal Server Error');
                });
            });
        },
        post: function(url, body) {
            return $q(function(resolve, reject) {
                $http.post(url, body).then(results => {
                    resolve(results);
                }, function(error) {
                    dealWithError(error);
                    reject(error);
                }).catch(function(exception) {
                    dealWithError(exception);
                    reject('Internal Server Error');
                });
            });
        },
        put: function(url, body) {
            return $q(function(resolve, reject) {
                $http.put(url, body).then(results => {
                    resolve(results);
                }, function(error) {
                    dealWithError(error);
                    reject(error);
                }).catch(function(exception) {
                    dealWithError(exception);
                    reject('Internal Server Error');
                });
            });
        },
        delete: function(url) {
            return $q(function(resolve, reject) {
                $http.delete(url).then(results => {
                    resolve(results);
                }, function(error) {
                    dealWithError(error);
                    reject(error);
                }).catch(function(exception) {
                    dealWithError(exception);
                    reject('Internal Server Error');
                });
            });
        }
    };
});
