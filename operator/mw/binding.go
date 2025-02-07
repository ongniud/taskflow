package mw

import (
	"github.com/ongniud/taskflow/model"
	"github.com/ongniud/taskflow/utils"
	"github.com/pkg/errors"
)

type BindingMW struct {
	Op model.IOperator
}

func NewBindingMW(op model.IOperator) *BindingMW {
	return &BindingMW{
		Op: op,
	}
}

func (op *BindingMW) Execute(ctx model.IOpContext) error {
	inputs := ctx.GetInputs()
	if err := utils.BindInputs(op.Op, inputs); err != nil {
		return errors.Wrap(err, "bind inputs err")
	}

	if err := op.Op.Execute(ctx); err != nil {
		return errors.Wrap(err, "bind inputs err")
	}

	outputs, err := utils.BindOutputs(op.Op)
	if err != nil {
		return errors.Wrap(err, "bind outputs err")
	}

	ctx.SetOutputs(outputs)
	return nil
}

func WithBinding(op model.IOperator) model.IOperator {
	return NewBindingMW(op)
}
