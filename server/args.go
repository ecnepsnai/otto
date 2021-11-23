package server

import (
	"fmt"
	"os"
	"time"

	"github.com/ecnepsnai/otto/server/environ"
	"github.com/ecnepsnai/secutil"
)

func preBootstrapArgs() {
	args := os.Args[1:]
	i := 0
	count := len(args)
	for i < count {
		arg := args[i]

		if arg == "-d" || arg == "--data-dir" {
			if i == count-1 {
				fmt.Fprintf(os.Stderr, "%s requires exactly 1 parameter\n", arg)
				printHelpAndExit()
			}

			value := args[i+1]
			dataDirectory = value
			i++
		} else if arg == "-b" || arg == "--bind-addr" {
			if i == count-1 {
				fmt.Fprintf(os.Stderr, "%s requires exactly 1 parameter\n", arg)
				printHelpAndExit()
			}

			value := args[i+1]
			bindAddress = value
			i++
		} else if arg == "--no-scheduler" {
			cronDisabled = true
		} else if arg == "-h" || arg == "--help" {
			printHelpAndExit()
		}

		i++
	}
}

func isVerbose() bool {
	args := os.Args[1:]
	for _, arg := range args {
		if arg == "-v" || arg == "--verbose" {
			return true
		}
	}
	return false
}

func postBootstrapArgs() {
	args := os.Args[1:]
	i := 0
	count := len(args)
	for i < count {
		arg := args[i]

		if arg == "--demo" {
			script, err := ScriptStore.NewScript(newScriptParameters{
				Name:        "Update Software",
				Executable:  "/bin/sh",
				Script:      "#!/bin/sh\n${YUM_CMD} -y update\n",
				Environment: []environ.Variable{},
			})
			if err != nil {
				panic(err.Message)
			}
			el7, err := GroupStore.NewGroup(newGroupParameters{
				Name:      "CentOS 7 Servers",
				ScriptIDs: []string{script.ID},
				Environment: []environ.Variable{
					environ.New("YUM_CMD", "yum"),
				},
			})
			if err != nil {
				panic(err.Message)
			}
			el8, err := GroupStore.NewGroup(newGroupParameters{
				Name:      "CentOS 8 Servers",
				ScriptIDs: []string{script.ID},
				Environment: []environ.Variable{
					environ.New("YUM_CMD", "dnf"),
				},
			})
			if err != nil {
				panic(err.Message)
			}
			el7Host, err := HostStore.NewHost(newHostParameters{
				Name:     "el7-host.example.com",
				Address:  "el7-host.example.com",
				Port:     12444,
				GroupIDs: []string{el7.ID},
			})
			if err != nil {
				panic(err.Message)
			}
			el8Host, err := HostStore.NewHost(newHostParameters{
				Name:     "el8-host.example.com",
				Address:  "el8-host.example.com",
				Port:     12444,
				GroupIDs: []string{el8.ID},
			})
			if err != nil {
				panic(err.Message)
			}
			schedule, err := ScheduleStore.NewSchedule(newScheduleParameters{
				Name:     "OS Updates",
				ScriptID: script.ID,
				Scope: ScheduleScope{
					GroupIDs: []string{
						el7.ID,
						el8.ID,
					},
				},
				Pattern: "0 0 * * *",
			})
			if err != nil {
				panic(err.Message)
			}
			x := int64(0)
			for x < 10 {
				report := ScheduleReport{
					ID:         newID(),
					ScheduleID: schedule.ID,
					HostIDs:    []string{el7Host.ID, el8Host.ID},
					Time: ScheduleReportTime{
						Start:          time.Unix(time.Now().Unix()-x*7400, 0),
						Finished:       time.Unix(time.Now().Unix()-x*7350, 0),
						ElapsedSeconds: 50.0,
					},
					Result: 0,
				}
				ScheduleReportStore.Table.Add(report)
				x++
			}
			_, err = RegisterRuleStore.NewRule(newRegisterRuleParams{
				Name: "CentOS 7 Servers",
				Clauses: []RegisterRuleClause{
					{
						Property: RegisterRulePropertyDistributionName,
						Pattern:  "CentOS Linux",
					},
					{
						Property: RegisterRulePropertyDistributionVersion,
						Pattern:  "7",
					},
				},
				GroupID: el7.ID,
			})
			if err != nil {
				panic(err.Message)
			}
			_, err = RegisterRuleStore.NewRule(newRegisterRuleParams{
				Name: "CentOS 8 Servers",
				Clauses: []RegisterRuleClause{
					{
						Property: RegisterRulePropertyDistributionName,
						Pattern:  "CentOS Linux",
					},
					{
						Property: RegisterRulePropertyDistributionVersion,
						Pattern:  "8",
					},
				},
				GroupID: el8.ID,
			})
			if err != nil {
				panic(err.Message)
			}
			o := Options
			o.Register.Enabled = true
			o.Register.Key = secutil.RandomString(6)
			o.Save()
		}

		i++
	}

	EventStore.ServerStarted(args)
}

func printHelpAndExit() {
	fmt.Printf("Usage %s [options]\n", os.Args[0])
	fmt.Printf("Options:\n")
	fmt.Printf("-d --data-dir <path>        Specify the absolute path to the data directory\n")
	fmt.Printf("-b --bind-addr <socket>     Specify the listen address for the web server\n")
	fmt.Printf("-v --verbose                Set the log level to debug\n")
	fmt.Printf("--no-scheduler              Disable all automatic tasks\n")
	os.Exit(1)
}
