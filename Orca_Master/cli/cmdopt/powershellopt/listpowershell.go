package powershellopt

import (
	"Orca_Master/cli/common"
	"Orca_Master/define/colorcode"
	"encoding/json"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type PowershellYaml struct {
	PowershellStructs []PowershellStruct `yaml:"powershell"`
}

type PowershellStruct struct {
	Name        string                `yaml:"name"`
	File        string                `yaml:"file"`
	Arch        string                `yaml:"arch"`
	Privilege   string                `yaml:"privilege"`
	Author      string                `yaml:"author"`
	Description string                `yaml:"description"`
	Cmd         string                `yaml:"cmd"`
	Options     []PowershellCmdOption `yaml:"options"`
	Loaded      bool
}

type PowershellCmdOption struct {
	Name        string `yaml:"name"`
	Value       string `yaml:"value"`
	Required    bool   `yaml:"required"`
	Flag        string `yaml:"flag"`
	Description string `yaml:"description"`
}

type PowershellLoaded struct {
	Name   string
	File   string
	Loaded bool
}

func InitPowershellLoaded(powershellYaml PowershellYaml) []PowershellLoaded {
	var powershellLoadeds []PowershellLoaded
	for _, ps := range powershellYaml.PowershellStructs {
		powershellLoadeds = append(powershellLoadeds, PowershellLoaded{
			Name:   ps.Name,
			File:   ps.File,
			Loaded: false,
		})
	}
	return powershellLoadeds
}

func GetPowershellFileName(name string, powershellStructs []PowershellStruct) string {
	for _, powershellStruct := range powershellStructs {
		if name == powershellStruct.Name {
			return powershellStruct.File
		}
	}
	return ""
}

func GetPowershellLoadeds(message string) []PowershellLoaded {
	var retData common.RetData
	var powershellLoadeds []PowershellLoaded
	err := json.Unmarshal([]byte(message), &retData)
	data := retData.Data.(string)
	err = json.Unmarshal([]byte(data), &powershellLoadeds)
	if err != nil {
		return nil
	}
	return powershellLoadeds
}

func ReadYamlFile(yamlPath string) PowershellYaml {
	var powershellYaml PowershellYaml
	yamlFile, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		return PowershellYaml{}
	}
	err = yaml.UnmarshalStrict(yamlFile, &powershellYaml)
	if err != nil {
		return PowershellYaml{}
	}
	return powershellYaml
}

func PowershellCmdParse(powershellYaml PowershellYaml, initCmd, url string) string {
	var baseOption = make(map[string]string)
	initCmd = strings.TrimSpace(initCmd)
	sa := strings.Split(initCmd, " ")
	var powershellStruct PowershellStruct
	cmd := ""
	for i, ps := range powershellYaml.PowershellStructs {
		if ps.Name == sa[0] {
			powershellStruct = powershellYaml.PowershellStructs[i]
			cmd = powershellStruct.Cmd
			break
		}
	}
	if cmd == "" {
		colorcode.PrintMessage(colorcode.SIGN_ERROR, "powershell script not find")
		return ""
	}

	r, _ := regexp.Compile(`{{\w+.\w+}}`)
	subStr := r.FindAllString(cmd, -1)
	for _, option := range subStr {
		option = strings.TrimSpace(option)
		option = strings.TrimLeft(option, "{{")
		option = strings.TrimRight(option, "}}")
		optionName, optionType, _ := strings.Cut(option, ".")
		baseOption[optionName] = optionType
	}

	for i := 1; i < len(sa); i++ {
		for _, option := range powershellStruct.Options {
			if option.Flag == sa[i] {
				if baseOption[option.Name] == "Flag" {
					oldflag := fmt.Sprintf("{{%s.Flag}}", option.Name)
					cmd = strings.Replace(cmd, oldflag, sa[i], -1)
				}
				if baseOption[option.Name] == "Value" {
					if i+1 > len(sa) {
						colorcode.PrintMessage(colorcode.SIGN_ERROR, "the command is missing parameters")
					}
					oldflag := fmt.Sprintf("{{%s.Value}}", option.Name)
					newFlag := fmt.Sprintf("%s %s", sa[i], sa[i+1])
					cmd = strings.Replace(cmd, oldflag, newFlag, -1)
				}
			}
		}
	}
	cmd = r.ReplaceAllString(cmd, "")
	cmd = strings.Replace(cmd, `$url$`, fmt.Sprintf("'%s'", url), -1)
	return cmd
}

func PrintScriptsTable(powershellYaml PowershellYaml) {
	var data [][]string
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"id", "name", "arch", "privilege", "author", "description", "loaded"})
	table.SetBorders(tablewriter.Border{Left: false, Top: false, Right: false, Bottom: false})
	table.SetColWidth(48)
	for i, powershellStruct := range powershellYaml.PowershellStructs {
		var loaded string
		if powershellStruct.Loaded {
			loaded = colorcode.Colorf(colorcode.COLOR_PURPLE, "loaded")
		}
		data = append(data, []string{strconv.Itoa(i + 1), powershellStruct.Name, powershellStruct.Arch, powershellStruct.Privilege, powershellStruct.Author, powershellStruct.Description, loaded})
	}

	for _, raw := range data {
		table.Append(raw)
	}
	table.Render()
}

func PrintOptionsTable(powershellYaml PowershellYaml, scriptName string) {
	var data [][]string
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"name", "value", "required", "description"})
	table.SetBorders(tablewriter.Border{Left: false, Top: false, Right: false, Bottom: false})
	table.SetColWidth(48)
	var ps PowershellStruct
	for _, powershellStruct := range powershellYaml.PowershellStructs {
		if scriptName == powershellStruct.Name {
			ps = powershellStruct
			break
		}
	}
	for _, option := range ps.Options {
		data = append(data, []string{option.Name, option.Value, strconv.FormatBool(option.Required), option.Description})
	}
	for _, raw := range data {
		table.Append(raw)
	}
	table.Render()
}
