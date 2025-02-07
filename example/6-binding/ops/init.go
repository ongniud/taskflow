package ops

import (
	"github.com/ongniud/taskflow/model"
	"github.com/ongniud/taskflow/model/config"
	"github.com/ongniud/taskflow/operator/mw"
	reg "github.com/ongniud/taskflow/registry"
)

func Init() error {
	if err := reg.RegisterOpBuilder("start", func(node *config.Node) (model.IOperator, error) {
		return &Start{}, nil
	}); err != nil {
		return err
	}
	if err := reg.RegisterOpBuilder("binding", func(node *config.Node) (model.IOperator, error) {
		return mw.WithBinding(&Binding{}), nil
	}); err != nil {
		return err
	}
	if err := reg.RegisterOpBuilder("end", func(node *config.Node) (model.IOperator, error) {
		return &End{}, nil
	}); err != nil {
		return err
	}
	return nil
}
