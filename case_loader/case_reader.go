package case_loader

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

type TestCases map[string]*TestCase

type TestCase struct {
	Name         string
	Method       string                    `yaml:"method"`
	RelativePath string                    `yaml:"relative_path"`
	Parameters   map[string]*ParameterCase `yaml:"parameters"`
}

type ParameterCase struct {
	Default   interface{}              `yaml:"default"`
	TestCases map[string][]interface{} `yaml:"test_cases"`
}

func ReadCasesFiles() TestCases {
	dir := "cases"
	cases, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	testCases := TestCases{}
	for _, caseFile := range cases {
		apiFilePath := filepath.Join(dir, caseFile.Name())
		bytes, err := ioutil.ReadFile(apiFilePath)
		if err != nil {
			panic(err)
		}
		partCases := TestCases{}
		err = yaml.Unmarshal(bytes, &partCases)
		if err != nil {
			panic(err)
		}
		for k, v := range partCases {
			testCases[k] = v
		}
		fmt.Println(`[ok] read test cases from `, apiFilePath)
	}
	return testCases
}
