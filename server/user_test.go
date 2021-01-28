package server

import "testing"

func TestAddGetUser(t *testing.T) {
	username := randomString(6)
	email := randomString(6)

	user, err := UserStore.NewUser(newUserParameters{
		Username: username,
		Email:    email,
		Password: randomString(6),
	})
	if err != nil {
		t.Fatalf("Error making new user: %s", err.Message)
	}
	if user == nil {
		t.Fatalf("No user returned")
	}

	if UserStore.UserWithUsername(username) == nil {
		t.Fatalf("No user with username")
	}

	if UserStore.UserWithEmail(email) == nil {
		t.Fatalf("No user with email")
	}
}

func TestEditUser(t *testing.T) {
	username := randomString(6)
	email := randomString(6)

	user, err := UserStore.NewUser(newUserParameters{
		Username: username,
		Email:    email,
		Password: randomString(6),
	})
	if err != nil {
		t.Fatalf("Error making new user: %s", err.Message)
	}

	_, err = UserStore.EditUser(user, editUserParameters{
		Email:   randomString(6),
		Enabled: true,
	})
	if err != nil {
		t.Fatalf("Error modifying user: %s", err.Message)
	}

	user = UserStore.UserWithUsername(username)
	if user == nil {
		t.Fatalf("Should return user")
	}

	if user.Email == email {
		t.Fatalf("Should update email")
	}
}

func TestDeleteUser(t *testing.T) {
	username := randomString(6)
	email := randomString(6)

	user, err := UserStore.NewUser(newUserParameters{
		Username: username,
		Email:    email,
		Password: randomString(6),
	})
	if err != nil {
		t.Fatalf("Error making new user: %s", err.Message)
	}

	if err := UserStore.DeleteUser(user); err != nil {
		t.Fatalf("Error deleting user: %s", err.Message)
	}

	if UserStore.UserWithUsername(username) != nil {
		t.Fatalf("Should not return user after deleting")
	}
}

func TestDuplicateUser(t *testing.T) {
	username := randomString(6)
	email := randomString(6)

	_, err := UserStore.NewUser(newUserParameters{
		Username: username,
		Email:    email,
		Password: randomString(6),
	})
	if err != nil {
		t.Fatalf("Error making new user: %s", err.Message)
	}

	// Duplicate username
	_, err = UserStore.NewUser(newUserParameters{
		Username: username,
		Email:    randomString(6),
		Password: randomString(6),
	})
	if err == nil {
		t.Fatalf("Should return error on duplicate username")
	}

	// Duplicate email
	_, err = UserStore.NewUser(newUserParameters{
		Username: randomString(6),
		Email:    email,
		Password: randomString(6),
	})
	if err == nil {
		t.Fatalf("Should return error on duplicate email")
	}
}
