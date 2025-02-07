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
	if err := reg.RegisterOpBuilder("request", func(node *config.Node) (model.IOperator, error) {
		return &Request{}, nil
	}); err != nil {
		return err
	}
	if err := reg.RegisterOpBuilder("response", func(node *config.Node) (model.IOperator, error) {
		return &Response{}, nil
	}); err != nil {
		return err
	}
	if err := reg.RegisterOpBuilder("SubGraphRequest", func(node *config.Node) (model.IOperator, error) {
		return &SubGraphRequest{}, nil
	}); err != nil {
		return err
	}
	if err := reg.RegisterOpBuilder("SubGraphResponse", func(node *config.Node) (model.IOperator, error) {
		return &SubGraphResponse{}, nil
	}); err != nil {
		return err
	}
	return nil
}
