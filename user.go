package nasfilesapi

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/webdav"
)

//User is a struct holding important details about every account
type User struct {
	UID      int       `json:"UID"`
	Username string    `json:"Username"`
	Email    string    `json:"Email"`
	Name     string    `json:"Name"`
	Password string    `json:",omitempty"`
	Admin    bool      `json:"Admin"`
	Created  time.Time `json:"Created"`

	Handler *webdav.Handler `json:"-"`
}

//UserService ...
type UserService interface {
	Add(u *User) error
	GetByUsername(username string) (*User, error)
	GetsAll() (*[]User, error)
	Delete(uid string) error
}

// SetPassword generates an hashed password from a string
func (u *User) SetPassword(password string) error {
	if len(password) <= 4 {
		return ErrPasswordTooShort
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update password
	u.Password = string(hashedPassword)

	return nil
}
