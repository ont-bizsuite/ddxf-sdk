package main

import (
	"encoding/hex"
	"fmt"
	"github.com/ont-bizsuite/ddxf-sdk"
	"github.com/ont-bizsuite/ddxf-sdk/example/base"
	"github.com/ont-bizsuite/ddxf-sdk/market_place_contract"
	"github.com/ont-bizsuite/ddxf-sdk/split_policy_contract"
	"github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/core/types"
	"github.com/ontio/ontology/core/utils"
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
	payer         *ontology_go_sdk.Account
	gasPrice      = uint64(500)
	tokenTemplate *market_place_contract.TokenTemplate
)

//3f2c66242810aacc4d033758c03f182fbf31df84  split
func main() {
	testNet := "http://106.75.224.136:20336"
	testNet = ddxf_sdk.TestNet
	//testNet = "http://113.31.112.154:20336"
	//testNet = ddxf_sdk.MainNet
	sdk := ddxf_sdk.NewDdxfSdk(testNet)
	//106.75.224.136
	wasmFile := "/Users/sss/dev/dockerData/rust_project/ddxf_market/output/marketplace.wasm"
	wasmFile = "/Users/sss/dev/dockerData/rust_project/ddxf_market/output/dtoken.wasm"
	wasmFile = "/Users/sss/dev/dockerData/rust_project/ddxf_market/output/data_id.wasm"
	//wasmFile = "/Users/sss/dev/dockerData/rust_project/ddxf_market/output/split_policy.wasm"
	//wasmFile = "/Users/sss/dev/rust_project/oep4-rust/output/oep_4.wasm"
	//wasmFile = "/Users/sss/dev/dockerData/rust_project/ddxf_market/output/open_kg.wasm"
	//wasmFile = "/Users/sss/dev/dockerData/rust_project/ddxf_market/output/accountant.wasm"
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
	payer, _ = wallet.GetAccountByAddress("AQCQ3Krh6qxeWKKRACNehA8kAATHxoQNWJ", pwd)

	if true {

		bs, err := sdk.GetOntologySdk().Native.OntId.GetDocumentJson("did:ont:TXvDhLqrqvAV6XUAmLEfWLjxmS1ESxbZBr")

		fmt.Println(string(bs))
		return
		wallet, err := sdk.GetOntologySdk().OpenWallet("./wallet.dat")
		if err != nil {
			fmt.Println(err)
			return
		}
		iden, err := wallet.NewDefaultSettingIdentity(pwd)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("iden:", iden.ID)
		txhash, err := sdk.GetOntologySdk().Native.OntId.RegIDWithPublicKey(500, 2000000, seller, iden.ID, seller)
		if err != nil {
			fmt.Println(err)
			return
		}
		evt, err := sdk.GetSmartCodeEvent(txhash.ToHexString())
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("RegIDWithPublicKey evt:", evt)

		//att := []*DDOAttribute{
		//	&DDOAttribute{
		//		Key:       []byte("key"),
		//		Value:     []byte("value"),
		//		ValueType: []byte{},
		//	},
		//}
		contractAddr, _ := common.AddressFromHexString("df04263aa6ff06bdaf6ba50d29c4cb2a188078cd")
		con := sdk.DefContract(contractAddr)

		iden2, err := wallet.NewDefaultSettingIdentity(pwd)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("iden2:", iden2.ID)
		rp := base.RegIdParam{
			Ontid: []byte(iden2.ID),
			Group: base.Group{
				Members:   [][]byte{[]byte(iden.ID)},
				Threshold: 1,
			},
			Signer: []base.Signer{
				base.Signer{
					Id:    []byte(iden.ID),
					Index: uint32(1),
				},
			},
			Attributes: []base.DDOAttribute{
				base.DDOAttribute{
					Key:       []byte("key"),
					Value:     []byte("value"),
					ValueType: []byte("ty"),
				},
			},
		}
		sink := common.NewZeroCopySink(nil)
		rp.Serialize(sink)

		bs, err = utils.BuildWasmContractParam([]interface{}{"reg_id_add_attribute_array", []interface{}{sink.Bytes()}})
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(hex.EncodeToString(bs))

		txhash, err = con.Invoke("reg_id_add_attribute_array", seller,
			[]interface{}{[]interface{}{sink.Bytes()}})

		if err != nil {
			fmt.Println(err)
			return
		}

		evt, err = sdk.GetSmartCodeEvent(txhash.ToHexString())
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(evt)
		return

		//tx, err := sdk.GetOntologySdk().Native.OntId.NewAddAttributesTransaction(500, 200000, iden.ID, att, seller.PublicKey)
		if err != nil {
			fmt.Println(err)
			return
		}
		var tx *types.MutableTransaction
		sdk.GetOntologySdk().SignToTransaction(tx, seller)
		txhash, err = sdk.GetOntologySdk().SendTransaction(tx)
		if err != nil {
			fmt.Println(err)
			return
		}
		evt, err = sdk.GetSmartCodeEvent(txhash.ToHexString())
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("AddAttributes evt:", evt)

		return

		data, er := sdk.GetOntologySdk().Native.OntId.GetDocumentJson("did:ont:AVFKrE54v1uSrB2c3uxkkcB4KnPpYm7Au6")
		if er != nil {
			fmt.Println(er)
			return
		}
		fmt.Println("data:", string(data))
		evt, _ = sdk.GetSmartCodeEvent("4096bc1c8d7337cb1527d4e959bda3cd500976cfcaf3344cea4055446bb9de8a")
		fmt.Println(evt)
		return
	}

	codeHex := common.ToHexString(code)
	contractAddr := common.AddressFromVmCode(code)
	fmt.Printf("contractAddr:%s, contractAddr:%s\n", contractAddr.ToBase58(), contractAddr.ToHexString())
	//oep init
	if false {
		contract := sdk.DefContract(contractAddr)
		txHash, err := contract.Invoke("init", seller, []interface{}{})
		if err != nil {
			fmt.Println("err:", err)
			return
		}
		evt, err := sdk.GetSmartCodeEvent(txHash.ToHexString())
		if err != nil {
			fmt.Println("err:", err)
			return
		}
		fmt.Println(evt)
		return
	}

	//return
	if true {
		deployContract(sdk, seller, codeHex)
		return
	}
	if false {
		kit := sdk.DefContract(contractAddr)
		txHash, err := kit.Invoke("init", seller, []interface{}{})
		if err != nil {
			fmt.Println("err", err)
			return
		}
		time.Sleep(6 * time.Second)
		evt, err := sdk.GetSmartCodeEvent(txHash.ToHexString())
		if err != nil {
			fmt.Println("err", err)
			return
		}
		fmt.Println("evt:", evt)
		return
	}

	sdk.SetMpContractAddress(contractAddr)
	if false {
		dtoken, _ := common.AddressFromHexString("466b94488bf2ad1b1eec0ae7e49e40708e71a35d")
		split, _ := common.AddressFromHexString("3f2c66242810aacc4d033758c03f182fbf31df84")
		sdk.SetGasPrice(0)
		txHash, err := sdk.DefMpKit().Init(seller, dtoken, split)
		if err != nil {
			fmt.Println("Init failed: ", err)
			return
		}
		showNotify(sdk, "init", txHash.ToHexString())
		return
	}

	if true {
		sdk.SetGasPrice(500)
		contractAddr, _ := common.AddressFromHexString("e01d500ed0c1719b7750367ae59b4b2d308d1ceb")
		txHash, err := sdk.DefDTokenKit().SetMpContractAddr(seller, contractAddr)
		if err != nil {
			fmt.Println(err)
			return
		}
		evt, err := sdk.GetSmartCodeEvent(txHash.ToHexString())
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(evt)
		return
	}

	if true {
		resourceIdBytes := []byte(strconv.Itoa(rand.Int()))
		dataId := ""
		tokenTemplate = &market_place_contract.TokenTemplate{
			DataID:     dataId,
			TokenHashs: []string{string(common.UINT256_EMPTY[:])},
		}

		if err = publish(sdk, resourceIdBytes); err != nil {
			fmt.Println("publish error: ", err)
			return
		}
		return
		if err := delete(sdk, resourceIdBytes); err != nil {
			fmt.Println("delete error: ", err)
			return
		}

		if err := update(sdk, resourceIdBytes); err != nil {
			fmt.Println("update error: ", err)
			return
		}

		//if err = buyAndUseToken(sdk, resourceIdBytes); err != nil {
		//	fmt.Println("buyAndUseToken error: ", err)
		//	return
		//}

		if err = buyDtoken(sdk, resourceIdBytes); err != nil {
			fmt.Println("buyDtoken error: ", err)
			return
		}

		if err = addAgents(sdk); err != nil {
			fmt.Println("addAgents error: ", err)
			return
		}

		if err = useTokenByAgent(sdk); err != nil {
			fmt.Println("useTokenByAgent error: ", err)
			return
		}

		if err = removeAgents(sdk); err != nil {
			fmt.Println("removeAgents error: ", err)
			return
		}

		if err = addTokenAgents(sdk); err != nil {
			fmt.Println("addTokenAgents error: ", err)
			return
		}
		if err = removeTokenAgents(sdk); err != nil {
			fmt.Println("removeTokenAgents error: ", err)
			return
		}

		err = useToken(sdk)
		if err != nil {
			fmt.Printf("useToken: %s\n", err)
			return
		}
	}
}

func addTokenAgents(sdk *ddxf_sdk.DdxfSdk) error {
	txHash, err := sdk.DefDTokenKit().AddTokenAgents(buyer,
		[]common.Address{agent.Address}, *tokenTemplate, 1)
	if err != nil {
		return err
	}
	return showNotify(sdk, "addTokenAgents", txHash.ToHexString())
}

func removeTokenAgents(sdk *ddxf_sdk.DdxfSdk) error {
	txHash, err := sdk.DefDTokenKit().RemoveTokenAgents(*tokenTemplate, buyer,
		[]common.Address{agent.Address})
	if err != nil {
		return err
	}
	return showNotify(sdk, "removeTokenAgents", txHash.ToHexString())
}

func removeAgents(sdk *ddxf_sdk.DdxfSdk) error {
	txHash, err := sdk.DefDTokenKit().RemoveAgents(buyer, []common.Address{agent.Address}, []market_place_contract.TokenTemplate{*tokenTemplate})
	if err != nil {
		return err
	}
	return showNotify(sdk, "removeAgents", txHash.ToHexString())
}

func useTokenByAgent(sdk *ddxf_sdk.DdxfSdk) error {
	txHash, err := sdk.DefDTokenKit().UseTokenByAgents(buyer.Address, agent, *tokenTemplate, 1)
	if err != nil {
		return err
	}
	return showNotify(sdk, "UseTokenByAgents", txHash.ToHexString())
}

func addAgents(sdk *ddxf_sdk.DdxfSdk) error {
	txHash, err := sdk.DefDTokenKit().AddAgents(buyer,
		[]common.Address{agent.Address}, 1, []market_place_contract.TokenTemplate{*tokenTemplate})
	if err != nil {
		fmt.Printf("AddAgents: %s\n", err)
		return err
	}
	return showNotify(sdk, "addAgents", txHash.ToHexString())
}
func useToken(sdk *ddxf_sdk.DdxfSdk) error {
	txHash, err := sdk.DefDTokenKit().UseToken(buyer, *tokenTemplate, 1)
	if err != nil {
		return err
	}
	return showNotify(sdk, "useToken", txHash.ToHexString())
}

func buyAndUseToken(sdk *ddxf_sdk.DdxfSdk, resourceIdBytes []byte) error {
	txHash, err := sdk.DefMpKit().BuyAndUseToken(buyer, payer, resourceIdBytes, 2, *tokenTemplate)
	if err != nil {
		return err
	}
	return showNotify(sdk, "buyAndUseToken", txHash.ToHexString())
}

func buyDtoken(sdk *ddxf_sdk.DdxfSdk, resourceIdBytes []byte) error {
	txHash, err := sdk.DefMpKit().BuyDtoken(buyer, payer, resourceIdBytes, 2)
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

func delete(sdk *ddxf_sdk.DdxfSdk, resourceIdBytes []byte) error {
	txHash, err := sdk.DefMpKit().Delete(seller, resourceIdBytes)
	if err != nil {
		fmt.Println(err)
		return err
	}
	evt, err := sdk.GetSmartCodeEvent(txHash.ToHexString())
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(evt)
	return nil
}
func update(sdk *ddxf_sdk.DdxfSdk, resourceIdBytes []byte) error {
	itemMeta := map[string]interface{}{
		"key": "value",
	}
	bs, err := ddxf.HashObject(itemMeta)
	if err != nil {
		return err
	}
	itemMetaHash, err := common.Uint256ParseFromBytes(bs[:])

	ddo := market_place_contract.ResourceDDO{
		Manager:      seller.Address,       // data owner id
		ItemMetaHash: itemMetaHash,         // required if len(Templates) > 1
		DTC:          common.ADDRESS_EMPTY, // can be empty
		MP:           common.ADDRESS_EMPTY, // can be empty
		Split:        common.ADDRESS_EMPTY,
	}

	item := market_place_contract.DTokenItem{
		Fee: market_place_contract.Fee{
			ContractType: 0,
			Count:        1,
		},
		ExpiredDate: uint64(time.Now().Unix()) + 10000,
		Stocks:      10000,
		Templates:   []*market_place_contract.TokenTemplate{tokenTemplate},
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

	txHash, err := sdk.DefMpKit().Update(seller, resourceIdBytes, ddo, item, sp)
	if err != nil {
		fmt.Printf("Publish error:%s\n", err)
		return err
	}
	fmt.Println("publish txHash: ", txHash.ToHexString())

	return showNotify(sdk, "update", txHash.ToHexString())
}

func publish(sdk *ddxf_sdk.DdxfSdk, resourceIdBytes []byte) error {
	itemMeta := map[string]interface{}{
		"key": "value",
	}
	bs, err := ddxf.HashObject(itemMeta)
	if err != nil {
		return err
	}
	itemMetaHash, err := common.Uint256ParseFromBytes(bs[:])

	ddo := market_place_contract.ResourceDDO{
		Manager:      seller.Address,       // data owner id
		ItemMetaHash: itemMetaHash,         // required if len(Templates) > 1
		DTC:          common.ADDRESS_EMPTY, // can be empty
		MP:           common.ADDRESS_EMPTY, // can be empty
		Split:        common.ADDRESS_EMPTY,
	}

	item := market_place_contract.DTokenItem{
		Fee: market_place_contract.Fee{
			ContractType: 0,
			Count:        1,
		},
		ExpiredDate: uint64(time.Now().Unix()) + 10000,
		Stocks:      10000,
		Templates:   []*market_place_contract.TokenTemplate{tokenTemplate},
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

	txHash, err := sdk.DefMpKit().Publish(seller, resourceIdBytes, ddo, item, sp)
	if err != nil {
		fmt.Printf("Publish error:%s\n", err)
		return err
	}
	fmt.Println("publish txHash: ", txHash.ToHexString())

	return showNotify(sdk, "publish", txHash.ToHexString())
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
