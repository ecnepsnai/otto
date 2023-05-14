package server

import (
	"github.com/ecnepsnai/ds"
	"github.com/ecnepsnai/limits"
	"github.com/ecnepsnai/secutil"
)

const systemUsername = "system"

func (s *userStoreObject) UserWithUsername(username string) (user *User) {
	s.Table.StartRead(func(tx ds.IReadTransaction) error {
		user = s.userWithUsername(tx, username)
		return nil
	})
	return
}

func (s *userStoreObject) userWithUsername(tx ds.IReadTransaction, username string) *User {
	object, err := tx.Get(username)
	if err != nil {
		log.PError("Error getting user by username", map[string]interface{}{
			"username": username,
			"error":    err.Error(),
		})
		return nil
	}
	if object == nil {
		return nil
	}
	user, ok := object.(User)
	if !ok {
		log.Panic("Invalid object type in user store")
	}
	return &user
}

func (s *userStoreObject) AllUsers() (users []User) {
	s.Table.StartRead(func(tx ds.IReadTransaction) error {
		users = s.allUsers(tx)
		return nil
	})
	return
}

func (s *userStoreObject) allUsers(tx ds.IReadTransaction) []User {
	objects, err := tx.GetAll(&ds.GetOptions{Sorted: true, Ascending: true})
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
	Password           string
	MustChangePassword bool
}

func (s *userStoreObject) NewUser(params newUserParameters) (user *User, err *Error) {
	s.Table.StartWrite(func(tx ds.IReadWriteTransaction) error {
		user, err = s.newUser(tx, params)
		return nil
	})
	return
}

func (s *userStoreObject) newUser(tx ds.IReadWriteTransaction, params newUserParameters) (*User, *Error) {
	if s.userWithUsername(tx, params.Username) != nil {
		log.Warn("User with username '%s' already exists", params.Username)
		return nil, ErrorUser("User with username '%s' already exists", params.Username)
	}

	hashedPassword, err := secutil.HashPassword([]byte(params.Password))
	if err != nil {
		log.Error("Error hasing user password: %s", err.Error())
		return nil, ErrorFrom(err)
	}

	user := User{
		Username:           params.Username,
		CanLogIn:           true,
		MustChangePassword: params.MustChangePassword,
	}

	if err := limits.Check(user); err != nil {
		return nil, ErrorUser(err.Error())
	}

	if err := tx.Add(user); err != nil {
		log.Error("Error adding new user '%s': %s", params.Username, err.Error())
		return nil, ErrorFrom(err)
	}
	ShadowStore.Set(user.Username, *hashedPassword)

	UserCache.Update(tx)
	log.Info("New user added: username='%s'", params.Username)
	return &user, nil
}

type editUserParameters struct {
	Password           string
	CanLogIn           bool
	MustChangePassword bool
}

func (s *userStoreObject) EditUser(user *User, params editUserParameters) (newUser *User, err *Error) {
	s.Table.StartWrite(func(tx ds.IReadWriteTransaction) error {
		newUser, err = s.editUser(tx, user, params)
		return nil
	})
	return
}

func (s *userStoreObject) editUser(tx ds.IReadWriteTransaction, user *User, params editUserParameters) (*User, *Error) {
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

	if err := tx.Update(*user); err != nil {
		log.Error("Error updating user '%s': %s", user.Username, err.Error())
		return nil, ErrorFrom(err)
	}

	UserCache.Update(tx)
	return user, nil
}

func (s *userStoreObject) ResetPassword(username string, newPassword []byte) (user *User, err *Error) {
	s.Table.StartWrite(func(tx ds.IReadWriteTransaction) error {
		user, err = s.resetPassword(tx, username, newPassword)
		return nil
	})
	return
}

func (s *userStoreObject) resetPassword(tx ds.IReadWriteTransaction, username string, newPassword []byte) (*User, *Error) {
	user := s.userWithUsername(tx, username)
	if user == nil {
		return nil, ErrorUser("no user with username %s", username)
	}

	passwordHash, err := secutil.HashPassword(newPassword)
	if err != nil {
		log.Error("Error hasing password for user: username='%s' error='%s'", username, err.Error())
		return nil, ErrorFrom(err)
	}

	user.MustChangePassword = false

	if err := tx.Update(*user); err != nil {
		log.Error("Error updating user '%s': %s", username, err.Error())
		return nil, ErrorFrom(err)
	}

	ShadowStore.Set(user.Username, *passwordHash)

	UserCache.Update(tx)
	return user, nil
}

func (s *userStoreObject) ResetAPIKey(username string) (apiKey *string, err *Error) {
	s.Table.StartWrite(func(tx ds.IReadWriteTransaction) error {
		apiKey, err = s.resetAPIKey(tx, username)
		return nil
	})
	return
}

func (s *userStoreObject) resetAPIKey(tx ds.IReadWriteTransaction, username string) (*string, *Error) {
	u := s.userWithUsername(tx, username)
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

func (s *userStoreObject) DeleteUser(user *User) (err *Error) {
	s.Table.StartWrite(func(tx ds.IReadWriteTransaction) error {
		err = s.deleteUser(tx, user)
		return nil
	})
	return
}

func (s *userStoreObject) deleteUser(tx ds.IReadWriteTransaction, user *User) *Error {
	if err := tx.Delete(*user); err != nil {
		log.Error("Error deleting user '%s': %s", user.Username, err.Error())
		return ErrorFrom(err)
	}
	ShadowStore.Delete(user.Username)

	UserCache.Update(tx)
	log.Warn("User deleted: username='%s'", user.Username)
	return nil
}

func (s *userStoreObject) DisableUser(username string) (err *Error) {
	s.Table.StartWrite(func(tx ds.IReadWriteTransaction) error {
		err = s.disableUser(tx, username)
		return nil
	})
	return
}

func (s *userStoreObject) disableUser(tx ds.IReadWriteTransaction, username string) *Error {
	u := UserStore.UserWithUsername(username)
	if u == nil {
		return ErrorUser("no user with username %s", username)
	}
	user := *u

	user.CanLogIn = false

	if err := tx.Update(user); err != nil {
		log.Error("Error updating user '%s': %s", username, err.Error())
		return ErrorFrom(err)
	}

	UserCache.Update(tx)

	log.PWarn("Disabled user", map[string]interface{}{
		"user": username,
	})
	return nil
}
