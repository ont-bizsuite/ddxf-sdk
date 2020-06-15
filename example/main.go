package main

import (
	"fmt"
	"github.com/ontio/ddxf-sdk"
	"github.com/ontio/ddxf-sdk/ddxf_contract"
	"github.com/ontio/ddxf-sdk/split_policy_contract"
	"github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology/common"
	"github.com/zhiqiangxu/ddxf"
	"io/ioutil"
	"time"
)

func main() {
	sdk := ddxf_sdk.NewDdxfSdk(ddxf_sdk.TestNet)
	wasmFile := "/Users/sss/dev/dockerData/rust_project/ddxf_market/output/ddxf.wasm"
	//wasmFile = "/Users/sss/dev/dockerData/rust_project/ddxf_market/output/dtoken.wasm"
	code, err := ioutil.ReadFile(wasmFile)
	if err != nil {
		fmt.Printf("error in ReadFile:%s\n", err)
		return
	}
	pwd := []byte("123456")
	ontSdk := sdk.GetOntologySdk()
	wallet, err := ontSdk.OpenWallet("./wallet.dat")
	if err != nil {
		fmt.Printf("error in ReadFile:%s\n", err)
		return
	}
	admin, _ := wallet.GetAccountByAddress("AYnhakv7kC9R5ppw65JoE2rt6xDzCjCTvD", pwd)
	codeHex := common.ToHexString(code)
	contractAddr := common.AddressFromVmCode(code)
	fmt.Printf("contractAddr:%s, contractAddr:%s\n", contractAddr.ToBase58(), contractAddr.ToHexString())
	if true {
		deployContract(sdk, admin, codeHex)
	}
	sdk.SetDDXFContractAddress(contractAddr)
}


func testDdxf(sdk *ddxf_sdk.DdxfSdk, wallet *ontology_go_sdk.Wallet, pwd []byte) {
	seller, _ := wallet.GetAccountByAddress("Aejfo7ZX5PVpenRj23yChnyH64nf8T1zbu", pwd)
	dataId := ""
	tokenTemplate := &ddxf_contract.TokenTemplate{
		DataID:     dataId,
		TokenHashs: []string{string(common.UINT256_EMPTY[:])},
	}
	trt := &ddxf_contract.TokenResourceTyEndpoint{
		TokenTemplate: tokenTemplate,
		ResourceType:  0,
		Endpoint:      "",
	}
	resourceIdBytes := []byte("")
	itemMeta := map[string]interface{}{
		"key": "value",
	}
	bs, err := ddxf.HashObject(itemMeta)
	if err != nil {
		return
	}
	itemMetaHash, err := common.Uint256ParseFromBytes(bs[:])

	ddo := ddxf_contract.ResourceDDO{
		TokenResourceTyEndpoints: []*ddxf_contract.TokenResourceTyEndpoint{trt}, // RT for tokens
		Manager:                  seller.Address,                                // data owner id
		ItemMetaHash:             itemMetaHash,                                  // required if len(Templates) > 1
		DTC:                      common.ADDRESS_EMPTY,                          // can be empty
		MP:                       common.ADDRESS_EMPTY,                          // can be empty
		Split:                    common.ADDRESS_EMPTY,
	}

	item := ddxf_contract.DTokenItem{
		Fee: ddxf_contract.Fee{
			ContractType: 0,
			Count:        1,
		},
		ExpiredDate: uint64(time.Now().Unix()) + 10000,
		Stocks:      10000,
		Templates:   []*ddxf_contract.TokenTemplate{tokenTemplate},
	}

	sp := split_policy_contract.SplitPolicyRegisterParam{
		AddrAmts: []*split_policy_contract.AddrAmt{},
		TokenTy:  split_policy_contract.ONG,
	}
	txHash, err := sdk.DefaultDDXFContract().Publish(seller, resourceIdBytes, ddo, item, sp)
	if err != nil {
		fmt.Printf("Publish error:%s\n", err)
		return
	}
	sdk.GetSmartCodeEvent(txHash.ToHexString())
}

func deployContract(sdk *ddxf_sdk.DdxfSdk, admin *ontology_go_sdk.Account, codeHex string) {
	txHash, err := sdk.DeployContract(admin, codeHex, "ddxf", "0.1.0", "lucas", "", "")
	if err != nil {
		fmt.Printf("DeployContract error:%s\n", err)
		return
	}
	evt, err := sdk.GetSmartCodeEvent(txHash.ToHexString())
	if err != nil {
		fmt.Printf("DeployContract GetSmartCodeEvent error:%s, txHash: %s\n", err, txHash.ToHexString())
		return
	}
	fmt.Println("evt:", evt)
}
