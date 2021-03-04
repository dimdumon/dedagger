package backend

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/pkg/errors"

	"github.com/svarogg/dedagger/model"

	"github.com/kaspanet/kaspad/domain/consensus/model/testapi"
)

func storeInterfaces(tc testapi.TestConsensus) []interface{} {
	return []interface{}{
		tc.AcceptanceDataStore(),
		tc.BlockHeaderStore(),
		tc.BlockRelationStore(),
		tc.BlockStatusStore(),
		tc.BlockStore(),
		tc.ConsensusStateStore(),
		tc.GHOSTDAGDataStore(),
		tc.HeaderTipsStore(),
		tc.MultisetStore(),
		tc.PruningStore(),
		tc.ReachabilityDataStore(),
		tc.UTXODiffStore(),
		tc.HeadersSelectedChainStore(),
	}
}

func extractStores(tc testapi.TestConsensus) ([]*model.Store, error) {
	storeInterfaces := storeInterfaces(tc)

	stores := make([]*model.Store, len(storeInterfaces))
	for i, storeInterface := range storeInterfaces {
		store, err := newStore(storeInterface)
		if err != nil {
			return nil, err
		}

		stores[i] = store
	}
	return stores, nil
}

func newStore(storeInterface interface{}) (*model.Store, error) {
	value := reflect.ValueOf(storeInterface)
	typeof := value.Type()
	methods, err := extractMethods(value, typeof)
	if err != nil {
		return nil, err
	}

	fullName := typeof.String()
	name := fullName[strings.Index(fullName, ".")+1:] // fullName has prefix of package name - remove it
	return &model.Store{
		Name:    name,
		Value:   value,
		Typeof:  typeof,
		Methods: methods,
	}, nil
}

func extractMethods(value reflect.Value, typeof reflect.Type) ([]*model.Method, error) {
	numMethods := typeof.NumMethod()
	methods := []*model.Method{}

	for i := 0; i < numMethods; i++ {
		reflectMethod := typeof.Method(i)
		if isFilteredMethod(reflectMethod.Name) {
			continue
		}

		methodValue := value.MethodByName(reflectMethod.Name)

		parameters, err := extractParameters(reflectMethod)
		if err != nil {
			return nil, err
		}

		method := &model.Method{
			Name:       reflectMethod.Name,
			Value:      methodValue,
			Parameters: parameters,
		}
		methods = append(methods, method)
	}

	return methods, nil
}

var filteredMethodRegexes = []string{
	"Commit",
	"Delete",
	"Discard",
	".*Stage.*",
	".*Updat.*",
}

func isFilteredMethod(methodName string) bool {
	methodNameBytes := []byte(methodName)

	for _, filteredMethodRegex := range filteredMethodRegexes {
		match, err := regexp.Match(filteredMethodRegex, methodNameBytes)
		if err != nil {
			panic(errors.Errorf("Error matching filteredMethod regex '%s' to '%s': %+v", filteredMethodRegex, methodName, err))
		}
		if match {
			return true
		}
	}

	return false
}

func extractParameters(method reflect.Method) ([]*model.Parameter, error) {
	numIn := method.Type.NumIn()
	parameters := []*model.Parameter{}

	for i := 2; i < numIn; i++ { // Start from 1 to skip the receiver and DatabaseContext
		parameter := &model.Parameter{
			Type: method.Type.In(i),
		}

		parameters = append(parameters, parameter)
	}
	return parameters, nil
}
