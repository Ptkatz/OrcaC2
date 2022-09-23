md out\master && md out\server && md out\puppet

go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,direct

cd Orca_Master\ && go build -o ..\out\master\Orca_Master.exe -ldflags "-s -w" && cd ..
cd Orca_Server\ && go build -o ..\out\server\Orca_Server.exe -ldflags "-s -w" && cd ..
cd Orca_Puppet\ && go build -o ..\out\puppet\Orca_Puppet.exe -ldflags "-s -w" && cd ..

xcopy /s /y Orca_Master\3rd_party\ out\master\3rd_party\
xcopy /s /y Orca_Server\db\ out\server\db\
xcopy /s /y Orca_Server\conf\ out\server\conf\

pause