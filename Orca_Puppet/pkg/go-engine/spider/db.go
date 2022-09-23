package spider

import (
	"database/sql"
	"Orca_Puppet/pkg/go-engine/common"
	"Orca_Puppet/pkg/go-engine/loggo"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type DB struct {
	gdb         *sql.DB
	gInsertStmt *sql.Stmt
	gSizeStmt   *sql.Stmt
	gLastStmt   *sql.Stmt
	gFindStmt   *sql.Stmt
	gDeleteStmt *sql.Stmt
	gSelectStmt *sql.Stmt
}

type JobDB struct {
	gdb            *sql.DB
	src            string
	gInsertJobStmt *sql.Stmt
	gSizeJobStmt   *sql.Stmt
	gPeekJobStmt   *sql.Stmt
	gDeleteJobStmt *sql.Stmt
	gHasJobStmt    *sql.Stmt
	gDropStmt      *sql.Stmt
}

type DoneDB struct {
	gdb             *sql.DB
	src             string
	gInsertDoneStmt *sql.Stmt
	gSizeDoneStmt   *sql.Stmt
	gDeleteDoneStmt *sql.Stmt
	gHasDoneStmt    *sql.Stmt
	gDropStmt       *sql.Stmt
}

type DBLinkInfo struct {
	Host  string
	Title string
	Name  string
	Url   string
}

func Load(dsn string, conn int, name string, expireday int) *DB {

	loggo.Info("mysql Load start %s", name)

	gdb, err := sql.Open("mysql", dsn)
	if err != nil {
		loggo.Error("open mysql fail %v", err)
		return nil
	}

	loggo.Info("mysql Load Open ok")

	err = gdb.Ping()
	if err != nil {
		loggo.Error("open mysql fail %v", err)
		return nil
	}

	loggo.Info("mysql Load Ping ok")

	gdb.SetConnMaxLifetime(0)
	gdb.SetMaxIdleConns(conn)
	gdb.SetMaxOpenConns(conn)

	ret := new(DB)

	ret.gdb = gdb

	_, err = gdb.Exec("CREATE DATABASE IF NOT EXISTS " + name + " ")
	if err != nil {
		loggo.Error("CREATE DATABASE fail %v", err)
		return nil
	}

	_, err = gdb.Exec("CREATE TABLE  IF NOT EXISTS " + name + ".link_info(" +
		"url VARCHAR(200)  NOT NULL ," +
		"title VARCHAR(200) NOT NULL," +
		"name VARCHAR(200) NOT NULL," +
		"time DATETIME NOT NULL," +
		"PRIMARY KEY(url)," +
		"INDEX `time`(`time`) USING BTREE" +
		");")
	if err != nil {
		loggo.Error("CREATE TABLE fail %v", err)
		return nil
	}

	loggo.Info("mysql Load CREATE TABLE ok")

	stmt, err := gdb.Prepare("insert IGNORE  into " + name + ".link_info(title, name, url, time) values(?, ?, ?, NOW())")
	if err != nil {
		loggo.Error("Prepare mysql fail %v", err)
		return nil
	}
	ret.gInsertStmt = stmt

	stmt, err = gdb.Prepare("select count(*) as ret from " + name + ".link_info")
	if err != nil {
		loggo.Error("hasDone Prepare mysql fail %v", err)
		return nil
	}
	ret.gSizeStmt = stmt

	stmt, err = gdb.Prepare("select title,name,url from " + name + ".link_info order by time desc limit 0, ?")
	if err != nil {
		loggo.Error("Prepare mysql fail %v", err)
		return nil
	}
	ret.gLastStmt = stmt

	stmt, err = gdb.Prepare("select title,name,url from (select title,name,url,time from " + name + ".link_info where name like ? or title like ?  limit 0,?) as A  order by time desc ")
	if err != nil {
		loggo.Error("Prepare mysql fail %v", err)
		return nil
	}
	ret.gFindStmt = stmt

	stmt, err = gdb.Prepare("delete from " + name + ".link_info where (TO_DAYS(NOW()) - TO_DAYS(time))>=" + strconv.Itoa(expireday))
	if err != nil {
		loggo.Error("Prepare mysql fail %v", err)
		return nil
	}
	ret.gDeleteStmt = stmt

	stmt, err = gdb.Prepare("SELECT title,name,url FROM " + name + ".link_info LIMIT ?, ?")
	if err != nil {
		loggo.Error("Prepare mysql fail %v", err)
		return nil
	}
	ret.gSelectStmt = stmt

	loggo.Info("mysql Load Prepare stmt ok")

	go deleteOldSpider(ret)

	return ret
}

func dropJob(db *JobDB) {
	db.gDropStmt.Exec()
}

func closeJob(db *JobDB) {
	db.gInsertJobStmt.Close()
	db.gSizeJobStmt.Close()
	db.gPeekJobStmt.Close()
	db.gDeleteJobStmt.Close()
	db.gHasJobStmt.Close()
	db.gDropStmt.Close()
	db.gdb.Close()
}

func loadJob(dsn string, conn int, src string) *JobDB {

	loggo.Info("Load Job start %v", src)

	dstURL, _ := url.Parse(src)
	host := dstURL.Host
	host = strings.ReplaceAll(host, ".", "_")
	host = strings.ReplaceAll(host, "-", "_")

	gdb, err := sql.Open("mysql", dsn)
	if err != nil {
		loggo.Error("open mysql fail %v", err)
		return nil
	}

	err = gdb.Ping()
	if err != nil {
		loggo.Error("open mysql fail %v", err)
		return nil
	}

	gdb.SetConnMaxLifetime(0)
	gdb.SetMaxIdleConns(conn)
	gdb.SetMaxOpenConns(conn)

	_, err = gdb.Exec("CREATE DATABASE IF NOT EXISTS spiderjob")
	if err != nil {
		loggo.Error("CREATE DATABASE fail %v", err)
		return nil
	}

	ret := new(JobDB)

	ret.gdb = gdb
	ret.src = src

	_, err = gdb.Exec("CREATE TABLE  IF NOT EXISTS spiderjob." + host + "(" +
		"src TEXT NOT NULL," +
		"url VARCHAR(200)  NOT NULL ," +
		"deps INT NOT NULL," +
		"time DATETIME NOT NULL," +
		"PRIMARY KEY(url));")
	if err != nil {
		loggo.Error("CREATE DATABASE fail %v", err)
		return nil
	}

	stmt, err := gdb.Prepare("insert IGNORE into spiderjob." + host + "(src, url, deps, time) values(?, ?, ?, NOW())")
	if err != nil {
		loggo.Error("Prepare Job fail %v", err)
		return nil
	}
	ret.gInsertJobStmt = stmt

	stmt, err = gdb.Prepare("select count(*) from spiderjob." + host + " where src = ?")
	if err != nil {
		loggo.Error("hasDone Job Prepare fail %v", err)
		return nil
	}
	ret.gSizeJobStmt = stmt

	stmt, err = gdb.Prepare("delete from spiderjob." + host + " where src = ? and url = ?")
	if err != nil {
		loggo.Error("Prepare Job fail %v", err)
		return nil
	}
	ret.gDeleteJobStmt = stmt

	stmt, err = gdb.Prepare("select url, deps from spiderjob." + host + " where src = ? limit 0, ?")
	if err != nil {
		loggo.Error("Prepare Job fail %v", err)
		return nil
	}
	ret.gPeekJobStmt = stmt

	stmt, err = gdb.Prepare("select url from spiderjob." + host + " where src = ? and url = ?")
	if err != nil {
		loggo.Error("Prepare Job fail %v", err)
		return nil
	}
	ret.gHasJobStmt = stmt

	stmt, err = gdb.Prepare("drop table spiderjob." + host + "")
	if err != nil {
		loggo.Error("Prepare Job fail %v", err)
		return nil
	}
	ret.gDropStmt = stmt

	num := getJobSize(ret)
	loggo.Info("Job size %v %v", src, num)

	return ret
}

func dropDone(db *DoneDB) {
	db.gDropStmt.Exec()
}

func closeDone(db *DoneDB) {
	db.gInsertDoneStmt.Close()
	db.gSizeDoneStmt.Close()
	db.gDeleteDoneStmt.Close()
	db.gHasDoneStmt.Close()
	db.gDropStmt.Close()
	db.gdb.Close()
}

func loadDone(dsn string, conn int, src string) *DoneDB {

	loggo.Info("Load Done start %v", src)

	dstURL, _ := url.Parse(src)
	host := dstURL.Host
	host = strings.ReplaceAll(host, ".", "_")
	host = strings.ReplaceAll(host, "-", "_")

	gdb, err := sql.Open("mysql", dsn)
	if err != nil {
		loggo.Error("open mysql fail %v", err)
		return nil
	}

	err = gdb.Ping()
	if err != nil {
		loggo.Error("open mysql fail %v", err)
		return nil
	}

	gdb.SetConnMaxLifetime(0)
	gdb.SetMaxIdleConns(conn)
	gdb.SetMaxOpenConns(conn)

	_, err = gdb.Exec("CREATE DATABASE IF NOT EXISTS spiderdone")
	if err != nil {
		loggo.Error("CREATE DATABASE fail %v", err)
		return nil
	}

	ret := new(DoneDB)
	ret.gdb = gdb
	ret.src = src

	_, err = gdb.Exec("CREATE TABLE  IF NOT EXISTS spiderdone." + host + "(" +
		"src TEXT NOT NULL," +
		"url VARCHAR(200)  NOT NULL," +
		"time DATETIME NOT NULL," +
		"PRIMARY KEY(url));")
	if err != nil {
		loggo.Error("CREATE DATABASE fail %v", err)
		return nil
	}

	////

	stmt, err := gdb.Prepare("insert IGNORE into spiderdone." + host + "(src, url, time) values(?, ?, NOW())")
	if err != nil {
		loggo.Error("Prepare fail %v", err)
		return nil
	}
	ret.gInsertDoneStmt = stmt

	stmt, err = gdb.Prepare("select count(*) from spiderdone." + host + " where src = ?")
	if err != nil {
		loggo.Error("hasDone Prepare fail %v", err)
		return nil
	}
	ret.gSizeDoneStmt = stmt

	stmt, err = gdb.Prepare("delete from spiderdone." + host + " where src = ?")
	if err != nil {
		loggo.Error("Prepare fail %v", err)
		return nil
	}
	ret.gDeleteDoneStmt = stmt

	stmt, err = gdb.Prepare("select url from spiderdone." + host + " where src = ? and url = ?")
	if err != nil {
		loggo.Error("Prepare fail %v", err)
		return nil
	}
	ret.gHasDoneStmt = stmt

	stmt, err = gdb.Prepare("drop table spiderdone." + host + "")
	if err != nil {
		loggo.Error("Prepare Job fail %v", err)
		return nil
	}
	ret.gDropStmt = stmt

	////

	num := getDoneSize(ret)
	loggo.Info("size %v %v", src, num)

	return ret
}

func popSpiderJob(db *JobDB, n int, stat *Stat) ([]string, []int) {

	defer common.Elapsed(func(d time.Duration) {
		stat.JobPopNum++
		stat.JobPopTotalTime += int64(d)
	})()

	var ret []string
	var retdeps []int

	b := time.Now()

	rows, err := db.gPeekJobStmt.Query(db.src, n)
	if err != nil {
		loggo.Error("popSpiderJob Query sqlite3 fail %v %v", db.src, err)
		return ret, retdeps
	}
	defer rows.Close()

	for rows.Next() {

		var url string
		var deps int
		err = rows.Scan(&url, &deps)
		if err != nil {
			loggo.Error("popSpiderJob Scan sqlite3 fail %v %v", db.src, err)
		}

		ret = append(ret, url)
		retdeps = append(retdeps, deps)
	}

	for i, url := range ret {
		db.gDeleteJobStmt.Exec(db.src, url)
		loggo.Info("popSpiderJob %v %v %v %s", db.src, url, retdeps[i], time.Now().Sub(b).String())
	}

	return ret, retdeps
}

func deleteSpiderDone(db *DoneDB) {
	db.gDeleteDoneStmt.Exec(db.src)
}

func insertSpiderJob(db *JobDB, url string, deps int, stat *Stat) {

	defer common.Elapsed(func(d time.Duration) {
		stat.JobInsertNum++
		stat.JobInsertTotalTime += int64(d)
	})()

	b := time.Now()

	_, err := db.gInsertJobStmt.Exec(db.src, url, deps)
	if err != nil {
		loggo.Error("insertSpiderJob insert sqlite3 fail %v %v", url, err)
	}

	loggo.Info("insertSpiderJob %v %s", url, time.Now().Sub(b).String())
}

func insertSpiderDone(db *DoneDB, url string, stat *Stat) {

	defer common.Elapsed(func(d time.Duration) {
		stat.DoneInsertNum++
		stat.DoneInsertTotalTime += int64(d)
	})()

	b := time.Now()

	_, err := db.gInsertDoneStmt.Exec(db.src, url)
	if err != nil {
		loggo.Error("insertSpiderDone insert sqlite3 fail %v %v", url, err)
	}

	loggo.Info("insertSpiderDone %v %s", url, time.Now().Sub(b).String())
}

func deleteOldSpider(db *DB) {
	defer common.CrashLog()

	for {
		b := time.Now()

		db.gDeleteStmt.Exec()

		loggo.Info("deleteOldSpider %v %s", GetSize(db), time.Now().Sub(b).String())

		time.Sleep(time.Hour)
	}
}

func InsertSpider(db *DB, title string, name string, url string) {

	_, err := db.gInsertStmt.Exec(title, name, url)
	if err != nil {
		loggo.Error("InsertSpider insert sqlite3 fail %v %v", url, err)
	}

	loggo.Info("InsertSpider %v %v %v", title, name, url)
}

func hasJob(db *JobDB, url string, stat *Stat) bool {
	defer common.Elapsed(func(d time.Duration) {
		stat.JobHasNum++
		stat.JobHasTotalTime += int64(d)
	})()

	var surl string
	err := db.gHasJobStmt.QueryRow(db.src, url).Scan(&surl)
	if err != nil {
		return false
	}
	return true
}

func hasDone(db *DoneDB, url string, stat *Stat) bool {
	defer common.Elapsed(func(d time.Duration) {
		stat.DoneHasNum++
		stat.DoneHasTotalTime += int64(d)
	})()

	var surl string
	err := db.gHasDoneStmt.QueryRow(db.src, url).Scan(&surl)
	if err != nil {
		return false
	}
	return true
}

func GetSize(db *DB) int {
	var ret int
	err := db.gSizeStmt.QueryRow().Scan(&ret)
	if err != nil {
		loggo.Error("GetSize fail %v", err)
	}
	return ret
}

func getJobSize(db *JobDB) int {
	var ret int
	err := db.gSizeJobStmt.QueryRow(db.src).Scan(&ret)
	if err != nil {
		loggo.Error("getJobSize fail %v %v", db.src, err)
	}
	return ret
}

func getDoneSize(db *DoneDB) int {
	var ret int
	err := db.gSizeDoneStmt.QueryRow(db.src).Scan(&ret)
	if err != nil {
		loggo.Error("getDoneSize fail %v %v", db.src, err)
	}
	return ret
}

type FindData struct {
	Title string
	Name  string
	URL   string
}

func Last(db *DB, n int) []FindData {

	var ret []FindData

	retmap := make(map[string]string)

	rows, err := db.gLastStmt.Query(n)
	if err != nil {
		loggo.Error("Last Query sqlite3 fail %v", err)
		return ret
	}
	defer rows.Close()

	for rows.Next() {

		var title string
		var name string
		var url string
		err := rows.Scan(&title, &name, &url)
		if err != nil {
			loggo.Error("Last Scan sqlite3 fail %v", err)
		}

		_, ok := retmap[url]
		if ok {
			continue
		}
		retmap[url] = name

		ret = append(ret, FindData{title, name, url})
	}

	return ret
}

func Find(db *DB, str string, max int) []FindData {

	var ret []FindData

	rows, err := db.gFindStmt.Query("%"+str+"%", "%"+str+"%", max)
	if err != nil {
		loggo.Error("Find Query sqlite3 fail %v", err)
		return ret
	}
	defer rows.Close()

	for rows.Next() {

		var title string
		var name string
		var url string
		err = rows.Scan(&title, &name, &url)
		if err != nil {
			loggo.Error("Find Scan sqlite3 fail %v", err)
		}

		ret = append(ret, FindData{title, name, url})
	}

	return ret
}

func Select(db *DB, offset int, count int) []FindData {

	var ret []FindData

	rows, err := db.gSelectStmt.Query(offset, count)
	if err != nil {
		loggo.Error("Select Query sqlite3 fail %v", err)
		return ret
	}
	defer rows.Close()

	for rows.Next() {

		var title string
		var name string
		var url string
		err = rows.Scan(&title, &name, &url)
		if err != nil {
			loggo.Error("Select Scan sqlite3 fail %v", err)
		}

		ret = append(ret, FindData{title, name, url})
	}

	return ret
}
