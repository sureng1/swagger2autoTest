package main

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"producerPy/case_loader"
	"producerPy/parser"
	"strings"
)

func main() {
	testCases := case_loader.ReadCasesFiles("case_loader/cases")
	apiList := parser.Parse("parser/fcst-platform.json")
	for name, testCase := range testCases {
		targetApi := FindApi(apiList, testCase)
		if targetApi == nil {
			panic(`[error] can not find the api you defied, name: ` + name + ` method: ` + testCase.Method)
		}
		err := generateRequestCase(targetApi, testCase)
		if err != nil {
			log.Println("[failed] generator yaml failed. reason:\n", err)
		}
	}
}

func FindApi(apiList []*parser.API, testCase *case_loader.TestCase) *parser.API {
	for i, api := range apiList {
		if api.RelativePath == testCase.RelativePath && api.Method == testCase.Method {
			return apiList[i]
		}
	}
	return nil
}

type Request struct {
	API API `yaml:"api"`
	Case string `yaml:"case"`
	Body struct{} `yaml:"body"`
	Level string `yaml:"level"`
	UrlInput string `yaml:"url_input"`
	Params string `yaml:"params"`
	StatusCode int `yaml:"status_code"`
	BusinessCode int `yaml:"business_code"`
	Reason string `yaml:"reason"`
}

type API struct {
	AliasName string `yaml:"alias_name"`
	Method string `yaml:"method"`
	SubUrl string `yaml:"sub_url"`
}

func api2Request(api *parser.API, testCase *case_loader.Case) Request {
	req := Request{}
	req.API.SubUrl = api.RelativePath
	req.API.Method = api.Method
	req.API.AliasName = "todo is what?"

	req.Body = struct{}{}
	req.BusinessCode = testCase.BusinessCode
	req.Case = testCase.CaseName
	req.Level = testCase.Level

	req.Reason = "ok(todo)"
	req.StatusCode = testCase.StatusCode


	req.Params = api.Params
	allProp := map[string]parser.Prop
	for _, param := range api.Params {
		for name, prop := range param.Props
		allProp param.Props
	}
	// todo
	//req.UrlInput = 1
	//req.API.SubUrl = 1

	return req
}

func generateRequestCase(api *parser.API, testCaseRef *case_loader.TestCase) error {
	for _, param := range api.Params {
		for param2Test, paramCase := range testCaseRef.Parameters {
			parser.SetProp(param.Props, param2Test, paramCase.Default)
		}
	}

	templateList := make([]Request, 0)
	for _, param := range api.Params {
		for param2Test, paramCase := range testCaseRef.Parameters {
			for _, testCase := range paramCase.TestCases {
				for _, testValue := range testCase.ValueList {
					old, _ := parser.SetPropAndGetOld(param.Props, param2Test, testValue)
					req := api2Request(api, testCase)
					templateList = append(templateList, req)
					parser.SetProp(param.Props, param2Test, old)
				}
			}
		}
	}

	b, err := yaml.Marshal(templateList)
	if err != nil {
		return err
	}
	f, err := os.Create("cases_generated/" + strings.ReplaceAll(api.RelativePath, "/", "_") +".yaml")
	if err != nil {
		return err
	}
	_, err = f.Write(b)
	return err
}
