package ops

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/ongniud/taskflow/example"
	"github.com/ongniud/taskflow/model"
)

type Response struct{}

func (op *Response) Execute(ctx model.IOpContext) error {
	reqObj, _ := ctx.GetParam("request")
	if reqObj == nil {
		fmt.Println("[Response] request not exist")
		return fmt.Errorf("request not exist")
	}
	req, ok := reqObj.(*example.Request)
	if !ok {
		fmt.Println("[Response] request type invalid")
		return fmt.Errorf("request type invalid")
	}

	reqStr, _ := json.Marshal(req)
	fmt.Println("[Response] req:", string(reqStr))

	rsp := example.Response{}

	for i := 0; i < 10; i++ {
		doc := &example.Doc{
			ID:    strconv.FormatInt(int64(i), 10),
			Typ:   "item",
			Score: float64(i) + 123.456,
			Sigs:  "7788",
		}
		rsp.Docs = append(rsp.Docs, doc)
	}

	rspStr, _ := json.Marshal(rsp)
	fmt.Println("[Response] rsp:", string(rspStr))

	ctx.SetGraphOutput("response", rsp)
	return nil
}
