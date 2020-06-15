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
	"math/rand"
	"strconv"
	"time"
)

var (
	admin         *ontology_go_sdk.Account
	seller        *ontology_go_sdk.Account
	buyer         *ontology_go_sdk.Account
	agent         *ontology_go_sdk.Account
	gasPrice      = uint64(0)
	tokenTemplate *ddxf_contract.TokenTemplate
)

func main() {
	sdk := ddxf_sdk.NewDdxfSdk(ddxf_sdk.LocalNet)
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
	admin, _ = wallet.GetAccountByAddress("AYnhakv7kC9R5ppw65JoE2rt6xDzCjCTvD", pwd)
	seller, _ = wallet.GetAccountByAddress("Aejfo7ZX5PVpenRj23yChnyH64nf8T1zbu", pwd)
	buyer, _ = wallet.GetAccountByAddress("AHhXa11suUgVLX1ZDFErqBd3gskKqLfa5N", pwd)
	agent, _ = wallet.GetAccountByAddress("ANb3bf1b67WP2ZPh5HQt4rkrmphMJmMCMK", pwd)

	codeHex := common.ToHexString(code)
	contractAddr := common.AddressFromVmCode(code)
	fmt.Printf("contractAddr:%s, contractAddr:%s\n", contractAddr.ToBase58(), contractAddr.ToHexString())
	if false {
		deployContract(sdk, admin, codeHex)
		return
	}
	sdk.SetDDXFContractAddress(contractAddr)
	if false {
		dtoken, _ := common.AddressFromHexString("49c2dc97ee58b2292e55499e1122c579fc0690e3")
		split, _ := common.AddressFromHexString("d1c175f1e485c4af5d3fb8906e850a5d569c4c48")
		txHash, err := sdk.DefDDXFKit().Init(admin, dtoken, split)
		if err != nil {
			fmt.Println("Init failed: ", err)
			return
		}
		showNotify(sdk, "init", txHash.ToHexString())
		return
	}
	if true {
		resourceIdBytes := []byte(strconv.Itoa(rand.Int()))
		publish(sdk, resourceIdBytes)
		buyDtoken(sdk, resourceIdBytes)

		if err = addAgents(sdk, resourceIdBytes); err != nil {
			fmt.Println("addAgents error: ", err)
			return
		}
		if err = useTokenByAgent(sdk, resourceIdBytes); err != nil {
			fmt.Println("useTokenByAgent error: ", err)
			return
		}
		if err = removeAgents(sdk, resourceIdBytes); err != nil {
			fmt.Println("removeAgents error: ", err)
			return
		}
		if err = addTokenAgents(sdk, resourceIdBytes); err != nil {
			fmt.Println("addTokenAgents error: ", err)
			return
		}
		if err = removeTokenAgents(sdk, resourceIdBytes); err != nil {
			fmt.Println("removeTokenAgents error: ", err)
			return
		}

		err = useToken(sdk, resourceIdBytes)
		if err != nil {
			fmt.Println("useToken: %s", err)
			return
		}
	}
}

func addTokenAgents(sdk *ddxf_sdk.DdxfSdk, resourceIdBytes []byte) error {
	txHash, err := sdk.DefDDXFKit().AddTokenAgents(resourceIdBytes, buyer,
		[]common.Address{agent.Address}, *tokenTemplate, 1)
	if err != nil {
		return err
	}
	return showNotify(sdk, "addTokenAgents", txHash.ToHexString())
}

func removeTokenAgents(sdk *ddxf_sdk.DdxfSdk, resourceIdBytes []byte) error {
	txHash, err := sdk.DefDDXFKit().RemoveTokenAgents(resourceIdBytes, *tokenTemplate, buyer,
		[]common.Address{agent.Address})
	if err != nil {
		return err
	}
	return showNotify(sdk, "removeTokenAgents", txHash.ToHexString())
}

func removeAgents(sdk *ddxf_sdk.DdxfSdk, resourceIdBytes []byte) error {
	txHash, err := sdk.DefDDXFKit().RemoveAgents(resourceIdBytes, buyer, []common.Address{agent.Address})
	if err != nil {
		return err
	}
	return showNotify(sdk, "removeAgents", txHash.ToHexString())
}

func useTokenByAgent(sdk *ddxf_sdk.DdxfSdk, resourceIdBytes []byte) error {
	txHash, err := sdk.DefDDXFKit().UseTokenByAgents(resourceIdBytes, buyer.Address, agent, *tokenTemplate, 1)
	if err != nil {
		return err
	}
	return showNotify(sdk, "UseTokenByAgents", txHash.ToHexString())
}

func addAgents(sdk *ddxf_sdk.DdxfSdk, resourceIdBytes []byte) error {
	txHash, err := sdk.DefDDXFKit().AddAgents(resourceIdBytes, buyer,
		[]common.Address{agent.Address}, 1)
	if err != nil {
		fmt.Println("AddAgents: %s", err)
		return err
	}
	return showNotify(sdk, "addAgents", txHash.ToHexString())
}
func useToken(sdk *ddxf_sdk.DdxfSdk, resourceIdBytes []byte) error {
	txHash, err := sdk.DefDDXFKit().UseToken(resourceIdBytes, buyer, *tokenTemplate, 1)
	if err != nil {
		return err
	}
	return showNotify(sdk, "useToken", txHash.ToHexString())
}
func buyDtoken(sdk *ddxf_sdk.DdxfSdk, resourceIdBytes []byte) error {
	txHash, err := sdk.DefDDXFKit().BuyDtoken(buyer, resourceIdBytes, 2)
	if err != nil {
		return err
	}
	return showNotify(sdk, "buyDtoken", txHash.ToHexString())
}

func showNotify(sdk *ddxf_sdk.DdxfSdk, method, txHash string) error {
	fmt.Printf("method: %s, txHash: %s\n", method, txHash)
	evt, err := sdk.GetSmartCodeEvent(txHash)
	if err != nil {
		return err
	}
	for _, notify := range evt.Notify {
		fmt.Printf("method: %s,evt: %v\n", method, notify)
	}
	return nil
}

func publish(sdk *ddxf_sdk.DdxfSdk, resourceIdBytes []byte) {
	dataId := ""
	tokenTemplate = &ddxf_contract.TokenTemplate{
		DataID:     dataId,
		TokenHashs: []string{string(common.UINT256_EMPTY[:])},
	}
	trt := &ddxf_contract.TokenResourceTyEndpoint{
		TokenTemplate: tokenTemplate,
		ResourceType:  0,
		Endpoint:      "",
	}
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
		AddrAmts: []*split_policy_contract.AddrAmt{
			&split_policy_contract.AddrAmt{
				To:          seller.Address,
				Percent:     5000,
				HasWithdraw: false,
			},
			&split_policy_contract.AddrAmt{
				To:          admin.Address,
				Percent:     5000,
				HasWithdraw: false,
			},
		},
		TokenTy: split_policy_contract.ONG,
	}

	txHash, err := sdk.DefDDXFKit().Publish(seller, resourceIdBytes, ddo, item, sp)
	if err != nil {
		fmt.Printf("Publish error:%s\n", err)
		return
	}
	fmt.Println("publish txHash: ", txHash.ToHexString())
	evt, err := sdk.GetSmartCodeEvent(txHash.ToHexString())
	if err != nil {
		fmt.Printf("Publish error:%s\n", err)
		return
	}
	fmt.Println("publish evt:", evt)
}

func deployContract(sdk *ddxf_sdk.DdxfSdk, admin *ontology_go_sdk.Account, codeHex string) {
	sdk.SetGasPrice(gasPrice)
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
