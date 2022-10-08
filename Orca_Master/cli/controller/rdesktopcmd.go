package controller

import (
	"Orca_Master/cli/cmdopt/screenstreamopt"
	"Orca_Master/cli/common"
	"Orca_Master/define/colorcode"
	"fmt"
	"github.com/desertbit/grumble"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/go-vgo/robotgo"
	"golang.org/x/image/colornames"
	_ "image/png"
)

var screenStreamCmd = &grumble.Command{
	Name:  "stream",
	Help:  "get screenshot stream",
	Usage: "screen stream [-h | --help]",
	Run: func(c *grumble.Context) error {
		if SelectId == -1 {
			colorcode.PrintMessage(colorcode.SIGN_ERROR, "please select the id first")
			return nil
		}
		pixelgl.Run(runScreenStream)
		colorcode.PrintMessage(colorcode.SIGN_NOTICE, "remote desktop interrupt")
		return nil
	},
}

func runScreenStream() {
	// 发送请求获取屏幕分辨率
	sx, sy := robotgo.GetScreenSize()
	sx, sy = screenstreamopt.GetScreenSize(SelectClientId, common.ClientId)
	if sx <= 0 || sy <= 0 {
		fmt.Errorf("screen size get failed")
		return
	}
	maxWidth := float64(sx)
	maxHeight := float64(sy)

	// 创建窗口
	title := "- Orca Remote Screen - " + SelectIp + " - "
	mouseState := "mouse:Off - "
	keyboardState := "keyboard:Off - "
	cfg := pixelgl.WindowConfig{
		Title:     title + mouseState + keyboardState,
		Bounds:    pixel.R(0, 0, maxWidth/1.6, maxHeight/1.6),
		VSync:     true,
		Resizable: true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		return
	}
	win.SetSmooth(true)

	mouseFlag := false
	keyboardFlag := false
	colorcode.PrintMessage(colorcode.SIGN_WARNING, "Tips: Press Ctrl+Alt+A to enable/disable mouse input")
	colorcode.PrintMessage(colorcode.SIGN_WARNING, "Tips: Press Ctrl+Alt+W to enable/disable keyboard input")
	for !win.Closed() {
		// 处理屏幕截图请求
		screenshotImg, err := screenstreamopt.SettleScreenRequest(SelectClientId, common.ClientId)
		if err != nil {
			fmt.Errorf(err.Error())
			return
		}
		pic := pixel.PictureDataFromImage(screenshotImg)
		sprite := pixel.NewSprite(pic, pic.Bounds())

		// 在窗口中渲染
		win.Clear(colornames.Skyblue)
		bounds := win.Bounds()
		width := bounds.Max.X
		height := bounds.Max.Y
		x := width / maxWidth
		y := height / maxHeight

		mat := pixel.IM
		mat = mat.Moved(win.Bounds().Center())
		mat = mat.ScaledXY(win.Bounds().Center(), pixel.V(x, y))
		sprite.Draw(win, mat)
		win.Update()

		if mouseFlag {
			mouseInfo := screenstreamopt.MouseInfo{
				MousePos:    screenstreamopt.GetMousePosition(win, x, y, maxHeight),
				PressButton: screenstreamopt.GetMousePressBtn(win),
				Scroll:      screenstreamopt.GetScroll(win, x, y),
				Toggle:      screenstreamopt.IsToggle(win),
			}
			// 发送鼠标信息
			if mouseInfo.PressButton != "" || mouseInfo.Scroll.ScrollX != 0 || mouseInfo.Scroll.ScrollY != 0 {
				screenstreamopt.SendMouseInfoMsg(SelectClientId, common.ClientId, mouseInfo)
			}
		}

		if keyboardFlag {
			press, release := screenstreamopt.CollectKeyStrokes(win)
			if press != nil || release != nil {
				keyboardInfo := screenstreamopt.KeyboardInfo{
					PressKey:   press,
					ReleaseKey: release,
				}
				// 发送键盘消息
				screenstreamopt.SendKeyboardInfoMsg(SelectClientId, common.ClientId, keyboardInfo)
			}
		}

		// 按下Ctrl+Alt+A可以启用/禁用鼠标输入
		if win.Pressed(pixelgl.KeyLeftControl) {
			if win.Pressed(pixelgl.KeyLeftAlt) {
				if win.JustPressed(pixelgl.KeyA) {
					mouseFlag = !mouseFlag
					if mouseFlag {
						mouseState = "mouse:On - "
						win.SetTitle(title + mouseState + keyboardState)
						colorcode.PrintMessage(colorcode.SIGN_NOTICE, "mouse is on!")
					} else {
						mouseState = "mouse:Off - "
						win.SetTitle(title + mouseState + keyboardState)
						colorcode.PrintMessage(colorcode.SIGN_NOTICE, "mouse is off!")
					}
				}
			}
		}

		// 按下Ctrl+Alt+D可以启用/禁用键盘输入
		if win.Pressed(pixelgl.KeyLeftControl) {
			if win.Pressed(pixelgl.KeyLeftAlt) {
				if win.JustPressed(pixelgl.KeyW) {
					keyboardFlag = !keyboardFlag
					if mouseFlag {
						keyboardState = "keyboard:On - "
						win.SetTitle(title + mouseState + keyboardState)
						colorcode.PrintMessage(colorcode.SIGN_NOTICE, "keyboard is on!")
					} else {
						keyboardState = "keyboard:Off - "
						win.SetTitle(title + mouseState + keyboardState)
						colorcode.PrintMessage(colorcode.SIGN_NOTICE, "keyboard is off!")
					}
				}
			}
		}
	}
}
