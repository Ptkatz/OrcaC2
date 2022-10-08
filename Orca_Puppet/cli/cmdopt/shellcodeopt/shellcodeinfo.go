//go:build windows
// +build windows

package shellcodeopt

type ShellcodeMetaInfo struct {
	Fid        string // 操作文件ID，随机生成的UUID
	LoadFunc   string // 加载器类型
	Pid        int    // 注入的pid
	SliceNum   int    // 基础分片数量
	SliceSize  int64  // 基础分片大小
	RemainSize int64  // 剩余分片大小
}
