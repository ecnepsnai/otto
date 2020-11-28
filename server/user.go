package server

import (
	"github.com/ecnepsnai/ds"
	"github.com/ecnepsnai/limits"
	"github.com/ecnepsnai/security"
)

// User describes a user object
type User struct {
	Username     string `ds:"primary" max:"32" min:"1"`
	Email        string `ds:"unique" max:"128" min:"1"`
	Enabled      bool
	PasswordHash security.HashedPassword `json:"-" structs:"-" ts:"-"`
}

func (s *userStoreObject) UserWithUsername(username string) (*User, *Error) {
	obj, err := s.Table.Get(username)
	if err != nil {
		log.Error("Error getting user with Username '%s': %s", username, err.Error())
		return nil, ErrorFrom(err)
	}
	if obj == nil {
		return nil, nil
	}
	user, k := obj.(User)
	if !k {
		log.Error("Object is not of type 'User'")
		return nil, ErrorServer("incorrect type")
	}

	return &user, nil
}

func (s *userStoreObject) UserWithEmail(email string) (*User, *Error) {
	obj, err := s.Table.GetUnique("Email", email)
	if err != nil {
		log.Error("Error getting user with email '%s': %s", email, err.Error())
		return nil, ErrorFrom(err)
	}
	if obj == nil {
		return nil, nil
	}
	user, k := obj.(User)
	if !k {
		log.Error("Object is not of type 'User'")
		return nil, ErrorServer("incorrect type")
	}

	return &user, nil
}

func (s *userStoreObject) AllUsers() ([]User, *Error) {
	objs, err := s.Table.GetAll(&ds.GetOptions{Sorted: true, Ascending: true})
	if err != nil {
		log.Error("Error getting all users: %s", err.Error())
		return nil, ErrorFrom(err)
	}
	if objs == nil || len(objs) == 0 {
		return []User{}, nil
	}

	users := make([]User, len(objs))
	for i, obj := range objs {
		user, k := obj.(User)
		if !k {
			log.Error("Object is not of type 'User'")
			return []User{}, ErrorServer("incorrect type")
		}
		users[i] = user
	}

	return users, nil
}

type newUserParameters struct {
	Username string
	Email    string
	Password string
}

func (s *userStoreObject) NewUser(params newUserParameters) (*User, *Error) {
	existingUser, err := s.UserWithEmail(params.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		log.Warn("User with email '%s' already exists", params.Email)
		return nil, ErrorUser("User with email '%s' already exists", params.Email)
	}

	hashedPassword, erro := security.HashPassword([]byte(params.Password))
	if erro != nil {
		log.Error("Error hasing user password: %s", erro.Error())
		return nil, ErrorFrom(erro)
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

	log.Info("Added new user '%s'", params.Email)
	return &user, nil
}

type editUserParameters struct {
	Email    string
	Enabled  bool
	Password string
}

func (s *userStoreObject) EditUser(user *User, params editUserParameters) (*User, *Error) {
	existingUser, err := s.UserWithEmail(params.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil && existingUser.Username != user.Username {
		log.Warn("User with email '%s' already exists", params.Email)
		return nil, ErrorUser("User with email '%s' already exists", params.Email)
	}

	user.Email = params.Email
	user.Enabled = params.Enabled
	if params.Password != "" {
		hashedPassword, erro := security.HashPassword([]byte(params.Password))
		if erro != nil {
			log.Error("Error hasing user password: %s", erro.Error())
			return nil, ErrorFrom(erro)
		}

		user.PasswordHash = *hashedPassword
	}

	if err := s.Table.Update(*user); err != nil {
		log.Error("Error updating user '%s': %s", params.Email, err.Error())
		return nil, ErrorFrom(err)
	}

	log.Info("Updating user '%s'", params.Email)
	return user, nil
}

func (s *userStoreObject) DeleteUser(user *User) *Error {
	if err := s.Table.Delete(*user); err != nil {
		log.Error("Error deleting user '%s': %s", user.Email, err.Error())
		return ErrorFrom(err)
	}

	log.Info("Deleting user '%s'", user.Email)
	return nil
}
