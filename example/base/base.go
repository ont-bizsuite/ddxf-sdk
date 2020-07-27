package base

import (
	"fmt"
	"github.com/ontio/ontology/common"
)

type Signer struct {
	Id    []byte
	Index uint32
}

func (this *Signer) Serialize(sink *common.ZeroCopySink) {
	sink.WriteVarBytes(this.Id)
	sink.WriteUint32(this.Index)
}

type Group struct {
	Members   [][]byte
	Threshold uint
}

func (this *Group) Serialize(sink *common.ZeroCopySink) {
	sink.WriteVarUint(uint64(len(this.Members)))
	for _, item := range this.Members {
		sink.WriteVarBytes(item)
	}
	sink.WriteUint32(uint32(this.Threshold))
}

type RegIdParam struct {
	Ontid      []byte
	Group      Group
	Signer     []Signer
	Attributes []DDOAttribute
}

func (this *RegIdParam) ToBytes() []byte {
	sink := common.NewZeroCopySink(nil)
	this.Serialize(sink)
	return sink.Bytes()
}

func (this *RegIdParam) Serialize(sink *common.ZeroCopySink) {
	sink.WriteVarBytes(this.Ontid)
	this.Group.Serialize(sink)
	sink.WriteVarUint(uint64(len(this.Signer)))
	for _, signer := range this.Signer {
		signer.Serialize(sink)
	}
	sink.WriteVarUint(uint64(len(this.Attributes)))
	for _, attr := range this.Attributes {
		attr.Serialize(sink)
	}
}

type DDOAttribute struct {
	Key       []byte
	Value     []byte
	ValueType []byte
}

//invoke wam contract
func (this *DDOAttribute) Serialize(sink *common.ZeroCopySink) {
	sink.WriteVarBytes(this.Key)
	sink.WriteVarBytes(this.Value)
	sink.WriteVarBytes(this.ValueType)
}

//parse from txpayload
func (this *DDOAttribute) Deserialize(source *common.ZeroCopySource) error {
	key, _, irr, eof := source.NextVarBytes()
	if irr || eof {
		return fmt.Errorf("irr: %s, eof: %s", irr, eof)
	}
	value, _, irr, eof := source.NextVarBytes()
	if irr || eof {
		return fmt.Errorf("irr: %s, eof: %s", irr, eof)
	}
	valueTy, _, irr, eof := source.NextVarBytes()
	if irr || eof {
		return fmt.Errorf("irr: %s, eof: %s", irr, eof)
	}
	this.Key = key
	this.Value = value
	this.ValueType = valueTy
	return nil
}
