package webkit

import (
	"github.com/dairaga/log"
	"github.com/gorilla/schema"
)

var _schema = schema.NewDecoder()

func init() {
	_schema.IgnoreUnknownKeys(true)
	_schema.ZeroEmpty(true)
	log.Debug("init schema")
}

// SchemaIgnorUnknownKeys ...
func SchemaIgnorUnknownKeys(ignore bool) {
	_schema.IgnoreUnknownKeys(ignore)
}

// SchemaZeroEmpty ...
func SchemaZeroEmpty(empty bool) {
	_schema.ZeroEmpty(empty)
}
