package ops

import (
	"github.com/ongniud/taskflow/model"
	"github.com/ongniud/taskflow/model/config"
	reg "github.com/ongniud/taskflow/registry"
)

func Init() error {
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
	if err := reg.RegisterOpBuilder("process", func(node *config.Node) (model.IOperator, error) {
		return &Process{}, nil
	}); err != nil {
		return err
	}
	return nil
}
