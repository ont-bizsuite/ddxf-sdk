package utils

import (
	"github.com/ont-bizsuite/ddxf-sdk/example/base"
	"github.com/ont-bizsuite/ddxf-sdk"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/core/utils"
	"fmt"
	"github.com/ontio/ontology-go-sdk"
	"encoding/hex"
)

func DataIdTest(sdk *ddxf_sdk.DdxfSdk, pwd []byte, seller *ontology_go_sdk.Account) {

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
	//contractAddr, _ := common.AddressFromHexString("df04263aa6ff06bdaf6ba50d29c4cb2a188078cd")
	//con := sdk.DefContract(contractAddr)

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

	tx, err := sdk.GetOntologySdk().Native.OntId.NewAddAttributesTransaction(500, 200000, iden.ID, nil, seller.PublicKey)
	if err != nil {
		fmt.Println(err)
		return
	}
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
}