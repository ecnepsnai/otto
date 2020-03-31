angular.module('otto').factory('$api', function($http, notify, $q) {
    return {
        get: function(url) {
            return $q(function(resolve, reject) {
                $http.get(url).then(results => {
                    resolve(results);
                }, function(error) {
                    var message = 'Internal Server Error';
                    if (error && error.data && error.data.message) {
                        message = error.data.message;
                    }
                    console.error(error);
                    notify.error(message);
                    reject(error);
                }).catch(function(exception) {
                    console.error(exception);
                    notify.error('Internal Server Error');
                    reject('Internal Server Error');
                });
            });
        },
        post: function(url, body) {
            return $q(function(resolve, reject) {
                $http.post(url, body).then(results => {
                    resolve(results);
                }, function(error) {
                    var message = 'Internal Server Error';
                    if (error && error.data && error.data.message) {
                        message = error.data.message;
                    }
                    console.error(error);
                    notify.error(message);
                    reject(error);
                }).catch(function(exception) {
                    console.error(exception);
                    notify.error('Internal Server Error');
                    reject('Internal Server Error');
                });
            });
        },
        put: function(url, body) {
            return $q(function(resolve, reject) {
                $http.put(url, body).then(results => {
                    resolve(results);
                }, function(error) {
                    var message = 'Internal Server Error';
                    if (error && error.data && error.data.message) {
                        message = error.data.message;
                    }
                    console.error(error);
                    notify.error(message);
                    reject(error);
                }).catch(function(exception) {
                    console.error(exception);
                    notify.error('Internal Server Error');
                    reject('Internal Server Error');
                });
            });
        },
        delete: function(url) {
            return $q(function(resolve, reject) {
                $http.delete(url).then(results => {
                    resolve(results);
                }, function(error) {
                    var message = 'Internal Server Error';
                    if (error && error.data && error.data.message) {
                        message = error.data.message;
                    }
                    console.error(error);
                    notify.error(message);
                    reject(error);
                }).catch(function(exception) {
                    console.error(exception);
                    notify.error('Internal Server Error');
                    reject('Internal Server Error');
                });
            });
        }
    };
});
