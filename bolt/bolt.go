package bolt

import (
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/fatih/color"

	"github.com/nasfiles/nasfilesapi"
	"github.com/nasfiles/nasfilesapi/terminal"
)

// Setup populates the Bolt database with buckets and all the required stuff
// for the api to work
func Setup(db *bolt.DB) {
	// create users bucket
	createBucketIfNotExists(db, "users")
}

// bucketsList retrieves the list of all buckets
func bucketsList(db *bolt.DB) ([]string, error) {
	var buckets []string

	if err := db.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, _ *bolt.Bucket) error {
			buckets = append(buckets, string(name))
			return nil
		})
	}); err != nil {
		return nil, err
	}

	return buckets, nil
}

// usersList retrieves the list of all users registered
func usersList(db *bolt.DB) ([]nasfilesapi.User, error) {
	users := []nasfilesapi.User{}

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("users"))

		err := b.ForEach(func(k, v []byte) error {
			user := nasfilesapi.User{}

			// decode user json information
			err := json.Unmarshal(v, &user)
			if err != nil {
				return err
			}

			// append to array of users
			users = append(users, user)

			return nil
		})

		return err
	})

	return users, err
}

func dumpConsoleLog(buckets []string, users []nasfilesapi.User) {
	// get terminal size
	width := terminal.Size()

	// Begin
	fmt.Printf("\n\n\n\n")
	terminal.LineSeparator("-", color.New(color.FgBlack).Add(color.BgRed), width)

	color.HiRed("BOLTDB DUMP")

	color.HiCyan("\n%d bucket(s)\n%d user(s)\n\n", len(buckets), len(users))

	// Bucket List
	color.HiYellow("Buckets List:\n")
	for _, bucket := range buckets {
		fmt.Printf(" - ")
		color.HiGreen(bucket)
	}

	// User List
	color.HiYellow("Users List:\n")
	for _, user := range users {
		// Username
		fmt.Printf(" - ")
		color.HiGreen("%s\n", user.Username)

		// Detailed user information
		fmt.Printf("\tName: %s\n\tEmail: %s", user.Name, user.Email)
	}

	terminal.LineSeparator("-", color.New(color.FgBlack).Add(color.BgRed), width)
	fmt.Printf("\n\n\n\n")
	// End
}

// Dump creates a JSON file with all the contents stored in the database in case you
// want to use a .json file as configuration
func Dump(db *bolt.DB, consoleLog bool) error {
	// Bucket List
	buckets, err := bucketsList(db)
	if err != nil {
		return err
	}

	// User list
	users, err := usersList(db)
	if err != nil {
		return err
	}

	// Print database information
	if consoleLog {
		dumpConsoleLog(buckets, users)
	}

	return nil
}

// createBucketIfNotExists creates a bucket if it doesn't exist
func createBucketIfNotExists(db *bolt.DB, bucketName string) error {
	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))

		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
