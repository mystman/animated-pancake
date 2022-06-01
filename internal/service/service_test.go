package service

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"testing"

	"mystman.com/animated-pancake/internal/data"
)

func TestService(t *testing.T) {

	log.SetOutput(ioutil.Discard)

	// t.Logf("temp dir is: %v", os.TempDir())
	// t.Logf("dbFilePath: %v", dbFilePath)

	dbFilePath := fmt.Sprintf("%sanim-pancake.db", os.TempDir())

	// Setup service
	repo := data.NewRepository(dbFilePath)
	svc := NewService(repo)

	// POST
	t.Run("Create a data entry", func(t *testing.T) {
		dt, err := svc.PostData(data.NetworkType, data.Data{Payload: getPayloadSample()})

		if err != nil {
			t.Fatalf("Failed to create Data with PostData: %v", err)
		}

		t.Logf("PostData entry %v", dt)
	})

	// POST & GET
	t.Run("Create and retrieve data entry", func(t *testing.T) {

		want, err1 := svc.PostData(data.NetworkType, data.Data{Payload: getPayloadSample()})
		got, err2 := svc.GetData(want.Metadata.ID)

		if err1 != nil {
			t.Fatalf("Failed to create Data with PostData: %v", err1)
		}
		if err2 != nil {
			t.Fatalf("Failed to retrieve Data with GetData: %v", err2)
		}

		t.Logf("GetData entry %#v", got)

		if !reflect.DeepEqual(got, want) {
			t.Fatalf("Post and Get data does not match | POST: %#v | GET: %#v", want, got)
		}

	})

	// GET
	t.Run("Get entry by non-existing ID", func(t *testing.T) {

		ID := "987654321"
		got, err := svc.GetData(ID)

		if err == nil {
			t.Errorf("Should have returned an error, but got %#v", got)
		}

		t.Logf("GetData entry %#v", got)

	})

	// GET all entries
	t.Run("Get all entries", func(t *testing.T) {

		got, err := svc.GetAllData("", "")
		if err != nil {
			t.Fatalf("Failed to retrieve all entries with GetAllData: %v", err)
		}
		t.Logf("Got all entries %#v", got)

	})

	defer os.Remove(dbFilePath)
}

// Helper function to get payload string
func getPayloadSample() string {
	return `{ "subnetName": "testSubnet", "ipRange": "10.2.1.0/24" }`
}
