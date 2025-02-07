package example

type Request struct {
	ReqId    string
	Scene    string
	Country  string
	Uid      int32
	DeviceID string
}

type Doc struct {
	Typ   string
	ID    string
	Score float64
	Sigs  string
}

type Response struct {
	Docs []*Doc
}
