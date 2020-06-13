package main

import "github.com/ontio/ddxf-sdk"

func main() {
	sdk := ddxf_sdk.NewDdxfSdk("http://")
	sdk.DefaultDDXFContract()
}
