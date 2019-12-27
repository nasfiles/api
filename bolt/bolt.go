package bolt

import (
	"github.com/boltdb/bolt"
)

//Setup populates the Bolt database with buckets and all the required stuff
//for the api to work
func Setup(db *bolt.DB) {
	// create users bucket
	createBucketIfNotExists(db, "users")
}

//Dump creates a JSON file with all the contents stored in the database in case you
//want to use a .json file as configuration
func Dump(db *bolt.DB) error {

	return nil
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
