package service

import (
	"log"

	"mystman.com/animated-pancake/internal/data"
)

// Service - is a struct to hold all service logic
type Service struct {
	Repo data.Repository
}

// PostData - creates a new entry
func (svc *Service) PostData(tp string, d data.Data) (data.Data, error) {
	return svc.Repo.PostData(tp, d)
}

// UpdateData - updates an existing entry
func (svc *Service) UpdateData(d data.Data) (data.Data, error) {
	return svc.Repo.UpdateData(d)
}

// GetData - gets an entry by ID
func (svc *Service) GetData(ID string) (data.Data, error) {
	return svc.Repo.GetData(ID)
}

// DeleteData - deletes and entry by ID
func (svc *Service) DeleteData(ID string) error {
	return svc.Repo.DeleteData(ID)
}

// GetAllData - gets entries matching the conditions
func (svc *Service) GetAllData(ID string, typ string) ([]data.Data, error) {
	return svc.Repo.GetAllData(ID, typ)
}

// NewService - creates and returns a new instace of a serive
func NewService(repo *data.Repository) *Service {
	log.Printf("Initializing new service")
	return &Service{
		Repo: *repo,
	}
}
