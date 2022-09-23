package htmlgen

import (
	"container/list"
	"fmt"
	"Orca_Puppet/pkg/go-engine/common"
	"Orca_Puppet/pkg/go-engine/loggo"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type HtmlGen struct {
	name        string
	path        string
	lastest     list.List
	lastestmax  int
	lastesttpl  string
	maxday      int
	subpagetpl  string
	cur         []string
	lasttime    time.Time
	lastsub     time.Time
	daypagetpl  string
	hourpagetpl string
}

func New(name string, path string, maxlastest int, maxday int,
	mainpagetpl string, subpagetpl string,
	daypagetpl string, hourpagetpl string) *HtmlGen {
	loggo.Info("Ini start %s", path)
	os.MkdirAll(path, os.ModePerm)
	os.MkdirAll(path+"/htmlgen/", os.ModePerm)
	hg := &HtmlGen{}
	hg.name = name
	hg.path = path
	hg.lastestmax = maxlastest
	hg.maxday = maxday
	hg.lasttime = time.Now()

	if len(mainpagetpl) > 0 {
		hg.lastesttpl = mainpagetpl
	} else {
		hg.lastesttpl = common.GetSrcDir() + "/htmlgen/" + "mainpage.tpl"
		if _, err := os.Stat(hg.lastesttpl); os.IsNotExist(err) {
			panic("no main page tpl at " + hg.lastesttpl)
		}
	}

	if len(subpagetpl) > 0 {
		hg.subpagetpl = subpagetpl
	} else {
		hg.subpagetpl = common.GetSrcDir() + "/htmlgen/" + "subpage.tpl"
		if _, err := os.Stat(hg.subpagetpl); os.IsNotExist(err) {
			panic("no main page tpl at " + hg.subpagetpl)
		}
	}

	if len(daypagetpl) > 0 {
		hg.daypagetpl = daypagetpl
	} else {
		hg.daypagetpl = common.GetSrcDir() + "/htmlgen/" + "daypage.tpl"
		if _, err := os.Stat(hg.daypagetpl); os.IsNotExist(err) {
			panic("no main page tpl at " + hg.daypagetpl)
		}
	}

	if len(hourpagetpl) > 0 {
		hg.hourpagetpl = hourpagetpl
	} else {
		hg.hourpagetpl = common.GetSrcDir() + "/htmlgen/" + "hourpage.tpl"
		if _, err := os.Stat(hg.hourpagetpl); os.IsNotExist(err) {
			panic("no main page tpl at " + hg.hourpagetpl)
		}
	}

	hg.deleteHtml()
	go func() {
		defer common.CrashLog()
		for {
			time.Sleep(time.Hour)
			hg.deleteHtml()
		}
	}()

	return hg
}

func (hg *HtmlGen) AddHtml(html string) error {
	now := time.Now()
	hg.addLatest(html)
	err := hg.saveLatest(now)
	if err != nil {
		return err
	}
	err = hg.saveDayTime(now)
	if err != nil {
		return err
	}
	err = hg.saveHourTime(now)
	if err != nil {
		return err
	}
	err = hg.saveSub(now, html)
	if err != nil {
		return err
	}
	loggo.Info("AddHtml %s", html)
	return nil
}

func (hg *HtmlGen) addLatest(s string) {
	hg.lastest.PushFront(s)
	if hg.lastest.Len() > hg.lastestmax {
		var last *list.Element
		for e := hg.lastest.Front(); e != nil; e = e.Next() {
			last = e
		}
		if last != nil {
			hg.lastest.Remove(last)
		}
	}
}

type mainpageLastest struct {
	Name string
}

type mainpageSub struct {
	Name string
}

type mainpage struct {
	Name    string
	Lastest []mainpageLastest
	Sub     []mainpageSub
}

func noescape(str string) template.HTML {
	return template.HTML(str)
}

func (hg *HtmlGen) savefile(data interface{}, des string, src string) error {

	file, err := os.Create(des)
	if err != nil {
		loggo.Error("os create %s", err)
		return err
	}
	defer file.Close()

	t := template.New("text")
	if err != nil {
		loggo.Error("template New %s", err)
		return err
	}

	t = t.Funcs(template.FuncMap{"noescape": noescape})

	srcfile, err := os.Open(src)
	if err != nil {
		loggo.Error("os Open %s", err)
		return err
	}
	defer srcfile.Close()

	var buffer [1024 * 1024]byte
	n, rerr := srcfile.Read(buffer[0:])
	if rerr != nil {
		loggo.Error("srcfile Read %s", err)
		return err
	}

	t, err = t.Parse(string(buffer[0:n]))
	if err != nil {
		loggo.Error("template Parse %s", err)
		return err
	}

	err = t.Execute(file, data)
	if err != nil {
		loggo.Error("template Execute %s", err)
		return err
	}

	return nil
}

func (hg *HtmlGen) saveLatest(now time.Time) error {
	mp := &mainpage{}
	mp.Name = hg.name
	for e := hg.lastest.Front(); e != nil; e = e.Next() {
		t := mainpageLastest{}
		t.Name = e.Value.(string)
		mp.Lastest = append(mp.Lastest, t)
	}

	for i := 0; i < hg.maxday; i++ {
		tt := time.Now().Add(-24 * time.Hour * time.Duration(i))
		t := mainpageSub{}
		t.Name = tt.Format("2006-01-02")
		mp.Sub = append(mp.Sub, t)
	}

	des := hg.path + "/" + "htmlgen.html"

	src := hg.lastesttpl

	return hg.savefile(mp, des, src)
}

type subpageData struct {
	Name string
}

type subpage struct {
	Name string
	Data []subpageData
}

func (hg *HtmlGen) saveSub(now time.Time, s string) error {

	cur := now.Format("2006-01-02_15-04")
	last := hg.lasttime.Format("2006-01-02_15-04")

	mustsave := false
	var old []string
	if cur != last {
		old = hg.cur
		hg.cur = make([]string, 0)
		hg.lasttime = now
		mustsave = true
	}
	hg.cur = append(hg.cur, s)

	if !mustsave {
		return nil
	}

	head := cur

	sp := &subpage{}
	sp.Name = head
	for i := len(old) - 1; i >= 0; i-- {
		t := subpageData{}
		t.Name = old[i]
		sp.Data = append(sp.Data, t)
	}

	des := hg.path + "/htmlgen" + "/" + head + ".html"

	src := hg.subpagetpl

	return hg.savefile(sp, des, src)
}

func (hg *HtmlGen) deleteHtml() {
	now := time.Now().Format("2006-01-02")
	nowt, _ := time.Parse("2006-01-02", now)
	nowunix := nowt.Unix()
	filepath.Walk(hg.path+"/htmlgen", func(path string, f os.FileInfo, err error) error {

		if f == nil || f.IsDir() {
			return nil
		}

		if !strings.HasSuffix(f.Name(), ".html") {
			return nil
		}

		date := f.Name()
		date = strings.TrimRight(date, ".html")
		index := strings.Index(date, "_")
		if index > 0 {
			date = date[0:index]
		}

		t, e := time.Parse("2006-01-02", date)
		if e != nil {
			loggo.Error("delete Parse file fail %v %v %v", f.Name(), date, err)
			return nil
		}
		tunix := t.Unix()
		if nowunix-tunix > int64(hg.maxday)*24*3600 {
			err := os.Remove(hg.path + "/htmlgen" + "/" + f.Name())
			if e != nil {
				loggo.Error("delete file fail %v %v", f.Name(), err)
				return nil
			}
		}

		return nil
	})
}

type subtimepageData struct {
	Name string
}

type timepage struct {
	Name string
	Data []subtimepageData
}

func (hg *HtmlGen) saveDayTime(now time.Time) error {
	day := time.Now().Format("2006-01-02")

	dp := &timepage{}
	dp.Name = day

	for i := 0; i < 24; i++ {
		t := subtimepageData{}
		t.Name = day + "_" + fmt.Sprintf("%02d", i)
		dp.Data = append(dp.Data, t)
	}

	des := hg.path + "/htmlgen" + "/" + day + ".html"

	src := hg.daypagetpl

	return hg.savefile(dp, des, src)
}

func (hg *HtmlGen) saveHourTime(now time.Time) error {
	hour := time.Now().Format("2006-01-02_15")

	dp := &timepage{}
	dp.Name = hour

	for i := 0; i < 60; i++ {
		t := subtimepageData{}
		t.Name = hour + "-" + fmt.Sprintf("%02d", i)
		dp.Data = append(dp.Data, t)
	}

	des := hg.path + "/htmlgen" + "/" + hour + ".html"

	src := hg.hourpagetpl

	return hg.savefile(dp, des, src)
}
