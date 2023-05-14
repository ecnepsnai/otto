package server

import "testing"

func TestAddGetUser(t *testing.T) {
	username := randomString(6)

	user, err := UserStore.NewUser(newUserParameters{
		Username: username,
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
}

func TestEditUser(t *testing.T) {
	username := randomString(6)

	user, err := UserStore.NewUser(newUserParameters{
		Username: username,
		Password: randomString(6),
	})
	if err != nil {
		t.Fatalf("Error making new user: %s", err.Message)
	}

	_, err = UserStore.EditUser(user, editUserParameters{
		CanLogIn: true,
	})
	if err != nil {
		t.Fatalf("Error modifying user: %s", err.Message)
	}

	user = UserStore.UserWithUsername(username)
	if user == nil {
		t.Fatalf("Should return user")
	}
}

func TestDeleteUser(t *testing.T) {
	username := randomString(6)

	user, err := UserStore.NewUser(newUserParameters{
		Username: username,
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

	_, err := UserStore.NewUser(newUserParameters{
		Username: username,
		Password: randomString(6),
	})
	if err != nil {
		t.Fatalf("Error making new user: %s", err.Message)
	}

	// Duplicate username
	_, err = UserStore.NewUser(newUserParameters{
		Username: username,
		Password: randomString(6),
	})
	if err == nil {
		t.Fatalf("Should return error on duplicate username")
	}
}

func TestResetUserPassword(t *testing.T) {
	username := randomString(6)

	user, err := UserStore.NewUser(newUserParameters{
		Username:           username,
		Password:           randomString(6),
		MustChangePassword: true,
	})
	if err != nil {
		t.Fatalf("Error making new user: %s", err.Message)
	}
	if user == nil {
		t.Fatalf("No user returned")
	}

	user, err = UserStore.ResetPassword(username, []byte(randomString(6)))
	if err != nil {
		t.Fatalf("Error changing password: %s", err.Message)
	}

	if user.MustChangePassword {
		t.Fatalf("Should not require password change after password change")
	}
}
