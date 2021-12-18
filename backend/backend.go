package backend

import (
	"github.com/kaspanet/kaspad/domain/consensus"
	"github.com/kaspanet/kaspad/domain/consensus/model/testapi"
	"github.com/kaspanet/kaspad/domain/dagconfig"
	"github.com/svarogg/dedagger/model"
)

type Backend struct {
	consensus testapi.TestConsensus
	Stores    map[string]*model.Store
}

func NewBackend(dataDir string, dagParams *dagconfig.Params) (backend *Backend, teardown func(), err error) {
	tc, teardown, err := initConsensus(dataDir, dagParams)
	if err != nil {
		return nil, nil, err
	}

	stores, err := extractStores(tc)
	if err != nil {
		teardown()
		return nil, nil, err
	}
	return &Backend{
		consensus: tc,
		Stores:    stores,
	}, teardown, nil
}

func initConsensus(dataDir string, dagParams *dagconfig.Params) (tc testapi.TestConsensus, teardown func(), err error) {
	factory := consensus.NewFactory()
	factory.SetTestDataDir(dataDir)
	config := consensus.Config{*dagParams, false, false, false}
	tc, tcTeardown, err := factory.NewTestConsensus(&config, "dedagger")
	if err != nil {
		return nil, nil, err
	}
	return tc, func() { tcTeardown(true) }, nil
}
