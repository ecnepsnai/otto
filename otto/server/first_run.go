package server

var DangerousSkipForcePasswordChangeForDefaultUser = ""

var defaultUser = newUserParameters{
	Username:           "admin",
	Password:           "admin",
	MustChangePassword: true,
	Permissions:        UserPermissionsMax(),
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
	RunLevel: ScriptRunLevelReadOnly,
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
		params := defaultUser
		if DangerousSkipForcePasswordChangeForDefaultUser == "true" {
			log.Error("[DANGER] Not forcing a password change for default user. This should only be used in a development environment.")
			params.MustChangePassword = false
		}
		user, err := UserStore.NewUser(params)
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
		if AutoRegisterOptions.DefaultGroupID == "" {
			o := AutoRegisterOptions
			o.DefaultGroupID = group.ID
			o.Save()
		}
	}
}
