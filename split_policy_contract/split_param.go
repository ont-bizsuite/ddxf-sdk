package split_policy_contract

import (
	"github.com/ontio/ontology/common"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/io"
)

type AddrAmt struct {
	To      common.Address
	percent uint32
}

func (this *AddrAmt) Serialize(sink *common.ZeroCopySink) {
	sink.WriteAddress(this.To)
	sink.WriteUint32(this.percent)
}

type SplitPolicy struct {
	AddrAmts []AddrAmt
	TokenTy  io.TokenType
}

func (this *SplitPolicy) ToBytes() []byte {
	sink := common.NewZeroCopySink(nil)
	this.Serialize(sink)
	sink.Bytes()
}

func (this *SplitPolicy) Serialize(sink *common.ZeroCopySink) {
	sink.WriteVarUint(uint64(len(this.AddrAmts)))
	sink.WriteByte(byte(this.TokenTy))
}
