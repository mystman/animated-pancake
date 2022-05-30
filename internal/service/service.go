package service

import "log"

// Service - is a struct to hold all service logic
type Service struct {
}

// NewService - creates and returns a new instace of a serice
func NewService() *Service {
	log.Printf("Initializing new service")
	return &Service{}
}
