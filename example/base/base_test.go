package base

import (
	"github.com/ontio/ontology/common"
	"testing"
	"fmt"
	"encoding/hex"
)

func TestSigner_Serialize(t *testing.T) {
	rp := RegIdParam{
		Ontid: []byte("123"),
		Group: Group{
			Members:   [][]byte{[]byte("123")},
			Threshold: 1,
		},
		Signer: []Signer{
			Signer{
				Id:    []byte("123"),
				Index: uint32(123),
			},
		},
		Attributes: []DDOAttribute{
			DDOAttribute{
				Key:       []byte("key"),
				Value:     []byte("value"),
				ValueType: []byte("ty"),
			},
		},
	}
	sink := common.NewZeroCopySink(nil)
	rp.Serialize(sink)
	fmt.Println(hex.EncodeToString(sink.Bytes()))
}
