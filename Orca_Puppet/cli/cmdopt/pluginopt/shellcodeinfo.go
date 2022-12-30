package pluginopt

type ShellcodeMetaInfo struct {
	Fid        string // 操作文件ID，随机生成的UUID
	Params     string // 参数
	SliceNum   int    // 基础分片数量
	SliceSize  int64  // 基础分片大小
	RemainSize int64  // 剩余分片大小
}
