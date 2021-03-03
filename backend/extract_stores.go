package backend

import (
	"reflect"

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

	return &model.Store{
		Name:    typeof.String(),
		Value:   value,
		Typeof:  typeof,
		Methods: methods,
	}, nil
}

func extractMethods(value reflect.Value, typeof reflect.Type) ([]*model.Method, error) {
	numMethods := typeof.NumMethod()
	methods := make([]*model.Method, numMethods)

	for i := 0; i < numMethods; i++ {
		reflectMethod := typeof.Method(i)
		methodValue := value.MethodByName(reflectMethod.Name)

		parameters, err := extractParameters(reflectMethod)
		if err != nil {
			return nil, err
		}

		methods[i] = &model.Method{
			Name:       reflectMethod.Name,
			Value:      methodValue,
			Parameters: parameters,
		}
	}

	return methods, nil
}

func extractParameters(method reflect.Method) ([]*model.Parameter, error) {
	numIn := method.Type.NumIn()
	parameters := make([]*model.Parameter, numIn-1) // `- 1` to skip the receiver

	for i := 1; i < numIn; i++ {
		parameter := &model.Parameter{
			Type: method.Type.In(i),
		}

		parameters[i-1] = parameter
	}
	return parameters, nil
}
