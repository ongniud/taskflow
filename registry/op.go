package registry

import (
	"fmt"

	"github.com/ongniud/taskflow/model"
	"github.com/ongniud/taskflow/model/config"
)

func GetOp(op string, node *config.Node) (model.IOperator, error) {
	builder := GetOpBuilder(op)
	if builder == nil {
		return nil, fmt.Errorf("builder[%s] not found", op)
	}
	iop, err := builder(node)
	return iop, err
}
