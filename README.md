# XWProxy
## 一、简介
XWProxy是我个人开发的一个http(https)代理小工具（取名困难症,
就把我名字的首字母xw加在Proxy单词前面了^_^）。 
在开发web项目过程中, 经常需要切换各种开发环境、测试环境等等。结合此工具可以很方便的
配置和切换各种环境

![工作原理](https://github.com/zh-five/XWProxy/blob/master/work.png)

XWProxy的一些特性：
- 支持http和https
- 支持类似系统hosts功能的配置文件，可以很方便的配置域名指向的ip地址
- 配置域名时支持*号通配符，降低配置的繁琐程度
- 修改配置文件后实时生效，不需要重启代理服务
- 支持同时代理多个端口，并使用不同的hosts配置
- 支持代理转发时`X-Forwarded-For`的设置，可以匿名代理、真实代理和伪装ip


## 二、应用场景
原则上负载不是特别高的代理需求XWProxy都能胜任，这里列举一些实际使用的例子：
- 在内网电脑上启动XWProxy服务（不会影响系统的hosts），为手机提供代理，测试网页或APP接口
- 本机启动XWProxy服务，结合浏览器代理插件（如 [SwitchyOmega](https://github.com/FelisCatus/SwitchyOmega) ）随时切换访问环境
- XWProxy启动多个端口配置多个环境的hosts，结合chrome浏览器不同用户间数据隔离的特性。创建多个本地账户并用插件配置到不一样的端口，从而达到不同的浏览器窗口访问不同的环境
- XWProxy启动在内网的某一个机器上，让团队成员（开发和测试）共用一套代理服务，保证环境配置的稳定性

## 三、安装使用
### 3.1 安装和启动XWProxy
（待完善）

### 3.2 浏览器代理配置
（待完善）

