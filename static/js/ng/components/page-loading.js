angular.module('otto').component('pageLoading', {
    bindings: {},
    controller: function($scope, $element, $attrs) {
        var that = this;
        this.$onInit = function () {
            $element.append('<i class="fas fa-spinner fa-pulse"></i><strong class="a-little-to-the-left">Loading...</strong>');
        };
    }
});