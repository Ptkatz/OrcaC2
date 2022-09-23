package console

import (
	"fmt"
	"Orca_Puppet/pkg/go-engine/common"
	"testing"
	"time"
)

func Test0001(t *testing.T) {
	c := NewConsole(true, 0, false)
	c.SetPretext("welcome:")
	if c != nil {
		c.Put("aaa")
		c.Put("aaa124")
		c.Put("aaa12312")
		time.Sleep(time.Second * 2)
		c.Stop()
	}
}

func Test0002(t *testing.T) {
	ci := ConsoleInput{}
	err := ci.Init()
	if err != nil {
		fmt.Println(err)
		return
	}

	go func() {
		defer common.CrashLog()

		for {
			e := ci.PollEvent()
			if e != nil {
				fmt.Println(e.Name())
			}
		}
	}()

	time.Sleep(time.Second * 2)
	ci.Stop()
}

func Test0003(t *testing.T) {

	eb := NewEditBox(3)
	eb.Input(&EventKey{key: KeyRune, ch: 'a'})
	eb.Input(&EventKey{key: KeyRune, ch: 'b'})
	eb.Input(&EventKey{key: KeyRune, ch: 'c'})
	fmt.Println(eb.GetText())
	if eb.GetText() != "abc" {
		t.Error("fail")
	}
}

func Test0004(t *testing.T) {

	eb := NewEditBox(3)
	eb.Input(&EventKey{key: KeyRune, ch: 'a'})
	eb.Input(&EventKey{key: KeyRune, ch: 'b'})
	eb.Input(&EventKey{key: KeyRune, ch: 'c'})
	eb.Input(&EventKey{key: KeyBackspace})
	fmt.Println(eb.GetText())
	if eb.GetText() != "ab" {
		t.Error("fail")
	}
}

func Test0005(t *testing.T) {

	eb := NewEditBox(3)
	eb.Input(&EventKey{key: KeyRune, ch: 'a'})
	eb.Input(&EventKey{key: KeyRune, ch: 'b'})
	eb.Input(&EventKey{key: KeyRune, ch: 'c'})
	eb.Input(&EventKey{key: KeyBackspace})
	eb.Input(&EventKey{key: KeyBackspace})
	eb.Input(&EventKey{key: KeyBackspace})
	eb.Input(&EventKey{key: KeyBackspace})
	fmt.Println(eb.GetText())
	if eb.GetText() != "" {
		t.Error("fail")
	}
}

func Test0006(t *testing.T) {

	eb := NewEditBox(3)
	eb.Input(&EventKey{key: KeyRune, ch: 'a'})
	eb.Input(&EventKey{key: KeyRune, ch: 'b'})
	eb.Input(&EventKey{key: KeyRune, ch: 'c'})
	eb.Input(&EventKey{key: KeyLeft})
	eb.Input(&EventKey{key: KeyBackspace})
	fmt.Println(eb.GetText())
	if eb.GetText() != "ac" {
		t.Error("fail")
	}
}

func Test0007(t *testing.T) {

	eb := NewEditBox(3)
	eb.Input(&EventKey{key: KeyRune, ch: 'a'})
	eb.Input(&EventKey{key: KeyRune, ch: 'b'})
	eb.Input(&EventKey{key: KeyRune, ch: 'c'})
	eb.Input(&EventKey{key: KeyLeft})
	eb.Input(&EventKey{key: KeyBS})
	eb.Input(&EventKey{key: KeyLeft})
	eb.Input(&EventKey{key: KeyDelete})
	fmt.Println(eb.GetText())
	if eb.GetText() != "c" {
		t.Error("fail")
	}
}

func Test0008(t *testing.T) {

	eb := NewEditBox(3)
	eb.Input(&EventKey{key: KeyRune, ch: 'a'})
	eb.Input(&EventKey{key: KeyRune, ch: 'b'})
	eb.Input(&EventKey{key: KeyRune, ch: 'c'})
	eb.Input(&EventKey{key: KeyLeft})
	eb.Input(&EventKey{key: KeyBS})
	eb.Input(&EventKey{key: KeyLeft})
	eb.Input(&EventKey{key: KeyDelete})
	eb.Input(&EventKey{key: KeyDelete})
	eb.Input(&EventKey{key: KeyDelete})
	fmt.Println(eb.GetText())
	if eb.GetText() != "" {
		t.Error("fail")
	}
}

func Test0009(t *testing.T) {

	eb := NewEditBox(3)
	eb.Input(&EventKey{key: KeyRune, ch: 'a'})
	eb.Input(&EventKey{key: KeyRune, ch: 'b'})
	eb.Input(&EventKey{key: KeyRune, ch: 'c'})
	eb.Input(&EventKey{key: KeyLeft})
	eb.Input(&EventKey{key: KeyRune, ch: 'd'})
	if eb.GetText() != "abdc" {
		t.Error("fail")
	}
	eb.Input(&EventKey{key: KeyLeft})
	eb.Input(&EventKey{key: KeyRune, ch: 'f'})
	if eb.GetText() != "abfdc" {
		t.Error("fail")
	}
	eb.Input(&EventKey{key: KeyRight})
	eb.Input(&EventKey{key: KeyRune, ch: 'c'})
	if eb.GetText() != "abfdcc" {
		t.Error("fail")
	}
	eb.Input(&EventKey{key: KeyEnter})
	if eb.GetText() != "" {
		t.Error("fail")
	}
	if eb.GetEnterText() != "abfdcc" {
		t.Error("fail")
	}
	eb.Input(&EventKey{key: KeyDown})
	if eb.GetText() != "" {
		t.Error("fail")
	}
	eb.Input(&EventKey{key: KeyUp})
	if eb.GetText() != "abfdcc" {
		t.Error("fail")
	}
	eb.Input(&EventKey{key: KeyLeft})
	eb.Input(&EventKey{key: KeyRune, ch: '1'})
	if eb.GetText() != "abfdc1c" {
		t.Error("fail")
	}
	eb.Input(&EventKey{key: KeyEnter})
	eb.Input(&EventKey{key: KeyUp})
	if eb.GetText() != "abfdc1c" {
		t.Error("fail")
	}
	eb.Input(&EventKey{key: KeyUp})
	if eb.GetText() != "abfdcc" {
		t.Error("fail")
	}
	eb.Input(&EventKey{key: KeyUp})
	if eb.GetText() != "abfdcc" {
		t.Error("fail")
	}
	eb.Input(&EventKey{key: KeyDown})
	if eb.GetText() != "abfdc1c" {
		t.Error("fail")
	}
	eb.Input(&EventKey{key: KeyRune, ch: '2'})
	if eb.GetText() != "abfdc1c2" {
		t.Error("fail")
	}
	eb.Input(&EventKey{key: KeyEnter})
	eb.Input(&EventKey{key: KeyUp})
	if eb.GetText() != "abfdc1c2" {
		t.Error("fail")
	}
	eb.Input(&EventKey{key: KeyUp})
	if eb.GetText() != "abfdc1c" {
		t.Error("fail")
	}
	eb.Input(&EventKey{key: KeyDown})
	if eb.GetText() != "abfdc1c2" {
		t.Error("fail")
	}
	eb.Input(&EventKey{key: KeyDown})
	eb.Input(&EventKey{key: KeyDown})
	eb.Input(&EventKey{key: KeyDown})
	if eb.GetText() != "" {
		t.Error("fail")
	}
	eb.Input(&EventKey{key: KeyUp})
	eb.Input(&EventKey{key: KeyUp})
	if eb.GetText() != "abfdc1c" {
		t.Error("fail")
	}
	eb.Input(&EventKey{key: KeyUp})
	if eb.GetText() != "abfdcc" {
		t.Error("fail")
	}
	eb.Input(&EventKey{key: KeyRune, ch: '3'})
	eb.Input(&EventKey{key: KeyEnter})
	eb.Input(&EventKey{key: KeyUp})
	if eb.GetText() != "abfdcc3" {
		t.Error("fail")
	}
	eb.Input(&EventKey{key: KeyUp})
	if eb.GetText() != "abfdc1c2" {
		t.Error("fail")
	}
	eb.Input(&EventKey{key: KeyUp})
	if eb.GetText() != "abfdc1c" {
		t.Error("fail")
	}
	eb.Input(&EventKey{key: KeyUp})
	if eb.GetText() != "abfdc1c" {
		t.Error("fail")
	}
}

func Test00010(t *testing.T) {

	eb := NewEditBox(3)
	eb.Input(&EventKey{key: KeyRune, ch: 'a'})
	eb.Input(&EventKey{key: KeyRune, ch: 'b'})
	eb.Input(&EventKey{key: KeyRune, ch: 'c'})
	fmt.Println(eb.GetShowText(true))

	eb.Input(&EventKey{key: KeyLeft})
	fmt.Println(eb.GetShowText(true))

	eb.Input(&EventKey{key: KeyLeft})
	fmt.Println(eb.GetShowText(true))

	eb.Input(&EventKey{key: KeyLeft})
	fmt.Println(eb.GetShowText(true))

	eb.Input(&EventKey{key: KeyLeft})
	fmt.Println(eb.GetShowText(true))

	eb.Input(&EventKey{key: KeyRune, ch: 'd'})
	fmt.Println(eb.GetShowText(true))
	fmt.Println(eb.GetShowText(false))
}

func Test00011(t *testing.T) {

	eb := NewEditBox(3)
	eb.Input(&EventKey{key: KeyRune, ch: 'a'})
	eb.Input(&EventKey{key: KeyRune, ch: 'b'})
	eb.Input(&EventKey{key: KeyRune, ch: 'c'})

	eb.Input(&EventKey{key: KeyLeft})
	fmt.Println(eb.GetShowText(false))

	eb.Input(&EventKey{key: KeyEnter})
	fmt.Println(eb.GetShowText(false))

}
