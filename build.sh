#!/bin/bash

#自定义加载 golang.org/x/
myget(){
    m=$1
    path="$GOPATH/src/"$m
    if [[ ! -d "$path" ]]; then
        url='https://'${m/golang.org\/x/github.com\/golang}'.git'
        echo $url
        git clone $url $path
    fi
}



#修改 $GOPATH
cd `dirname $0`

CURDIR=`pwd`
OLDGOPATH="$GOPATH"
export GOPATH="$CURDIR"
#------

#get
echo '加载包...'
myget golang.org/x/sys
go get github.com/fsnotify/fsnotify

echo '加载完成!'
echo
#exit;

echo '编译 ...'

rm -f xwproxy_*
echo 'mac 64...'
GOOS=darwin GOARCH=amd64 go build  -o xwproxy_mac64  main.go
zip -r xwproxy_mac64.zip xwproxy_mac64

echo 'linux 64...'
GOOS=linux GOARCH=amd64 go build  -o xwproxy_linux64  main.go
tar zcvf xwproxy_linux64.tar.gz xwproxy_linux64

echo 'windows 64 ...'
GOOS=windows GOARCH=amd64 go build  -o xwproxy_win64.exe  main.go

echo 'windows 32 ...'
GOOS=windows GOARCH=386 go build  -o xwproxy_win32.exe  main.go

echo '编译完成!'
echo

#复原 $GOPATH
export GOPATH="$OLDGOPATH"


