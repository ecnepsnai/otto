angular.module('otto').factory('title', function() {
    return {
        set: function(val) {
            document.title = val + ' - Otto';
        }
    };
});