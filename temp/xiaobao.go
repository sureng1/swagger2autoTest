package temp

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
//	{{ $case.CaseName }} = [{{ $case.CaseValue }}]
//{{ end }}
//
//{{ .Template }}
//
//class {{.APINameCamo}}Generator:
//    def __init__(self):
//        curPath = os.path.abspath(os.path.dirname(__file__))
//        self.saveYamlFile = os.path.join(curPath,"../../../../cases/spcdn/cdn_console/App/case_data/{{.APIName}}.yaml")
//    def gen_case(self):
//        arr = []
//        cnt = 0
//
//		{{ range $idx,$case := .TestCases }}
//        for appid in {{$case.CaseName}}:
//// 把这个返回值放到你想要的地方 {{firstSubString $case.CaseName}}
//// 也可以用另外一个函数的返回值来做个循环
//			{{ range $i,$value := subStrings $case.CaseName }}
//				{{ $value }}
//			{{ end }}
//
//// end
//            cnt += 1
//            tmp = copy.deepcopy(template)
//            tmp['url_input']['appid'] = appid
//            tmp['case'] = "case-%d: appid-合法参数-"%(cnt) + str(appid)
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
//	APIName     string // get_hello
//	APINameCamo string // getHello
//}
//
//func main() {
//	cases := []Cases{
//		{
//			CaseName:  "apiid200_class",
//			CaseValue: `[-1,0,1.3,'a','$%','1_2&','👌']`,
//		},
//		{
//			CaseName:  "apiid201_class",
//			CaseValue: `[-1,0,1.3,'a','$%','1_2&','👌']`,
//		},
//		{
//			CaseName:  "apiid500_class",
//			CaseValue: `[-1,0,1.3,'a','$%','1_2&','👌']`,
//		},
//	}
//
//	api1 := &Xiaobai{}
//	api1.TestCases = cases
//	api1.Template = `template = {
//    "api": { "alias_name": "app_detail", "method": "GET", "sub_url": "console/api/v1/app/${appid}/detail" },
//    "case": "参数：appid - 正常参数",
//    "body":{},
//    "level": "P2",
//    "url_input": {
//        "x-request-id":"",
//        "appid":""
//    },
//    "params":{},
//    "status_code":200,
//    "business_code":0,
//    "reason":"OK"
//}`
//	api1.APIName = `get_xiao_hei`
//	api1.APINameCamo = `getXiaoHEI`
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
//	t, err := template.New(`j`).Funcs(funcs).Parse(autoTest)
//	if err != nil {
//		panic(err)
//	}
//	file, err := os.Create(fmt.Sprintf(`/tmp/%s.py`, api1.APIName))
//	if err != nil {
//		panic(err)
//	}
//	err = t.Execute(file, api1)
//	if err != nil {
//		panic(err)
//	}
//}
