#! /bin/bash
set -e

REPO_ROOT=$(dirname "${BASH_SOURCE[0]}")/

# 换源
cp /etc/apt/sources.list /etc/apt/sources.list.bak
rm -f /etc/apt/sources.list
cat > /etc/apt/sources.list <<EOF
deb http://mirrors.tuna.tsinghua.edu.cn/ubuntu/ hirsute main restricted

deb http://mirrors.tuna.tsinghua.edu.cn/ubuntu/ hirsute-updates main restricted

deb http://mirrors.tuna.tsinghua.edu.cn/ubuntu/ hirsute universe

deb http://mirrors.tuna.tsinghua.edu.cn/ubuntu/ hirsute-updates universe

deb http://mirrors.tuna.tsinghua.edu.cn/ubuntu/ hirsute multiverse

deb http://mirrors.tuna.tsinghua.edu.cn/ubuntu/ hirsute-updates multiverse

deb http://mirrors.tuna.tsinghua.edu.cn/ubuntu/ hirsute-backports main restricted universe multiverse

deb http://mirrors.tuna.tsinghua.edu.cn/ubuntu/ hirsute-security main restricted

deb http://mirrors.tuna.tsinghua.edu.cn/ubuntu/ hirsute-security multiverse
EOF

apt-get update

# 下载wget
apt install wget

# 下载Go
cd $REPO_ROOT/../
wget https://go.dev/dl/go1.17.9.linux-amd64.tar.gz -t 3 -O go1.17.9.linux-amd64.tar.gz
tar -zxf go1.17.9.linux-amd64.tar.gz
rm -rf /usr/local/go
mv go /usr/local/go
echo 'export PATH=/usr/local/go/bin:$PATH' >> ~/.bashrc
echo 'export PATH=/root/go/bin:$PATH' >> ~/.bashrc
#source ~/.bashrc
#go env -w GOPROXY=https://goproxy.cn
#cd $REPO_ROOT
#

