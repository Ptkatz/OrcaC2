package screenstreamopt

import (
	"Orca_Master/cli/common"
	"Orca_Master/define/config"
	"Orca_Master/tools/crypto"
	"encoding/json"
	"github.com/faiface/pixel/pixelgl"
)

type KeyboardInfo struct {
	PressKey   []string
	ReleaseKey []string
}

var KeyboardBtnMap = map[pixelgl.Button]string{
	pixelgl.Key0:            "0",
	pixelgl.Key1:            "1",
	pixelgl.Key2:            "2",
	pixelgl.Key3:            "3",
	pixelgl.Key4:            "4",
	pixelgl.Key5:            "5",
	pixelgl.Key6:            "6",
	pixelgl.Key7:            "7",
	pixelgl.Key8:            "8",
	pixelgl.Key9:            "9",
	pixelgl.KeyA:            "a",
	pixelgl.KeyB:            "b",
	pixelgl.KeyC:            "c",
	pixelgl.KeyD:            "d",
	pixelgl.KeyE:            "e",
	pixelgl.KeyF:            "f",
	pixelgl.KeyG:            "g",
	pixelgl.KeyH:            "h",
	pixelgl.KeyI:            "i",
	pixelgl.KeyJ:            "j",
	pixelgl.KeyK:            "k",
	pixelgl.KeyL:            "l",
	pixelgl.KeyM:            "m",
	pixelgl.KeyN:            "n",
	pixelgl.KeyO:            "o",
	pixelgl.KeyP:            "p",
	pixelgl.KeyQ:            "q",
	pixelgl.KeyR:            "r",
	pixelgl.KeyS:            "s",
	pixelgl.KeyT:            "t",
	pixelgl.KeyU:            "u",
	pixelgl.KeyV:            "v",
	pixelgl.KeyW:            "w",
	pixelgl.KeyX:            "x",
	pixelgl.KeyY:            "y",
	pixelgl.KeyZ:            "z",
	pixelgl.KeySpace:        "space",
	pixelgl.KeyEscape:       "escape",
	pixelgl.KeyEnter:        "enter",
	pixelgl.KeyTab:          "tab",
	pixelgl.KeyBackspace:    "backspace",
	pixelgl.KeyInsert:       "insert",
	pixelgl.KeyDelete:       "delete",
	pixelgl.KeyRight:        "right",
	pixelgl.KeyLeft:         "left",
	pixelgl.KeyDown:         "down",
	pixelgl.KeyUp:           "up",
	pixelgl.KeyPageUp:       "pageup",
	pixelgl.KeyPageDown:     "pagedown",
	pixelgl.KeyHome:         "home",
	pixelgl.KeyEnd:          "end",
	pixelgl.KeyCapsLock:     "capslock",
	pixelgl.KeyNumLock:      "num_lock",
	pixelgl.KeyPrintScreen:  "printscreen",
	pixelgl.KeyF1:           "f1",
	pixelgl.KeyF2:           "f2",
	pixelgl.KeyF3:           "f3",
	pixelgl.KeyF4:           "f4",
	pixelgl.KeyF5:           "f5",
	pixelgl.KeyF6:           "f6",
	pixelgl.KeyF7:           "f7",
	pixelgl.KeyF8:           "f8",
	pixelgl.KeyF9:           "f9",
	pixelgl.KeyF10:          "f10",
	pixelgl.KeyF11:          "f11",
	pixelgl.KeyF12:          "f12",
	pixelgl.KeyF13:          "f13",
	pixelgl.KeyF14:          "f14",
	pixelgl.KeyF15:          "f15",
	pixelgl.KeyF16:          "f16",
	pixelgl.KeyF17:          "f17",
	pixelgl.KeyF18:          "f18",
	pixelgl.KeyF19:          "f19",
	pixelgl.KeyF20:          "f20",
	pixelgl.KeyF21:          "f21",
	pixelgl.KeyF22:          "f22",
	pixelgl.KeyF23:          "f23",
	pixelgl.KeyF24:          "f24",
	pixelgl.KeyLeftShift:    "lshift",
	pixelgl.KeyLeftControl:  "lctrl",
	pixelgl.KeyLeftAlt:      "lalt",
	pixelgl.KeyLeftSuper:    "cmd",
	pixelgl.KeyRightShift:   "rshift",
	pixelgl.KeyRightControl: "rctrl",
	pixelgl.KeyRightAlt:     "ralt",
	pixelgl.KeyRightSuper:   "rcmd",
	pixelgl.KeyMenu:         "menu",
	pixelgl.KeyMinus:        "-",
	pixelgl.KeyEqual:        "=",
	pixelgl.KeySlash:        "/",
	pixelgl.KeyPeriod:       ".",
	pixelgl.KeyComma:        ",",
	pixelgl.KeyApostrophe:   "'",
	pixelgl.KeySemicolon:    ";",
	pixelgl.KeyBackslash:    "\\",
	pixelgl.KeyLeftBracket:  "[",
	pixelgl.KeyRightBracket: "]",
	pixelgl.KeyGraveAccent:  "`",
}

// 获取按下与释放按键
func CollectKeyStrokes(win *pixelgl.Window) (pressed []string, release []string) {
	for k, v := range KeyboardBtnMap {
		key := v
		if win.JustPressed(k) || win.Repeated(k) {
			pressed = append(pressed, key)
		}
		if win.JustReleased(k) {
			release = append(release, key)
		}
	}
	return
}

// 发送键盘消息
func SendKeyboardInfoMsg(clientId, sendUserId string, keyboardInfo KeyboardInfo) common.HttpRetData {
	msg := "keyboardAction"
	keyboardInfoStr, err := json.Marshal(keyboardInfo)
	decKeyboardInfo, err := crypto.Encrypt(keyboardInfoStr, []byte(config.AesKey))
	if err != nil {
		return common.HttpRetData{}
	}
	data := decKeyboardInfo
	return common.SendSuccessMsg(clientId, sendUserId, msg, data, "")
}
