# ddxf-sdk


This sdk is used to interact with the [ddxf contract suite](https://github.com/ont-bizsuite/ddxf-contract-suite).


First of all, create a ddxf sdk instance and prepare the account:

```golang
import (
    osdk "github.com/ontio/ontology-go-sdk"
)

// ontology testnet
addr := "http://polaris2.ont.io:20336"
sdk := NewDdxfSdk(addr)

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

type DataIdInfo struct {
	DataId       string
	DataMetaHash common.Uint256
	DataHash     common.Uint256
	Owners       []*OntIdIndex
}

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


wallet, _ := osdk.NewOntologySdk().OpenWallet("./wallet.dat")

account, _ = wallet.GetAccountByAddress("account_address", []byte("password"))

resource := uuid.NewV4().String()

itemMeta := map[string]interface{}{
    "attr1": "value1",
    "attr2": "value2",
}

itemMetaHash, _ = ddxf.HashObject(input.Item)
ddo := ResourceDDO{
		Manager:      account.Address,
		ItemMetaHash: itemMetaHash,
	}

templates := make([]*market_place_contract.TokenTemplate, 0)
for i := 0; i < len(dataMetas); i++ {
    var dataMetaHash [sha256.Size]byte
    dataMetaHash = dataMetaHashArray[i]
    u, _ := common2.Uint256ParseFromBytes(dataMetaHash[:])
    dataId := res[u.ToHexString()]
    tt := &market_place_contract.TokenTemplate{
        DataID:     dataId.(string),
        TokenHashs: []string{"1"},
        Endpoint:   "aaaa",
    }
    templates = append(templates, tt)
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