package client

type BeaconClient interface {
	CreateProbe()
}

func NewBeaconClient() BeaconClient {
	panic("not implemented")
}
