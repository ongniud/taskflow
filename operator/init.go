package operator

import (
	"github.com/ongniud/taskflow/model"
	"github.com/ongniud/taskflow/model/config"
	"github.com/ongniud/taskflow/operator/ops"
	reg "github.com/ongniud/taskflow/registry"
)

func Init() error {
	if err := reg.RegisterOpBuilder("debug", func(node *config.Node) (model.IOperator, error) {
		return &ops.Debug{}, nil
	}); err != nil {
		return err
	}
	if err := reg.RegisterOpBuilder("request", func(node *config.Node) (model.IOperator, error) {
		return &ops.Request{}, nil
	}); err != nil {
		return err
	}
	if err := reg.RegisterOpBuilder("response", func(node *config.Node) (model.IOperator, error) {
		return &ops.Response{}, nil
	}); err != nil {
		return err
	}
	if err := reg.RegisterOpBuilder("SubGraphRequest", func(node *config.Node) (model.IOperator, error) {
		return &ops.SubGraphRequest{}, nil
	}); err != nil {
		return err
	}
	if err := reg.RegisterOpBuilder("SubGraphResponse", func(node *config.Node) (model.IOperator, error) {
		return &ops.SubGraphResponse{}, nil
	}); err != nil {
		return err
	}
	return nil
}
