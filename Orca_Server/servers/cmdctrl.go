package servers

import (
	"Orca_Server/cli/cmdopt/fileopt"
	"Orca_Server/cli/cmdopt/hostopt"
	"Orca_Server/cli/cmdopt/portopt/portcrackopt"
	"Orca_Server/cli/cmdopt/portopt/portscanopt"
	"Orca_Server/cli/cmdopt/proxyopt"
	"Orca_Server/cli/cmdopt/sshopt"
	"Orca_Server/cli/common/setchannel"
	"Orca_Server/define/colorcode"
	"Orca_Server/define/config"
	"Orca_Server/define/retcode"
	"Orca_Server/pkg/go-engine/loggo"
	"Orca_Server/pkg/go-engine/proxy"
	"Orca_Server/setting"
	"Orca_Server/sqlmgmt"
	"Orca_Server/tools/crypto"
	"Orca_Server/tools/qqwry"
	"Orca_Server/tools/util"
	"encoding/json"
	"fmt"
	"github.com/4dogs-cn/TXPortMap/pkg/output"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func hostCmd(decData string) {
	var hostInfo hostopt.HostInfo
	err := json.Unmarshal([]byte(decData), &hostInfo)
	if err != nil {
		return
	}
	sqlmgmt.AddHost(hostInfo)
	data, _ := crypto.Encrypt([]byte(decData), []byte(setting.CommonSetting.CryptoKey))
	SendMessage2System(hostInfo.SystemId, "Server", retcode.ONLINE_MESSAGE_CODE, "online", data)
	return
}

func ipSearchCmd(data string, clientId string) {
	qqwry.IPData.FilePath = config.GeoIpDB
	startTime := time.Now().UnixNano()
	res := qqwry.IPData.InitIPData()
	if v, ok := res.(error); ok {
		logrus.Printf(v.Error())
	}
	endTime := time.Now().UnixNano()
	logrus.Printf("IP 库加载完成 共加载:%d 条 IP 记录, 所花时间:%.1f ms", qqwry.IPData.IPNum, float64(endTime-startTime)/1000000)
	q := qqwry.NewQQwry()
	r := q.Find(data)
	retData := fmt.Sprintf("%s %s", r.Country, r.Area)
	SendMessage2Client(clientId, "Server", retcode.SUCCESS, "ipSearch_ret", &retData)
}

func listCmd(clientId string) {
	hosts := sqlmgmt.ListHosts()
	hostsStr := string(hosts)
	SendMessage2Client(clientId, "Server", retcode.SUCCESS, "listHosts_ret", &hostsStr)
}

func proxyServerStartCmd(decData, clientId string) {
	loggo.Ini(loggo.Config{
		Level:     1,
		Prefix:    "orca",
		MaxDay:    3,
		NoLogFile: true,
		NoPrint:   false,
	})
	var proxyServerParam proxyopt.ProxyServerParam
	err := json.Unmarshal([]byte(decData), &proxyServerParam)
	if err != nil {
		return
	}
	key := proxyServerParam.Key
	protos := []string{proxyServerParam.Proto}
	listenaddrs := []string{proxyServerParam.Listen}
	defConfig := proxy.DefaultConfig()
	defConfig.Key = key
	defConfig.Encrypt = setting.CommonSetting.CryptoKey
	server, err := proxy.NewServer(defConfig, protos, listenaddrs)
	var proxyServer = proxyopt.ProxyServer{
		Uid:    util.GenUUID(),
		Server: server,
		Param:  proxyServerParam,
	}
	if err != nil {
		retData := fmt.Sprintf("main NewServer fail %s", err.Error())
		SendMessage2Client(clientId, "Server", retcode.FAIL, "proxyServerStart_ret", &retData)
		return
	}
	proxyopt.ProxyServerLists = append(proxyopt.ProxyServerLists, proxyServer)
	retData := "proxy server start!"
	SendMessage2Client(clientId, "Server", retcode.SUCCESS, "proxyServerStart_ret", &retData)
}

func proxyServerListCmd(clientId string) {
	proxyList := proxyopt.ProxyServerLists
	marshal, _ := json.Marshal(proxyList)
	data, _ := crypto.Encrypt(marshal, []byte(setting.CommonSetting.CryptoKey))
	SendMessage2Client(clientId, "Server", retcode.SUCCESS, "proxyServerList_ret", &data)
}

func proxyServerCloseCmd(decData string) {
	index := -1
	closeFlag := false
	for i, server := range proxyopt.ProxyServerLists {
		if decData == server.Uid {
			server.Server.Close()
			index = i
			closeFlag = true
		}
	}
	if closeFlag {
		proxyopt.ProxyServerLists = append(proxyopt.ProxyServerLists[:index], proxyopt.ProxyServerLists[index+1:]...)
	}
}

func proxyClientStartCmd(decData, clientId string) {
	fmt.Println(decData)
	var proxyClientLists []proxyopt.ProxyClient
	json.Unmarshal([]byte(decData), &proxyClientLists)
	for _, client := range proxyClientLists {
		sqlmgmt.AddProxyClient(client)
	}
}

func proxyClientListCmd(clientId string) {
	proxyLists := sqlmgmt.ListProxy()
	proxyListsStr := string(proxyLists)
	fmt.Println(proxyListsStr)
	SendMessage2Client(clientId, "Server", retcode.SUCCESS, "proxyClientList_ret", &proxyListsStr)
}

func proxyClientCloseCmd(decData string) {
	sqlmgmt.DelProxyRecordByUid(decData)
}

func sshConnTestCmd(clientId, decData string) {
	var sshStruct sshopt.SshStruct
	json.Unmarshal([]byte(decData), &sshStruct)
	username := sshStruct.SSHUser
	password := sshStruct.SSHPwd
	ip, port, _ := strings.Cut(sshStruct.SSHHost, ":")
	client := sshopt.NewSSHClient(username, password, ip, port)
	err := client.Connect()
	if err != nil {
		outputMsg := colorcode.OutputMessage(colorcode.SIGN_FAIL, err.Error())
		retData, _ := crypto.Encrypt([]byte(outputMsg), []byte(setting.CommonSetting.CryptoKey))
		SendMessage2Client(clientId, "Server", retcode.FAIL, "sshTestConn_ret", &retData)
	} else {
		outputMsg := colorcode.OutputMessage(colorcode.SIGN_SUCCESS, "ssh connection is successful")
		retData, _ := crypto.Encrypt([]byte(outputMsg), []byte(setting.CommonSetting.CryptoKey))
		SendMessage2Client(clientId, "Server", retcode.SUCCESS, "sshTestConn_ret", &retData)
	}
}

func sshRunCmd(clientId, decData string) {
	var sshRun sshopt.SshRunStruct
	json.Unmarshal([]byte(decData), &sshRun)
	username := sshRun.SshStruct.SSHUser
	password := sshRun.SshStruct.SSHPwd
	ip, port, _ := strings.Cut(sshRun.SshStruct.SSHHost, ":")
	client := sshopt.NewSSHClient(username, password, ip, port)

	cmd := sshRun.Command
	backinfo, err := client.Run(cmd)
	if err != nil {
		message := fmt.Sprintf("failed to run shell: %v", err)
		outputMsg := colorcode.OutputMessage(colorcode.SIGN_FAIL, message)
		retData, _ := crypto.Encrypt([]byte(outputMsg), []byte(setting.CommonSetting.CryptoKey))
		SendMessage2Client(clientId, "Server", retcode.FAIL, "sshRun_ret", &retData)
		return
	}
	message := fmt.Sprintf("'%v' back info: \n%v", cmd, backinfo)
	outputMsg := colorcode.OutputMessage(colorcode.SIGN_SUCCESS, message)
	retData, _ := crypto.Encrypt([]byte(outputMsg), []byte(setting.CommonSetting.CryptoKey))
	SendMessage2Client(clientId, "Server", retcode.SUCCESS, "sshRun_ret", &retData)
}

func sshUploadCmd(clientId, decData string) {
	m, exist := setchannel.GetFileSliceDataChan(clientId)
	if !exist {
		m = make(chan interface{})
		setchannel.AddFileSliceDataChan(clientId, m)
	}
	var sshUploadStruct sshopt.SshUploadStruct
	err := json.Unmarshal([]byte(decData), &sshUploadStruct)
	if err != nil {
		return
	}
	username := sshUploadStruct.SshStruct.SSHUser
	password := sshUploadStruct.SshStruct.SSHPwd
	ip, port, _ := strings.Cut(sshUploadStruct.SshStruct.SSHHost, ":")
	client := sshopt.NewSSHClient(username, password, ip, port)

	fileMetaInfo := sshUploadStruct.FileMetaInfo
	saveFile := fileMetaInfo.SaveFileName
	sliceNum := fileMetaInfo.SliceNum

	var fileByte []byte
	// 循环获取分片数据
	for i := 0; i < sliceNum+1; i++ {
		select {
		case metaData := <-m:
			fileByte = append(fileByte, metaData.([]byte)...)
		case <-time.After(5 * time.Second):
			setchannel.DeleteFileSliceDataChan(clientId)
			return
		}
	}
	_, err = client.UploadFile(fileByte, saveFile)
	if err != nil {
		message := fmt.Sprintf("upload failed: %v", err)
		outputMsg := colorcode.OutputMessage(colorcode.SIGN_FAIL, message)
		retData, _ := crypto.Encrypt([]byte(outputMsg), []byte(setting.CommonSetting.CryptoKey))
		SendMessage2Client(clientId, "Server", retcode.FAIL, "sshUpload_ret", &retData)
		return
	} else {
		outputMsg := colorcode.OutputMessage(colorcode.SIGN_SUCCESS, "file upload success")
		retData, _ := crypto.Encrypt([]byte(outputMsg), []byte(setting.CommonSetting.CryptoKey))
		SendMessage2Client(clientId, "Server", retcode.FAIL, "sshUpload_ret", &retData)
	}
}

func sshDownloadCmd(clientId, decData string) {
	var sshDownloadStruct sshopt.SshDownloadStruct
	json.Unmarshal([]byte(decData), &sshDownloadStruct)
	username := sshDownloadStruct.SshStruct.SSHUser
	password := sshDownloadStruct.SshStruct.SSHPwd
	ip, port, _ := strings.Cut(sshDownloadStruct.SshStruct.SSHHost, ":")
	client := sshopt.NewSSHClient(username, password, ip, port)
	requestFile := sshDownloadStruct.RequestFile
	remoteFile := requestFile.DestFileName
	saveFile := requestFile.SaveFileName
	fileInfo, err := client.LStateFile(remoteFile)
	if err != nil {
		message := fmt.Sprintf("download failed: %v", err)
		outputMsg := colorcode.OutputMessage(colorcode.SIGN_FAIL, message)
		retData, _ := crypto.Encrypt([]byte(outputMsg), []byte(setting.CommonSetting.CryptoKey))
		SendMessage2Client(clientId, "Server", retcode.FAIL, "sshDownload_ret", &retData)
		return
	}
	if fileInfo.IsDir() {
		outputMsg := colorcode.OutputMessage(colorcode.SIGN_FAIL, "the requested file is a directory")
		retData, _ := crypto.Encrypt([]byte(outputMsg), []byte(setting.CommonSetting.CryptoKey))
		SendMessage2Client(clientId, "Server", retcode.FAIL, "sshDownload_ret", &retData)
		return
	}

	// 发送文件元信息
	fileMetaInfo := sshopt.GetFileMetaInfo(fileInfo, saveFile)
	sliceNum := fileMetaInfo.SliceNum
	sliceSize := fileMetaInfo.SliceSize
	remainSize := fileMetaInfo.RemainSize
	metaInfo, err := json.Marshal(fileMetaInfo)
	data, _ := crypto.Encrypt(metaInfo, []byte(setting.CommonSetting.CryptoKey))
	SendMessage2Client(clientId, "Server", retcode.SUCCESS, "fileMetaInfo", &data)
	// 发送文件分片
	if client.SshClient == nil {
		if err = client.Connect(); err != nil {
			return
		}
	}
	pUploadFile, err := client.SftpClient.Open(remoteFile)
	if err != nil {
		return
	}
	defer pUploadFile.Close()
	defer client.SshClient.Close()
	defer client.SftpClient.Close()
	for i := 0; i < sliceNum; i++ {
		sliceData := make([]byte, sliceSize)
		_, err = pUploadFile.Read(sliceData)
		if err != nil && err != io.EOF {
			return
		}
		encData, _ := crypto.Encrypt(sliceData, []byte(setting.CommonSetting.CryptoKey))
		encMsg, _ := crypto.Encrypt([]byte("sliceData"), []byte(setting.CommonSetting.CryptoKey))
		SendMessage2Client(clientId, "Server", retcode.SUCCESS, encMsg, &encData)
	}
	// 处理最后一个分片
	sliceData := make([]byte, remainSize)
	_, err = pUploadFile.Read(sliceData)
	if err != nil && err != io.EOF {
		return
	}
	encData, _ := crypto.Encrypt(sliceData, []byte(setting.CommonSetting.CryptoKey))
	encMsg, _ := crypto.Encrypt([]byte("sliceData"), []byte(setting.CommonSetting.CryptoKey))
	SendMessage2Client(clientId, "Server", retcode.SUCCESS, encMsg, &encData)
}

func sshTunnelStartCmd(clientId, decData string) {
	var err error
	var sshTunnelStruct sshopt.SshTunnelStruct
	json.Unmarshal([]byte(decData), &sshTunnelStruct)
	username := sshTunnelStruct.SshStruct.SSHUser
	password := sshTunnelStruct.SshStruct.SSHPwd
	ip, port, _ := strings.Cut(sshTunnelStruct.SshStruct.SSHHost, ":")
	target := sshTunnelStruct.Target
	source := sshTunnelStruct.Source
	client := sshopt.NewSSHClient(username, password, ip, port)
	err = client.Connect()
	if err != nil {
		outputMsg := colorcode.OutputMessage(colorcode.SIGN_FAIL, err.Error())
		retData, _ := crypto.Encrypt([]byte(outputMsg), []byte(setting.CommonSetting.CryptoKey))
		SendMessage2Client(clientId, "Server", retcode.FAIL, "sshTunnelStart_ret", &retData)
		return
	}

	sshTunnel := client.Cli2Tunnel(target, source)
	sshTunnel.Log = log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds)
	go func() {
		err = sshTunnel.Start()
	}()
	if err != nil {
		outputMsg := colorcode.OutputMessage(colorcode.SIGN_FAIL, "tunnel open failed")
		retData, _ := crypto.Encrypt([]byte(outputMsg), []byte(setting.CommonSetting.CryptoKey))
		SendMessage2Client(clientId, "Server", retcode.FAIL, "sshTunnelStart_ret", &retData)
		return
	}
	for _, recordList := range sshopt.SshTunnelRecordLists {
		if recordList.SshTunnelBaseRecord.Target == target && recordList.SshTunnelBaseRecord.Source == source {
			outputMsg := colorcode.OutputMessage(colorcode.SIGN_FAIL, "the tunnel is repeated")
			retData, _ := crypto.Encrypt([]byte(outputMsg), []byte(setting.CommonSetting.CryptoKey))
			SendMessage2Client(clientId, "Server", retcode.FAIL, "sshTunnelStart_ret", &retData)
			return
		}
	}
	message := fmt.Sprintf("tunnel open successfully: %s --> %s", target, source)
	outputMsg := colorcode.OutputMessage(colorcode.SIGN_SUCCESS, message)
	retData, _ := crypto.Encrypt([]byte(outputMsg), []byte(setting.CommonSetting.CryptoKey))
	SendMessage2Client(clientId, "Server", retcode.SUCCESS, "sshTunnelStart_ret", &retData)
	sshTunnelBaseRecord := sshopt.SshTunnelBaseRecord{
		Uid:      util.GenUUID(),
		ClientId: "Server",
		Source:   source,
		Target:   target,
	}
	sshTunnelRecord := sshopt.SshTunnelRecord{
		SSHTunnel:           sshTunnel,
		SshTunnelBaseRecord: sshTunnelBaseRecord,
	}
	sshopt.SshTunnelRecordLists = append(sshopt.SshTunnelRecordLists, sshTunnelRecord)
	sqlmgmt.AddSshTunnel(sshTunnelBaseRecord)
}

func sshTunnelListCmd(clientId string) {
	sshTunnelLists := sqlmgmt.ListSshTunnel()
	sshTunnelListsStr := string(sshTunnelLists)
	fmt.Println(sshTunnelListsStr)
	SendMessage2Client(clientId, "Server", retcode.SUCCESS, "sshTunnelList_ret", &sshTunnelListsStr)
}

func sshTunnelCloseCmd(decData string) {
	index := -1
	closeFlag := false
	for i, record := range sshopt.SshTunnelRecordLists {
		if decData == record.SshTunnelBaseRecord.Uid {
			record.SSHTunnel.Close()
			index = i
			closeFlag = true
		}
	}
	if closeFlag {
		sshopt.SshTunnelRecordLists = append(sshopt.SshTunnelRecordLists[:index], sshopt.SshTunnelRecordLists[index+1:]...)
	}
	sqlmgmt.DelSshTunnelRecordByUid(decData)
}

func sshTunnelAddCmd(decData string) {
	var sshTunnelBaseRecord sshopt.SshTunnelBaseRecord
	json.Unmarshal([]byte(decData), &sshTunnelBaseRecord)
	sqlmgmt.AddSshTunnel(sshTunnelBaseRecord)
}

func sshTunnelDelCmd(decData string) {
	sqlmgmt.DelSshTunnelRecordByUid(decData)
}

func portScanCmd(clientId, decData string) {
	portscanopt.ResultEvents = make([]*output.ResultEvent, 0)
	var scanCmdMsg portscanopt.ScanCmdMsg
	json.Unmarshal([]byte(decData), &scanCmdMsg)
	portscanopt.CheckSliceValue(&scanCmdMsg.CmdIps)
	portscanopt.CheckSliceValue(&scanCmdMsg.CmdPorts)
	portscanopt.CheckSliceValue(&scanCmdMsg.ExcIps)
	portscanopt.CheckSliceValue(&scanCmdMsg.ExcPorts)
	portscanopt.Init(scanCmdMsg.CmdIps, scanCmdMsg.CmdPorts, scanCmdMsg.CmdT1000, scanCmdMsg.CmdRandom, scanCmdMsg.NumThreads, scanCmdMsg.Limit, scanCmdMsg.ExcIps, scanCmdMsg.ExcPorts, "", false, true, "", "", scanCmdMsg.Tout, scanCmdMsg.Nbtscan)
	engine := portscanopt.CreateEngine()
	// 命令行参数错误
	if err := engine.Parser(); err != nil {
		outputMsg, _ := json.Marshal(err.Error())
		retData, _ := crypto.Encrypt(outputMsg, []byte(setting.CommonSetting.CryptoKey))
		SendMessage2Client(clientId, "Server", retcode.FAIL, "portScan_ret", &retData)
		return
	}
	engine.Run()
	// 等待扫描任务完成
	engine.Wg.Wait()
	outputMsg, _ := json.Marshal(portscanopt.ResultEvents)
	retData, _ := crypto.Encrypt(outputMsg, []byte(setting.CommonSetting.CryptoKey))
	SendMessage2Client(clientId, "Server", retcode.SUCCESS, "portScan_ret", &retData)
	if portscanopt.Writer != nil {
		portscanopt.Writer.Close()
	}
}

func portCrackCmd(clientId, decData string) {
	var options *portcrackopt.Options
	json.Unmarshal([]byte(decData), &options)
	newRunner, err := portcrackopt.NewRunner(options)
	if err != nil {
		msg := fmt.Sprintf("Could not create runner: %v", err)
		outputMsg := colorcode.OutputMessage(colorcode.SIGN_FAIL, msg)
		retData, _ := crypto.Encrypt([]byte(outputMsg), []byte(setting.CommonSetting.CryptoKey))
		SendMessage2Client(clientId, "Server", retcode.FAIL, "portCrack_ret", &retData)
		return
	}
	outputMsg := newRunner.Run()
	retData, _ := crypto.Encrypt([]byte(outputMsg), []byte(setting.CommonSetting.CryptoKey))
	SendMessage2Client(clientId, "Server", retcode.SUCCESS, "portCrack_ret", &retData)
}

func fileUploadCmd(clientId, decData string) {
	// 获取id对应的管道
	m, exist := setchannel.GetFileSliceDataChan(clientId)
	if !exist {
		m = make(chan interface{})
		setchannel.AddFileSliceDataChan(clientId, m)
	}
	defer setchannel.DeleteFileSliceDataChan(clientId)
	// 获取文件元信息
	var fileMetaInfo fileopt.FileMetaInfo
	err := json.Unmarshal([]byte(decData), &fileMetaInfo)
	if err != nil {
		return
	}
	if !fileopt.IsDir("files") {
		err = os.Mkdir("files", 0666)
		if err != nil {
			fmt.Errorf("%s", err)
		}
	}
	saveFile, _ := filepath.Abs("files/" + fileMetaInfo.SaveFileName)
	sliceNum := fileMetaInfo.SliceNum
	md5sum := fileMetaInfo.Md5sum

	fmt.Println(saveFile, sliceNum, md5sum)
	pSaveFile, _ := os.OpenFile(saveFile, os.O_CREATE|os.O_RDWR, 0600)
	defer pSaveFile.Close()

	// 循环获取分片数据
	for i := 0; i < sliceNum+1; i++ {
		select {
		case metaData := <-m:
			pSaveFile.Write(metaData.([]byte))
		case <-time.After(5 * time.Second):
			SendMessage2Client(clientId, "Server", retcode.FAIL, "file upload failed", nil)
			setchannel.DeleteFileSliceDataChan(clientId)
			return
		}
	}
	saveFileMd5 := util.GetFileMd5Sum(saveFile)
	if md5sum == saveFileMd5 {
		data := colorcode.OutputMessage(colorcode.SIGN_SUCCESS, "file upload success")
		outputMsg, _ := crypto.Encrypt([]byte(data), []byte(setting.CommonSetting.CryptoKey))
		SendMessage2Client(clientId, "Server", retcode.SUCCESS, "fileUpload_ret", &outputMsg)
	} else {
		data := colorcode.OutputMessage(colorcode.SIGN_FAIL, "file upload failed")
		outputMsg, _ := crypto.Encrypt([]byte(data), []byte(setting.CommonSetting.CryptoKey))
		SendMessage2Client(clientId, "Server", retcode.FAIL, "fileUpload_ret", &outputMsg)
	}
}
