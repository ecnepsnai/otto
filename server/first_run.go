package server

var defaultUser = newUserParameters{
	Username: "admin",
	Email:    "admin@localhost",
	Password: "admin",
}

var defaultGroup = newGroupParameters{
	Name: "Otto Clients",
}

func atLeastOneUser() bool {
	users, err := UserStore.AllUsers()
	if err != nil {
		panic(err)
	}
	return len(users) > 0
}

func atLeastOneGroup() bool {
	groups, err := GroupStore.AllGroups()
	if err != nil {
		panic(err)
	}
	return len(groups) > 0
}

func checkFirstRun() {
	if !atLeastOneUser() {
		log.Warn("Creating default user")
		user, err := UserStore.NewUser(defaultUser)
		if err != nil {
			log.Fatal("Unable to make default user: %s", err.Message)
		}
		EventStore.UserAdded(user, "system")
	}

	if !atLeastOneGroup() {
		log.Warn("Creating default group")
		group, err := GroupStore.NewGroup(defaultGroup)
		if err != nil {
			log.Fatal("Unable to make default group: %s", err.Message)
		}
		EventStore.GroupAdded(group, "system")
	}
}
