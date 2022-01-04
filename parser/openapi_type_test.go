package parser

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestToJson(t *testing.T) {
	str := NewString("hi")
	str.Name = "name"
	array := NewArray()
	array.AddProp(str)

	subObj := NewObject()
	subObj.Props["name"] = str
	subObj.Props["list"] = array

	p := NewObject()
	p.Props["mans"] = array
	p.Props["sub"] = subObj
	p.Props["name"] = str
	b, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}
