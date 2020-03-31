angular.module('otto').component('createNew', {
    bindings: {
        link: '<'
    },
    controller: function($scope, $element, $attrs) {
        var that = this;
        this.$onInit = function () {
            var out = '<div class="mb-2">';
            out += '<a href="' + that.link + '" class="btn btn-sm btn-outline-primary"><i class="fas fa-plus"></i> Create New</a>';
            out += '</div>';
            $element.html(out);
        };
    }
});