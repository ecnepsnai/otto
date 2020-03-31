angular.module('otto').component('pageWrapper', {
    bindings: {
        pageTitle: '@'
    },
    controller: function($scope, $element, $attrs) {
        var that = this;
        this.$onInit = function () {
            var title = '<div class="page-title"><div class="container">';
            title += '<p>' + that.pageTitle + '</p>';
            title += '</div></div>';
            $element.before(title);
            $element.wrap('<div class="container"></div>');
        };
    }
});