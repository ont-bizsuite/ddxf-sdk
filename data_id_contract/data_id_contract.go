package data_id_contract

import (
	"github.com/ontio/ddxf-sdk/base_contract"
	"github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology/common"
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

func (this *DataIdContractKit) SetContractAddress(dataId common.Address) {
	this.contractAddress = dataId
}

func (this *DataIdContractKit) RegisterDataIdInfo(info DataIdInfo,
	seller *ontology_go_sdk.Account) (common.Uint256, error) {
	return this.bc.Invoke(this.contractAddress, seller, "registerDataId",
		[]interface{}{info.ToBytes()})
}

func (this *DataIdContractKit) GetDataIdInfo(dataId string) (*DataIdInfo, error) {
	res, err := this.bc.PreInvoke(this.contractAddress, "getDataIdInfo",
		[]interface{}{dataId})
	if err != nil {
		return nil, err
	}
	data, err := res.ToByteArray()
	if err != nil {
		return nil, err
	}
	info := &DataIdInfo{}
	err = info.FromBytes(data)
	if err != nil {
		return nil, err
	}
	return info, nil
}
