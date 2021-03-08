package converters

import (
	"encoding/json"

	"github.com/kaspanet/kaspad/domain/consensus/model/externalapi"
)

func RenderOutput(output interface{}) (string, error) {
	switch outputObj := output.(type) {
	case *externalapi.DomainHash:
		return outputObj.String(), nil
	default:
		bytes, err := json.MarshalIndent(output, "", "\t")
		if err != nil {
			return "", err
		}
		return string(bytes), nil
	}
}
