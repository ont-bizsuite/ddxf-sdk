package split_policy_contract

import (
	"fmt"
	"github.com/ontio/ontology/common"
)

// TokenType def
type TokenType byte

const (
	// ONT token
	ONT TokenType = iota
	// ONG token
	ONG
	// OEP4 token
	OEP4
	// OEP5 token
	OEP5
	// OEP8 token
	OEP8
	// OEP68 token
	OEP68
)

type AddrAmt struct {
	To      common.Address
	Percent uint32
}

func (this *AddrAmt) Serialize(sink *common.ZeroCopySink) {
	sink.WriteAddress(this.To)
	sink.WriteUint32(this.Percent)
}
func (this *AddrAmt) Deserialize(source *common.ZeroCopySource) error {
	addr, eof := source.NextAddress()
	if eof {
		return fmt.Errorf("[AddrAmt] read to failed, eof: %v", eof)
	}
	p, eof := source.NextUint32()
	if eof {
		return fmt.Errorf("[AddrAmt] read percent failed, eof: %v", eof)
	}
	this.To = addr
	this.Percent = p
	return nil
}

type SplitPolicyRegisterParam struct {
	AddrAmts     []*AddrAmt
	TokenTy      TokenType
	ContractAddr common.Address // option
}

func (this *SplitPolicyRegisterParam) ToBytes() []byte {
	sink := common.NewZeroCopySink(nil)
	this.Serialize(sink)
	return sink.Bytes()
}

func (this *SplitPolicyRegisterParam) Serialize(sink *common.ZeroCopySink) {
	sink.WriteVarUint(uint64(len(this.AddrAmts)))
	for _, v := range this.AddrAmts {
		v.Serialize(sink)
	}
	sink.WriteByte(byte(this.TokenTy))
	if this.ContractAddr != common.ADDRESS_EMPTY {
		sink.WriteBool(true)
		sink.WriteAddress(this.ContractAddr)
	} else {
		sink.WriteBool(false)
	}
}
func (this *SplitPolicyRegisterParam) Deserialize(source *common.ZeroCopySource) error {
	l, _, irregular, eof := source.NextVarUint()
	if irregular || eof {
		return fmt.Errorf("read AddrAmts length failed,irregular: %v,eof: %v", irregular, eof)
	}
	aa := make([]*AddrAmt, l)
	for i := 0; i < int(l); i++ {
		err := aa[i].Deserialize(source)
		if err != nil {
			return err
		}
	}
	ty, eof := source.NextByte()
	if eof {
		return fmt.Errorf("[SplitPolicyRegisterParam] read TokenTy failed, eof: %v", eof)
	}
	this.TokenTy = TokenType(ty)
	boo, irregular, eof := source.NextBool()
	if irregular || eof {
		return fmt.Errorf("[SplitPolicyRegisterParam] read ContractAddr failed, irregular:%v,eof: %v", irregular, eof)
	}
	if boo {
		this.ContractAddr, eof = source.NextAddress()
		if eof {
			return fmt.Errorf("[SplitPolicyRegisterParam] read ContractAddr failed, eof: %v", eof)
		}
	}
	return nil
}

func (this *SplitPolicyRegisterParam) FromBytes(data []byte) error {
	source := common.NewZeroCopySource(data)
	return this.Deserialize(source)
}
