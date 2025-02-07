package registry

import (
	"fmt"
	"sync"

	"github.com/ongniud/taskflow/model"
	"github.com/ongniud/taskflow/model/config"
)

var (
	builders sync.Map
)

// TODO:
// comments:
// why we need builder here, because op is state relate and not reusable, each execute need new one.
// is we add a reset here maybe can reuse it?

type OpBuilder func(node *config.Node) (model.IOperator, error)

func RegisterOpBuilder(op string, builder OpBuilder) error {
	if op == "" || builder == nil {
		return fmt.Errorf("arg null")
	}
	_, ok := builders.Load(op)
	if ok {
		return fmt.Errorf("op builder already exists")
	}
	builders.Store(op, builder)
	return nil
}

func GetOpBuilder(op string) OpBuilder {
	v, ok := builders.Load(op)
	if !ok {
		return nil
	}
	b, ok := v.(OpBuilder)
	if !ok {
		return nil
	}
	return b
}
