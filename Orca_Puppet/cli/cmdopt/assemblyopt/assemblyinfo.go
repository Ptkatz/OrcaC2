//go:build windows
// +build windows

package assemblyopt

import (
	"github.com/Ne0nd0g/go-clr"
	"sync"
)

type AssemblyMetaInfo struct {
	Fid        string // 操作文件ID，随机生成的UUID
	FileName   string // 程序名
	SliceNum   int    // 基础分片数量
	SliceSize  int64  // 基础分片大小
	RemainSize int64  // 剩余分片大小
}

type Assembly struct {
	name       string
	version    string
	methodInfo *clr.MethodInfo
}

type Results struct {
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
}

var Mutex sync.Mutex
var Assemblies = make(map[string]Assembly)
var runtimeHost *clr.ICORRuntimeHost
var redirected bool
var patched bool
