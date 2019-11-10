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
go get github.com/icattlecoder/godaemon
#go get github.com/kballard/go-shellquote
#go get github.com/openatx/androidutils
#go get github.com/sevlyar/go-daemon
#go get github.com/openatx/androidutils

echo '加载完成!'
echo
#exit;

echo '编译 ...'
go build  -o xwproxy  main.go
echo '编译完成!'
echo

#复原 $GOPATH
export GOPATH="$OLDGOPATH"


#./xwproxy -c ./hosts.txt
#./xwproxy -h