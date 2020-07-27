package utils

import (
	"github.com/ont-bizsuite/ddxf-sdk/example/base"
	"github.com/ont-bizsuite/ddxf-sdk"
	"github.com/ontio/ontology/common"
	"fmt"
	"github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology/core/utils"
	"encoding/hex"
)

func DataIdTest(sdk *ddxf_sdk.DdxfSdk, pwd []byte, seller *ontology_go_sdk.Account, contractAddr common.Address) {

	bs, err := sdk.GetOntologySdk().Native.OntId.GetDocumentJson("did:ont:TXvDhLqrqvAV6XUAmLEfWLjxmS1ESxbZBr")

	fmt.Println(string(bs))
	wallet, err := sdk.GetOntologySdk().OpenWallet("./wallet.dat")
	if err != nil {
		fmt.Println(err)
		return
	}
	var iden *ontology_go_sdk.Identity
	members := make([][]byte,0)
	for i:=0;i<16;i++ {
		iden, err = wallet.NewDefaultSettingIdentity(pwd)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("iden:", iden.ID)
		txhash, err := sdk.GetOntologySdk().Native.OntId.RegIDWithPublicKey(2500, 2000000, seller, iden.ID, seller)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("txhash:",txhash.ToHexString())
		evt, err := sdk.GetSmartCodeEvent(txhash.ToHexString())
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("RegIDWithPublicKey evt:", evt)
		members = append(members, []byte(iden.ID))
	}


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
			Members:   members,
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


	bs, err = utils.BuildWasmContractParam([]interface{}{"reg_id_add_attribute_array", []interface{}{[]interface{}{sink.Bytes()}}})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("param:", hex.EncodeToString(bs))

	txhash, err := con.Invoke("reg_id_add_attribute_array",seller, []interface{}{[]interface{}{sink.Bytes()}})
	if err != nil {
		fmt.Println(err)
		return
	}


	evt, err := sdk.GetSmartCodeEvent(txhash.ToHexString())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(evt)
	return

	tx, err := sdk.GetOntologySdk().Native.OntId.NewAddAttributesTransaction(2500, 200000, iden.ID, nil, seller.PublicKey)
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