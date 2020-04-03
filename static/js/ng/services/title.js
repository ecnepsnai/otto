angular.module('otto').factory('title', () => {
    return {
        set: function(val) {
            document.title = val + ' - Otto';
        }
    };
});