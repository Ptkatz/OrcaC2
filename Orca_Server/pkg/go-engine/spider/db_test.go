package spider

import (
	"fmt"
	"github.com/go-sql-driver/mysql"
	"testing"
)

func Test0001(t *testing.T) {

	dbconfig := mysql.NewConfig()
	dbconfig.User = "root"
	dbconfig.Passwd = "123123"
	dbconfig.Addr = "192.168.0.106:4406"
	dbconfig.Net = "tcp"

	f := loadJob(dbconfig.FormatDSN(), 10, "http://www.baidu.com")
	if f == nil {
		return
	}
	s := Stat{}
	insertSpiderJob(f, "aaa", 1, &s)
	insertSpiderJob(f, "aaaa", 1, &s)
	insertSpiderJob(f, "aaaaa", 1, &s)
	fmt.Println(hasJob(f, "aaa", &s))
	fmt.Println(hasJob(f, "aaba", &s))
	u, d := popSpiderJob(f, 1, &s)
	fmt.Println(u)
	fmt.Println(d)
}
