package data

import (
	"log"
	"time"

	bolt "go.etcd.io/bbolt"
)

//BucketName - name of the BoltDB bucket
const BucketName = "pancake"

// Repository - interface for the data storage
type Repository interface {
	PostData(tp string, d Data) (Data, error)
	UpdateData(ID string, d Data) error
	GetData(ID string) (Data, error)
	DeleteData(ID string) error
	GetAllData(ID string, typ string) ([]Data, error)
}

// Repo - is an implementetion of the Repository interface
type Repo struct {
	db         *bolt.DB
	bucketName string
}

// NewRepo - creates a new instance Repository
func NewRepo(dbFilePath string) *Repo {

	if len(dbFilePath) == 0 {
		log.Fatalf("Repository file location is required")
	}
	log.Printf("Initializing repository at %v", dbFilePath)

	boltDB, err := InitBoltDB(dbFilePath, BucketName)
	if err != nil {
		log.Fatalf("Initializing Bolt DB failed: %v", err)
	}

	return &Repo{
		db:         boltDB,
		bucketName: BucketName,
	}

}

// PostData - takes Data and stores it as a new entity
func (r *Repo) PostData(tp string, d Data) (Data, error) {

	tmpData := Data{
		Metadata: Metadata{
			LastUpdated: getISOTimestamp(),
			Type:        tp,
		},
		Payload: d.Payload,
	}

	dta, err := putData(r.db, tmpData)

	if err != nil {
		return Data{}, err
	}

	log.Printf("Created new entry %v: %v", dta.Metadata.ID, dta)

	return dta, nil
}

// GetData - based on an ID it returns Data
func (r *Repo) GetData(ID string) (Data, error) {

	dta, err := getDataByID(r.db, ID)
	if err != nil {
		return Data{}, err
	}

	return dta, nil
}

// UpdateData - takes Data and stores it as a new entity
func (r *Repo) UpdateData(ID string, d Data) error {

	d.Metadata.LastUpdated = getISOTimestamp()

	return putDataByID(r.db, ID, d)
}

// DeleteData - based on the ID it removes a stored Data
func (r *Repo) DeleteData(ID string) error {
	return deleteDataByID(r.db, ID)
}

// GetAllData - based on the parameters it returns the matching Data as a slice
func (r *Repo) GetAllData(ID string, typ string) ([]Data, error) {

	dt, err := getAllEntires(r.db)

	if err != nil {
		return []Data{}, err
	}

	// TODO: filter out not requied OR delegate the filtering to getAllEntires()

	log.Printf("GetAllData retrieved %v entries", len(dt))

	return dt, nil
}

func getISOTimestamp() string {
	return time.Now().Format(time.RFC3339)
}
