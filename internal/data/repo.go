package data

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	bolt "go.etcd.io/bbolt"
)

//bucketName - name of the BoltDB bucket
const BucketName = "pancake"

// Repository - interface for the data storage
type Repository interface {
	PostData(tp string, d Data) (Data, error)
	UpdateData(d Data) (Data, error)
	GetData(ID string) (Data, error)
	DeleteData(ID string) error
	GetAllData(ID string, typ string) ([]Data, error)
}

// Repo - is an implementetion of the Repository interface
type Repo struct {
	db         *bolt.DB
	bucketName string
}

// PostData - takes Data and stores it as a new entity
func (r *Repo) PostData(tp string, d Data) (Data, error) {

	dta := Data{
		Metadata: Metadata{
			// ID:          r.newID(),
			LastUpdated: time.Now().Format(time.RFC3339),
			Type:        tp,
		},
		Payload: d.Payload,
	}

	var ID string

	err := r.db.Update(func(tx *bolt.Tx) error {
		// Retrieve the users bucket.
		// This should be created when the DB is first opened.
		b := tx.Bucket([]byte(r.bucketName))

		// Generate ID for the user.
		// This returns an error only if the Tx is closed or not writeable.
		// That can't happen in an Update() call so I ignore the error check.
		id, _ := b.NextSequence()
		ID = fmt.Sprint(id)
		dta.Metadata.ID = ID

		// Marshal user data into bytes.
		buf, err := json.Marshal(dta)
		if err != nil {
			return err
		}

		// Persist bytes to users bucket.
		return b.Put([]byte(dta.Metadata.ID), buf)
	})

	if err != nil {
		return Data{}, err
	}

	// Testing all
	// TestDisplayAllEnteries(r.db)
	got, err := getDataByID(r.db, ID)
	if err != nil {
		return Data{}, err
	}

	log.Printf("Created new entry %v: %v", ID, string(got))

	return dta, nil
}

// UpdateData - takes Data and stores it as a new entity
func (r *Repo) UpdateData(d Data) (Data, error) {
	return Data{}, nil
}

// GetData - based on an ID it returns Data
func (r *Repo) GetData(ID string) (Data, error) {
	return Data{}, nil
}

// DeleteData - based on the ID it removes a stored Data
func (r *Repo) DeleteData(ID string) error {
	return nil
}

// GetAllData - based on the parameters it returns the matching Data as a slice
func (r *Repo) GetAllData(ID string, typ string) ([]Data, error) {
	return []Data{}, nil
}

// NewRepository - creates a new instance Repository
func NewRepository(dbFilePath string) *Repository {

	if len(dbFilePath) == 0 {
		log.Fatalf("Repository file location is required")
	}
	log.Printf("Initializing repository at %v", dbFilePath)

	boltDB, err := InitBoltDB(dbFilePath, BucketName)
	if err != nil {
		log.Fatalf("Initializing Bolt DB failed: %v", err)
	}

	r := &Repo{
		db:         boltDB,
		bucketName: BucketName,
	}

	rp := Repository(r)

	return &rp
}
