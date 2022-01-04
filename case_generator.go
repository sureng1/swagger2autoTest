package main

import (
	"log"
	"producerPy/case_loader"
	"producerPy/parser"
	"strings"
)

func main() {
	testCases := case_loader.ReadCasesFiles()
	apiList := parser.Parse("parser/fcst-platform.json")
	for name, testCase := range testCases {
		targetApi := FindApi(apiList, testCase)
		if targetApi == nil {
			panic(`[error] can not find the api you defied, name:` + name + `method:` + testCase.Method)
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

func generateRequestCase(api *parser.API, testCase *case_loader.TestCase) error {
	for _, param := range api.Params {
		for name, value := range testCase.Parameters {
			setProp(param.Props, name, value.Default)
		}
	}

	for _, param := range api.Params {
		for propNameTestCase, paramTestCase := range testCase.Parameters {
			for caseName, caseValue := range paramTestCase.TestCases {
				parent := findPropParent(param.Props, propNameTestCase)
				propName := getPropName(propNameTestCase)
				j
				defauleValue := parent.Props[propName]
				parent.Props[propName] = parser.NewOverwrite(caseValue)
				generate(api, caseName)
				parent.Props[propName] = defauleValue
			}
		}
	}
}

func setProp(prop parser.Prop, key string, value interface{}) bool {
	_, ok := setPropAndGetOld(prop, key, value)
	return ok
}

func setPropAndGetOld(prop parser.Prop, key string, value interface{}) (parser.Prop, bool) {
	if value == nil {
		return nil, false
	}
	parent := findPropParent(prop, key)
	propName := getPropName(key)
	old := parent.Props[propName]
	parent.Props[propName] = parser.NewOverwrite(value)
	return old, true
}

func findPropParent(prop parser.Prop, key string) *parser.Object {
	props := strings.Split(key, ".")
	// find prop to modify
	temp := prop
	var exists bool
	for i, propName := range props {
		if i == len(props)-1 {
			break
		}
		temp, exists = temp.(*parser.Object).Props[propName]
		if !exists {
			return nil
		}
	}
	return temp.(*parser.Object)
}

func getPropName(key string) string { // a.b.c
	props := strings.Split(key, ".")
	return props[len(props)-1]
}
