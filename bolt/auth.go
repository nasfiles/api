package bolt

import (
	"encoding/json"

	"github.com/boltdb/bolt"
	"golang.org/x/crypto/bcrypt"

	"github.com/nasfiles/nasfilesapi"
)

// AuthService implements all the methods to manage authentication
type AuthService struct {
	DB *bolt.DB
}

// Login authenticates a user given its username and password guess
func (s *AuthService) Login(guess *nasfilesapi.AuthGuess) error {
	user := nasfilesapi.User{}

	// Fetch user information
	err := s.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("users"))

		// Get user info
		userBytes := b.Get([]byte(guess.Username))
		if len(userBytes) == 0 {
			return nasfilesapi.ErrUserNotFound
		}

		// Decode JSON information
		if e := json.Unmarshal(userBytes, &user); e != nil {
			return e
		}

		return nil
	})

	// // There was some error fetching the user information
	if err != nil {
		return err
	}

	// Check if password guess matches
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(guess.Password))
	if err != nil {
		return nasfilesapi.ErrWrongPassword
	}

	return nil
}
