package split_policy_contract

import (
	"testing"

	"encoding/hex"
	"fmt"
	"github.com/ont-bizsuite/ddxf-sdk/base_contract"
	"github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology/common"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"math/rand"
	"strconv"
	"time"
)

var (
	splitPolicyKit *SplitPolicyKit
	wallet         *ontology_go_sdk.Wallet
	pwd            = []byte("123456")
	ontSdk         *ontology_go_sdk.OntologySdk
	payer          *ontology_go_sdk.Account
	seller         *ontology_go_sdk.Account
	admin          *ontology_go_sdk.Account
	gasPrice       = uint64(0)
	gasLimit       = uint64(20000000)
	testNet        = "http://polaris1.ont.io:20336"
	localHost      = "http://127.0.0.1:20336"
)

func TestMain(m *testing.M) {
	ontSdk = ontology_go_sdk.NewOntologySdk()
	ontSdk.NewRpcClient().SetAddress(localHost)
	var err error
	wallet, err = ontSdk.OpenWallet("./wallet.dat")
	if err != nil {
		fmt.Printf("error in ReadFile:%s\n", err)
		return
	}
	payer, _ = wallet.GetAccountByAddress("AYnhakv7kC9R5ppw65JoE2rt6xDzCjCTvD", pwd)
	seller, _ = wallet.GetAccountByAddress("Aejfo7ZX5PVpenRj23yChnyH64nf8T1zbu", pwd)
	admin, _ = wallet.GetAccountByAddress("AYnhakv7kC9R5ppw65JoE2rt6xDzCjCTvD", pwd)
	wasmFile := "/Users/sss/dev/dockerData/rust_project/ddxf_market/output/split_policy.wasm"
	code, err := ioutil.ReadFile(wasmFile)
	if err != nil {
		fmt.Printf("error in ReadFile:%s\n", err)
		return
	}
	//only need execute once
	if true {
		txHash, err := ontSdk.WasmVM.DeployWasmVMSmartContract(gasPrice, 200000000, admin,
			hex.EncodeToString(code), "", "", "", "", "")
		if err != nil {
			fmt.Println("err:", err)
			return
		}
		time.Sleep(10 * time.Second)
		evt, err := ontSdk.GetSmartContractEvent(txHash.ToHexString())
		if err != nil {
			fmt.Println("err:", err)
			return
		}
		fmt.Println("evts: ", evt)
	}
	contractAddress := common.AddressFromVmCode(code)
	fmt.Println("contractAddress: ", contractAddress.ToHexString())
	bc := base_contract.NewBaseContract(ontSdk, 20000000, gasPrice, payer)
	splitPolicyKit = NewSplitPolicyKit(contractAddress, bc)
	m.Run()
}

func TestSplitPolicyKit_Register(t *testing.T) {
	randIn := rand.Int()
	key := []byte(strconv.Itoa(randIn))
	rp := SplitPolicyRegisterParam{
		AddrAmts: []*AddrAmt{
			&AddrAmt{
				To:      seller.Address,
				Percent: 5000,
			},
			&AddrAmt{
				To:      admin.Address,
				Percent: 5000,
			},
		},
		TokenTy:      ONG,
		ContractAddr: common.ADDRESS_EMPTY,
	}
	fmt.Println("param: ", hex.EncodeToString(rp.ToBytes()))
	txHash, err := splitPolicyKit.Register(key, rp, seller)
	time.Sleep(10 * time.Second)
	assert.Nil(t, err)
	evt, err := ontSdk.GetSmartContractEvent(txHash.ToHexString())
	assert.Nil(t, err)
	fmt.Println("evt: ", evt)

	regParam, err := splitPolicyKit.GetRegisterParam(key)
	assert.Nil(t, err)
	fmt.Println("regParam: ", regParam)

	txHash, err = splitPolicyKit.Withdraw(key, seller)
	assert.Nil(t, err)
	time.Sleep(10 * time.Second)
	evt, err = ontSdk.GetSmartContractEvent(txHash.ToHexString())
	assert.Nil(t, err)
	fmt.Println("evt:", evt)
}
