package spider

import (
	"github.com/PuerkitoBio/goquery"
	"Orca_Puppet/pkg/go-engine/common"
	"Orca_Puppet/pkg/go-engine/loggo"
	"Orca_Puppet/pkg/go-engine/shell"
	"math"
	"net/url"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Config struct {
	Threadnum    int
	Buffersize   int
	Deps         int
	FocusSpider  bool
	Crawlfunc    string // simple,puppeteer
	CrawlTimeout int
	CrawlRetry   int
}

type PageLinkInfo struct {
	UI   URLInfo
	Name string
}

type PageInfo struct {
	UI    URLInfo
	Title string
	Son   []PageLinkInfo
}

type URLInfo struct {
	Url  string
	Deps int
}

var gonce sync.Once

func Ini() {
	gonce.Do(ini)
}

func ini() {
	if runtime.GOOS == "linux" {
		go startChrome()
		go getChrome()
		for i := 0; i < 60; i++ {
			if len(gSpiderData.chromeWSEndpoint) > 0 {
				break
			}
			time.Sleep(time.Second)
		}

		if len(gSpiderData.chromeWSEndpoint) <= 0 {
			panic("spider start chrome fail")
		}

		loggo.Info("spider start chrome %v", gSpiderData.chromeWSEndpoint)
	}
}

type SpiderData struct {
	chromeWSEndpoint string
}

var gSpiderData SpiderData

func GetChromeWSEndpoint() string {
	return gSpiderData.chromeWSEndpoint
}

func getChrome() {
	defer common.CrashLog()

	for {
		ret, _ := shell.Run(common.GetNodeDir()+"/get_chrome.sh", true, common.GetNodeDir())
		ret = strings.TrimSpace(ret)
		if len(ret) > 0 {
			if ret != gSpiderData.chromeWSEndpoint {
				gSpiderData.chromeWSEndpoint = ret
				loggo.Info("spider get chromeWSEndpoint %v", gSpiderData.chromeWSEndpoint)
			}
		}
		time.Sleep(time.Second)
	}
}

func startChrome() {
	defer common.CrashLog()

	for {
		shell.RunTimeout(common.GetNodeDir()+"/close_chrome.sh", true, 60)
		loggo.Info("spider restart chrome ")
		shell.Run(common.GetNodeDir()+"/start_chrome.sh", true, common.GetNodeDir())
		loggo.Info("spider exit chrome ")
		time.Sleep(time.Second)
	}
}

type Stat struct {
	CrawBePushJobNum int

	CrawChannelNum  int
	CrawFunc        string
	CrawNum         int
	CrawRetrtyNum   int
	CrawOKNum       int64
	CrawFailNum     int
	CrawOKTotalTime int64
	CrawOKAvgTime   int64

	ParseChannelNum int
	ParseNum        int
	ParseValidNum   int
	ParseSpawnNum   int
	ParseFinishNum  int
	ParseTooDeepNum int
	ParseJobNum     int

	SaveChannelNum int
	SaveNum        int

	InsertNum       int64
	InsertTotalTime int64
	InsertAvgTime   int64

	JobInsertNum       int64
	JobInsertTotalTime int64
	JobInsertAvgTime   int64
	JobPopNum          int64
	JobPopTotalTime    int64
	JobPopAvgTime      int64
	JobHasNum          int64
	JobHasTotalTime    int64
	JobHasAvgTime      int64

	DoneInsertNum       int64
	DoneInsertTotalTime int64
	DoneInsertAvgTime   int64
	DoneHasNum          int64
	DoneHasTotalTime    int64
	DoneHasAvgTime      int64
}

type Content struct {
	Dsn   string
	Conn  int
	Crawl func(pg *PageInfo, doc *goquery.Document) *PageInfo
	Parse func(hosturl string, pg *PageInfo, save chan<- interface{}) bool
	Save  func(result interface{})
}

func Start(ctx *Content, config Config, url string, stat *Stat) {
	loggo.Info("Spider Start  %v", url)

	dsn := ctx.Dsn
	conn := ctx.Conn

	jbd := loadJob(dsn, conn, url)
	if jbd == nil {
		loggo.Error("Spider job loadJob fail %v", url)
		return
	}
	dbd := loadDone(dsn, conn, url)
	if dbd == nil {
		loggo.Error("Spider job loadDone fail %v", url)
		return
	}

	old := getJobSize(jbd)
	if old == 0 {
		insertSpiderJob(jbd, url, 0, stat)
		deleteSpiderDone(dbd)
	}

	old = getJobSize(jbd)
	if old == 0 {
		loggo.Error("Spider job no jobs %v", url)
		return
	}

	crawl := make(chan *URLInfo, config.Buffersize)
	parse := make(chan *PageInfo, config.Buffersize)
	save := make(chan interface{}, config.Buffersize)

	entry, deps := popSpiderJob(jbd, int(math.Min(float64(old), float64(config.Buffersize))), stat)
	if len(entry) == 0 {
		loggo.Error("Spider job no jobs %v", url)
		return
	}

	for i, u := range entry {
		crawl <- &URLInfo{u, deps[i]}
	}

	var jobsCrawlerTotal int32
	var jobsCrawlerFail int32

	var wg sync.WaitGroup
	var running int32

	for i := 0; i < config.Threadnum; i++ {
		wg.Add(3)
		go Crawler(&running, &wg, jbd, dbd, config, crawl, parse, &jobsCrawlerTotal, &jobsCrawlerFail,
			config.Crawlfunc, config.CrawlTimeout, config.CrawlRetry, stat, ctx)
		go Parser(&running, &wg, jbd, dbd, config, crawl, parse, save, url, stat, ctx)
		go Saver(&running, &wg, save, stat, ctx)
	}

	for {
		tmpurls, tmpdeps := popSpiderJob(jbd, config.Buffersize, stat)
		if len(tmpurls) == 0 {
			time.Sleep(time.Second)
			run := atomic.LoadInt32(&running)
			if run == 0 {
				time.Sleep(time.Second)
				tmpurls, tmpdeps = popSpiderJob(jbd, config.Buffersize, stat)
				if len(tmpurls) == 0 && run == 0 &&
					len(crawl) == 0 && len(parse) == 0 && len(save) == 0 {
					break
				}
			}
		}

		for i, url := range tmpurls {
			stat.CrawBePushJobNum++
			crawl <- &URLInfo{url, tmpdeps[i]}
		}
	}

	loggo.Info("Spider jobs done %v crawl %v, failed %v", url, jobsCrawlerTotal, jobsCrawlerFail)

	crawl <- nil
	parse <- nil
	save <- nil
	wg.Wait()

	close(crawl)
	close(parse)
	close(save)

	loggo.Info("Spider end %v %v", url, getDoneSize(dbd))

	dropJob(jbd)
	dropDone(dbd)

	closeJob(jbd)
	closeDone(dbd)
}

func Crawler(running *int32, group *sync.WaitGroup, jbd *JobDB, dbd *DoneDB, config Config, crawl <-chan *URLInfo, parse chan<- *PageInfo,
	jobsCrawlerTotal *int32, jobsCrawlerTotalFail *int32, crawlfunc string, crawlTimeout int, crawlRetry int, stat *Stat, ctx *Content) {
	defer common.CrashLog()

	defer group.Done()

	loggo.Info("Crawler start")
	for job := range crawl {
		if job == nil {
			break
		}
		atomic.AddInt32(running, 1)

		stat.CrawChannelNum = len(crawl)
		//loggo.Info("receive crawl job %v", job)

		ok := hasDone(dbd, job.Url, stat)
		if !ok {
			insertSpiderDone(dbd, job.Url, stat)
			if job.Deps < config.Deps {
				atomic.AddInt32(jobsCrawlerTotal, 1)
				var pg *PageInfo
				b := time.Now()
				stat.CrawNum++
				stat.CrawFunc = crawlfunc
				for t := 0; t < crawlRetry; t++ {
					stat.CrawRetrtyNum++
					if crawlfunc == "simple" {
						pg = simplecrawl(job, crawlTimeout, ctx)
					} else if crawlfunc == "puppeteer" {
						pg = puppeteercrawl(job, crawlTimeout, ctx)
					}
					if pg != nil {
						break
					}
				}
				if pg != nil {
					stat.CrawOKNum++
					stat.CrawOKTotalTime += int64(time.Now().Sub(b))
					loggo.Info("crawl job ok %v %v %v %s", job.Url, pg.Title, len(pg.Son), time.Now().Sub(b).String())
					parse <- pg
				} else {
					stat.CrawFailNum++
					atomic.AddInt32(jobsCrawlerTotalFail, 1)
				}
			}
		}

		atomic.AddInt32(running, -1)
	}
	loggo.Info("Crawler end")
}

func Parser(running *int32, group *sync.WaitGroup, jbd *JobDB, dbd *DoneDB, config Config, crawl chan<- *URLInfo, parse <-chan *PageInfo, save chan<- interface{},
	hosturl string, stat *Stat, ctx *Content) {
	defer common.CrashLog()

	defer group.Done()

	loggo.Info("Parser start")

	for job := range parse {
		if job == nil {
			break
		}
		atomic.AddInt32(running, 1)

		stat.ParseChannelNum = len(parse)
		//loggo.Info("receive parse job %v %v", job.Title, job.UI.Url)

		stat.ParseNum++

		srcURL, err := url.Parse(job.UI.Url)
		if err != nil {
			atomic.AddInt32(running, -1)
			continue
		}

		stat.ParseValidNum++

		ok := ctx.Parse(hosturl, job, save)
		if ok {
			stat.ParseFinishNum++
		}

		for _, s := range job.Son {
			sonurl := s.UI.Url

			stat.ParseSpawnNum++

			if strings.HasPrefix(sonurl, "#") {
				continue
			}

			if sonurl == "/" {
				continue
			}

			if strings.Contains(sonurl, "javascript:") {
				continue
			}

			if s.UI.Deps >= config.Deps {
				stat.ParseTooDeepNum++
				continue
			}

			sonu, err := url.Parse(sonurl)
			if err != nil {
				continue
			}

			sonu = srcURL.ResolveReference(sonu)
			sonurl = sonu.String()

			var tmp *URLInfo

			finded := hasDone(dbd, sonurl, stat)
			if !finded {
				if config.FocusSpider {
					dstParams := strings.Split(sonu.Host, ".")
					srcParams := strings.Split(srcURL.Host, ".")

					if len(dstParams) >= 2 && len(srcParams) >= 2 &&
						dstParams[len(dstParams)-1] == srcParams[len(srcParams)-1] &&
						dstParams[len(dstParams)-2] == srcParams[len(srcParams)-2] {
						tmp = &URLInfo{sonurl, s.UI.Deps}
					}
				} else {
					tmp = &URLInfo{sonurl, s.UI.Deps}
				}
			}

			if tmp != nil {
				hasJob := hasJob(jbd, tmp.Url, stat)
				if !hasJob {
					stat.ParseJobNum++

					insertSpiderJob(jbd, tmp.Url, tmp.Deps, stat)

					//loggo.Info("parse spawn job %v %v %v", job.UI.Url, sonurl, getJobSize(src))
				}
			}
		}
		atomic.AddInt32(running, -1)
	}
	loggo.Info("Parser end")
}

func Saver(running *int32, group *sync.WaitGroup, save <-chan interface{}, stat *Stat, ctx *Content) {
	defer common.CrashLog()

	defer group.Done()

	loggo.Info("Saver start")

	for job := range save {
		if job == nil {
			break
		}
		atomic.AddInt32(running, 1)

		stat.SaveChannelNum = len(save)
		//loggo.Info("receive save job %v %v %v", job.Title, job.Name, job.Url)

		stat.SaveNum++

		stat.InsertNum++
		b := time.Now()
		ctx.Save(job)
		stat.InsertTotalTime += int64(time.Since(b))

		atomic.AddInt32(running, -1)
	}

	loggo.Info("Saver end")
}
