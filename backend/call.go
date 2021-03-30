package backend

import (
	"reflect"

	consensusmodel "github.com/kaspanet/kaspad/domain/consensus/model"

	"github.com/svarogg/dedagger/model"
)

func (be *Backend) Call(method *model.Method, parameters []reflect.Value) []reflect.Value {
	in := []reflect.Value{}
	stagingAreaIndex := 0
	if method.Value.Type().NumIn() > 0 && method.Value.Type().In(0).String() == "model.DBReader" {
		valueOfDatabaseContext := reflect.ValueOf(be.consensus.DatabaseContext())
		in = append(in, valueOfDatabaseContext)

		stagingAreaIndex++
	}

	if method.Value.Type().NumIn() > stagingAreaIndex && method.Value.Type().In(stagingAreaIndex).String() == "*model.StagingArea" {
		stagingArea := consensusmodel.NewStagingArea()
		valueOfStagingArea := reflect.ValueOf(stagingArea)
		in = append(in, valueOfStagingArea)
	}
	in = append(in, parameters...)
	return method.Value.Call(in)
}
