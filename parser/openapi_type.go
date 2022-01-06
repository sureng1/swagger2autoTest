package parser

import (
	"encoding/json"
	"errors"
	"fmt"
)

const (
	String_T = "string"
	Object_T = "object"
	Array_T  = "array"
	Int_T    = "integer"
	Bool_T   = "boolean"
	Number_T = "number"
)

type Prop interface {
	DefaultValue() interface{}
	GetType() string
	GetIn() string // todo query,path args, indict to replace api path
	GetProp(name string) (Prop, error)
	GetName() string
	SetValue(value interface{})
	json.Marshaler
}

type BasicProp struct {
	APIBelong   *API
	In          string
	Name        string
	Description string
	Type        string
	IsRequired  bool // todo 可以获取它是否是需要的参数
	Value       interface{}
}

func (b *BasicProp) GetIn() string {
	return b.In
}

func (b *BasicProp) SetValue(value interface{}) {
	b.Value = value
}

func (b *BasicProp) GetType() string {
	return b.Type
}

func (b *BasicProp) GetName() string {
	return b.Name
}

func (b *BasicProp) GetProp(name string) (Prop, error) {
	return nil, errors.New(`no subProperty for this type`)
}

func (b *BasicProp) DefaultValue() interface{} {
	panic("must be over write")
}

func (b *BasicProp) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.Value)
}

func NewBasicProp() *BasicProp {
	return &BasicProp{}
}

type Object struct {
	*BasicProp
	Props map[string]Prop
}

func (o *Object) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.Props)
}

func (o *Object) GetProp(name string) (Prop, error) {
	prop, exists := o.Props[name]
	if !exists {
		return nil, errors.New(fmt.Sprintf(`property %s not found in object type`, name))
	}
	return prop, nil
}

func NewObject() *Object {
	b := NewBasicProp()
	m := map[string]Prop{}
	return &Object{
		BasicProp: b,
		Props:     m,
	}
}

type Array struct {
	*BasicProp
	PropsArr []Prop
}

func (a *Array) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.PropsArr)
}

func (a *Array) AddProp(prop Prop) {
	a.PropsArr = append(a.PropsArr, prop)
}

func (a *Array) Get(name string) (Prop, error) {
	for _, prop := range a.PropsArr {
		if prop.GetName() == name {
			return prop, nil
		}
	}
	return nil, errors.New(fmt.Sprintf(`property %s not found`, name))
}

func NewArray() *Array {
	return &Array{BasicProp: NewBasicProp(), PropsArr: make([]Prop, 0)}
}

type String struct {
	*BasicProp
}

func NewString(v string) *String {
	b := NewBasicProp()
	b.Value = v
	return &String{b}
}

type Int struct {
	*BasicProp
}

func (str *Int) MarshalJSON() ([]byte, error) {
	return json.Marshal(0)
}

func NewInt(v int) *Int {
	b := NewBasicProp()
	b.Value = v
	return &Int{b}
}

type Bool struct {
	*BasicProp
}

func NewBool(v bool) *Bool {
	b := NewBasicProp()
	b.Value = v
	return &Bool{b}
}

type Number struct {
	*BasicProp
}

func NewNumber(v float64) *Number {
	b := NewBasicProp()
	b.Value = v
	return &Number{BasicProp: b}
}
