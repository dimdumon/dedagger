package converters

import (
	"encoding/json"
	"reflect"

	"github.com/kaspanet/kaspad/domain/consensus/model/externalapi"
	"github.com/svarogg/dedagger/model"
)

func ParseParameter(parameter *model.Parameter, valueString string) (reflect.Value, error) {
	switch parameter.Type.String() {
	case "*externalapi.DomainHash":
		hash, err := externalapi.NewDomainHashFromString(valueString)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(hash), nil
	default:
		valueInterface := reflect.New(parameter.Type).Interface()
		err := json.Unmarshal([]byte(valueString), valueInterface)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(valueInterface), nil
	}
}
