angular.module('otto').controller('scriptExecute', function($scope, $api, $q, $script, $group, $route, $location, title) {
    var $ctrl = this;
    $ctrl.stage = 'input';
    $ctrl.hostIDMap = {};
    $ctrl.title = 'Execute Script';
    title.set($ctrl.title);
    $ctrl.selectedHosts = {};
    $ctrl.selectedGroups = {};

    $ctrl.loadData = () => {
        $ctrl.loading = true;
        var scriptID = $route.current.params.id;
        return $q((resolve, reject) => {
            $q.all({
                groups: $script.getGroups(scriptID),
                hosts: $script.getHosts(scriptID),
                script: $script.get(scriptID)
            }, reject).then((results) => {
                var promises = {};
                results.groups.forEach((group) => {
                    promises[group.ID] = $group.getHosts(group.ID);
                });
                results.hosts.forEach((host) => {
                    $ctrl.hostIDMap[host.HostID] = host.HostName;
                });
                $q.all(promises).then((groupHosts) => {
                    $ctrl.loading = false;
                    resolve({
                        groups: results.groups,
                        hosts: results.hosts,
                        script: results.script,
                        groupHosts: groupHosts,
                    });
                }, reject);
            });
        });
    };

    $scope.$watch('$ctrl.allGroups', (val) => {
        if (val === undefined) {
            return;
        }

        $ctrl.groups.forEach((group) => {
            $ctrl.selectedGroups[group.ID] = val;
        });
    });
    $scope.$watch('$ctrl.allHosts', (val) => {
        if (val === undefined) {
            return;
        }

        $ctrl.hosts.forEach((host) => {
            $ctrl.selectedHosts[host.HostID] = val;
        });
    });
    $scope.$watch('$ctrl.selectedGroups', () => {
        var allEnabled = true;
        if (Object.keys($ctrl.selectedGroups).length === 0) {
            return;
        }
        Object.keys($ctrl.selectedGroups).forEach((groupID) => {
            if (!$ctrl.selectedGroups[groupID]) {
                allEnabled = false;
            }
        });
        $ctrl.allGroups = allEnabled;
    }, true);
    $scope.$watch('$ctrl.selectedHosts', () => {
        var allEnabled = true;
        if (Object.keys($ctrl.selectedHosts).length === 0) {
            return;
        }
        Object.keys($ctrl.selectedHosts).forEach((hostID) => {
            if (!$ctrl.selectedHosts[hostID]) {
                allEnabled = false;
            }
        });
        $ctrl.allHosts = allEnabled;
    }, true);

    $ctrl.loadData().then((results) => {
        $ctrl.hosts = results.hosts;
        $ctrl.groups = results.groups;
        $ctrl.script = results.script;
        $ctrl.title = 'Execute Script: ' + $ctrl.script.Name;
        title.set($ctrl.title);
        $ctrl.loaded = true;

        var host = $location.search().host;
        var group = $location.search().group;
        if (group) {
            $ctrl.selectedGroups[group] = true;
            $ctrl.execute(true);
        } else if (host) {
            $ctrl.selectedHosts[host] = true;
            $ctrl.execute(true);
        }
    });

    $ctrl.executeHosts = () => {
        var hosts = {};

        var selectedHosts = $ctrl.hosts.filter((h) => {
            return $ctrl.selectedHosts[h.HostID];
        });

        var selectedGroups = $ctrl.groups.filter((g) => {
            return $ctrl.selectedGroups[g.ID];
        });

        selectedHosts.forEach((h) => {
            hosts[h.HostID] = true;
        });
        selectedGroups.forEach((g) => {
            $ctrl.hosts.forEach((h) => {
                if (h.GroupID === g.ID) {
                    hosts[h.HostID] = true;
                }
            });
        });

        return Object.keys(hosts);
    };

    $ctrl.execute = (valid) => {
        $ctrl.executeHosts();

        if (!valid) {
            return;
        }

        var selectedHosts = $ctrl.executeHosts();
        if (selectedHosts.length === 0) {
            return;
        }

        var executions = [];
        selectedHosts.forEach((host) => {
            executions.push({
                Action: 'run_script',
                HostID: host,
                ScriptID: $ctrl.script.ID,
            });
        });

        $ctrl.executed = 0;
        $ctrl.total = executions.length;
        $ctrl.stage = 'executing';
        $ctrl.startExecution(executions);
    };

    $ctrl.startExecution = (executions) => {
        $ctrl.executePercent = 0;
        $ctrl.searchProgressStyle = {width: $ctrl.executePercent + '%'};

        $ctrl.results = [];

        function updateExecuteProgress() {
            $ctrl.executed++;
            $ctrl.executePercent = ($ctrl.executed / $ctrl.total) * 100;
            $ctrl.searchProgressStyle = {width: $ctrl.executePercent + '%'};
        }

        var promises = [];
        executions.forEach((execution) => {
            promises.push($q((resolve) => {
                $api.put('/api/request', execution).then((response) => {
                    var result = response.data.data;
                    result.Host = $ctrl.hostIDMap[execution.HostID];

                    var keys = Object.keys(result.Environment).sort();
                    var environmentListSorted = [];
                    keys.forEach((key) => {
                        environmentListSorted.push({
                            Key: key,
                            Value: result.Environment[key],
                        });
                    });
                    result.Environment = environmentListSorted;

                    updateExecuteProgress();
                    resolve(result);
                }, (response) => {
                    var result = {error: response.data.error};
                    result.Host = $ctrl.hostIDMap[execution.HostID];
                    updateExecuteProgress();
                    resolve(result);
                });
            }));
        });

        $q.all(promises).then((results) => {
            $ctrl.results = results;
            $ctrl.stage = 'results';
        });
    };

    $ctrl.borderClass = (result) => {
        if (result.error || !result.Result.Success) {
            return {
                'border-danger': true,
            };
        }
        return {
            'border-danger': false,
        };
    };

    $ctrl.headerClass = (result) => {
        if (result.error || !result.Result.Success) {
            return {
                'text-white': true,
                'bg-danger': true,
            };
        }
        return {
            'text-white': false,
            'bg-danger': false,
        };
    };
});
