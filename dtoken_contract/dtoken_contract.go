package dtoken_contract

import (
	"fmt"
	"github.com/ont-bizsuite/ddxf-sdk/base_contract"
	"github.com/ont-bizsuite/ddxf-sdk/market_place_contract"
	"github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology/common"
)

type DTokenKit struct {
	bc              *base_contract.BaseContract
	contractAddress common.Address
}

func NewDTokenKit(contractAddress common.Address, bc *base_contract.BaseContract) *DTokenKit {
	return &DTokenKit{
		contractAddress: contractAddress,
		bc:              bc,
	}
}

func (this *DTokenKit) SetContractAddr(addr common.Address) {
	this.contractAddress = addr
}

func (this *DTokenKit) SetMpContractAddr(admin *ontology_go_sdk.Account, addr common.Address) (common.Uint256, error) {
	return this.bc.Invoke(this.contractAddress, admin, "setMpContract",
		[]interface{}{addr})
}

func (this *DTokenKit) GetMpContractAddr() (common.Address, error) {
	res, err := this.bc.PreInvoke(this.contractAddress, "getMpContract", []interface{}{})
	if err != nil {
		return common.ADDRESS_EMPTY, err
	}
	bs, err := res.ToByteArray()
	if err != nil {
		return common.ADDRESS_EMPTY, err
	}
	source := common.NewZeroCopySource(bs)
	addr, eof := source.NextAddress()
	if eof {
		return common.ADDRESS_EMPTY, fmt.Errorf("read address failed, eof: %v", eof)
	}
	return addr, nil
}

func (this *DTokenKit) CreateTokenTemplate(creator *ontology_go_sdk.Account,
	tt market_place_contract.TokenTemplate) (common.Uint256, error) {
	sink := common.NewZeroCopySink(nil)
	tt.Serialize(sink)
	return this.bc.Invoke(this.contractAddress, creator, "createTokenTemplate", []interface{}{creator.Address, sink.Bytes()})
}

func (this *DTokenKit) AuthorizeTokenTemplate(creator *ontology_go_sdk.Account, tokenTemplateId []byte,
	authorizeAddr common.Address) (common.Uint256, error) {
	return this.bc.Invoke(this.contractAddress, creator, "authorizeTokenTemplate",
		[]interface{}{tokenTemplateId, authorizeAddr})
}

func (this *DTokenKit) GetAuthorizedAddr(tokenTemplateId []byte) ([]common.Address, error) {
	res, err := this.bc.PreInvoke(this.contractAddress, "getAuthorizedAddr", []interface{}{tokenTemplateId})
	if err != nil {
		return nil, err
	}
	arr, err := res.ToArray()
	if err != nil {
		return nil, err
	}
	addrs := make([]common.Address, 0)
	for _, item := range arr {
		bs, err := item.ToByteArray()
		if err != nil {
			return nil, err
		}
		addr, err := common.AddressParseFromBytes(bs)
		if err != nil {
			return nil, err
		}
		addrs = append(addrs, addr)
	}
	return addrs, nil
}

func (this *DTokenKit) GetTokenIdByTemplateId(templateId []byte) ([]byte, error) {
	res, err := this.bc.PreInvoke(this.contractAddress, "getTokenIdByTemplateId", []interface{}{templateId})
	if err != nil {
		return nil, err
	}
	return res.ToByteArray()
}
func (this *DTokenKit) GetTemplateIdByTokenId(tokenId []byte) ([]byte, error) {
	res, err := this.bc.PreInvoke(this.contractAddress, "getTemplateIdByTokenId", []interface{}{tokenId})
	if err != nil {
		return nil, err
	}
	return res.ToByteArray()
}

func (this *DTokenKit) GenerateDToken(acc *ontology_go_sdk.Account, tokenTemplateId []byte, n int) (common.Uint256, error) {
	return this.bc.Invoke(this.contractAddress, acc, "generateDToken", []interface{}{acc.Address, tokenTemplateId, n})
}

func (this *DTokenKit) BalanceOf(addr common.Address, tokenId []byte) (uint64, error) {
	res, err := this.bc.PreInvoke(this.contractAddress, "balanceOf", []interface{}{addr, tokenId})
	if err != nil {
		return 0, err
	}
	in, err := res.ToInteger()
	if err != nil {
		return 0, err
	}
	return in.Uint64(), nil
}

func (this *DTokenKit) DeleteToken(acc *ontology_go_sdk.Account, tokenId []byte) (common.Uint256, error) {
	return this.bc.Invoke(this.contractAddress, acc, "deleteToken", []interface{}{acc.Address, tokenId})
}

func (this *DTokenKit) GetAgentBalance(agent common.Address) (uint64, error) {
	res, err := this.bc.PreInvoke(this.contractAddress, "getAgentBalance", []interface{}{agent})
	if err != nil {
		return 0, err
	}
	in, err := res.ToInteger()
	if err != nil {
		return 0, err
	}
	return in.Uint64(), nil
}

func (this *DTokenKit) UseToken(buyer *ontology_go_sdk.Account,
	tokenId []byte, n int) (common.Uint256, error) {
	return this.bc.Invoke(this.contractAddress, buyer, "useToken",
		[]interface{}{buyer.Address, tokenId, n})
}

func (this *DTokenKit) UseTokenByAgents(tokenOwner common.Address,
	agent *ontology_go_sdk.Account, tokenId []byte, n int) (common.Uint256, error) {
	return this.bc.Invoke(this.contractAddress, agent, "useTokenByAgent",
		[]interface{}{tokenOwner, agent.Address, tokenId, n})
}

func (this *DTokenKit) SetAgents(account *ontology_go_sdk.Account,
	agents []common.Address, n int) (common.Uint256, error) {
	return this.bc.Invoke(this.contractAddress, account, "setAgents",
		[]interface{}{account.Address, agents, n})
}

func (this *DTokenKit) SetTokenAgents(account *ontology_go_sdk.Account,
	agents []common.Address, tokenId []byte, n int) (common.Uint256, error) {
	return this.bc.Invoke(this.contractAddress, account, "setTokenAgents",
		[]interface{}{account.Address, agents, tokenId, n})
}

func (this *DTokenKit) AddAgents(account *ontology_go_sdk.Account,
	agents []common.Address, n int, tokenIds [][]byte) (common.Uint256, error) {
	tokenIdParam := make([]interface{}, 0)
	for _, item := range tokenIds {
		tokenIdParam = append(tokenIdParam, item)
	}
	return this.bc.Invoke(this.contractAddress, account, "addAgents",
		[]interface{}{account.Address, parseAddressArr(agents), n, tokenIdParam})
}

func parseAddressArr(addrs []common.Address) []interface{} {
	agentArr := make([]interface{}, len(addrs))
	for i := 0; i < len(addrs); i++ {
		agentArr[i] = addrs[i]
	}
	return agentArr
}

func (this *DTokenKit) AddTokenAgents(account *ontology_go_sdk.Account,
	agents []common.Address, tokenId []byte, n int) (common.Uint256, error) {
	return this.bc.Invoke(this.contractAddress, account, "addTokenAgents",
		[]interface{}{account.Address, tokenId, parseAddressArr(agents), n})
}

func (this *DTokenKit) RemoveAgents(account *ontology_go_sdk.Account,
	agents []common.Address, tokenIds [][]byte) (common.Uint256, error) {
	tokenIdParam := make([]interface{}, 0)
	for _, item := range tokenIds {
		tokenIdParam = append(tokenIdParam, item)
	}
	return this.bc.Invoke(this.contractAddress, account, "removeAgents",
		[]interface{}{account.Address, parseAddressArr(agents), tokenIdParam})
}

func (this *DTokenKit) RemoveTokenAgents(tokenId []byte, account *ontology_go_sdk.Account,
	agents []common.Address) (common.Uint256, error) {
	return this.bc.Invoke(this.contractAddress, account, "removeTokenAgents",
		[]interface{}{account.Address, tokenId, parseAddressArr(agents)})
}
