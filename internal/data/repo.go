package data

import (
	"log"
	"strconv"
	"sync"
	"time"
)

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
	m      map[int]Data
	lastID int //temporary
	l      sync.RWMutex
}

// Temporary -  thread safe ID generation for testing
func (r *Repo) newID() string {
	r.l.Lock()
	defer r.l.Unlock()

	r.lastID++

	return strconv.Itoa(r.lastID)
}

// PostData - takes Data and stores it as a new entity
func (r *Repo) PostData(tp string, d Data) (Data, error) {

	dta := Data{
		Metadata: Metadata{
			ID:          r.newID(),
			LastUpdated: time.Now().Format(time.RFC3339),
			Type:        tp,
		},
		Payload: d.Payload,
	}
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
	log.Printf("Initializing new repository at %v", dbFilePath)

	InitBoltDB(dbFilePath)

	mp := make(map[int]Data)

	r := &Repo{
		m: mp,
	}

	rp := Repository(r)

	return &rp
}
