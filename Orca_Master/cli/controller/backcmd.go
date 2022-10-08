package controller

import (
	"Orca_Master/cli/cmdopt/sshopt"
	"github.com/desertbit/grumble"
)

var backCmd = &grumble.Command{
	Name:  "back",
	Help:  "back to the main menu",
	Usage: "back [-h | --help]",
	Run: func(c *grumble.Context) error {
		BackMainMenu()
		return nil
	},
}

func BackMainMenu() {
	App.SetPrompt(InitPrompt)
	SelectClientId = ""
	SelectId = -1
	SelectIp = ""
	SelectVer = ""
	sshopt.MySsh.Node = SelectClientId
	RemoveCommand()
	sshopt.InitSshOption()
}
