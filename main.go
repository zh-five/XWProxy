package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"xwproxy"
)

//程序主动启动后台运行的标识参数
const DAEMON_ARG = "DAEMON"

func main() {
	// 参数处理
	file := flag.String("c", "", "配置文件(不存在时将尝试创建默认配置文件")
	isStop := flag.Bool("s", false, "停止程序")
	isTest := flag.Bool("t", false, "测试配置文件")
	isDaemon := flag.Bool("d", false, "后台运行")
	usage()      //修改默认的帮助信息
	flag.Parse() //解析

	//配置文件路径检查和处理
	if *file == "" {
		flag.Usage()
		return
	}
	absFile, err := filepath.Abs(*file)
	if err != nil {
		fmt.Println("配置文件路径错误:", err)
		return
	}

	//停止
	if *isStop {
		toStop(absFile)
		return
	}

	//配置文件不存在的处理
	if !checkFile(absFile) {
		return
	}

	//解析配置文件
	fCfg := &xwproxy.FileConfig{File: absFile, IsDebug: *isTest || !*isDaemon}
	ok := fCfg.Parse() //解析
	if !ok {
		fmt.Println("配置文件有错误, 退出!")
		return
	}

	//测试
	if *isTest {
		return
	}

	//后台运行()
	if *isDaemon && !inSlice(DAEMON_ARG, flag.Args()) {
		toDaemon(absFile)
		return
	}

	//启动代理服务
	xwproxy.Run(fCfg)

	//监视配置文件变动
	fCfg.Watch()

}

//修改默认的帮助信息
func usage() {
	flag.Usage = func() {
		desc := `
XWProxy : http(https)代理工具
版本 0.10 
项目主页 <https://github.com/zh-five/XWProxy>
问题反馈 <https://github.com/zh-five/XWProxy/issues>
`
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintln(flag.CommandLine.Output(), desc)
	}
}

//切片中是否有某个值
func inSlice(val string, s []string) bool {
	for _, v := range s {
		if val == v {
			return true
		}
	}

	return false
}

//后台启动
func toDaemon(confFile string) {
	f, _ := filepath.Abs(os.Args[0]) //绝对路径
	args := append(os.Args[1:], DAEMON_ARG, ">", "/dev/null", "2>&1", "&")

	//启动
	cmd := exec.Command(f, args ...)
	cmd.Start()

	//保存pid
	ioutil.WriteFile(confFile+".pid", []byte(strconv.Itoa(cmd.Process.Pid)), 0666)

	fmt.Println("启动后台进程:", cmd.Process.Pid)
}

//停止后台程序
func toStop(confFile string) {
	pid, err := ioutil.ReadFile(confFile+".pid")
	if err != nil {
		fmt.Println("停止后台程序失败:", err)
		return
	}

	cmd := exec.Command("kill", string(pid))
	cmd.Start()
}

//检查配置文件, 不存在则尝试创建
func checkFile(absFile string) bool {
	_,err := os.Stat(absFile)
	if err == nil {
		return true
	}

	//错误原因不是不存在
	if os.IsExist(err) {
		fmt.Println(err)
		return false
	}

	//
	fmt.Print("指定的配置文件不存在", absFile, "\n是否尝试创建默认的配置文件[y/n]:")
	reader := bufio.NewReader(os.Stdin)
	y, _ := reader.ReadString('\n')
	if strings.Trim(y, " \n\r\t") != "y" {
		fmt.Println("放弃创建默认配置文件, 退出")
		return false
	}

	err = ioutil.WriteFile(absFile,confFileData(), 0666)
	if err == nil {
		fmt.Println("已经成功创建默认配置文件:", absFile)
	} else {
		fmt.Println("创建默认配置文件失败:", err)
	}

	return false
}

func confFileData() []byte{
	text := `############################################################################
#  xwproxy 代理工具配置文件说明
#  1.'#'开头到行尾是注释内容,解析时被忽略
#  2.但'#!'开头的行为选项配置行
#  3.一个addr选项对应一个代理配置, 至少一个, 可配置多个
#  4.修改配置文件后, 实时生效. 但addr选项除外, 有增删或修改addr时, 应重启服务
############################################################################


# 监听地址选项, 必须. (一个代理配置的开始)
# 包含ip地址和端口号
#! addr = 127.0.0.1:8033

# 转发ip选项, 可选, 默认为0.
# 影响转发请求时head里'X-Forwarded-For'的设置, 有三种取值:
# 0  : 不设置'X-Forwarded-For'
# 1  : 按照真实情况设置'X-Forwarded-For'
# ip : 设置为指定的ip
#! forwardedIP = 0


# 以下是指定host的配置, 格式兼容系统的hosts文件
# *可用于表示任何非点号(.)的1个或多个字符.
# 匹配时从上到下检查, 遇到第一个合格时停止检查. 建议严格的条件靠前放置
# 以下是几种配置示例(注意:删除ip前'#'才能生效)

#192.168.6.33		example.com
#192.168.6.24		a.example.com b.example.com
#192.168.6.34		xw.a.example.com
#192.168.6.33		xw.*.example.com


# 第2个代理配置开始
##! addr = 127.0.0.1:8024
##! forwardedIP = 1

#192.168.6.33 abc.com
`

	return []byte(text)
}