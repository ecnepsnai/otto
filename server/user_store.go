package server

import (
	"github.com/ecnepsnai/ds"
	"github.com/ecnepsnai/limits"
	"github.com/ecnepsnai/security"
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
	Username string
	Email    string
	Password string
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

	hashedPassword, err := security.HashPassword([]byte(params.Password))
	if err != nil {
		log.Error("Error hasing user password: %s", err.Error())
		return nil, ErrorFrom(err)
	}

	user := User{
		Username:     params.Username,
		Email:        params.Email,
		Enabled:      true,
		PasswordHash: *hashedPassword,
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
	Email    string
	Enabled  bool
	Password string
}

func (s *userStoreObject) EditUser(user *User, params editUserParameters) (*User, *Error) {
	if existingUser := s.UserWithEmail(params.Email); existingUser != nil && existingUser.Username != user.Username {
		log.Warn("User with email '%s' already exists", params.Email)
		return nil, ErrorUser("User with email '%s' already exists", params.Email)
	}

	user.Email = params.Email
	user.Enabled = params.Enabled
	if params.Password != "" {
		hashedPassword, err := security.HashPassword([]byte(params.Password))
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

func (s *userStoreObject) DeleteUser(user *User) *Error {
	if err := s.Table.Delete(*user); err != nil {
		log.Error("Error deleting user '%s': %s", user.Email, err.Error())
		return ErrorFrom(err)
	}

	UpdateUserCache()

	log.Warn("User deleted: username='%s' email='%s'", user.Username, user.Email)
	return nil
}