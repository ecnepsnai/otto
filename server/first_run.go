package server

import "os"

func checkFirstRun() {
	users, err := UserStore.AllUsers()
	if err != nil {
		panic(err)
	}
	if len(users) > 0 {
		return
	}

	_, err = UserStore.NewUser(newUserParameters{
		Username: "admin",
		Email:    "admin@local",
		Password: "admin",
	})
	if err != nil {
		log.Error("Unable to make default user: %s", err.Message)
		os.Exit(1)
	}
}
