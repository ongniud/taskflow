package utils

import (
	"github.com/ongniud/taskflow/model"
)

func FieldsToMap(fields []*model.FieldData) map[string]any {
	res := make(map[string]any, len(fields))
	for _, f := range fields {
		if f != nil && f.Data != nil {
			field := f.Field.Mapping
			if field == "" {
				field = f.Field.Name
			}
			res[field] = f.Data.Val
		}
	}
	return res
}
