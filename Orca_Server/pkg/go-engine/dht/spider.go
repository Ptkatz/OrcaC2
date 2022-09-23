package dht

import (
	"Orca_Server/pkg/go-engine/common"
	"Orca_Server/pkg/go-engine/loggo"
	"Orca_Server/pkg/go-engine/tmap"
	"encoding/hex"
	"encoding/json"
	"github.com/shiyanhui/dht"
)

var gdht *tmap.TMysql
var gcb func(infohash string, name string)

func Load(dsn string, conn int) error {

	loggo.Info("mysql dht Load start")

	gdht = tmap.NewTMysql(dsn, conn, "dht", 30)

	err := gdht.Load()
	if err != nil {
		return err
	}

	loggo.Info("mysql dht Load ok")

	go Crawl()

	return nil
}

func SetCallback(cb func(infohash string, name string)) {
	gcb = cb
}

type file struct {
	Path   []interface{} `json:"path"`
	Length int           `json:"length"`
}

type bitTorrent struct {
	InfoHash string `json:"infohash"`
	Name     string `json:"name"`
	Files    []file `json:"files,omitempty"`
	Length   int    `json:"length,omitempty"`
}

func OnCrawl(w *dht.Wire) {
	defer common.CrashLog()

	for resp := range w.Response() {
		loggo.Info("OnCrawl resp bytes %v", len(resp.MetadataInfo))

		metadata, err := dht.Decode(resp.MetadataInfo)
		if err != nil {
			continue
		}
		info := metadata.(map[string]interface{})

		if _, ok := info["name"]; !ok {
			continue
		}

		bt := bitTorrent{
			InfoHash: hex.EncodeToString(resp.InfoHash),
			Name:     info["name"].(string),
		}

		if v, ok := info["files"]; ok {
			files := v.([]interface{})
			bt.Files = make([]file, len(files))

			for i, item := range files {
				f := item.(map[string]interface{})
				bt.Files[i] = file{
					Path:   f["path"].([]interface{}),
					Length: f["length"].(int),
				}
			}
		} else if _, ok := info["length"]; ok {
			bt.Length = info["length"].(int)
		}

		data, err := json.Marshal(bt)
		if err == nil {
			loggo.Info("Crawl %s", data)

			InsertSpider(bt.InfoHash, bt.Name)
		}
	}
}

func InsertSpider(infohash string, name string) {
	if gdht.Insert(infohash, name) != nil {
		return
	}
	if gcb != nil {
		gcb(infohash, name)
	}
}

func GetSize() int {
	return gdht.GetSize()
}

func Crawl() {
	defer common.CrashLog()

	w := dht.NewWire(65536, 1024, 256)
	go OnCrawl(w)
	go func() {
		defer common.CrashLog()
		w.Run()
	}()

	config := dht.NewCrawlConfig()
	config.OnAnnouncePeer = func(infoHash, ip string, port int) {
		w.Request([]byte(infoHash), ip, port)
	}
	d := dht.New(config)

	go func() {
		defer common.CrashLog()
		d.Run()
	}()
}

type FindData struct {
	Infohash string
	Name     string
}

func Last(n int) []FindData {
	var ret []FindData
	find := gdht.Last(n)
	for _, v := range find {
		ret = append(ret, FindData{v.Name, v.Value})
	}
	return ret
}

func Find(str string, max int) []FindData {
	var ret []FindData
	find := gdht.FindValue(str, max)
	for _, v := range find {
		ret = append(ret, FindData{v.Name, v.Value})
	}
	return ret
}
