package service

import (
	"fmt"
	"os"
	"testing"

	"mystman.com/animated-pancake/internal/data"
)

func TestRepo(t *testing.T) {

	// t.Logf("temp dir is: %v", os.TempDir())

	dbFilePath := fmt.Sprintf("%sanim-pancake.db", os.TempDir())

	// t.Logf("dbFilePath: %v", dbFilePath)

	// starting the service
	repo := data.NewRepository(dbFilePath)
	svc := NewService(repo)

	// Test
	dt, err := svc.PostData(data.NetworkType, data.Data{Payload: "-=IPSUM LOREM=-"})

	if err != nil {
		t.Fatalf("Failed to PostData: %v", err)
	}

	t.Logf("PostData entry %v", dt)

	// defer os.Remove(dbFilePath)
}
