package parser

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	apiList := Parse("fcst-platform.json")
	fmt.Println(len(apiList))
	for _, api := range apiList {
		fmt.Println(api.Method+`      api:---------------------------------`, api.RelativePath)
		for _, param := range api.Params {
			b, err := param.Props.MarshalJSON()
			if err != nil {
				panic(err)
			}
			fmt.Println(string(b))
		}
	}
}
