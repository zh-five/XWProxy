#!/bin/bash


# 目录
cd `dirname $0`
mkdir bin 
echo
#exit;

echo '编译 ...'

rm -f xwproxy_*

if [ $# != 1 ]; then
    go build -o ./bin/xwproxy main.go
else
    echo 'mac 64...'
    GOOS=darwin GOARCH=amd64 go build  -o ./bin/xwproxy_mac64  main.go
    #zip -r xwproxy_mac64.zip xwproxy_mac64

    echo 'linux 64...'
    GOOS=linux GOARCH=amd64 go build  -o ./bin/xwproxy_linux64  main.go
    #tar zcvf xwproxy_linux64.tar.gz xwproxy_linux64

    echo 'windows 64 ...'
    GOOS=windows GOARCH=amd64 go build  -o ./bin/xwproxy_win64.exe  main.go

    echo 'windows 32 ...'
    GOOS=windows GOARCH=386 go build  -o ./bin/xwproxy_win32.exe  main.go

fi

echo '编译完成!'
echo



