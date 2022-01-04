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
	Default   interface{} `yaml:"default_value"`
	TestCases []*Case     `yaml:"test_cases"`
}

type Case struct {
	CaseName     string        `yaml:"case_name"`
	StatusCode   int           `yaml:"status_code"`
	BusinessCode int           `yaml:"business_code"`
	Level        string        `yaml:"level"`
	ValueList    []interface{} `yaml:"value_list"`
}

func ReadCasesFiles(dir string) TestCases {
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
