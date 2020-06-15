package ddxf_sdk

import (
	"fmt"
	"github.com/ontio/ddxf-sdk/base_contract"
	"github.com/ontio/ddxf-sdk/data_id_contract"
	"github.com/ontio/ddxf-sdk/ddxf_contract"
	"github.com/ontio/ontology-go-sdk"
	common2 "github.com/ontio/ontology-go-sdk/common"
	"github.com/ontio/ontology/common"
	"time"
)

const (
	dDXFContractAddress   = "90982cd1d33ec7b33bffe54b289f5acaf02815a8"
	dataIdContractAddress = "e854316627dfc44bef9c0eb583e941804d0716d5"
)

const (
	MainNet  = "http://dappnode1.ont.io:20336"
	TestNet  = "http://polaris1.ont.io:20336"
	LocalNet = "http://127.0.0.1:20336"
)

const (
	defaultGasPrice = 500
	defaultGasLimit = 31200000
)

type DdxfSdk struct {
	sdk                   *ontology_go_sdk.OntologySdk
	bc                    *base_contract.BaseContract
	rpc                   string
	defaultDdxfContract   *ddxf_contract.DDXFKit
	defaultDataIdContract *data_id_contract.DataIdKit
	gasPrice              uint64
	gasLimit              uint64
}

func NewDdxfSdk(addr string) *DdxfSdk {
	sdk := ontology_go_sdk.NewOntologySdk()
	sdk.NewRpcClient().SetAddress(addr)
	return &DdxfSdk{
		sdk:      sdk,
		rpc:      addr,
		gasPrice: defaultGasPrice,
		gasLimit: defaultGasLimit,
		bc:       base_contract.NewBaseContract(sdk, defaultGasLimit, defaultGasPrice, nil),
	}
}

func (sdk *DdxfSdk) SetPayer(payer *ontology_go_sdk.Account) {
	sdk.bc.SetPayer(payer)
}

func (sdk *DdxfSdk) SetGasLimit(gasLimit uint64) {
	sdk.gasLimit = gasLimit
	sdk.bc.SetGasLimit(gasLimit)
}

func (sdk *DdxfSdk) SetGasPrice(gasPrice uint64) {
	sdk.gasPrice = gasPrice
	sdk.bc.SetGasPrice(gasPrice)
}

func (sdk *DdxfSdk) GetOntologySdk() *ontology_go_sdk.OntologySdk {
	return sdk.sdk
}

func (sdk *DdxfSdk) DefaultDataIdContract() *data_id_contract.DataIdKit {
	if sdk.defaultDataIdContract == nil {
		contractAddress, _ := common.AddressFromHexString(dataIdContractAddress)
		sdk.defaultDataIdContract = data_id_contract.NewDataIdContractKit(contractAddress,
			sdk.bc)
	}
	return sdk.defaultDataIdContract
}

func (sdk *DdxfSdk) DefDDXFKit() *ddxf_contract.DDXFKit {
	if sdk.defaultDdxfContract == nil {
		contractAddress, _ := common.AddressFromHexString(dDXFContractAddress)
		sdk.defaultDdxfContract = ddxf_contract.NewDDXFContractKit(contractAddress,
			sdk.bc)
	}
	return sdk.defaultDdxfContract
}

func (sdk *DdxfSdk) SetDDXFContractAddress(ddxf common.Address) {
	sdk.DefDDXFKit().SetContractAddress(ddxf)
}

func (sdk *DdxfSdk) GetSmartCodeEvent(txHash string) (*common2.SmartContactEvent, error) {
	for i := 0; i < 10; i++ {
		event, err := sdk.sdk.GetSmartContractEvent(txHash)
		if event != nil {
			return event, err
		}
		if err != nil {
			return nil, err
		}
		if event == nil {
			time.Sleep(3 * time.Second)
		}
	}
	return nil, fmt.Errorf("GetSmartCodeEvent timeout, txhash: %s", txHash)
}

func (this *DdxfSdk) DeployContract(signer *ontology_go_sdk.Account, code,
	name, version, author, email, desc string) (common.Uint256, error) {
	return this.sdk.WasmVM.DeployWasmVMSmartContract(this.gasPrice,
		this.gasLimit, signer, code, name,
		version, author, email, desc)
}
