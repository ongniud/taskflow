package ops

import (
	"fmt"
	"testing"

	"github.com/ongniud/taskflow/model"
	"github.com/ongniud/taskflow/utils"

	jsoniter "github.com/json-iterator/go"
)

func TestBindInputs(t *testing.T) {
	inputs := map[string]any{
		"request_id":  "abc:def:ghi",
		"user_id":     int64(1234567),
		"session_ids": []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		"HostIds":     []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
	}

	debug := &Debug{}

	var op model.IOperator
	op = debug
	if err := utils.BindInputs(op, inputs); err != nil {
		t.Fatal(err)
	}

	debugStr, _ := jsoniter.MarshalToString(debug)
	fmt.Println(debugStr)
}

//func TestBindOutputs(t *testing.T) {
//	debug := &Debug{}
//	debug.Req = proto.String("abc:def:ghi")
//	debug.Ctx = proto.String("ctx")
//	var op model.IOperator
//	op = debug
//	outputs, err := utils.BindOutputs(op)
//	if err != nil {
//		t.Fatal(err)
//	}
//	debugStr, _ := json.Marshal(outputs)
//	fmt.Println(string(debugStr))
//}
