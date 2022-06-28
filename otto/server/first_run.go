package server

var defaultUser = newUserParameters{
	Username:           "admin",
	Email:              "admin@localhost",
	Password:           "admin",
	MustChangePassword: true,
}

var defaultGroup = newGroupParameters{
	Name: "Otto Agents",
}

var defaultScript = newScriptParameters{
	Name:       "Hello World",
	Executable: "/bin/sh",
	Script:     "#!/bin/sh\n\necho \"Hello world, my name is $(hostname)\"\n",
	RunAs: RunAs{
		Inherit: true,
	},
}

func atLeastOneUser() bool {
	users := UserStore.AllUsers()
	return len(users) > 0
}

func atLeastOneGroup() bool {
	groups := GroupStore.AllGroups()
	return len(groups) > 0
}

func atLeastOneScript() bool {
	scripts := ScriptStore.AllScripts()
	return len(scripts) > 0
}

func checkFirstRun() {
	if !atLeastOneUser() {
		log.Warn("Creating default user")
		user, err := UserStore.NewUser(defaultUser)
		if err != nil {
			log.Fatal("Unable to make default user: %s", err.Message)
		}
		EventStore.UserAdded(user, systemUsername)
	}

	defaultScriptId := ""
	if !atLeastOneScript() {
		log.Warn("Creating default script")
		script, err := ScriptStore.NewScript(defaultScript)
		if err != nil {
			log.Fatal("Unable to make default script: %s", err.Message)
		}
		EventStore.ScriptAdded(script, systemUsername)
		defaultScriptId = script.ID
	}

	if !atLeastOneGroup() {
		log.Warn("Creating default group")
		if defaultScriptId != "" {
			defaultGroup.ScriptIDs = []string{defaultScriptId}
		}

		group, err := GroupStore.NewGroup(defaultGroup)
		if err != nil {
			log.Fatal("Unable to make default group: %s", err.Message)
		}
		EventStore.GroupAdded(group, systemUsername)
		if Options.Register.DefaultGroupID == "" {
			o := Options
			o.Register.DefaultGroupID = group.ID
			o.Save()
		}
	}
}
