package config

// Metadata holds static service information.
type Metadata struct {
	ServiceName string
	Version     string
}

func GetMetadata() Metadata {
	return Metadata{
		ServiceName: "Fibre_Rate_Limit_Service",
		Version:     "1.0.0",
	}
}
