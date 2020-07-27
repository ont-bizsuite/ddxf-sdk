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
func (this *Signer) Deserialize(source *common.ZeroCopySource) error {
	id, _, irr, eof := source.NextVarBytes()
	if irr || eof {
		return fmt.Errorf("irr: %v, eof:%v", irr, eof)
	}
	ind, eof := source.NextUint32()
	if eof {
		return fmt.Errorf("eof: %v", eof)
	}
	this.Index = ind
	this.Id = id
	return nil
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
func (this *Group) Deserialize(source *common.ZeroCopySource) error {
	l, _, irr, eof := source.NextVarUint()
	if irr || eof {
		return fmt.Errorf("irr: %v, eof: %v", irr, eof)
	}
	members := make([][]byte, 0)
	for i := 0; i < int(l); i++ {
		bs, _, irr, eof := source.NextVarBytes()
		if irr || eof {
			return fmt.Errorf("irr: %v, eof: %v", irr, eof)
		}
		members = append(members, bs)
	}
	th, eof := source.NextUint32()
	if eof {
		return fmt.Errorf("eof: %v", eof)
	}
	this.Members = members
	this.Threshold = uint(th)
	return nil
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

func (this *RegIdParam) Deserialize(source *common.ZeroCopySource) error {
	ontId, _, irr, eof := source.NextVarBytes()
	if irr || eof {
		return fmt.Errorf("irr: %v,eof: %v", irr, eof)
	}
	err := this.Group.Deserialize(source)
	if err != nil {
		return err
	}
	l, _, irr, eof := source.NextVarUint()
	if irr || eof {
		return fmt.Errorf("irr: %v,eof: %v", irr, eof)
	}
	for i := 0; i < int(l); i++ {
		signer := &Signer{}
		err = signer.Deserialize(source)
		if err != nil {
			return err
		}
		this.Signer = append(this.Signer, *signer)
	}
	l, _, irr, eof = source.NextVarUint()
	if irr || eof {
		return fmt.Errorf("irr: %v,eof: %v", irr, eof)
	}
	for i := 0; i < int(l); i++ {
		attri := &DDOAttribute{}
		err = attri.Deserialize(source)
		if err != nil {
			return err
		}
		this.Attributes = append(this.Attributes, *attri)
	}
	this.Ontid = ontId
	return nil
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
		return fmt.Errorf("parse key failed, irr: %v, eof: %v", irr, eof)
	}
	value, _, irr, eof := source.NextVarBytes()
	if irr || eof {
		return fmt.Errorf("parse value failed, irr: %v, eof: %v", irr, eof)
	}
	valueTy, _, irr, eof := source.NextVarBytes()
	if irr || eof {
		return fmt.Errorf("parse valueType failed, irr: %v, eof: %v", irr, eof)
	}
	this.Key = key
	this.Value = value
	this.ValueType = valueTy
	return nil
}
