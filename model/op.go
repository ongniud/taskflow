package model

type IOperator interface {
	Execute(ctx IOpContext) error
}
