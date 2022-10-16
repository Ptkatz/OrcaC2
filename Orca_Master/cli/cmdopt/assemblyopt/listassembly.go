package assemblyopt

import (
	"Orca_Master/cli/common"
	"Orca_Master/define/colorcode"
	"encoding/json"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type AssemblyYaml struct {
	AssemblyStructs []AssemblyStruct `yaml:"assembly"`
}

type AssemblyStruct struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Author      string `yaml:"author"`
	Loaded      string
}

func ReadYamlFile(yamlPath string) AssemblyYaml {
	var assemblyYaml AssemblyYaml
	yamlFile, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		return AssemblyYaml{}
	}
	err = yaml.UnmarshalStrict(yamlFile, &assemblyYaml)
	if err != nil {
		return AssemblyYaml{}
	}
	return assemblyYaml
}

func PrintTable(assemblyStructs []AssemblyStruct) {
	var data [][]string
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"id", "name", "description", "author", "loaded"})
	table.SetBorders(tablewriter.Border{Left: false, Top: false, Right: false, Bottom: false})
	table.SetColWidth(48)
	for i, assemblyStruct := range assemblyStructs {
		data = append(data, []string{strconv.Itoa(i + 1), assemblyStruct.Name, assemblyStruct.Description, assemblyStruct.Author, colorcode.Colorf(colorcode.COLOR_PURPLE, assemblyStruct.Loaded)})
	}

	for _, raw := range data {
		table.Append(raw)
	}
	table.Render()
}

func SettleLoadedAssembly(msg string) []AssemblyStruct {
	var assemblyNames []string
	// 打印程序集列表
	exist := false
	yamlFile, _ := filepath.Abs("3rd_party/windows/csharp/assembly.yaml")
	assemblyYaml := ReadYamlFile(yamlFile)
	_, _, data := common.SettleRetDataBt(msg)
	json.Unmarshal(data, &assemblyNames)
	for i, _ := range assemblyYaml.AssemblyStructs {
		for _, name := range assemblyNames {
			if strings.ToLower(assemblyYaml.AssemblyStructs[i].Name+".exe") == name {
				assemblyYaml.AssemblyStructs[i].Loaded = "loaded"
				exist = true
			}
		}
	}
	if !exist {
		for i, name := range assemblyNames {
			n := len(assemblyYaml.AssemblyStructs)
			x := n + i - 1
			assemblyYaml.AssemblyStructs[x].Name = name
			assemblyYaml.AssemblyStructs[x].Author = "Unknown"
			assemblyYaml.AssemblyStructs[x].Loaded = "loaded"
		}
	}
	return assemblyYaml.AssemblyStructs
}
