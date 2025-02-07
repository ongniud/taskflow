package model

import (
	"github.com/ongniud/taskflow/model/config"
)

type Data struct {
	Typ string
	Val interface{}
}

type FieldData struct {
	Field *config.Field
	Data  *Data
}
