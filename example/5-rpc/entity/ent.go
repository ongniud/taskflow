package entity

type Request struct {
	ReqId    string
	Scene    string
	Country  string
	Uid      int
	DeviceID string
}

type Response struct {
	Identity string
}
