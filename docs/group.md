# Groups

Otto organizes scripts and hosts into Groups. Groups allow you to categorize servers as you see fit, for example you may wish to have a group for CentOS server and another for Ubuntu servers.

Scripts are associated with groups. When you run a script, it must be associated with a group. This is required so that the correct set of environment variables are used.

There must always be at least one group on an Otto server. On first launch, a default "Otto Clients" group is created. This can be renamed, or deleted as long as another group is added.

Groups can not be deleted if they have any hosts belonging to them, if they are used in a schedule, or if they are used in the host registration configuration.