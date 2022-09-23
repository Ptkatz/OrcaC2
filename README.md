# OrcaC2

## 简介

`OrcaC2`是一款基于Websocket加密通信的多功能C&C框架，使用Golang实现。

由三部分组成：`Orca_Server`(服务端)、`Orca_Master`(控制端)、`Orca_Puppet`(被控端)。

<img src="https://i.imgur.com/zQzBHAS.jpeg" alt="logo.jpg" style="zoom: 50%;" />

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
- 远程Shellcode、PE加载（支持的注入方式：CreateThread、CreateRemoteThread、BananaPhone、RtlCreateUserThread）（控制端为Windows系统）
- 正/反向代理、socks5正/反向代理（支持的协议：tcp、rudp(可靠udp)、ricmp(可靠icmp)、rhttp(可靠http)、kcp、quic）
- 多协程端口扫描（指纹识别端口信息）
- 多协程端口爆破（支持ftp、ssh、wmi、wmihash、smb、mssql、oracle、mysql、rdp、postgres、redis、memcached、mongodb）
- 远程ssh命令执行/文件上传/文件下载
- 远程smb命令执行(无回显)/文件上传（通过rpc服务执行命令，类似wmiexec；通过ipc$上传文件，类似psexec）



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



**由衷感谢以上项目的作者/团队对开源的贡献与支持**



## TODO

- [ ] 支持Websocket SSL
- [ ] Dump Hash
- [ ] Powershell模块加载
- [ ] 完善Linux-memfd无文件执行
- [ ] 内网中间人攻击
- [ ] Linux系统的屏幕截图
- [ ] 基于VNC的远程桌面
- [ ] Webshell管理
- [ ] WireGuard搭建隧道接入内网
- [ ] 对MacOS系统更多支持
- [ ] 根据payload生成被控端加载器
- [ ] 使用C实现远程加载器加载被控端，解决被控端体积过大问题
- [ ] 多端口监听器
- [ ] ...



> 由于最近要学Rust，这个项目最近不会有太多更新，见谅



## 已知Bug

- 使用`assembly invoke`功能调用部分C#程序时会出错，在工作中务必先进行试验，建议使用`3rd_party`下的C#程序
- 在linux下使用隐藏执行(`-hide`)时，调用`pty`功能时程序崩溃！
- 利用smb命令执行（`smb exec`）上线时，无法使用屏幕截图与屏幕控制功能



## 免责声明

本工具仅面向**合法授权**的企业安全建设行为，如您需要测试本工具的可用性，请自行搭建靶机环境。

在使用本工具进行检测时，您应确保该行为符合当地的法律法规，并且已经取得了足够的授权。***\*请勿对非授权目标进行扫描。\****

如您在使用本工具的过程中存在任何非法行为，您需自行承担相应后果，本人将不承担任何法律及连带责任。

在安装并使用本工具前，请您***\*务必审慎阅读、充分理解各条款内容\****，限制、免责条款或者其他涉及您重大权益的条款可能会以加粗、加下划线等形式提示您重点注意。 除非您已充分阅读、完全理解并接受本协议所有条款，否则，请您不要安装并使用本工具。您的使用行为或者您以其他任何明示或者默示方式表示接受本协议的，即视为您已阅读并同意本协议的约束。
