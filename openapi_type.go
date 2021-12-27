package main

import (
	"encoding/json"
	"strings"
)

const (
	String_T = "string"
	Object_T = "object"
	Array_T = "array"
)

type Prop interface {
	DefaultValue() interface{}
	GetType() string
	GetIn() string // todo query,path args, indict to replace api path
	json.Marshaler
}

type BasicProp struct {
	APIBelong *API
	PropMap map[string]interface{}
	In string
	Name string
	Description string
	Type string
	IsRequired bool// todo 可以获取它是否是需要的参数
}

func (b BasicProp) GetIn() string {
	return b.In
}

// todo
func (b BasicProp) Set(ArgName string, TestCase Prop) {
	args := strings.Split(ArgName,".")
	if len(args) <= 0 {
		panic(`没有这个参数，检查输入路径:`+ArgName)
	}
	argToSet := args[0]
	if b.PropMap == nil {
		return
	}
	if _, ok := b.PropMap[argToSet] ;!ok{
		if !ok {
			return
		}
	}
	b.PropMap[argToSet] = TestCase
}

func (b BasicProp) GetType() string {
	return b.Type
}

func (b BasicProp) DefaultValue() interface{}{
	panic("must be over write")
}

func (b BasicProp) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.PropMap)
}

type Object struct {
	BasicProp
}

type Array struct {
	BasicProp
}

type String struct {
	BasicProp
}

func (str *String) DefaultValue() interface{} {
	return ""
}