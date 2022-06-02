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
func getDataByID(db *bolt.DB, ID string) ([]byte, error) {
	log.Printf("Retrieving data by ID: %v", ID)

	var entry []byte

	err := db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(BucketName))
		if b == nil {
			return fmt.Errorf("Getting bucket failed")
		}

		entry = b.Get([]byte(ID))

		return nil
	})

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
	log.Printf("Deleting data by ID: %v", ID)

	err := db.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(BucketName))
		if b == nil {
			return fmt.Errorf("Getting bucket failed")
		}

		return b.Delete([]byte(ID))
	})
	return err
}

//==================================================================
// InitSQLiteDB - Initializes an SQLite DB from a file for testing
//==================================================================

// Imports for SQLite:
// _ "github.com/jmoiron/sqlx"
// _ "github.com/mattn/go-sqlite3"

/*
func InitDB() (*sqlx.DB, error) {

	location := "/usr/share/pancake-data/pancake.db"

	db, err := sqlx.Open("sqlite3", location)

	if err != nil {
		log.Fatalf("Error opening location %v: %v", location, err)
	}

	//	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging DB %v: %v", db, err)
	}

	log.Printf("Pinging DB %v was successfull", db)

	schema := `CREATE TABLE IF NOT EXISTS data (
		ID INTEGER PRIMARY KEY AUTOINCREMENT,
		lastUpdated TEXT,
		type TEXT,
		data TEXT);`

	result, err := db.Exec(schema)
	if err != nil {
		log.Fatalf("Schema creation failed: %v", err)
	}
	log.Printf("Schema creation was successfull: %v", result)

	for i := 0; i < 10; i++ {
		entry := `INSERT INTO data (lastUpdated, type, data) VALUES (?, ?, ?)`
		db.MustExec(entry, "donno", NetworkType, "{'value':'nonsese'}")
	}

	log.Println("Insert done")

	rows, err := db.Query("SELECT ID, lastUpdated, type, data FROM data")

	// iterate over each row
	for rows.Next() {
		var ID int
		var lastUpdated, typ, data string
		err = rows.Scan(&ID, &lastUpdated, &typ, &data)

		log.Printf("Rows data: %v | %v | %v | %v", ID, lastUpdated, typ, data)
	}
	// check the error from rows
	err = rows.Err()
	if err != nil {
		log.Printf("Error from the rows: %v", err)
	}

	return db, nil
}
*/
