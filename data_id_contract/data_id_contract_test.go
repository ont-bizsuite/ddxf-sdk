package data_id_contract

import (
	"testing"

	"fmt"
	"github.com/ont-bizsuite/ddxf-sdk/base_contract"
	"github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology/common"
	"math/rand"
	"strconv"
	"time"
)

var (
	dataId  *DataIdKit
	account *ontology_go_sdk.Account
	sdk     *ontology_go_sdk.OntologySdk
)

func TestMain(m *testing.M) {
	contractAddr, _ := common.AddressFromHexString("03d07edc239d6c8cd99543bdef7a4bc407f9d44c")
	sdk = ontology_go_sdk.NewOntologySdk()
	sdk.NewRpcClient().SetAddress("http://127.0.0.1:20336")
	bc := base_contract.NewBaseContract(sdk, 200000000, 0, nil)
	dataId = NewDataIdContractKit(contractAddr, bc)

	wallet, _ := sdk.OpenWallet("./wallet.dat")
	account, _ = wallet.NewDefaultSettingAccount([]byte("123456"))
	txhash, err := sdk.Native.OntId.RegIDWithPublicKey(0, 20000000, account,
		"did:ont:"+account.Address.ToBase58(), account)
	if err != nil {
		fmt.Println(err)
		return
	}
	handleEvent(txhash)
	m.Run()
}

func handleEvent(txHash common.Uint256) {
	time.Sleep(6 * time.Second)
	evt, err := sdk.GetSmartContractEvent(txHash.ToHexString())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(evt)
}

func TestDataIdKit_RegisterDataIdInfo(t *testing.T) {
	info := DataIdInfo{
		DataId:       strconv.Itoa(rand.Int()),
		DataMetaHash: common.UINT256_EMPTY,
		DataHash:     common.UINT256_EMPTY,
		Owners: []*OntIdIndex{
			&OntIdIndex{
				OntId: "did:ont:" + account.Address.ToBase58(),
				index: 1,
			},
		},
	}
	txhash, err := dataId.RegisterDataIdInfo(info, account)
	if err != nil {
		fmt.Println(err)
		return
	}
	handleEvent(txhash)
}

func TestDataIdKit_RegisterDataIdInfoArray(t *testing.T) {
	info := DataIdInfo{
		DataId:       strconv.Itoa(rand.Int()),
		DataMetaHash: common.UINT256_EMPTY,
		DataHash:     common.UINT256_EMPTY,
		Owners: []*OntIdIndex{
			&OntIdIndex{
				OntId: "did:ont:" + account.Address.ToBase58(),
				index: 1,
			},
		},
	}
	txhash, err := dataId.RegisterDataIdInfoArray([]DataIdInfo{info}, account)
	if err != nil {
		fmt.Println(err)
		return
	}
	handleEvent(txhash)
}

func TestDataIdKit_BuildRegisterDataIdInfoArrayTx(t *testing.T) {

}
