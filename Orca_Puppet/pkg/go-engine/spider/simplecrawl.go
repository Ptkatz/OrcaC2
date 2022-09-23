package spider

import (
	"crypto/tls"
	"github.com/PuerkitoBio/goquery"
	"github.com/axgle/mahonia"
	"Orca_Puppet/pkg/go-engine/loggo"
	"net/http"
	"strings"
	"time"
)

func simplecrawl(ui *URLInfo, crawlTimeout int, ctx *Content) *PageInfo {

	url := ui.Url
	loggo.Info("start simple crawl %v", url)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   time.Duration(crawlTimeout) * time.Second,
	}
	defer client.CloseIdleConnections()

	res, err := client.Get(url)
	if err != nil {
		loggo.Info("simple crawl http Get fail %v %v", url, err)
		return nil
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		loggo.Info("simple crawl http StatusCode fail %v %v", url, res.StatusCode)
		return nil
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		loggo.Info("simple crawl http NewDocumentFromReader fail %v %v", url, err)
		return nil
	}

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
			//loggo.Info("simple simple crawl title %v", pg.Title)
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
			//loggo.Info("simple simple crawl link %v %v %v %v", i, pg.Title, name, href)

			if len(href) > 0 {
				pgl := PageLinkInfo{URLInfo{href, ui.Deps + 1}, name}
				pg.Son = append(pg.Son, pgl)
			}
		}
	})

	pg = ctx.Crawl(pg, doc)

	//if len(pg.Son) == 0 {
	//	html, _ := doc.Html()
	//	loggo.Info("simple simple crawl no link %v html:\n%v", url, html)
	//}

	return pg
}
