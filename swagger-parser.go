package main

import (
	"context"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"net/http"
)

type API struct {
	RelativePath string
	Method string

	In string // 参数的传递方式
	Prop Prop
	Params openapi3.Parameters // 用来debug
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

func (a *API)AddBasicProp(prop *BasicProp ) {
	a.BasicProps = append(a.BasicProps, prop)
}

func main() {
	ctx := context.Background()
	loader := &openapi3.Loader{Context: ctx}
	doc, err := loader.LoadFromFile("fcst-platform.json")
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
		api := GenerateAPI(http.MethodConnect,pathIten.Connect)
		addNewAPI(api, path)

		api = GenerateAPI(http.MethodDelete,pathIten.Delete)
		addNewAPI(api, path)

		api = GenerateAPI(http.MethodGet,pathIten.Get)
		addNewAPI(api, path)

		api = GenerateAPI(http.MethodHead,pathIten.Head)
		addNewAPI(api, path)

		api = GenerateAPI(http.MethodOptions,pathIten.Options)
		addNewAPI(api, path)

		api = GenerateAPI(http.MethodPatch,pathIten.Patch)
		addNewAPI(api, path)

		api = GenerateAPI(http.MethodPost,pathIten.Post)
		addNewAPI(api, path)

		api = GenerateAPI(http.MethodPut,pathIten.Put)
		addNewAPI(api, path)
	}
	for _,api := range APIs {
		fmt.Printf("%+v\n",api)
	}
}

func GenerateAPI(method string, operation *openapi3.Operation) *API {
	if operation == nil {
		return nil
	}
	api := &API{Method: method}
	api.OperationHandler(operation)
	return api
}

func (a *API)OperationHandler(item *openapi3.Operation) {
	a.Params = item.Parameters

	for _, paramRef := range item.Parameters {
		param := paramRef.Value
		// 单层无嵌套
		if len(param.ExtensionProps.Extensions) > 0{
			t := fmt.Sprintf("%s",param.ExtensionProps.Extensions["type"])
			prop := CreatePropWithType(t)
			prop.
			prop.Name:        param.Name,
				In:          param.In,
				Description: param.Description,
				Parents:     nil,
				Type:    ,
			}
			// todo prop.IsRequired
			a.AddBasicProp(prop)
		}
		parent := []string{param.Name}
		a.SchemaRefHandler(param.Schema, parent, "")
	}
}


func (a *API)SchemaRefHandler(ref *openapi3.SchemaRef, parent []string, firstParentType string) map[string]interface{}{
	if ref == nil {
		return
	}
	schema := ref.Value
	if schema.Type != Object_T && schema.Type != Array_T {
		// 没有嵌套的参数, 此时只有一个基本类型，"integer"，"string"
		name := parent[len(parent)-1]
		parent = parent[:len(parent)-1]
		prop := &BasicProp{
			Name:        name,
			In:          "body",
			Parents:     parent,
			PropType:    schema.Type,
		}
		a.AddBasicProp(prop)
		return
	}
	for key, prop := range schema.Properties { // properties
		tempParent := append(parent, key)
		a.SchemaRefHandler(prop, tempParent, schema.Type)
	}
	a.SchemaRefHandler(schema.Items, parent, schema.Type)
	a.SchemaRefHandler(schema.AdditionalProperties, parent, schema.Type)
}

func CreatePropWithType(t string) Prop {
	switch t {
	case String_T:
		return &String{}
	case Array_T:
		return &Array{}
	case Object_T:
		return &Object{}
	}
}