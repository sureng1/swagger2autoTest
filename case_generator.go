package main

import (
	"encoding/json"
	"log"
	"net/url"
	"os"
	"producerPy/case_loader"
	"producerPy/parser"
	"strings"

	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

func init() {
	zap.ReplaceGlobals(zap.NewExample())
}

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
	API          API      `yaml:"api"`
	Case         string   `yaml:"case"`
	Body         struct{} `yaml:"body"`
	Level        string   `yaml:"level"`
	UrlInput     string   `yaml:"url_input"`
	Params       string   `yaml:"params"`
	StatusCode   int      `yaml:"status_code"`
	BusinessCode int      `yaml:"business_code"`
	Reason       string   `yaml:"reason"`
}

type API struct {
	AliasName string `yaml:"alias_name"`
	Method    string `yaml:"method"`
	SubUrl    string `yaml:"sub_url"`
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

	toString := func(m interface{}) string {
		b, err := json.Marshal(m)
		if err != nil {
			log.Print(m)
			panic(err)
		}
		return string(b)
	}

	req.Params = toString(api.BodyParams)
	req.UrlInput = toString(api.PathParams)

	query := url.Values{}
	for queryName, queryValue := range api.QueryParams.Props {
		query.Add(queryName, toString(queryValue))
	}
	u := url.URL{}
	u.Path = api.RelativePath
	u.RawQuery = query.Encode()
	req.API.SubUrl = u.String()

	return req
}

func generateRequestCase(api *parser.API, testCaseRef *case_loader.TestCase) error {
		for param2Test, paramCase := range testCaseRef.Parameters { // 尝试负值默认值给三种类型的参数
			parser.SetProp(api.BodyParams, param2Test, paramCase.Default)
			parser.SetProp(api.PathParams, param2Test, paramCase.Default)
			parser.SetProp(api.QueryParams, param2Test, paramCase.Default)
		}

	templateList := make([]Request, 0)
	for param2Test, paramCase := range testCaseRef.Parameters {
		for _, testCase := range paramCase.TestCases {
			for _, testValue := range testCase.ValueList {
				trySetAndGenerate := func (p parser.Prop) {
					old, err := parser.SetPropAndGetOld(p, param2Test, testValue)
					if err != nil {
						return
					}
					req := api2Request(api, testCase)
					templateList = append(templateList, req)
					parser.SetProp(p, param2Test, old)
				}
				trySetAndGenerate(api.BodyParams)
				trySetAndGenerate(api.QueryParams)
				trySetAndGenerate(api.PathParams)
			}
		}
	}

	b, err := yaml.Marshal(templateList)
	if err != nil {
		return err
	}
	f, err := os.Create("cases_generated/" + strings.ReplaceAll(api.RelativePath, "/", "_") + ".yaml")
	if err != nil {
		return err
	}
	_, err = f.Write(b)
	return err
}
