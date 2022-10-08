md out\master && md out\server && md out\master\puppet && md out\master\stub

go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,direct

cd Orca_Master\ && go build -o ..\out\master\Orca_Master_win_x64.exe -ldflags "-s -w" && cd ..
cd Orca_Server\ && go build -o ..\out\server\Orca_Server_win_x64.exe -ldflags "-s -w" && cd ..
cd Orca_Puppet\ && go build -o ..\out\master\puppet\Orca_Puppet_win_x64.exe -ldflags "-s -w" && cd ..
cd Orca_Loader\windows && gcc stub.c lib\x64\winhttp.lib -s -w -o stub_win_x64.exe && cd ..\..

xcopy /s /y Orca_Master\3rd_party\ out\master\3rd_party\
xcopy /s /y Orca_Server\db\ out\server\db\
xcopy /s /y Orca_Server\conf\ out\server\conf\
move /y Orca_Loader\windows\stub_win_x64.exe out\master\stub\stub_win_x64.exe

pause
