package bolt

import (
	"encoding/json"
	"errors"

	"github.com/boltdb/bolt"

	"github.com/nasfiles/nasfilesapi"
)

//UserService implements all the methods to manage users
type UserService struct {
	DB *bolt.DB
}

//Add user
func (s *UserService) Add(u *nasfilesapi.User) error {
	err := s.DB.Update(func(tx *bolt.Tx) error {
		// select user bucket
		b := tx.Bucket([]byte("users"))

		// Encode User struct to json
		userJSON, e := json.Marshal(*u)
		if e != nil {
			return e
		}

		// put new data into the bucket
		err := b.Put([]byte(u.Username), userJSON)
		return err
	})

	if err != nil {
		return err
	}

	return nil
}

//GetByUsername retrieves a user from the database given its uid
func (s *UserService) GetByUsername(username string) (*nasfilesapi.User, error) {
	var u *nasfilesapi.User

	err := s.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("users"))

		userBytes := b.Get([]byte(username))

		if len(userBytes) == 0 {
			return errors.New("Not found")
		}

		if e := json.Unmarshal(userBytes, &u); e != nil {
			return e
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return u, nil
}

//Delete user
func (s *UserService) Delete(uid string) error {
	err := s.DB.Update(func(tx *bolt.Tx) error {
		// select user bucket
		b := tx.Bucket([]byte("users"))

		// delete user with given uid
		err := b.Delete([]byte(uid))
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
