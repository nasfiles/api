package bolt

import (
	"fmt"

	"github.com/nasfiles/api"
	"github.com/nasfiles/api/terminal"

	"github.com/boltdb/bolt"
	"github.com/fatih/color"
)

//Setup populates the Bolt database with buckets and all the required stuff
//for the api to work
func Setup(db *bolt.DB) {
	// create users bucket
	createBucketIfNotExists(db, "users")
}

//Dump creates a JSON file with all the contents stored in the database in case you
//want to use a .json file as configuration
func Dump(db *bolt.DB, consoleLog bool) error {
	// Bucket List
	buckets, err := bucketsList(db)
	if err != nil {
		return err
	}

	if consoleLog {
		dumpConsoleLog(buckets, []api.User{})
	}

	return nil
}

func dumpConsoleLog(buckets []string, users []api.User) {
	// get terminal size
	width := terminal.TerminalSize()

	// Begin
	fmt.Printf("\n\n\n\n")
	terminal.LineSeparator("-", color.New(color.FgBlack).Add(color.BgRed), width)

	color.HiRed("BOLTDB DUMP")

	color.HiCyan("\n%d bucket(s)\n%d user(s)\n\n", len(buckets), len(users))

	// Bucket List
	color.HiYellow("Buckets List:\n")
	for _, bucket := range buckets {
		fmt.Printf("  - ")
		color.HiGreen(bucket)
	}

	// User List
	color.HiYellow("Users List:\n")
	for _, user := range users {
		fmt.Printf("  - ")
		color.HiGreen(user.Name)
	}

	terminal.LineSeparator("-", color.New(color.FgBlack).Add(color.BgRed), width)
	fmt.Printf("\n\n\n\n")
	// End
}

//BucketsList
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

//createBucketIfNotExists creates a bucket if it doesn't exist
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
