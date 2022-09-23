#bin/bash
# 编译Master所安装的依赖库
apt update -y && apt upgrade -y
apt-get install -y libxt-dev libxcursor-dev libxrandr-dev libxinerama-dev libglx-dev libglx-dev xorg-dev

export GO111MODULE=on
export GOPROXY=https://goproxy.cn

mkdir -p out/master out/server out/puppet 
cd Orca_Master && go build -o ../out/master/Orca_Master && cd ..
cd Orca_Server && go build -o ../out/server/Orca_Server && cd ..
cd Orca_Puppet && go build -o ../out/puppet/Orca_Puppet && cd ..

cp -r Orca_Master/3rd_party out/master/3rd_party
cp -r Orca_Server/db out/server/db
cp -r Orca_Server/conf out/server/conf
