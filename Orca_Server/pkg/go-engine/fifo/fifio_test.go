package fifo

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

	f, err := NewFIFO(dbconfig.FormatDSN(), 10, "aa")
	if err != nil {
		fmt.Println(err)
		return
	}
	f.Write("aa")
	f.Write("bb")
	f.Write("cc")

	fmt.Println(f.GetSize())

	for f.GetSize() > 0 {
		datas, _ := f.Read(10)
		for _, d := range datas {
			fmt.Println(d)
		}
	}
}

func Test0002(t *testing.T) {

	f, err := NewFIFOLocal("aa", 2)
	if err != nil {
		fmt.Println(err)
		return
	}
	f.Write("aa")
	f.Write("bb")
	f.Write("cc")

	fmt.Println(f.GetSize())

	for f.GetSize() > 0 {
		datas, _ := f.Read(1)
		for _, d := range datas {
			fmt.Println(d)
		}
	}
}
