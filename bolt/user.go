package bolt

import (
	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/nasfiles/api"
)

//UserService implements all the methods to manage users
type UserService struct {
	DB *bolt.DB
}

//Add user
func (s *UserService) Add(u *api.User) error {
	err := s.DB.Update(func(tx *bolt.Tx) error {
		// select user bucket
		b := tx.Bucket([]byte("users"))

		// Encode User struct to json
		userJSON, e := json.Marshal(*u)
		if e != nil {
			return e
		}

		// put new data into the bucket
		err := b.Put([]byte(u.UID), userJSON)
		return err
	})

	if err != nil {
		return err
	}

	return nil
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
