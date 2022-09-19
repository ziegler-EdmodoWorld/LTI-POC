package engine

import (
	"fmt"
	"strings"
	"time"

	"github.com/mlhoyt/ramsql/engine/log"
	"github.com/mlhoyt/ramsql/engine/parser"
)

func updateValues(r *Relation, row int, values map[string]interface{}) error {
	for i := range r.table.attributes {
		val, ok := values[r.table.attributes[i].name]
		if !ok {
			switch onUpdateVal := r.table.attributes[i].onUpdateValue.(type) {
			case func() interface{}:
				val = (func() interface{})(onUpdateVal)()
			default:
				continue
			}
		} else {
			log.Debug("Type of '%s' is '%s'\n", r.table.attributes[i].name, r.table.attributes[i].typeName)
			switch strings.ToLower(r.table.attributes[i].typeName) {
			case "timestamp", "localtimestamp":
				switch valVal := val.(type) {
				case func() interface{}:
					val = (func() interface{})(valVal)()
				case time.Time:
					// format time.Time into parsable string
					val = valVal.Format(parser.DateLongFormat)
				case string:
					if valVal == "current_timestamp" || valVal == "now()" {
						val = time.Now().Format(parser.DateLongFormat)
					}
				}
			}
		}
		r.rows[row].Values[i] = fmt.Sprintf("%v", val)
	}

	return nil
}
