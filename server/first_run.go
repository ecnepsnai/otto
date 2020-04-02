package server

import "os"

var defaultUser = newUserParameters{
	Username: "admin",
	Email:    "admin@localhost",
	Password: "admin",
}

func checkFirstRun() {
	users, err := UserStore.AllUsers()
	if err != nil {
		panic(err)
	}
	if len(users) > 0 {
		return
	}

	if _, err := UserStore.NewUser(defaultUser); err != nil {
		log.Error("Unable to make default user: %s", err.Message)
		os.Exit(1)
	}
}
