package temp

var autoTest = `# coding=utf8

import copy
import os
import yaml

{{ range $idx,$case := .TestCases }}
	{{ $case.CaseName }} = [{{ $case.CaseValue }}]
{{ end }}

appid = 1
user.appid = 2

template = {
   "api": { "alias_name": "app_detail", "method": "GET", "sub_url": "console/api/v1/app/${appid}/detail" },
   "case": "参数：appid - 正常参数",
   "body":{},
   "level": "P2",
   "url_input": {
       "x-request-id":"",
       "appid":""
   },
   "params":{},
   "status_code":200,
   "business_code":0,
   "reason":"OK"
}

class {{.APINameCamo}}Generator:
   def __init__(self):
       curPath = os.path.abspath(os.path.dirname(__file__))
       self.saveYamlFile = os.path.join(curPath,"../../../../cases/spcdn/cdn_console/App/case_data/{{.APIName}}.yaml")
   def gen_case(self):
       arr = []
       cnt = 0

		{{ range $idx,$case := .TestCases }}
       for appid in {{$case.CaseName}}:
// 把这个返回值放到你想要的地方 {{firstSubString $case.CaseName}}
// 也可以用另外一个函数的返回值来做个循环
			{{ range $i,$value := subStrings $case.CaseName }}
				{{ $value }}
			{{ end }}

// end
           cnt += 1
           tmp = copy.deepcopy(template)
           tmp['url_input']['appid'] = appid
           tmp['case'] = "case-%d: appid-合法参数-"%(cnt) + str(appid)
           tmp['status_code'] = 200
           tmp['business_code'] = 0
           tmp['reason'] = "ok"
           tmp['level'] = "P1"
           arr.append(tmp)

		{{ end }}
       with open(self.saveYamlFile, "w", ) as f:
           print(self.saveYamlFile)
           yaml.dump(arr, f, allow_unicode=True)

if __name__ == "__main__":
   {{.APINameCamo}}Generator().gen_case()
`
