package ddxf_contract

import (
	"github.com/ontio/ontology/common"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/param"
)

type DataMetaInfo struct {
	DataMeta     map[string]interface{} `json:"dataMeta"`
	DataMetaHash string                 `json:"dataMetaHash"`
	ResourceType byte                   `json:"resourceType"`
	Fee          param.Fee              `json:"fee"`
	Stock        uint32                 `json:"stock"`
	ExpiredDate  uint64                 `json:"expiredDate"`
	DataEndpoint string                 `json:"dataEndpoint"`
	DataHash     string                 `json:"dataHash"`
	DataId       string                 `json:"dataId"`
}

type TokenMetaInfo struct {
	TokenMeta     map[string]interface{} `json:"tokenMeta"`
	TokenMetaHash string                 `json:"tokenMetaHash"`
	DataMetaHash  string                 `json:"dataMetaHash"`
	TokenEndpoint string                 `json:"tokenEndpoint"`
}

type ProductInfoOnChain struct {
	ResourceDdo *param.ResourceDDO
	DtokenItem  *param.DTokenItem
}

type ResourceIdAndN struct {
	ResourceId []byte
	N          int
}

func (this *ProductInfoOnChain) Deserialize(source *common.ZeroCopySource) error {
	ddo := &param.ResourceDDO{}
	err := ddo.Deserialize(source)
	if err != nil {
		return err
	}
	item := &param.DTokenItem{}
	err = item.Deserialize(source)
	if err != nil {
		return err
	}
	this.ResourceDdo = ddo
	this.DtokenItem = item
	return nil
}

func (this *ProductInfoOnChain) FromBytes(data []byte) error {
	source := common.NewZeroCopySource(data)
	return this.Deserialize(source)
}
