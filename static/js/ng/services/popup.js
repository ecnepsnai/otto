angular.module('otto').factory('popup', function($window, $q, $compile, $rootScope, rand) {
    var newPopup = function(options) {
        return $q(function(resolve) {
            options.id = 'popup-' + rand.ID();
            var cls = options.class || '';

            var popupHTML = '<div class="modal fade" id="' + options.id + '" tabindex="-1" role="dialog"><div class="modal-dialog + ' + cls + '" role="document"><div class="modal-content">';
            popupHTML += options.template;
            popupHTML += '</div></div></div>';


            var modalData;
            $scope = $rootScope.$new(true);
            $scope.popupOptions = options;
            $scope.popupData = options.data;
            $scope.popupResolve = function(data) {
                modalData = data;
            };

            $('body').append($compile(popupHTML)($scope));

            var $modal = $('#' + options.id);
            $scope.popupElement = $modal;
            $modal.modal();
            $modal.on('hidden.bs.modal', function() {
                resolve(modalData);
            });
        });
    };
    return {
        /**
         * Show a popup
         * @param {object} options The options for the popup
         * Options:
         * - class: optional class for the popup dialog
         * - template: angularJS template
         * - data: data to pass to the popup scope
         */
        new: newPopup,
        /**
         * Show an alert popup
         * @param {string} title The title of the alert
         * @param {string} body The body of the alert
         */
        alert: function(title, body) {
            return newPopup({
                template: '<alert-popup></alert-popup>',
                data: {
                    title: title,
                    body: body
                }
            });
        },
        /**
         * Show a confirmation popup
         * @param {string} title The title of the popup
         * @param {string} body The body of the popup
         * @param {[]string} choices Buttons. Maximum 2. First is considered no and second is yes.
         */
        confirm: function(title, body, choices) {
            return newPopup({
                template: '<confirm-popup></confirm-popup>',
                data: {
                    title: title,
                    body: body,
                    choices: choices,
                }
            });
        }
    };
});