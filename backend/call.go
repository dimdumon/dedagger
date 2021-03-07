package backend

import (
	"reflect"

	"github.com/svarogg/dedagger/model"
)

func (be *Backend) Call(method *model.Method, parameters []reflect.Value) []reflect.Value {
	valueOfDatabaseContext := reflect.ValueOf(be.consensus.DatabaseContext())
	in := append([]reflect.Value{valueOfDatabaseContext}, parameters...)
	return method.Value.Call(in)
}
