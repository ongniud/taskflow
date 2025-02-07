package ops

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/ongniud/taskflow/example/entity"
	"github.com/ongniud/taskflow/model"
)

type Response struct{}

func (op *Response) Execute(ctx model.IOpContext) error {
	reqObj, _ := ctx.GetParam("request")
	if reqObj == nil {
		fmt.Println("[Response] request not exist")
		return fmt.Errorf("request not exist")
	}
	req, ok := reqObj.(*entity.Request)
	if !ok {
		fmt.Println("[Response] request type invalid")
		return fmt.Errorf("request type invalid")
	}

	reqStr, _ := json.Marshal(req)
	fmt.Println("[Response] req:", string(reqStr))

	rsp := entity.Response{}

	for i := 0; i < 10; i++ {
		doc := &entity.Doc{
			Title:  strconv.FormatInt(int64(i), 10),
			Text:   "item",
			Author: "7788",
		}
		rsp.Docs = append(rsp.Docs, doc)
	}

	rspStr, _ := json.Marshal(rsp)
	fmt.Println("[Response] rsp:", string(rspStr))

	ctx.SetGraphOutput("response", rsp)
	return nil
}
