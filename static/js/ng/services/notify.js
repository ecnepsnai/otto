angular.module('otto').factory('notify', function($window, $q, $compile, $rootScope, rand) {
    var visibleNotifications = [];
    var notify = (options) => {
        return $q((resolve) => {
            options.id = 'notify-' + rand.ID();
            var cls = options.class || '';

            var notifyHtml = '<div class="notify alert-dismissible fade show alert alert-' + cls + '" id="' + options.id + '">';
            if (options.title) {
                notifyHtml += '<h5 class="alert-heading"><i class="' + options.icon + '"></i><span>' + options.title + '</span></h5>' + options.body;
            } else {
                notifyHtml += '<i class="' + options.icon + '"></i><span>' + options.body + '</span>';
            }
            notifyHtml +='<button type="button" class="close" data-dismiss="alert" aria-label="Close"><span aria-hidden="true">&times;</span></button></div>';
            $('body').append(notifyHtml);

            var $notify = $('#' + options.id);
            var idx = visibleNotifications.push($notify) - 1;
            calculateHeightOffset();
            $notify.alert();
            setTimeout(() => {
                $notify.alert('close');
            }, 2000);
            $notify.on('closed.bs.alert', () => {
                visibleNotifications.splice(visibleNotifications.indexOf($notify), 1);
                calculateHeightOffset();
                resolve();
            });
        });
    };
    var calculateHeightOffset = () => {
        if (visibleNotifications.length > 1) {
            var height = 16;
            for (var i = 0; i < visibleNotifications.length; i++) {
                visibleNotifications[i][0].style.top = height + 'px';
                height += visibleNotifications[i][0].clientHeight + 10;
            }
        }
    };

    window.addEventListener('message', (event) => {
        if (event.data.indexOf('notify:') === 0) {
            var components = event.data.split(':');
            var level = components[1];
            var title = components[2];
            components.splice(0, 3);
            var body = components.join(':');

            if (level === 'error') {
                notify({
                    class: 'danger',
                    title: title,
                    body: body,
                    icon: 'fas fa-exclamation-circle'
                });
            }
        }
    }, false);

    return {
        show: notify,
        success: (body, title) => {
            return notify({
                class: 'success',
                title: title,
                body: body,
                icon: 'fas fa-check-circle'
            });
        },
        info: (body, title) => {
            return notify({
                class: 'primary',
                title: title,
                body: body,
                icon: 'fas fa-info-circle'
            });
        },
        error: (body, title) => {
            return notify({
                class: 'danger',
                title: title,
                body: body,
                icon: 'fas fa-exclamation-circle'
            });
        },
        warning: (body, title) => {
            return notify({
                class: 'warning',
                title: title,
                body: body,
                icon: 'fas fa-times-circle'
            });
        },
        common: {
            saved: () => {
                return notify({
                    class: 'success',
                    body: 'Changes Applied',
                    icon: 'fas fa-check-circle'
                });
            }
        }
    };
});