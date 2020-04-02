angular.module('otto').factory('title', () => {
    return {
        set: (val) => {
            document.title = val + ' - Otto';
        }
    };
});