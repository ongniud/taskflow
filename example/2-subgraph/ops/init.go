package ops

import (
	"github.com/ongniud/taskflow/model"
	"github.com/ongniud/taskflow/model/config"
	reg "github.com/ongniud/taskflow/registry"
)

func Init() error {
	if err := reg.RegisterOpBuilder("debug", func(node *config.Node) (model.IOperator, error) {
		return &Debug{}, nil
	}); err != nil {
		return err
	}
	return nil
}
