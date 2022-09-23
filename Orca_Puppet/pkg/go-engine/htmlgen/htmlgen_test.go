package htmlgen

import (
	"testing"
	"time"
)

func Test0001(t *testing.T) {
	hg := New("test", "./", 10, 10,
		"./mainpage.tpl", "./subpage.tpl",
		"./daypage.tpl", "./hourpage.tpl")
	for i := 0; i < 20; i++ {
		hg.AddHtml("aa")
		time.Sleep(time.Second)
		hg.AddHtml("bb")
		time.Sleep(time.Second)
		hg.AddHtml("aaa")
		time.Sleep(time.Second)
		hg.AddHtml("啊啊")
		time.Sleep(time.Second)
		hg.AddHtml("3阿斯发a")
		time.Sleep(time.Second)
		hg.AddHtml("asfa")
		time.Sleep(time.Second)
	}
}
