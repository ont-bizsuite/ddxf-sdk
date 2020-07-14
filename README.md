# ddxf-sdk


This sdk is used to interact with the [ddxf contract suite](https://github.com/ont-bizsuite/ddxf-contract-suite).

目前，ddxf-sdk封装了`dtoken`合约、`marketplace`合约和`split_policy`合约的调用实现。


## dtoken合约接口设计

1. 初始化

先创建 `ddxf-sdk`实例, 需要指定连接的节点地址, gasPrice、gasLimit会使用默认值。
```
testNet = ddxf_sdk.TestNet
sdk := ddxf_sdk.NewDdxfSdk(testNet)
```

如果自定义gasPrice、gasLimit、payer(交易手续费的支付地址)
```
sdk.SetGasPrice(2500)
sdk.SetGasLimit(20000)
sdk.SetPayer(account)
dtoken := sdk.DefDTokenKit()
```

创建 `Dtoken`实例
```
dtoken := sdk.DefDTokenKit()
```

2. `CreateTokenTemplate` 

创建`TokenTemplate`, 该方法会生成TokenTemplateId, 用来表示链上唯一的TokenTemplate

接口定义
```
func (this *DTokenKit) CreateTokenTemplate(creator *ontology_go_sdk.Account,
	tt market_place_contract.TokenTemplate) (common.Uint256, error)
```

调用示例
```
tt := market_place_contract.TokenTemplate{
		DataID:"",
		TokenHashs:[]string{""},
		Endpoint:"",
		TokenName:"name",
		TokenSymbol:"symbol",
	}
dtoken.CreateTokenTemplate(creator, tt)
```

参数解释:
* `creator` `TokenTemplate`创建者Account,需要该Account对交易签名
* `TokenTemplate` 
   * `DataID` 物的ontid, 也是物的链上唯一标志
   * `TokenHashs` TokenHash数组
   * `Endpoint` TokenMeta的访问接口
   * `TokenName` 根据TokenTemplate创建Dtoken时生成的Token的Name
   * `TokenSymbol` 根据TokenTemplate创建Dtoken时生成的Token的Symbol

返回值解释
* `common.Uint256` 交易的hash
* `error` 交易错误信息

Event设计
该交易会推出 `Event`, 其数据结构如下
```
[string, string, bytearray, bytearray]
```
第一个参数表示调用的合约的方法名
第二个参数表示creator, base58编码的地址字符串
第三个参数表示TokenTemplate序列化后的字节数组
第四个参数表示TokenTemplateId


3. `AuthorizeTokenTemplate` 
TokenTemplate创建者授权别的地址可以根据该TokenTemplate生成DToken

接口定义
```
func (this *DTokenKit) AuthorizeTokenTemplate(creator *ontology_go_sdk.Account, 
tokenTemplateId []byte,authorizeAddr common.Address) (common.Uint256, error) 
```
调用示例
```
dtoken.AuthorizeTokenTemplate(creator, tokenTemplateId, authorizeAddr)
```
参数解释
* `creator` TokenTemplate的创建者
* `tokenTemplateId` TokenTemplate链上的唯一性标志
* `authorizeAddr` 被授权的地址

4. `GenerateDToken` 
生成DToken， 该交易会推出Event, 其中含有生成的 `TokenId`字段, 生成的DToken都在传入的地址里面

接口设计
```
func (this *DTokenKit) GenerateDToken(acc *ontology_go_sdk.Account, 
tokenTemplateId []byte, n int) (common.Uint256, error) 
```
调用示例
```
dtoken.GenerateDToken(acc, tokenTemplateId, n)
```
参数解释
* `acc` TokenTemplate的创建者或者被授权者
* `tokenTemplateId` TokenTemplate链上的唯一性标志
* `n` 生成的DToken的数量

Event设计
该交易会推出 `Event`, 其数据结构如下
```
[string, string, bytearray, int, bytearray]
```
第一个参数表示调用的合约的方法名
第二个参数输入的参数acc的地址
第三个参数表示TokenTemplateId
第四个参数表示n
第五个参数表示生成的TokenId

5. BalanceOf
查询给定地址的DToken余额

接口设计
```
func (this *DTokenKit) BalanceOf(addr common.Address, tokenId []byte) (uint64, error) 
```

调用示例
```
dtoken.BalanceOf(addr, tokenId)
```

参数解释
* addr 查询余额的地址
* tokenId TokenID

6. `UseToken`

使用Token
```
func (this *DTokenKit) UseToken(buyer *ontology_go_sdk.Account,
	tokenId []byte, n int) (common.Uint256, error) 
```
调用示例
```
dtoken.UseToken(buyer, tokenId, n)
```
参数解释
* `buyer` DToken的购买者
* `tokenId` TokenId
* `n` 使用的数量


## Marketplace合约接口设计

1. 初始化Marketplace合约实例

先创建 `ddxf-sdk`实例, 需要指定连接的节点地址, gasPrice、gasLimit会使用默认值。
```
testNet = ddxf_sdk.TestNet
sdk := ddxf_sdk.NewDdxfSdk(testNet)
```

如果自定义gasPrice、gasLimit、payer(交易手续费的支付地址)
```
sdk.SetGasPrice(2500)
sdk.SetGasLimit(20000)
sdk.SetPayer(account)
```

创建 `mp`实例
```
mp := sdk.DefMpKit()
```

2. 发布商品

用户将商品发布到链上

接口设计
```
func (this *MpKit) Publish(seller *ontology_go_sdk.Account, resourceId []byte, 
ddo ResourceDDO, item DTokenItem,splitPolicyParam split_policy_contract.SplitPolicyRegisterParam)
	 (common.Uint256, error)
```

调用示例
```
mp.Publish(seller, resourceId, ddo, item)
```

参数解释
* `seller` 卖家的Account
* `resourceId` 用来标志链上商品的唯一性
* `ddo` 
   * `Manager` 是seller的地址
   * `ItemMetaHash` ItemMeta的hash
   * `DTC` DToken合约地址
   * `Accountant` Accountant合约地址
   * `Split` 分润的合约地址
* `item`
   * `Fee` 
      * `ContractAddr` fee的合约地址
      * `ContractType` fee合约类型，目前支持ont,ong,oep4类型的合约
      * `Count`  手续费
   * `ExpiredDate` 过期时间
   * `Stocks` 库存数量
   * `Sold` 已经销售的数量
   * `TokenTemplateIds` TokenTemplateId数组
* `splitPolicyParam`
   * `AddrAmts` 地址和数量数组
   * `TokenTy` Token类型
   * `ContractAddr` Token的合约地址




3. `BuyDtoken`
购买商品，用户可以购买已经发布的商品，购买后，会生成的`TokenId`, 生成的DToken会登记在买家的地址，

接口设计
```
func (this *MpKit) BuyDtoken(buyer, payer *ontology_go_sdk.Account, 
resourceId []byte,n int) (common.Uint256, error) 
```

参数解释
* `buyer` 购买方的Account
* `payer` 付钱方的Account
* `resourceId` 商品的唯一性标志
* `n` 购买的数量




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

dataID, _ := ontology_go_sdk.GenerateID()
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