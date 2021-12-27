////+build ignore
package main
//
//import (
//	"fmt"
//	"os"
//	"strings"
//	"text/template"
//)
//
//var autoTest = `# coding=utf8
//
//import copy
//import os
//import yaml
//
//{{ range $idx,$case := .TestCases }}
//{{ $case.CaseName }} = {{ $case.CaseValue }}
//{{ end }}
//
//{{ .Template }}
//
//class {{.APINameCamo}}Generator:
//    def __init__(self):
//        curPath = os.path.abspath(os.path.dirname(__file__))
//        self.saveYamlFile = os.path.join(curPath,"{{ .SaveYamlPath }}")
//    def gen_case(self):
//        arr = []
//        cnt = 0
//
//		{{ range $idx,$case := .TestCases }}
//        for {{firstSubString $case.CaseName}} in {{$case.CaseName}}:
//            cnt += 1
//            tmp = copy.deepcopy(template)
//            tmp['url_input']['{{firstSubString $case.CaseName}}'] = {{firstSubString $case.CaseName}}
//            tmp['case'] = "case-%d: {{firstSubString $case.CaseName}}  -合法参数-"%(cnt) + str({{firstSubString $case.CaseName}})
//            tmp['status_code'] = 200
//            tmp['business_code'] = 0
//            tmp['reason'] = "ok"
//            tmp['level'] = "P1"
//            arr.append(tmp)
//
//		{{ end }}
//        with open(self.saveYamlFile, "w", ) as f:
//            print(self.saveYamlFile)
//            yaml.dump(arr, f, allow_unicode=True)
//
//if __name__ == "__main__":
//    {{.APINameCamo}}Generator().gen_case()
//`
//
//type Cases struct {
//	CaseName  string
//	CaseValue string
//}
//
//type Xiaobai struct {
//	TestCases   []Cases
//	Template    string
//	APIMethod   string
//	APIName     string // get_hello
//	APINameCamo string // getHello
//	APIUrl      string
//	SaveYamlPath string
//}
//
////请求方法也要是变量 get post
//
//func main() {
//	cases := []Cases{
//		{
//			CaseName:  "campaignID_200",
//			CaseValue: `[102]`,
//		},
//		//{
//		//	CaseName:  "region_200",
//		//	CaseValue: `["sg","br"]`,
//		//},
//	}
//
//	api1 := &Xiaobai{}
//	api1.TestCases = cases
//	api1.APIMethod = `GET`
//	api1.APIName = `get_strategy-reference`
//	api1.APINameCamo = `getStrategyReference`
//	api1.APIUrl = `strategy-reference`
//	api1.SaveYamlPath = `../../../../cases/fcst/Overviwe/case_data/get_strategy-reference.yaml`
//	s :=  fmt.Sprintf(`template = {
//    "api": { "alias_name": '%s', "method": '%s', "sub_url": "/platform/api/v1/%s" },
//    "case": "参数：campaignID - 正常参数",
//    "body":{},
//    "level": "P2",
//    "url_input": {
//		"campaignID": 102,
//	},
//    "params":{},
//    "status_code":200,
//    "business_code":0,
//    "reason":"OK"
//}`,api1.APIName,api1.APIMethod,api1.APIUrl)
//	api1.Template = s
//
//	/*
//		定义你希望的函数，然后注册到模板中，
//		在模板中调用这个函数
//	*/
//	firstSubStringFunc := func(s string) string {
//		return strings.Split(s, "_")[0]
//	}
//	subStringsFunc := func(s string) []string {
//		return strings.Split(s, "_")
//	}
//	funcs := map[string]interface{}{
//		"firstSubString": firstSubStringFunc,
//		"subStrings":     subStringsFunc,
//	}
//
//	t, err := template.New(`j`).Funcs(funcs).Parse(autoTest)
//	if err != nil {
//		panic(err)
//	}
//	file, err := os.Create(fmt.Sprintf(`/Users/zhesun/Documents/gitRepository/mt-autotest/controller/fcst/Overviwe/%s.py`, api1.APIName))
//	if err != nil {
//		panic(err)
//	}
//	err = t.Execute(file, api1)
//	if err != nil {
//		panic(err)
//	}
//}
