package data_id_contract

import (
	"github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology/common"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/seller/server"
	"github.com/ontio/ddxf-sdk/base_contract"
)

type DataIdContractKit struct {
	bc              *base_contract.BaseContract
	contractAddress common.Address
}

func NewDataIdContractKit(
	contractAddress common.Address,
	bc *base_contract.BaseContract) *DataIdContractKit {
	return &DataIdContractKit{
		contractAddress: contractAddress,
		bc:              bc,
	}
}

func (this *DataIdContractKit) RegisterDataId(info server.DataIdInfo,
	seller *ontology_go_sdk.Account) (common.Uint256, error) {
	return this.bc.Invoke(this.contractAddress, seller, "registerDataId",
		[]interface{}{info.ToBytes()})
}
