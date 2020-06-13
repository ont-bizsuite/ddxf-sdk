package split_policy_contract

import (
	"github.com/ontio/ddxf-sdk/base_contract"
	"github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology/common"
)

type SplitPolicyKit struct {
	bc              *base_contract.BaseContract
	contractAddress common.Address
}

func NewSplitPolicyKit(contractAddress common.Address, bc *base_contract.BaseContract) *SplitPolicyKit {
	return &SplitPolicyKit{
		contractAddress: contractAddress,
		bc:              bc,
	}
}

func (this *SplitPolicyKit) SetContractAddress(addr common.Address) {
	this.contractAddress = addr
}

func (this *SplitPolicyKit) Register(key []byte, rp SplitPolicyRegisterParam,
	signer *ontology_go_sdk.Account) (common.Uint256, error) {
	return this.bc.Invoke(this.contractAddress, signer, "register",
		[]interface{}{key, rp.ToBytes()})
}

func (this *SplitPolicyKit) GetRegisterParam(key []byte) (*SplitPolicyRegisterParam, error) {
	res, err := this.bc.PreInvoke(this.contractAddress, "getRegisterParam", []interface{}{key})
	if err != nil {
		return nil, err
	}
	data, err := res.ToByteArray()
	if err != nil {
		return nil, err
	}
	spp := &SplitPolicyRegisterParam{}
	err = spp.FromBytes(data)
	if err != nil {
		return nil, err
	}
	return spp, nil
}
