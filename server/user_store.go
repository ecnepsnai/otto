package server

import (
	"github.com/ecnepsnai/ds"
	"github.com/ecnepsnai/limits"
	"github.com/ecnepsnai/secutil"
)

func (s *userStoreObject) UserWithUsername(username string) *User {
	user, present := GetUserCache()[username]
	if !present {
		return nil
	}
	return &user
}

func (s *userStoreObject) UserWithEmail(email string) *User {
	objects, err := s.Table.GetUnique("Email", email)
	if err != nil {
		log.Error("Error getting user: email='%s' error='%s'", email, err.Error())
		return nil
	}
	if objects == nil {
		return nil
	}
	user, k := objects.(User)
	if !k {
		log.Fatal("Error getting user: email='%s' error='%s'", email, "invalid type")
	}

	return &user
}

func (s *userStoreObject) AllUsers() []User {
	objects, err := s.Table.GetAll(&ds.GetOptions{Sorted: true, Ascending: true})
	if err != nil {
		log.Error("Error listing all users: error='%s'", err.Error())
	}
	if objects == nil || len(objects) == 0 {
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
		PasswordHash:       *hashedPassword,
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

	UpdateUserCache()

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

		user.PasswordHash = *hashedPassword
	}

	if err := s.Table.Update(*user); err != nil {
		log.Error("Error updating user '%s': %s", params.Email, err.Error())
		return nil, ErrorFrom(err)
	}

	UpdateUserCache()

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

	user.PasswordHash = *passwordHash
	user.MustChangePassword = false

	if err := s.Table.Update(*user); err != nil {
		log.Error("Error updating user '%s': %s", username, err.Error())
		return nil, ErrorFrom(err)
	}

	UpdateUserCache()
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
	user.APIKey = *hashedKey
	err = UserStore.Table.Update(user)
	if err != nil {
		return nil, ErrorFrom(err)
	}
	UpdateUserCache()
	log.Info("API key reset: username='%s'", username)
	return &apiKey, nil
}

func (s *userStoreObject) DeleteUser(user *User) *Error {
	if err := s.Table.Delete(*user); err != nil {
		log.Error("Error deleting user '%s': %s", user.Email, err.Error())
		return ErrorFrom(err)
	}

	UpdateUserCache()

	log.Warn("User deleted: username='%s' email='%s'", user.Username, user.Email)
	return nil
}
