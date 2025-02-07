package ops

import (
	"testing"
)

func TestDebugBindInputs(t *testing.T) {
	//inputs := map[string]any{
	//	"request_id":  "abc:def:ghi",
	//	"user_id":     int64(1234567),
	//	"session_ids": []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
	//	"HostIds":     []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
	//}
	//
	//debug := &Debug2{}
	//debugWithBinding := common.WithBinding(debug)
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
