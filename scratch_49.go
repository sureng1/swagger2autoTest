package main

//import (
//	"fmt"
//	"os"
//	"text/template"
//)
//
//var autoTest = `# coding=utf8
//
//import copy
//import os
//import yaml
//
//{{ range $idx,$case := .TestCases }} #什么意思
//	{{ $case.CaseName }} = {{ $case.CaseValue }}
//{{ end }}
//
//{{ .Template }}
//
//class {{.APINameCamo}}Generator:
//   def __init__(self):
//       curPath = os.path.abspath(os.path.dirname(__file__))
//       self.saveYamlFile = os.path.join(curPath,"../../../../cases/fcst/Overviwe/case_data/{{.APIName}}.yaml")
//   def gen_case(self):
//       arr = []
//       cnt = 0
//
//		{{ range $idx,$case := .TestCases }}
//       for appid in {{$case.CaseName}}:  # appid应该是变量，与CaseName对应
//           cnt += 1
//           tmp = copy.deepcopy(template)
//           tmp['url_input']['appid'] = appid
//           tmp['case'] = "case-%d: appid-合法参数-"%(cnt) + str(appid)
//           tmp['status_code'] = 200
//           tmp['business_code'] = 0
//           tmp['reason'] = "ok"
//           tmp['level'] = "P1"
//           arr.append(tmp)
//
//		{{ end }}
//       with open(self.saveYamlFile, "w", ) as f:
//           print(self.saveYamlFile)
//           yaml.dump(arr, f, allow_unicode=True)
//
//if __name__ == "__main__":
//   {{.APINameCamo}}Generator().gen_case()
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
//	APIName     string // get_hello
//	APINameCamo string // getHello
//}
//
//func main() {
//	cases := []Cases{
//		{
//			CaseName:  "campaignID_200",
//			CaseValue: `[102]`,
//		},
//		{
//			CaseName:  "behavior_200",
//			CaseValue: `["Place Order","Browse Microsite"]`,
//		},
//	}
//
//	api1 := &Xiaobai{}
//	api1.TestCases = cases
//	api1.Template = `template = {
//   "api": { "alias_name": "target_list", "method": "GET", "sub_url": "/platform/api/v1/target" },
//   "case": "参数：campaignID - 正常参数",
//   "body":{},
//   "level": "P2",
//   "url_input": {
//       "x-request-id":"",
//       "appid":""
//   },
//   "params":{},
//   "status_code":200,
//   "business_code":0,
//   "reason":"OK"
//}`
//	api1.APIName = `get_target`
//	api1.APINameCamo = `getTarget`
//
//	t, err := template.New(`j`).Parse(autoTest)
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
