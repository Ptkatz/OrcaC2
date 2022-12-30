package screenstreamopt

import (
	"Orca_Master/cli/common"
	"Orca_Master/define/config"
	"Orca_Master/tools/crypto"
	"encoding/json"
	"github.com/faiface/pixel/pixelgl"
)

type Pos struct {
	X int
	Y int
}

// 鼠标滚动
type Scroll struct {
	ScrollX int
	ScrollY int
}

type MouseInfo struct {
	MousePos    Pos
	PressButton string
	Scroll      Scroll
	Toggle      bool
}

var MouseBtnMap = map[pixelgl.Button]string{
	pixelgl.MouseButtonLeft:   "left",
	pixelgl.MouseButtonRight:  "right",
	pixelgl.MouseButtonMiddle: "center",
}

// 发送鼠标信息
func SendMouseInfoMsg(clientId, sendUserId string, mouseInfo MouseInfo) common.HttpRetData {
	msg := "mouseAction"
	mouseInfoStr, err := json.Marshal(mouseInfo)
	decMouseInfo, err := crypto.Encrypt(mouseInfoStr, []byte(config.AesKey))
	if err != nil {
		return common.HttpRetData{}
	}
	data := decMouseInfo
	return common.SendSuccessMsg(clientId, sendUserId, msg, data, "")
}

// 获取鼠标位置
func GetMousePosition(win *pixelgl.Window, x, y, maxHeight float64) Pos {
	mousePosition := win.MousePosition()

	moveX := int(mousePosition.X / x)
	moveY := int(maxHeight - (mousePosition.Y / y))
	pos := Pos{
		X: moveX,
		Y: moveY,
	}
	return pos
}

// 获取点击按钮
func GetMouseJustPressBtn(win *pixelgl.Window) string {
	var btn string
	if win.JustPressed(pixelgl.MouseButtonLeft) {
		btn = MouseBtnMap[pixelgl.MouseButtonLeft]
	} else if win.JustPressed(pixelgl.MouseButtonRight) {
		btn = MouseBtnMap[pixelgl.MouseButtonRight]
	} else if win.JustPressed(pixelgl.MouseButtonMiddle) {
		btn = MouseBtnMap[pixelgl.MouseButtonMiddle]
	}
	return btn
}

// 获取按下按钮
func GetMousePressBtn(win *pixelgl.Window) string {
	var btn string
	if win.Pressed(pixelgl.MouseButtonLeft) {
		btn = MouseBtnMap[pixelgl.MouseButtonLeft]
	} else if win.Pressed(pixelgl.MouseButtonRight) {
		btn = MouseBtnMap[pixelgl.MouseButtonRight]
	} else if win.Pressed(pixelgl.MouseButtonMiddle) {
		btn = MouseBtnMap[pixelgl.MouseButtonMiddle]
	}
	return btn
}

// 判断鼠标是否拖动
func IsToggle(win *pixelgl.Window) bool {
	jbtn := GetMouseJustPressBtn(win)
	btn := GetMousePressBtn(win)
	if btn != "" {
		if btn == "left" && jbtn == "" {
			return true
		}
	}
	return false
}

// 获取鼠标滚轮滚动信息
func GetScroll(win *pixelgl.Window, x, y float64) Scroll {
	scrollVec := win.MouseScroll()
	scrollX := int(scrollVec.X / x)
	scrollY := int(scrollVec.Y / y)
	scroll := Scroll{
		ScrollX: scrollX,
		ScrollY: scrollY,
	}
	return scroll
}
