# Schedules

Schedules are a way to run a script on groups or individual hosts on a specified frequency.

# Creating a Schedule

To create a schedule you must first already have a script configured and hosts that belong to groups with that script.

In the web interface navigate to Schedules and click "Create New"

Select the script that you wish to run from the dropdown then select a frequency.

You can use a pre-defined frequency template of:

- Every hour
- Every 4 hours
- Every day at midnight
- Every monday at midnight

Or define your own schedule using [cron patterns](https://pkg.go.dev/github.com/ecnepsnai/cron#pkg-overview).

**All schedules run in the UTC timezone, not the timezone of the server, agent, or your local computer.**

Last, select the individual hosts or groups that you want this script to run.

# Monitoring a Schedule

A history of the runs of the schedule is maintained and you can view the previous runs on the web interface.
