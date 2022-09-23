package shellopt

import (
	"github.com/olekukonko/tablewriter"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strconv"
)

type CmdCheatSheetYaml struct {
	CmdCheatSheetStructs []CmdCheatSheetStruct `yaml:"cheatsheet"`
}

type CmdCheatSheetStruct struct {
	Cmd         string `yaml:"cmd"`
	Description string `yaml:"description"`
	Privilege   string `yaml:"privilege"`
}

func ReadYamlFile(yamlPath string) CmdCheatSheetYaml {
	var cmdCheatSheet CmdCheatSheetYaml
	yamlFile, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		return CmdCheatSheetYaml{}
	}
	err = yaml.UnmarshalStrict(yamlFile, &cmdCheatSheet)
	if err != nil {
		return CmdCheatSheetYaml{}
	}
	return cmdCheatSheet
}

func PrintTable(cmdCheatSheetStructs []CmdCheatSheetStruct) {
	var data [][]string
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"id", "cmd", "description", "privilege"})
	table.SetBorders(tablewriter.Border{Left: false, Top: false, Right: false, Bottom: false})
	table.SetColWidth(48)
	for i, cmdCheatSheetStruct := range cmdCheatSheetStructs {
		data = append(data, []string{strconv.Itoa(i + 1), cmdCheatSheetStruct.Cmd, cmdCheatSheetStruct.Description, cmdCheatSheetStruct.Privilege})
	}

	for _, raw := range data {
		table.Append(raw)
	}
	table.Render()
}
