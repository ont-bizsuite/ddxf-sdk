package ddxf_sdk

import (
	"github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ddxf-sdk/ddxf_contract"
	"github.com/ontio/ddxf-sdk/data_id_contract"
	"github.com/ontio/ddxf-sdk/base_contract"
)

const (
	dDXFContractAddress   = "90982cd1d33ec7b33bffe54b289f5acaf02815a8"
	dataIdContractAddress = "e854316627dfc44bef9c0eb583e941804d0716d5"
)

const (
	defaultGasPrice = 500
	defaultGasLimit = 20000000
)

type DdxfSdk struct {
	bc                    *base_contract.BaseContract
	rpc                   string
	defaultDdxfContract   *ddxf_contract.DDXFContractKit
	defaultDataIdContract *data_id_contract.DataIdContractKit
}

func NewDdxfSdk(addr string) *DdxfSdk {
	sdk := ontology_go_sdk.NewOntologySdk()
	sdk.NewRpcClient().SetAddress(addr)
	return &DdxfSdk{
		rpc: addr,
		bc:  base_contract.NewBaseContract(sdk, defaultGasLimit, defaultGasPrice, nil),
	}
}

func (sdk *DdxfSdk) SetPayer(payer *ontology_go_sdk.Account) {
	sdk.bc.SetPayer(payer)
}

func (sdk *DdxfSdk) SetGasLimit(gasLimit uint64) {
	sdk.bc.SetGasLimit(gasLimit)
}

func (sdk *DdxfSdk) SetGasPrice(gasPrice uint64) {
	sdk.bc.SetGasPrice(gasPrice)
}

func (sdk *DdxfSdk) GetOntologySdk() *ontology_go_sdk.OntologySdk {
	return sdk.GetOntologySdk()
}

func (sdk *DdxfSdk) DefaultDataIdContract() *data_id_contract.DataIdContractKit {
	if sdk.defaultDataIdContract == nil {
		contractAddress, _ := common.AddressFromHexString(dataIdContractAddress)
		sdk.defaultDataIdContract = data_id_contract.NewDataIdContractKit(contractAddress,
			sdk.bc)
	}
	return sdk.defaultDataIdContract
}

func (sdk *DdxfSdk) DefaultDDXFContract() *ddxf_contract.DDXFContractKit {
	if sdk.defaultDdxfContract == nil {
		contractAddress, _ := common.AddressFromHexString(dDXFContractAddress)
		sdk.defaultDdxfContract = ddxf_contract.NewDDXFContractKit(contractAddress,
			sdk.bc)
	}
	return sdk.defaultDdxfContract
}
