package dtoken_contract

import (
	"fmt"
	"github.com/ontio/ddxf-sdk/base_contract"
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

func (this *DTokenKit) SetDDXFContractAddr(admin *ontology_go_sdk.Account, addr common.Address) (common.Uint256, error) {
	return this.bc.Invoke(this.contractAddress, admin, "setDdxfContract",
		[]interface{}{addr})
}

func (this *DTokenKit) GetDDXFContractAddr() (common.Address, error) {
	res, err := this.bc.PreInvoke(this.contractAddress, "getDdxfContract", []interface{}{})
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
