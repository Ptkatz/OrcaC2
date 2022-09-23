package ms

// 此文件提供SMB、NT相关的错误代码以及错误信息
// https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-erref/1bc92ddf-b79e-413c-bbaa-99a5281a6c90

const (
	STATUS_SUCCESS                  = 0x00000000
	STATUS_MORE_PROCESSING_REQUIRED = 0xC0000016
	STATUS_ACCESS_DENIED            = 0xC0000022
	STATUS_LOGON_FAILURE            = 0xC000006D
	STATUS_BAD_NETWORK_NAME         = 0xC00000CC
	STATUS_USER_SESSION_DELETED     = 0xC0000203 // 该问题似乎可以随机化ClientGuid解决
	STATUS_FILE_CLOSED              = 0xC0000128
	STATUS_PIPE_DISCONNECTED        = 0xC00000B0
	STATUS_INVALID_PARAMETER        = 0xC000000D
	STATUS_OBJECT_NAME_NOT_FOUND    = 0xC0000034
	STATUS_PIPE_BROKEN              = 0xC000014B
	NCA_S_FAULT_NDR                 = 0x0000006f7
	ERROR_FILE_NOT_FOUND            = 2
)

var StatusMap = map[uint32]string{
	STATUS_SUCCESS:                  "Requested operation succeeded.",
	STATUS_MORE_PROCESSING_REQUIRED: "More Processing Required",
	STATUS_ACCESS_DENIED:            "A process has requested access to an object but has not been granted those access rights.",
	STATUS_LOGON_FAILURE:            "Authentication failed.",
	STATUS_BAD_NETWORK_NAME:         "The specified share name cannot be found on the remote server.",
	STATUS_USER_SESSION_DELETED:     "STATUS_USER_SESSION_DELETED.",
	STATUS_FILE_CLOSED:              "An I/O request other than close and several other special case operations was attempted using a file object that had already been closed.",
	STATUS_PIPE_DISCONNECTED:        "The specified named psipe is in the disconnected state.",
	STATUS_INVALID_PARAMETER:        "An invalid parameter was passed to a service or function.",
	STATUS_OBJECT_NAME_NOT_FOUND:    "The object name is not found.",
	STATUS_PIPE_BROKEN:              "The pipe operation has failed because the other end of the pipe has been closed.",
	NCA_S_FAULT_NDR:                 "The stub data is invalid, doesn't match with the IDL definition",
	ERROR_FILE_NOT_FOUND:            "The system cannot find the file specified.",
}
