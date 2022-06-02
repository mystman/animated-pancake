package service

import (
	"encoding/json"
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
	repo := data.NewRepo(dbFilePath)
	svc := NewService(repo)

	// POST
	t.Run("Create a data entry", func(t *testing.T) {
		_, err := svc.PostData(data.NetworkType, data.Data{Payload: getPayloadSample()})

		if err != nil {
			t.Fatalf("Failed to create Data with PostData: %v", err)
		}

		// t.Logf("PostData entry %v", dt)
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

		// t.Logf("GetData entry %#v", got)

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

		// t.Logf("GetData entry %#v", got)
	})

	// GET all entries
	t.Run("Get all entries", func(t *testing.T) {

		got, err := svc.GetAllData("", "")
		if err != nil {
			t.Fatalf("Failed to retrieve all entries with GetAllData: %v", err)
		}

		if len(got) <= 0 {
			t.Fatalf("No entries in the list: %v", err)
		}

		// t.Logf("Got all entries %#v", got)
	})

	// POST & DELETE
	t.Run("Create and and delete data entry", func(t *testing.T) {

		want, err1 := svc.PostData(data.NetworkType, data.Data{Payload: getPayloadSample()})
		ret, err2 := svc.GetData(want.Metadata.ID)

		if err1 != nil {
			t.Fatalf("Failed to create Data with PostData: %v", err1)
		}
		if err2 != nil {
			t.Fatalf("Failed to retrieve Data with GetData: %v", err2)
		}

		err3 := svc.DeleteData(ret.Metadata.ID)
		if err3 != nil {
			t.Fatalf("Failed to delete Data: %v", err3)
		}

		_, err4 := svc.GetData(ret.Metadata.ID)
		if err4 == nil {
			t.Errorf("Deleted data with ID %v is not removed: %v", ret.Metadata.ID, err4)
		}

	})

	// POST & UPDATE
	t.Run("Create and update a data entry", func(t *testing.T) {

		want, err1 := svc.PostData(data.NetworkType, data.Data{Payload: getPayloadSample()})
		got, err2 := svc.GetData(want.Metadata.ID)

		if err1 != nil {
			t.Fatalf("Failed to create Data with PostData: %v", err1)
		}
		if err2 != nil {
			t.Fatalf("Failed to retrieve Data with GetData: %v", err2)
		}

		// // Store original
		// gotPayload := []byte(fmt.Sprintf("%v", got.Payload))
		// var expectedPayload map[string]interface{}
		// json.Unmarshal(gotPayload, &expectedPayload)

		// Change data
		tmpPayload := make(map[string]string)
		tmpPayload["subnetName"] = "CHANGEDNET"
		tmpPayload["ipRange"] = "1.2.3.4/24"

		tmpJSONPayload, _ := json.Marshal(tmpPayload)
		got.Payload = string(tmpJSONPayload)

		// Perform update
		err3 := svc.UpdateData(got.Metadata.ID, got)
		if err3 != nil {
			t.Fatalf("Failed to Update data: %v", err3)
		}

		// Read back updated data
		got2, err4 := svc.GetData(want.Metadata.ID)
		if err4 != nil {
			t.Fatalf("Failed to retrieve back update Data %v", err4)
		}

		got2PayloadBs := []byte(fmt.Sprintf("%v", got2.Payload))
		var got2Payload map[string]interface{}
		json.Unmarshal(got2PayloadBs, &got2Payload)

		if tmpPayload["subnetName"] != got2Payload["subnetName"] || tmpPayload["ipRange"] != got2Payload["ipRange"] {
			t.Fatalf("Post and Get data does not match |\n EXPECTED:\n%v\n\nGOT2:\n%v\n", got, got2)
		}

	})

	// Cleanup (just in case)
	defer os.Remove(dbFilePath)
}

// Helper function to get payload string
func getPayloadSample() string {
	return `{ "subnetName": "testSubnet", "ipRange": "10.2.1.0/24" }`
}
