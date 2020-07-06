# ddxf-sdk


This sdk is used to interact with the [ddxf contract suite](https://github.com/ont-bizsuite/ddxf-contract-suite).




// addr for node address
sdk := NewDdxfSdk(addr string)
sdk.DefMpKit().Publish(seller *ontology_go_sdk.Account, resourceId []byte, ddo ResourceDDO, item DTokenItem,
	splitPolicyParam split_policy_contract.SplitPolicyRegisterParam)


sdk.DefMpKit().Update(seller *ontology_go_sdk.Account, resourceId []byte, ddo ResourceDDO, item DTokenItem,
	splitPolicyParam split_policy_contract.SplitPolicyRegisterParam)

sdk.DefMpKit().Delete(seller *ontology_go_sdk.Account, resourceId []byte)



sdk.DefMpKit().BuyDtoken(buyer, payer *ontology_go_sdk.Account, resourceId []byte,
	n int) 