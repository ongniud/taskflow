package entity

type Request struct {
	ReqId    string
	Scene    string
	Country  string
	Uid      int32
	DeviceID string
}

type Doc struct {
	Title  string
	Author string
	Text   string
}

type Response struct {
	Docs []*Doc
}
