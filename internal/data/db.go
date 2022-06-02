package data

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	bolt "go.etcd.io/bbolt"
)

//InitBoltDB - initialize a BoltDB, create a bucket (if not exist) and return a reference for the repo
func InitBoltDB(path string, bucketName string) (*bolt.DB, error) {

	// Initializing the DB
	db, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}

	// Creating a bucket
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return fmt.Errorf("create bucket failed: %s", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}

// getDataByID - returns an entry or an error
func getDataByID(db *bolt.DB, ID string) (Data, error) {
	log.Printf("Retrieving data by ID: %v", ID)

	var entry Data

	err := db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(BucketName))
		if b == nil {
			return fmt.Errorf("Getting bucket failed")
		}

		bs := b.Get([]byte(ID))

		err := json.Unmarshal(bs, &entry)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return Data{}, err
	}

	return entry, err
}

// getDataByID - returns all entries or an error
func getAllEntires(db *bolt.DB) ([]Data, error) {
	log.Printf("Retrieving all entries")

	entries := make([]Data, 0)

	err := db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(BucketName))
		if b == nil {
			return fmt.Errorf("Getting bucket failed")
		}

		c := b.Cursor()

		var tmpData Data
		for k, v := c.First(); k != nil; k, v = c.Next() {
			err2 := json.Unmarshal(v, &tmpData)
			if err2 != nil {
				// on fail, just skip for now
				continue
			}
			entries = append(entries, tmpData)
		}

		return nil
	})

	return entries, err
}

// deleteDataByID - deletes an entry
func deleteDataByID(db *bolt.DB, ID string) error {
	log.Printf(" data by ID: %v", ID)

	err := db.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(BucketName))
		if b == nil {
			return fmt.Errorf("Getting bucket failed")
		}

		return b.Delete([]byte(ID))
	})
	return err
}

// putData - creates a new entry
func putData(db *bolt.DB, d Data) (Data, error) {
	log.Print("Creatting new data entry")

	var ID string

	err := db.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(BucketName))

		id, _ := b.NextSequence()
		ID = fmt.Sprint(id)
		d.Metadata.ID = ID

		buf, err := json.Marshal(d)
		if err != nil {
			return err
		}

		return b.Put([]byte(d.Metadata.ID), buf)
	})

	if err != nil {
		return Data{}, err
	}

	// Read back and return the entry
	dt, err := getDataByID(db, ID)
	if err != nil {
		return Data{}, err
	}
	return dt, nil
}

// putDataByID - updates an entry
func putDataByID(db *bolt.DB, ID string, d Data) error {
	log.Printf("Updating data by ID: %v", ID)

	err := db.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(BucketName))

		d.Metadata.ID = ID

		buf, err := json.Marshal(d)
		if err != nil {
			return err
		}

		return b.Put([]byte(d.Metadata.ID), buf)
	})

	return err
}
