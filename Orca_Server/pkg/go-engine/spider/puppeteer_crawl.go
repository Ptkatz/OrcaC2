package spider

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/axgle/mahonia"
	"Orca_Server/pkg/go-engine/common"
	"Orca_Server/pkg/go-engine/loggo"
	"Orca_Server/pkg/go-engine/node"
	"golang.org/x/net/html"
	"path/filepath"
	"strings"
)

func puppeteercrawl(ui *URLInfo, crawlTimeout int, ctx *Content) *PageInfo {

	url := ui.Url
	loggo.Info("start puppeteer crawl %v %v", url, gSpiderData.chromeWSEndpoint)

	tmp := common.GetNodeDir() + "/" + "puppeteer_crawl.js"
	tmp = filepath.Clean(tmp)
	tmp = filepath.ToSlash(tmp)

	ret := node.Run(tmp, true, crawlTimeout, url, gSpiderData.chromeWSEndpoint)
	if len(ret) <= 0 {
		loggo.Info("puppeteer crawl http fail %v %v", url, gSpiderData.chromeWSEndpoint)
		return nil
	}

	r := strings.NewReader(ret)

	root, err := html.Parse(r)
	if err != nil {
		loggo.Info("puppeteer crawl html Parse fail %v %v", url, err)
		return nil
	}

	// Load the HTML document
	doc := goquery.NewDocumentFromNode(root)

	gb2312 := false
	doc.Find("META").Each(func(i int, s *goquery.Selection) {
		content, ok := s.Attr("content")
		if ok {
			if strings.Contains(content, "gb2312") {
				gb2312 = true
			}
		}
	})

	pg := &PageInfo{}
	pg.UI = *ui
	doc.Find("title").Each(func(i int, s *goquery.Selection) {
		if pg.Title == "" {
			pg.Title = s.Text()
			pg.Title = strings.TrimSpace(pg.Title)
			if gb2312 {
				enc := mahonia.NewDecoder("gbk")
				pg.Title = enc.ConvertString(pg.Title)
			}
			//loggo.Info("puppeteer crawl title %v", pg.Title)
		}
	})

	// Find the items
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		name := s.Text()
		href, ok := s.Attr("href")
		if ok {
			href = strings.TrimSpace(href)
			name = strings.TrimSpace(name)
			name = strings.Replace(name, "\n", " ", -1)
			if gb2312 {
				enc := mahonia.NewDecoder("gbk")
				href = enc.ConvertString(href)
				name = enc.ConvertString(name)
			}
			//loggo.Info("puppeteer crawl link %v %v %v %v", i, pg.Title, name, href)

			if len(href) > 0 {
				pgl := PageLinkInfo{URLInfo{href, ui.Deps + 1}, name}
				pg.Son = append(pg.Son, pgl)
			}
		}
	})

	pg = ctx.Crawl(pg, doc)

	//if len(pg.Son) == 0 {
	//	html, _ := doc.Html()
	//	loggo.Info("puppeteer crawl no link %v html:\n%v", url, html)
	//}

	return pg
}
