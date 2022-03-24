package server

import (
	"github.com/ecnepsnai/ds"
	"github.com/ecnepsnai/limits"
	"github.com/ecnepsnai/secutil"
)

const systemUsername = "system"

func (s *userStoreObject) UserWithUsername(username string) *User {
	user, present := UserCache.ByUsername(username)
	if !present {
		return nil
	}
	return &user
}

func (s *userStoreObject) UserWithEmail(email string) *User {
	user, present := UserCache.ByEmail(email)
	if !present {
		return nil
	}
	return &user
}

func (s *userStoreObject) AllUsers() []User {
	objects, err := s.Table.GetAll(&ds.GetOptions{Sorted: true, Ascending: true})
	if err != nil {
		log.Error("Error listing all users: error='%s'", err.Error())
	}
	if len(objects) == 0 {
		return []User{}
	}

	users := make([]User, len(objects))
	for i, obj := range objects {
		user, k := obj.(User)
		if !k {
			log.Fatal("Error listing all users: error='%s'", "invalid type")
		}
		users[i] = user
	}

	return users
}

type newUserParameters struct {
	Username           string `max:"32" min:"1"`
	Email              string `max:"128" min:"1"`
	Password           string
	MustChangePassword bool
}

func (s *userStoreObject) NewUser(params newUserParameters) (*User, *Error) {
	if s.UserWithUsername(params.Username) != nil {
		log.Warn("User with username '%s' already exists", params.Username)
		return nil, ErrorUser("User with username '%s' already exists", params.Username)
	}
	if s.UserWithEmail(params.Email) != nil {
		log.Warn("User with email '%s' already exists", params.Email)
		return nil, ErrorUser("User with email '%s' already exists", params.Email)
	}

	hashedPassword, err := secutil.HashPassword([]byte(params.Password))
	if err != nil {
		log.Error("Error hasing user password: %s", err.Error())
		return nil, ErrorFrom(err)
	}

	user := User{
		Username:           params.Username,
		Email:              params.Email,
		CanLogIn:           true,
		MustChangePassword: params.MustChangePassword,
	}

	if err := limits.Check(user); err != nil {
		return nil, ErrorUser(err.Error())
	}

	if err := s.Table.Add(user); err != nil {
		log.Error("Error adding new user '%s': %s", params.Email, err.Error())
		return nil, ErrorFrom(err)
	}
	ShadowStore.Set(user.Username, *hashedPassword)

	UserCache.Update()
	log.Info("New user added: username='%s' email='%s'", params.Username, params.Email)
	return &user, nil
}

type editUserParameters struct {
	Email              string `max:"128" min:"1"`
	Password           string
	CanLogIn           bool
	MustChangePassword bool
}

func (s *userStoreObject) EditUser(user *User, params editUserParameters) (*User, *Error) {
	if existingUser := s.UserWithEmail(params.Email); existingUser != nil && existingUser.Username != user.Username {
		log.Warn("User with email '%s' already exists", params.Email)
		return nil, ErrorUser("User with email '%s' already exists", params.Email)
	}

	user.Email = params.Email
	user.CanLogIn = params.CanLogIn
	user.MustChangePassword = params.MustChangePassword
	if params.Password != "" {
		hashedPassword, err := secutil.HashPassword([]byte(params.Password))
		if err != nil {
			log.Error("Error hasing user password: %s", err.Error())
			return nil, ErrorFrom(err)
		}

		ShadowStore.Set(user.Username, *hashedPassword)
	}

	if err := s.Table.Update(*user); err != nil {
		log.Error("Error updating user '%s': %s", params.Email, err.Error())
		return nil, ErrorFrom(err)
	}

	UserCache.Update()
	return user, nil
}

func (s *userStoreObject) ResetPassword(username string, newPassword []byte) (*User, *Error) {
	user := s.UserWithUsername(username)
	if user == nil {
		return nil, ErrorUser("no user with username %s", username)
	}

	passwordHash, err := secutil.HashPassword(newPassword)
	if err != nil {
		log.Error("Error hasing password for user: username='%s' error='%s'", username, err.Error())
		return nil, ErrorFrom(err)
	}

	user.MustChangePassword = false

	if err := s.Table.Update(*user); err != nil {
		log.Error("Error updating user '%s': %s", username, err.Error())
		return nil, ErrorFrom(err)
	}

	ShadowStore.Set(user.Username, *passwordHash)

	UserCache.Update()
	return user, nil
}

func (s *userStoreObject) ResetAPIKey(username string) (*string, *Error) {
	u := UserStore.UserWithUsername(username)
	if u == nil {
		return nil, ErrorUser("no user with username %s", username)
	}
	user := *u

	apiKey := newAPIKey()
	hashedKey, err := secutil.HashPassword([]byte(apiKey))
	if err != nil {
		log.Error("Error hashing API key: %s", err.Error())
		return nil, ErrorFrom(err)
	}

	ShadowStore.Set("api_"+user.Username, *hashedKey)
	log.Info("API key reset: username='%s'", username)
	return &apiKey, nil
}

func (s *userStoreObject) DeleteUser(user *User) *Error {
	if err := s.Table.Delete(*user); err != nil {
		log.Error("Error deleting user '%s': %s", user.Email, err.Error())
		return ErrorFrom(err)
	}
	ShadowStore.Delete(user.Username)

	UserCache.Update()
	log.Warn("User deleted: username='%s' email='%s'", user.Username, user.Email)
	return nil
}

func (s *userStoreObject) DisableUser(username string) *Error {
	u := UserStore.UserWithUsername(username)
	if u == nil {
		return ErrorUser("no user with username %s", username)
	}
	user := *u

	user.CanLogIn = false

	if err := s.Table.Update(user); err != nil {
		log.Error("Error updating user '%s': %s", username, err.Error())
		return ErrorFrom(err)
	}

	UserCache.Update()

	log.PWarn("Disabled user", map[string]interface{}{
		"user": username,
	})
	return nil
}
