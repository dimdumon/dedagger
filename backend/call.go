package backend

import (
	"reflect"

	"github.com/svarogg/dedagger/model"
)

func (be *Backend) Call(method *model.Method, parameters []reflect.Value) []reflect.Value {
	in := []reflect.Value{}
	if method.Value.Type().NumIn() > 0 && method.Value.Type().In(0).String() == "model.DBReader" {
		valueOfDatabaseContext := reflect.ValueOf(be.consensus.DatabaseContext())
		in = append(in, valueOfDatabaseContext)
	}
	in = append(in, parameters...)
	return method.Value.Call(in)
}
