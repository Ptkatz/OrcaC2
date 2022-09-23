package common

import (
	"archive/zip"
	"encoding/json"
	"Orca_Server/pkg/go-engine/loggo"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

var gEngineDir string
var gNodeDir string
var gDataDir string
var gSrcDir string

var gonce sync.Once

// 引擎全局的初始化
// 如：引擎各个目录的实际地址
// 如：执行引擎初始化脚本
// 如：解压引擎数据文件
// 其他各个模块的初始化代码放到各自模块init中
func Ini() {
	gonce.Do(ini)
}

func ini() {

	var GOPATH string
	for i := 0; i < 100; i++ {
		GOPATH = os.Getenv("GOPATH")
		if len(GOPATH) > 0 {
			loggo.Info("GOPATH %v", GOPATH)
			break
		} else {
			loggo.Info("GOPATH empty, wait... %v", i)
			time.Sleep(time.Second)
		}
	}
	if len(GOPATH) <= 0 {
		panic("need GOPATH")
	}

	gpaths := strings.Split(GOPATH, ";")

	dir := ""
	for _, p := range gpaths {
		tmp := p + "/src/github.com/esrrhs/go-engine/"
		tmp = filepath.Clean(tmp)
		tmp = filepath.ToSlash(tmp)
		if _, err := os.Stat(tmp); !os.IsNotExist(err) {
			dir = tmp
			break
		}
	}

	if len(dir) <= 0 {
		panic("need install go-engine in GOPATH " + GOPATH)
		return
	}

	gEngineDir = dir
	loggo.Info("gEngineDir %v", gEngineDir)

	sysType := runtime.GOOS
	loggo.Info("sysType %v", sysType)
	gNodeDir = gEngineDir + "/node/" + sysType + "/"
	gNodeDir = filepath.Clean(gNodeDir)
	gNodeDir = filepath.ToSlash(gNodeDir)
	if _, err := os.Stat(gNodeDir); os.IsNotExist(err) {
		panic("need install node in go-engine " + gNodeDir)
		return
	}
	loggo.Info("gNodeDir %v", gNodeDir)

	gDataDir = gEngineDir + "/data/"
	gDataDir = filepath.Clean(gDataDir)
	gDataDir = filepath.ToSlash(gDataDir)
	if _, err := os.Stat(gDataDir); os.IsNotExist(err) {
		panic("need install data in go-engine " + gDataDir)
		return
	}
	loggo.Info("gDataDir %v", gDataDir)

	gSrcDir = gEngineDir + "/src/"
	gSrcDir = filepath.Clean(gSrcDir)
	gSrcDir = filepath.ToSlash(gSrcDir)
	if _, err := os.Stat(gSrcDir); os.IsNotExist(err) {
		panic("need install src in go-engine " + gSrcDir)
		return
	}
	loggo.Info("gSrcDir %v", gSrcDir)

	loadConfig()
	extract()
	initShell()
}

type EngineConfiguration struct {
	Extract []string
	Init    []string
}

var gEngineConfiguration EngineConfiguration

func initShell() {
	sysType := runtime.GOOS

	for _, f := range gEngineConfiguration.Init {
		if strings.Contains(f, sysType) {
			initShellFile(f)
		}
	}
}

func initShellFile(f string) {

	file := gEngineDir + "/" + f
	file = filepath.Clean(file)
	file = filepath.ToSlash(file)

	loggo.Info("InitShellFile %v", file)

	cmd := exec.Command("bash", file, gEngineDir)
	out, err := cmd.CombinedOutput()
	outstr := string(out)
	if err != nil {
		panic(err)
	}

	loggo.Info("InitShellFile ok %v %v", file, outstr)
}

func loadConfig() {
	file := gEngineDir + "/" + "config.json"
	file = filepath.Clean(file)
	file = filepath.ToSlash(file)
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	decoder := json.NewDecoder(f)
	conf := EngineConfiguration{}
	err = decoder.Decode(&conf)
	if err != nil {
		panic(err)
	}
	gEngineConfiguration = conf
}

func extractFile(file string) {
	file = filepath.Clean(file)
	file = filepath.ToSlash(file)
	if _, err := os.Stat(file); os.IsNotExist(err) {
		err = decompress(file + ".zip")
		if err != nil {
			panic("extractFile file fail " + file)
			return
		}
	}
	loggo.Info("extractFile %v", file)
}

func extract() {
	for _, f := range gEngineConfiguration.Extract {
		extractFile(gEngineDir + "/" + f)
	}
}

func GetEngineDir() string {
	return gEngineDir
}

func GetDataDir() string {
	return gDataDir
}

func GetNodeDir() string {
	return gNodeDir
}

func GetSrcDir() string {
	return gSrcDir
}

func decompress(tarFile string) error {
	dest := filepath.Dir(tarFile)

	srcFile, err := os.Open(tarFile)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	zipFile, err := zip.OpenReader(srcFile.Name())
	if err != nil {
		loggo.Error("Unzip File Error：%v", err)
		return err
	}
	defer zipFile.Close()

	for _, innerFile := range zipFile.File {
		info := innerFile.FileInfo()
		if info.IsDir() {
			err = os.MkdirAll(innerFile.Name, os.ModePerm)
			if err != nil {
				loggo.Error("Unzip File Error : %v", err)
				return err
			}
			continue
		}
		srcFile, err := innerFile.Open()
		if err != nil {
			loggo.Error("Unzip File Error : %v", err)
			continue
		}
		defer srcFile.Close()

		tmp := dest + "/" + innerFile.Name
		tmp = filepath.Clean(tmp)
		tmp = filepath.ToSlash(tmp)

		newFile, err := os.Create(tmp)
		if err != nil {
			loggo.Error("Unzip File Error : %v", err)
			continue
		}
		io.Copy(newFile, srcFile)
		newFile.Close()
	}
	return nil
}
