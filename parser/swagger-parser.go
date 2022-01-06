package parser

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

type API struct {
	RelativePath string
	Method       string

	BodyParams  *Object
	PathParams  *Object
	QueryParams *Object
	ParamsRaw   openapi3.Parameters // 用来debug
}

/*
map[string]interface{}
{
}
[1,2,3]
""

核心：通过json输出各种各样的case
问题，path，query的这种怎么解决？
* 把path或者api这个参数带下去，为每个类型实现一个set方法

测试用例怎么嵌入?

测试用例描述方式
 req.man.id = [1,2,3]

obj = find(req,man.id)
for case range idsInput:
obj = case.Value
*/

func NewAPI() *API {
	return &API{
		BodyParams:  NewObject(),
		PathParams:  NewObject(),
		QueryParams: NewObject(),
	}
}

// Parse ... apiFile 相对路径
func Parse(apiFile string) []*API {
	ctx := context.Background()
	loader := &openapi3.Loader{Context: ctx}
	doc, err := loader.LoadFromFile(apiFile)
	if err != nil {
		panic(err)
	}
	APIs := make([]*API, 0, 128)
	addNewAPI := func(api *API, relativePath string) {
		if api == nil {
			return
		}
		api.RelativePath = relativePath
		APIs = append(APIs, api)
	}
	for path, pathIten := range doc.Paths {
		api := GenerateAPI(http.MethodConnect, pathIten.Connect)
		addNewAPI(api, path)

		api = GenerateAPI(http.MethodDelete, pathIten.Delete)
		addNewAPI(api, path)

		api = GenerateAPI(http.MethodGet, pathIten.Get)
		addNewAPI(api, path)

		api = GenerateAPI(http.MethodHead, pathIten.Head)
		addNewAPI(api, path)

		api = GenerateAPI(http.MethodOptions, pathIten.Options)
		addNewAPI(api, path)

		api = GenerateAPI(http.MethodPatch, pathIten.Patch)
		addNewAPI(api, path)

		api = GenerateAPI(http.MethodPost, pathIten.Post)
		addNewAPI(api, path)

		api = GenerateAPI(http.MethodPut, pathIten.Put)
		addNewAPI(api, path)
	}
	return APIs
}

func GenerateAPI(method string, operation *openapi3.Operation) *API {
	if operation == nil {
		return nil
	}
	api := NewAPI()
	api.Method = method
	api.OperationHandler(operation)
	return api
}

func (a *API) OperationHandler(item *openapi3.Operation) {
	a.ParamsRaw = item.Parameters

	// 处理一个个param
	for _, paramRef := range item.Parameters {
		a.ParameterRefHandler(paramRef)
	}
}

func (a *API) ParameterRefHandler(ref *openapi3.ParameterRef) {
	if ref.Value == nil {
		return
	}
	rawParameter := ref.Value

	setPropInfo := func(p *Object) {
		p.In = rawParameter.In
		p.Description = rawParameter.Description
	}

	switch rawParameter.In {
	case "body":
		setPropInfo(a.BodyParams)
		newObj := SchemaRefHandler(rawParameter.Schema).(*Object)
		for k, v := range newObj.Props {
			a.BodyParams.Props[k] = v
		}
	case "query":
		setPropInfo(a.QueryParams)
		newParamName := rawParameter.Name
		newParam := ExtensionPropsHandler(rawParameter.ExtensionProps)
		if _, isObject := newParam.(*Object); isObject {
			panic("这宗情况需要处理一下，结构体参数出现在这两种类型中")
		}
		a.QueryParams.Props[newParamName] = newParam
	case "path":
		setPropInfo(a.PathParams)
		newParamName := rawParameter.Name
		newParam := ExtensionPropsHandler(rawParameter.ExtensionProps)
		if _, isObject := newParam.(*Object); isObject {
			panic("这宗情况需要处理一下，结构体参数出现在这两种类型中")
		}
		a.PathParams.Props[newParamName] = newParam
	default:
		log.Printf("%+v\n", rawParameter)
		panic("unexpect case")
	}
}

func ExtensionPropsHandler(extendProps openapi3.ExtensionProps) Prop {
	hasExtensionProps := len(extendProps.Extensions) != 0
	if !hasExtensionProps {
		return nil
	}

	msg, ok := extendProps.Extensions["type"].(json.RawMessage)
	if !ok {
		panic(`not json raw message`)
	}
	typeStringBytes, err := msg.MarshalJSON()
	if err != nil {
		panic(err)
	}
	t := strings.Trim(string(typeStringBytes), "\"")
	type Items struct {
		Type string `json:"type"`
	}

	switch t {
	case String_T:
		return NewString("")
	case Array_T:
		items := Items{}
		msg := extendProps.Extensions["items"].(json.RawMessage)
		b, _ := msg.MarshalJSON()
		json.Unmarshal(b, &items)
		arr := NewArray()
		var element Prop
		switch items.Type {
		case String_T:
			element = NewString("")
		case Int_T:
			element = NewInt(0)
		case Bool_T:
			element = NewBool(false)
		default:
			panic("not type:" + t)
		}
		arr.AddProp(element)
		return arr
	case Object_T:
		panic("额外参数的object需要被处理")
		return NewObject()
	case Int_T:
		return NewInt(0)
	case Bool_T:
		return NewBool(false)
	default:
		panic("not type:" + t)
	}
}

func SchemaRefHandler(ref *openapi3.SchemaRef) Prop {
	if ref == nil {
		panic("unexpect ref is nil")
	}
	schema := ref.Value
	if schema == nil {
		panic("unexpect schema is nil")
	}

	switch schema.Type {
	case Object_T:
		obj := NewObject()
		for name, subRef := range schema.Properties {
			obj.Props[name] = SchemaRefHandler(subRef)
		}
		return obj
	case Array_T:
		arr := NewArray()
		subProp := SchemaRefHandler(schema.Items)
		arr.AddProp(subProp)
		return arr
	case String_T:
		str := NewString("")
		return str
	case Bool_T:
		boolean := NewBool(false)
		return boolean
	case Int_T:
		integer := NewInt(0)
		return integer
	case Number_T:
		number := NewNumber(0.0)
		return number
	default:
		panic(`type not support yet`)
	}
}
