angular.module('otto').component('trueFalse', {
    bindings: {
        value: '<',
        trueText: '<',
        falseText: '<',
    },
    controller: function($scope, $element, $attrs) {
        var that = this;
        this.$onInit = () => {
            var cls;
            var text;

            if (that.value === true || that.value === 'true') {
                cls = 'badge-success';
                text = that.trueText || 'Enabled';
            } else {
                cls = 'badge-danger';
                text = that.falseText || 'Disabled';
            }

            $element.html('<span class="badge badge-pill ' + cls + '">' + text + '</span>');
        };
    }
});
