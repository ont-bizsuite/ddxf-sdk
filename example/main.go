package main

import (
	"encoding/hex"
	"fmt"
	"github.com/ont-bizsuite/ddxf-sdk"
	"github.com/ont-bizsuite/ddxf-sdk/example/utils"
	"github.com/ont-bizsuite/ddxf-sdk/market_place_contract"
	"github.com/ont-bizsuite/ddxf-sdk/split_policy_contract"
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
	payer         *ontology_go_sdk.Account
	gasPrice      = uint64(2500)
	tokenTemplate *market_place_contract.TokenTemplate
)

//3f2c66242810aacc4d033758c03f182fbf31df84  split
func main() {
	testNet := "http://106.75.224.136:20336"
	testNet = ddxf_sdk.TestNet
	//testNet = "http://172.168.3.47:20336"
	//testNet = "http://113.31.112.154:20336"
	//testNet = ddxf_sdk.MainNet
	sdk := ddxf_sdk.NewDdxfSdk(testNet)
	//106.75.224.136
	wasmFile := "/Users/sss/dev/dockerData/rust_project/ddxf_market/output/marketplace.wasm"
	//wasmFile = "/Users/sss/dev/dockerData/rust_project/ddxf_market/output/dtoken.wasm"
	//wasmFile = "/Users/sss/dev/dockerData/rust_project/ddxf_market/output/data_id.wasm"
	//wasmFile = "/Users/sss/dev/dockerData/rust_project/ddxf_market/output/split_policy.wasm"
	//wasmFile = "/Users/sss/dev/rust_project/oep4-rust/output/oep_4.wasm"
	wasmFile = "/Users/sss/dev/dockerData/rust_project/ddxf_market/output/open_kg.wasm"
	//wasmFile = "/Users/sss/dev/dockerData/rust_project/ddxf_market/output/accountant.wasm"
	//wasmFile = "/Users/sss/dev/dockerData/rust_project/vote/output/vote.wasm"
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

	codeHex := common.ToHexString(code)
	contractAddr := common.AddressFromVmCode(code)
	fmt.Printf("contractAddr:%s, contractAddr:%s\n", contractAddr.ToBase58(), contractAddr.ToHexString())
	//return

	if false {
		utils.DataIdTest(sdk, pwd, seller, contractAddr)
		return
	}

	if false {
		//dtoken b57362d82c658fb2927e793406ee8decf32ff2e9
		// mp ac820725a3bc5337b6044aad38dc3293032eae99
		// openkg dcd823e05f330a0a838730754bcc2f7e7cf0af57
		deployContract(sdk, seller, codeHex)
		return
	}

	if false {
		sdk.SetMpContractAddress(contractAddr)
		dtoken, _ := common.AddressFromHexString("b57362d82c658fb2927e793406ee8decf32ff2e9")
		split, _ := common.AddressFromHexString("f024034fe7e5ea69c53cede4774bd1dad566234f")
		sdk.SetGasPrice(2500)
		txHash, err := sdk.DefMpKit().Init(seller, dtoken, split)
		if err != nil {
			fmt.Println("Init failed: ", err)
			return
		}
		showNotify(sdk, "init", txHash.ToHexString())
		return
	}
	//openkg
	if true {
		con := sdk.DefContract(contractAddr)
		if false {
			mp, _ := common.AddressFromHexString("ac820725a3bc5337b6044aad38dc3293032eae99")
			dtoken, _ := common.AddressFromHexString("b57362d82c658fb2927e793406ee8decf32ff2e9")
			utils.Init(sdk, con, seller, mp, dtoken)
			return
		}
		if false {
			dtoken, _ := common.AddressFromHexString("3343753265152550e5a1741cea946436744ab442")
			utils.SetDtokenContractAddr(sdk, con, seller, dtoken)
			return
		}
		if false {
			mp, _ := common.AddressFromHexString("5fbcadf08b14aa737de8af429483dc4fb1ae13d3")
			utils.SetMpContractAddr(sdk, con, seller, mp)
			return
		}
		if true {
			resource_id := []byte("637088084609811033")
			template_id, _ := hex.DecodeString("30")
			utils.BuyAndUseToken(sdk, con, resource_id, 1, buyer, payer, template_id)
			return
		}
	}
	if false {
		sdk.DefDTokenKit().SetContractAddr(contractAddr)
		if false {
			utils.CreateTokenTemplate(sdk, seller)
			return
		}
		if false {
			utils.GenerateDtoken(sdk, seller)
			return
		}
		tokenId, _ := hex.DecodeString("31")
		if true {
			utils.BalanceOf(sdk, buyer.Address, tokenId)
			return
		}
		if err = addAgents(sdk, tokenId); err != nil {
			fmt.Println("addAgents error: ", err)
			return
		}

		//if err = useTokenByAgent(sdk, tokenId); err != nil {
		//	fmt.Println("useTokenByAgent error: ", err)
		//	return
		//}

		if err = removeAgents(sdk, tokenId); err != nil {
			fmt.Println("removeAgents error: ", err)
			return
		}

		if err = addTokenAgents(sdk, tokenId); err != nil {
			fmt.Println("addTokenAgents error: ", err)
			return
		}
		if err = removeTokenAgents(sdk); err != nil {
			fmt.Println("removeTokenAgents error: ", err)
			return
		}

		err = useToken(sdk, tokenId)
		if err != nil {
			fmt.Printf("useToken: %s\n", err)
			return
		}
		return
	}

	if true {
		sdk.DefMpKit().SetContractAddress(contractAddr)
		resourceId := strconv.Itoa(rand.Int())
		fmt.Println("resourceId:", resourceId)
		resourceIdBytes := []byte(resourceId)
		dataId := ""
		tokenTemplate = &market_place_contract.TokenTemplate{
			DataID:     dataId,
			TokenHashs: []string{string(common.UINT256_EMPTY[:])},
		}

		if err = publish(sdk, resourceIdBytes); err != nil {
			fmt.Println("publish error: ", err)
			return
		}

		//if err := delete(sdk, resourceIdBytes); err != nil {
		//	fmt.Println("delete error: ", err)
		//	return
		//}
		//return

		//if err := update(sdk, resourceIdBytes); err != nil {
		//	fmt.Println("update error: ", err)
		//	return
		//}

		//if err = buyAndUseToken(sdk, resourceIdBytes); err != nil {
		//	fmt.Println("buyAndUseToken error: ", err)
		//	return
		//}

		//if err = buyDToken(sdk, resourceIdBytes); err != nil {
		//	fmt.Println("buyDToken error: ", err)
		//	return
		//}
	}
}

func addTokenAgents(sdk *ddxf_sdk.DdxfSdk, tokenId []byte) error {
	txHash, err := sdk.DefDTokenKit().AddTokenAgents(buyer,
		[]common.Address{agent.Address}, tokenId, []int{1})
	if err != nil {
		return err
	}
	return showNotify(sdk, "addTokenAgents", txHash.ToHexString())
}

func removeTokenAgents(sdk *ddxf_sdk.DdxfSdk) error {
	txHash, err := sdk.DefDTokenKit().RemoveTokenAgents([]byte(""), buyer,
		[]common.Address{agent.Address})
	if err != nil {
		return err
	}
	return showNotify(sdk, "removeTokenAgents", txHash.ToHexString())
}

func removeAgents(sdk *ddxf_sdk.DdxfSdk, tokenId []byte) error {
	txHash, err := sdk.DefDTokenKit().RemoveAgents(buyer, []common.Address{agent.Address}, [][]byte{tokenId})
	if err != nil {
		return err
	}
	return showNotify(sdk, "removeAgents", txHash.ToHexString())
}

func useTokenByAgent(sdk *ddxf_sdk.DdxfSdk, tokenId []byte) error {

	txHash, err := sdk.DefDTokenKit().UseTokenByAgents(buyer.Address, agent, tokenId, 1)
	if err != nil {
		return err
	}
	return showNotify(sdk, "UseTokenByAgents", txHash.ToHexString())
}

func addAgents(sdk *ddxf_sdk.DdxfSdk, tokenId []byte) error {
	txHash, err := sdk.DefDTokenKit().AddAgents(buyer,
		[]common.Address{agent.Address}, []int{1}, [][]byte{tokenId})
	if err != nil {
		fmt.Printf("AddAgents: %s\n", err)
		return err
	}
	return showNotify(sdk, "addAgents", txHash.ToHexString())
}

func useToken(sdk *ddxf_sdk.DdxfSdk, tokenId []byte) error {
	txHash, err := sdk.DefDTokenKit().UseToken(seller, tokenId, 1)
	if err != nil {
		return err
	}
	return showNotify(sdk, "useToken", txHash.ToHexString())
}

//func buyAndUseToken(sdk *ddxf_sdk.DdxfSdk, resourceIdBytes []byte) error {
//	txHash, err := sdk.DefMpKit().BuyAndUseToken(buyer, payer, resourceIdBytes, 2, *tokenTemplate)
//	if err != nil {
//		return err
//	}
//	return showNotify(sdk, "buyAndUseToken", txHash.ToHexString())
//}

func buyDToken(sdk *ddxf_sdk.DdxfSdk, resourceIdBytes []byte) error {
	txHash, err := sdk.DefMpKit().BuyDToken(buyer, payer, resourceIdBytes, 10)
	if err != nil {
		return err
	}
	return showNotify(sdk, "buyDToken", txHash.ToHexString())
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
	return showNotify(sdk, "delete", txHash.ToHexString())
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
		DTC:          []common.Address{},   // can be empty
		Accountant:   common.ADDRESS_EMPTY, // can be empty
		Split:        common.ADDRESS_EMPTY,
	}

	tokenTemplateId, _ := hex.DecodeString("30")

	item := market_place_contract.DTokenItem{
		Fee: market_place_contract.Fee{
			ContractType: 0,
			Count:        1,
		},
		ExpiredDate:      uint64(time.Now().Unix()) + 10000,
		Stocks:           10000,
		TokenTemplateIds: []string{string(tokenTemplateId)},
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
		fmt.Printf("update error:%s\n", err)
		return err
	}
	fmt.Println("update txHash: ", txHash.ToHexString())

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
		DTC:          []common.Address{},   // can be empty
		Accountant:   common.ADDRESS_EMPTY, // can be empty
		Split:        common.ADDRESS_EMPTY,
	}

	tokenTemplateId, _ := hex.DecodeString("30")
	item := market_place_contract.DTokenItem{
		Fee: market_place_contract.Fee{
			ContractType: 0,
			Count:        1,
		},
		ExpiredDate:      uint64(time.Now().Unix()) + 10000,
		Stocks:           10000,
		TokenTemplateIds: []string{string(tokenTemplateId)},
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
	name := "ontology-vote"
	desc := "smart contract for ontology vote"
	txHash, err := sdk.DeployContract(admin, codeHex, name, "0.1.1", "lucas", "", desc)
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
