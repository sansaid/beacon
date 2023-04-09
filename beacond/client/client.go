package client

import (
	"beacon/beacond/server"
)

type BeaconClient interface {
	CreateProbe(string, string) error
	ListProbes() server.ListProbesResponse
	DeleteProbe(string, string) error
	Describe() server.BeaconDescribeResponse
	SetURL(string)
	URL() string
}

type Beacon struct {
	url string
}

type BeaconOption func(...interface{}) func(BeaconClient) BeaconClient

func NewBeaconClient(opts ...BeaconOption) BeaconClient {
	beaconClient := new(Beacon)

	beaconClient.SetURL("http://localhost:3232")

	for _, opt := range opts {
		opt(beaconClient)
	}

	return beaconClient
}

func (b *Beacon) CreateProbe(namespace string, repo string) error {
	// http.Post(fmt.Sprintf("%s%s?namespace=%s&repo=%s", b.URL(), "/probe", namespace, repo), "application/json", &bytes.Buffer{})
	panic("not implemented")
}

func (b *Beacon) ListProbes() server.ListProbesResponse {
	panic("not implemented")
}

func (b *Beacon) DeleteProbe(namespace string, repo string) error {
	panic("not implemented")
}

func (b *Beacon) Describe() server.BeaconDescribeResponse {
	panic("not implemented")
}

func (b *Beacon) SetURL(url string) {
	b.url = url
}

func (b *Beacon) URL() string {
	return b.url
}

func WithUrl(url string) func(BeaconClient) BeaconClient {
	return func(b BeaconClient) BeaconClient {
		b.SetURL(url)

		return b
	}
}
