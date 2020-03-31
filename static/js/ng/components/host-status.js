angular.module('otto').component('hostStatus', {
    bindings: {
        heartbeat: '<'
    },
    controller: function($element) {
        var that = this;
        this.$onInit = function () {
            var heartbeat = that.heartbeat;
            if (!heartbeat) {
                heartbeat = {
                    Unknown: true,
                };
            }

            var out = '';
            if (heartbeat.Unknown) {
                out = '<span class="badge badge-pill badge-secondary"><i class="fas fa-question-circle"></i> Unknown</span>';
            } else if (heartbeat.IsReachable) {
                out = '<span class="badge badge-pill badge-success"><i class="fas fa-check-circle"></i> Reachable</span>';
            } else {
                out = '<span class="badge badge-pill badge-danger"><i class="fas fa-times-circle"></i> Unreachable</span>';
            }
            $element.html(out);
        };
    }
});