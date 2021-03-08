package converters

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/kaspanet/kaspad/domain/consensus/utils/hashes"

	"github.com/kaspanet/kaspad/domain/consensus/model/externalapi"
)

func jsonMarshal(output interface{}) (string, error) {
	bytes, err := json.MarshalIndent(output, "", "\t")
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func RenderOutput(output interface{}) (string, error) {
	// check for nil or interface nil:
	// https://stackoverflow.com/a/50487104/474270
	if output == nil || (reflect.ValueOf(output).Kind() == reflect.Ptr && reflect.ValueOf(output).IsNil()) {
		return "nil", nil
	}
	switch outputObj := output.(type) {
	case *externalapi.DomainHash:
		return outputObj.String(), nil
	case externalapi.BlockHeader:
		return renderBlockHeader(outputObj)
	case *externalapi.DomainBlock:
		return renderBlock(outputObj)
	case error:
		return fmt.Sprintf("%+v", outputObj), nil
	default:
		return jsonMarshal(output)
	}
}

func renderBlock(block *externalapi.DomainBlock) (string, error) {
	jsonable := jsonableBlock(block)

	return jsonMarshal(jsonable)
}

func jsonableBlock(block *externalapi.DomainBlock) interface{} {
	return struct {
		Header       interface{}
		Transactions []*externalapi.DomainTransaction
	}{
		Header:       jsonableBlockHeader(block.Header),
		Transactions: block.Transactions,
	}
}

func renderBlockHeader(blockHeader externalapi.BlockHeader) (string, error) {
	jsonable := jsonableBlockHeader(blockHeader)

	return jsonMarshal(jsonable)
}

func jsonableBlockHeader(blockHeader externalapi.BlockHeader) interface{} {
	return struct {
		Version              uint16
		ParentHashes         []string
		HashMerkleRoot       string
		AcceptedIDMerkleRoot string
		UTXOCommitment       string
		TimeInMilliseconds   int64
		Bits                 uint32
		Nonce                uint64
	}{
		Version:              blockHeader.Version(),
		ParentHashes:         hashes.ToStrings(blockHeader.ParentHashes()),
		HashMerkleRoot:       blockHeader.HashMerkleRoot().String(),
		AcceptedIDMerkleRoot: blockHeader.AcceptedIDMerkleRoot().String(),
		UTXOCommitment:       blockHeader.UTXOCommitment().String(),
		TimeInMilliseconds:   blockHeader.TimeInMilliseconds(),
		Bits:                 blockHeader.Bits(),
		Nonce:                blockHeader.Nonce(),
	}
}
