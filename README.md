# ddxf-sdk


This sdk is used to interact with the [ddxf contract suite](https://github.com/ont-bizsuite/ddxf-contract-suite).


First of all, create a ddxf sdk instance and prepare the account:

```golang
import (
    osdk "github.com/ontio/ontology-go-sdk"
    dsdk "github.com/ont-bizsuite/ddxf-sdk"
)

// ontology testnet
addr := "http://polaris2.ont.io:20336"
sdk := dsdk.NewDdxfSdk(addr)

wallet, _ := osdk.NewOntologySdk().OpenWallet("./wallet.dat")

account, _ = wallet.GetAccountByAddress("account_address", []byte("password"))
```

Then create an `dataid` for the goods:

```golang
import (
    "github.com/ont-bizsuite/ddxf-sdk/data_id_contract"
    uuid "github.com/satori/go.uuid"
    "github.com/zhiqiangxu/ddxf"
)

dataID := uuid.NewV4().String()
dataMeta := map[string]interface{}{
    "dataMetaAttr":"value",
}
dataMetaHash, _ := ddxf.HashObject(dataMeta)
dataIdInfo := data_id_contract.DataIdInfo{
    DataId:dataID, DataMetaHash:dataMetaHash, 
    DataHash:     common.UINT256_EMPTY,
		Owners: []*OntIdIndex{
			&OntIdIndex{
				OntId: "did:ont:" + account.Address.ToBase58(),
				index: 1,
			},
        },
}
sdk.DefDataIdKit().RegisterDataIdInfo(dataIdInfo, account)


```

Then call `sdk.DefMpKit().Publish` to publish the goods onto ontology:

```golang
import (
    osdk "github.com/ontio/ontology-go-sdk"
    uuid "github.com/satori/go.uuid"
    "github.com/zhiqiangxu/ddxf"
    "crypto/sha256"
    "github.com/ont-bizsuite/ddxf-sdk/split_policy_contract"
    "github.com/ont-bizsuite/ddxf-sdk/market_place_contract"
)



resource := uuid.NewV4().String()

itemMeta := map[string]interface{}{
    "attr1": "value1",
    "attr2": "value2",
}

itemMetaHash, _ = ddxf.HashObject(itemMeta)
ddo := ResourceDDO{
		Manager:      account.Address,
		ItemMetaHash: itemMetaHash,
	}

item := market_place_contract.DTokenItem{
		Fee: market_place_contract.Fee{
			ContractAddr: addr,
			ContractType: split_policy_contract.OEP4,
			Count:        1,
		},
		ExpiredDate: uint64(time.Now().Unix()) + uint64(time.Hour*24*30),
		Stocks:      10000,
		Templates:   []*market_place_contract.TokenTemplate{
            &market_place_contract.TokenTemplate{
                DataID:     dataID,
                TokenHashs: []string{"1"},
                Endpoint:   "your ddxf endpoint",
            },
        },
    }
    
splitPolicyParam := split_policy_contract.SplitPolicyRegisterParam{
		AddrAmts: []*split_policy_contract.AddrAmt{
			&split_policy_contract.AddrAmt{
				To:      account.Address,
				Percent: 100,
			},
		},
		TokenTy: split_policy_contract.ONG,
	}    
sdk.DefMpKit().Publish(account, resourceId, ddo, item, splitPolicyParam)
```


After that, the buyer can buy the above dtoken:


```golang
sdk.DefMpKit().BuyDtoken(buyer, payer *ontology_go_sdk.Account, resourceId, n int)
```