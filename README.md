# XWProxy

## 一、简介

XWProxy是一个http(https)代理软件。
在开发web项目过程中, 经常需要切换各种开发环境、测试环境等等。结合此工具可以很方便的
配置和切换各种环境, 以及用于手机app测试.

![工作原理](https://github.com/zh-five/XWProxy/blob/master/work.png)

XWProxy的一些特性：

- 支持http和https
- 支持类似系统hosts功能的配置文件，可以很方便的配置域名指向的ip地址
- 配置域名时支持*号通配符，降低配置的繁琐程度
- 修改配置文件后实时生效，不需要重启代理服务
- 支持同时代理多个端口，并使用不同的hosts配置
- 支持代理转发时 `X-Forwarded-For`的设置，可以匿名代理、真实代理和伪装ip
- 支持macOS、Linux和Windows系统

## 二、应用场景

原则上负载不是特别高的代理需求XWProxy都能胜任，这里列举一些实际使用的例子：

- 在内网电脑上启动XWProxy服务（不会影响系统的hosts），为手机提供代理，测试网页或APP接口
- 本机启动XWProxy服务，结合浏览器代理插件（如 [SwitchyOmega](https://github.com/FelisCatus/SwitchyOmega) ）随时切换访问环境
- XWProxy启动多个端口配置多个环境的hosts，结合chrome浏览器不同用户间数据隔离的特性。创建多个本地账户并用插件配置到不一样的端口，从而达到不同的浏览器窗口访问不同的环境
- XWProxy启动在内网的某一个机器上，让团队成员（开发和测试）共用一套代理服务，保证环境配置的稳定性

## 三、安装使用

## 3.1 编译安装

以macOS或Linux为例

```bash
git clone https://github.com/zh-five/XWProxy
cd XWProxy

# 编译
sh ./build.sh

# 创建默认配置文件(配置文件不存在时, 创建默认配置)
./bin/xwproxy -c /your_path/xwp.conf

# 运行
./bin/xwproxy -c /your_path/xwp.conf

# 查看帮助文档
$ ./bin/xwproxy
Usage of ./xwproxy:
  -c string
    	必须.配置文件(指定文件不存在时将尝试创建默认配置文件
  -d	可选,是否后台运行
  -log string
    	可选.日志文件,后台运行时有效,无则不记录
  -t	检查指定配置文件是否有错误

XWProxy : http(https)代理工具
版本 1.0
项目主页 <https://github.com/zh-five/XWProxy>
问题反馈 <https://github.com/zh-five/XWProxy/issues>

```

## 3.2 直接下载编译好的可执行文件

下载地址: [https://github.com/zh-five/XWProxy/releases](https://github.com/zh-five/XWProxy/releases)

选择系统对应的可执行文件下载即可

## 3.3 启动

```bash
# 各系统下的可执行文件名称可能有些差别, 请注意替换, 参数是一样的

#1.创建配置文件
$ ./xwproxy -c proxy.txt
指定的配置文件不存在/data/git/github/zh-five/XWProxy/proxy.txt
是否尝试创建默认的配置文件[y/n]:y
已经成功创建默认配置文件: /data/git/github/zh-five/XWProxy/proxy.txt

#2.前台执行(若要后台运行,则需加上 -d 参数)
$ ./xwproxy -c proxy.txt
配置文件解析成功: /data/git/github/zh-five/XWProxy/proxy.txt
2020/06/16 16:47:33 HttpPxoy to runing on 127.0.0.1:8033

```

## 3.4 修改配置文件

自动创建的默认配置文件大约如下, 你可以把你的一个测试环境的hosts配置(如 `abc.com 192.168.6.33`)加入到文件中
代理会实时生效. 使用以下命令访问, 会请求到 `192.168.6.33`

`curl -x 127.0.0.1:8033 'http://abc.com'`

```txt
############################################################################
#  xwproxy 代理工具配置文件说明
#  1.'#'开头的行为注释
#  2.一个@addr选项对应一个代理配置, 至少一个, 可配置多个
#  3.修改配置文件后, 实时生效. 但@addr选项除外, 有增删或修改@addr时, 应重启服务
############################################################################


# 监听地址选项, 必须. (一个代理配置的开始)
#只允许本机访问
@addr = 127.0.0.1:8033  
#不限ip可以访问
#@addr = :8033  


# 转发ip选项, 可选, 默认为0.
# 影响转发请求时head里'X-Forwarded-For'的设置, 有三种取值:
# 0  : 不设置'X-Forwarded-For'
# 1  : 按照真实情况设置'X-Forwarded-For'
# 127.0.0.1 : 可以指定的为任意ip
@forwardedIP = 0


# 以下是指定host的配置, 格式兼容系统的hosts文件
# *可用于表示任何非点号(.)的1个或多个字符.
# 匹配时从上到下检查, 遇到第一个合格时停止检查. 建议严格的条件靠前放置
# 以下是几种配置示例(注意:删除ip前'#'才能生效)

#192.168.6.33		example.com
#192.168.6.24		a.example.com b.example.com
#192.168.6.34		xw.a.example.com
#192.168.6.33		xw.*.example.com


# 第2个代理配置开始
#@addr = 127.0.0.1:8024
#@forwardedIP = 1

#192.168.6.33 abc.com
```

### 3.4 各种使用场景下的代理配置

**本机管理各种开发测试环境**

1. 把系统的http代理配置到8033端口, 为加入到配置文件的地址不受影响. 切换环境时直接修改配置,实时生效(可行,但不推荐)
2. 代理程序配置多个端口, 每个端口配置不同的环境, 如8033是开发环境, 8034是测试环境等等.
   然后chrome浏览器安装代理插件SwitchyOmega, 设置多个用户, 每个用户配置不一样的代理端口. 这样要切换环境时, 使用不同用户的chrome浏览器窗口访问即可,不用频繁修改配置了.

**手机app测试**

注意要修改配置文件里的addr为 `:8033`, 然后手机连接在同一个局域网, 然后在手机上设置http代理即可

**多人共同测试**

为保证多人共同测试web服务时环境的一致性, 在一台机器上启动一个代理服务, 大家共一个代理服务, 则访问到的环境都是一致.代理服务的配置同上.
