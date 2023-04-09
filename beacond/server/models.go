package server

type BaseResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

type BeaconDescribeResponse struct {
	Registry string   `json:"registry"`
	Probes   []string `json:"probes"`
	Runtime  string   `json:"runtime"`
}

type ListProbesResponse struct {
	Probes []string `json:"probes"`
}
