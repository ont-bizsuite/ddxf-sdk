package utils

import (
	"encoding/hex"
	"fmt"
	"github.com/ont-bizsuite/ddxf-sdk"
	"github.com/ont-bizsuite/ddxf-sdk/market_place_contract"
	"github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology/common"
)

func CreateTokenTemplate(sdk *ddxf_sdk.DdxfSdk, seller *ontology_go_sdk.Account) {
	tt := market_place_contract.TokenTemplate{
		DataID:      "",
		TokenHashs:  []string{""},
		Endpoint:    "",
		TokenName:   "name",
		TokenSymbol: "symbol",
	}
	txhash, err := sdk.DefDTokenKit().CreateTokenTemplate(common.ADDRESS_EMPTY,seller, tt)
	if err != nil {
		fmt.Println("CreateTokenTemplate error: ", err)
		return
	}
	showNotify(sdk, "CreateTokenTemplate", txhash.ToHexString())
}

func AuthorizeTokenTemplate(sdk *ddxf_sdk.DdxfSdk, templateId []byte, seller *ontology_go_sdk.Account, buyer []common.Address) {
	txhash, err := sdk.DefDTokenKit().AuthorizeTokenTemplate(common.ADDRESS_EMPTY,seller, templateId, buyer)
	if err != nil {
		fmt.Println("CreateTokenTemplate error: ", err)
		return
	}
	showNotify(sdk, "AuthorizeTokenTemplate", txhash.ToHexString())
}

func GenerateDtoken(sdk *ddxf_sdk.DdxfSdk, seller *ontology_go_sdk.Account) {
	tokenTemplateId, _ := hex.DecodeString("3238")
	txhash, err := sdk.DefDTokenKit().GenerateDToken(common.ADDRESS_EMPTY,seller, tokenTemplateId, 1000)
	if err != nil {
		fmt.Println(err)
		return
	}
	showNotify(sdk, "GenerateDtoken", txhash.ToHexString())
}

func BalanceOf(sdk *ddxf_sdk.DdxfSdk, addr common.Address, tokenId []byte) {
	res, err := sdk.DefDTokenKit().BalanceOf(common.ADDRESS_EMPTY,addr, tokenId)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
}
