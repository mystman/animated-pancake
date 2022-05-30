package data

// Metadata - additional metadata stored with the request
// It is for internal use only; will not be displayed for clients
type Metadata struct {
	ID          string
	LastUpdated string
	Type        string
}

// Data - data to be stored in the database
// Consists the actual payload and metadata
type Data struct {
	Metadata Metadata
	Payload  interface{}
}

const (
	//NetworkType -
	NetworkType string = "network"

	//FirewallType -
	FirewallType = "firewall"

	//ApplicationType -
	ApplicationType = "application"

	//DataType -
	DataType = "data"
)
