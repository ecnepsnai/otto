package server

var defaultUser = newUserParameters{
	Username: "admin",
	Email:    "admin@localhost",
	Password: "admin",
}

var defaultGroup = newGroupParameters{
	Name: "Otto Clients",
}

func isFirstRun() bool {
	numberOfUsers := 0
	numberOfGroups := 0

	users, err := UserStore.AllUsers()
	if err != nil {
		panic(err)
	}
	numberOfUsers = len(users)

	groups, err := GroupStore.AllGroups()
	if err != nil {
		panic(err)
	}
	numberOfGroups = len(groups)

	return numberOfUsers == 0 && numberOfGroups == 0
}

func checkFirstRun() {
	if !isFirstRun() {
		return
	}

	if _, err := UserStore.NewUser(defaultUser); err != nil {
		log.Fatal("Unable to make default user: %s", err.Message)
	}

	if _, err := GroupStore.NewGroup(defaultGroup); err != nil {
		log.Fatal("Unable to make default group: %s", err.Message)
	}
}
