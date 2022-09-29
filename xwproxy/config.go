package xwproxy

//解析host文件

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
	"sync"
)

type oneHost struct {
	ip     string //对应ip
	domain string //域名
	mold   int    //类型. 0全字符串对比; 1后缀对比; 2前缀对比; 3正则对比
}

type FileConfig struct {
	File    string                       //配置文件路径
	IsDebug bool                         //是否输出调试信息
	options map[string]map[string]string //参数选项 [addr] => {key => val}
	hosts   map[string][]*oneHost        //指定地址转发的域名信息, 类似系统的hosts文件 [addr] => [*oneHost ...]
	err     bool                         //解析是否有错误
	addr    string                       //解析过程中用于记录当前的addr
	sync.RWMutex
}

//解析hosts文件, 每次解析都先清空原内容
//解析有错误返回false, 无则返回true
func (cfg *FileConfig) Parse() bool {
	cfg.Lock()
	defer cfg.Unlock()

	//初始化选项列表
	cfg.init()

	//读取
	buf, err := ioutil.ReadFile(cfg.File)
	if err != nil {
		cfg.debug(err)
		return false
	}
	str := string(buf)

	//解析
	lines := strings.Split(str, "\n") //分割为行
	n := len(lines)
	for i := 0; i < n; i++ {
		l := strings.Trim(lines[i], " \t\r")
		if l == "" {
			continue
		}

		if l[0] == '#' { //注释
			continue
		}

		if len(l) > 0 && l[0] == '@' { //选项配置行 '@addr = xxx'
			cfg.parseOption(l)
			continue
		}

		//解析host行
		cfg.parseHost(l)
	}

	//检查必须选项配置
	//检查addr
	if len(cfg.options) == 0{
		cfg.err = true
		cfg.debug("配置文件错误, 无任何有效代理:", cfg.File)
		return false
	}

	//解析结论
	if cfg.err == true {
		cfg.debug("配置文件解析有错误:", cfg.File)
	} else {
		cfg.debug("配置文件解析成功:", cfg.File)
	}

	return !cfg.err
}

func (cfg *FileConfig) Watch() {
	watch(cfg)
}

//获取选项配置
func (cfg *FileConfig) GetOption(addr, key string) string {
	//未找到对应代理的配置(可能启动后修改了配置文件)
	options, ok := cfg.options[addr]
	if !ok{
		return ""
	}

	val, ok := options[key]
	if ok {
		return val
	}

	return ""
}

//依据hosts文件获取对应ip, 无则返回空字符串
func (cfg *FileConfig) GetIP(addr, domain string) string {
	cfg.RLock()
	defer cfg.RUnlock()

	//未找到对应代理的配置(可能启动后修改了配置文件)
	hosts, ok := cfg.hosts[addr]
	if !ok{
		return ""
	}

	num := len(hosts)
	for i := 0; i < num; i++ {
		if hosts[i].mold == 0 {
			if hosts[i].domain == domain {
				return hosts[i].ip
			}
		} else if hosts[i].mold == 1 { //后缀对比
			if strings.HasSuffix(domain, hosts[i].domain) {
				return hosts[i].ip
			}
		} else if hosts[i].mold == 2 { //前缀对比
			if strings.HasPrefix(domain, hosts[i].domain) {
				return hosts[i].ip
			}
		} else { //正则对比
			if ok, _ := regexp.MatchString(hosts[i].domain, domain); ok {
				return hosts[i].ip
			}
		}
	}

	return ""
}

//解析前初始化
func (cfg *FileConfig) init() {
	cfg.options = make(map[string]map[string]string)
	cfg.hosts = make(map[string][]*oneHost)
	cfg.err = false
	cfg.addr = ""
}

//默认配置列表
func (cfg *FileConfig) defaultOptions() map[string]string {
	m := make(map[string]string)
	m["forwardedIP"] = "0"  //默认匿名

	return m
}

//解析option行 : '#! addr = 127.0.0.1 '
func (cfg *FileConfig) parseOption(l string) {
	s := strings.Split(l[1:], "=")
	if len(s) != 2 {
		cfg.err = true
		cfg.debug("选项配置格式错误, 忽略:", l)
		return
	}
	key := strings.Trim(s[0], " \t")
	val := strings.Trim(s[1], " \t")
	if key == "" || val == "" {
		cfg.err = true
		cfg.debug("选项错误, key和val都不能为空:", l)
		return
	}

	//初始化一个代理配置
	if key == "addr" {
		cfg.options[val] = cfg.defaultOptions()
		cfg.hosts[val] = []*oneHost{}
		cfg.addr = val
		return
	}

	//检查是否已经addr了
	if cfg.addr == "" {
		cfg.err = true
		panic("第一个有效配置必须是'addr'选项: " + l)
		return
	}
	options := cfg.options[cfg.addr]

	//无效配置
	if _, ok := options[key]; !ok {
		cfg.err = true
		cfg.debug("未知选项, 忽略:", l)
		return
	}

	options[key] = val
	//cfg.debug("选项ok:", key, "=", s[1])
}

//解析一行host配置
func (cfg *FileConfig) parseHost(l string) {
	tmp := strings.SplitN(l, "#", 2)[0] //删除右侧注释,如果有
	tmp = strings.Trim(tmp, " \t\r")    //删除首尾空白

	//按空白分割
	arr := regexp.MustCompile(`\s+`).Split(tmp, -1)
	num := len(arr) //分段数量
	if num < 2 {
		cfg.err = true
		cfg.debug("host行格式错误:", l)
		return
	}

	if cfg.addr == "" {
		cfg.err = true
		panic("第一个有效配置必须是 'addr' 选项: " + l)
		return
	}

	for j := 1; j < num; j++ {
		cfg.addHost(arr[0], arr[j]) //添加一个域名
	}
}

//插入一个域名配置
func (cfg *FileConfig) addHost(ip string, domain string) {
	idx := strings.IndexByte(domain, '*')
	if idx == -1 { //没有 '*'
		cfg.hosts[cfg.addr] = append(cfg.hosts[cfg.addr], &oneHost{ip, domain, 0})
		return
	}

	//只有一个 *
	if 1 == strings.Count(domain, "*") {
		if idx == 0 { //开头是 *
			cfg.hosts[cfg.addr] = append(cfg.hosts[cfg.addr], &oneHost{ip, domain[1:], 1})
			return
		} else if idx == len(domain)-1 { //末尾是一个 *
			cfg.hosts[cfg.addr] = append(cfg.hosts[cfg.addr], &oneHost{ip, domain[:len(domain)-2], 2})
			return
		}
	}

	//无法按前缀后缀匹配, 使用正则
	tmp := regexp.QuoteMeta(domain)
	tmp = strings.Replace(tmp, `\*`, `[^.]+`, -1)
	reg := `^` + tmp + `$`
	cfg.hosts[cfg.addr] = append(cfg.hosts[cfg.addr], &oneHost{ip, reg, 3})
}

func (cfg *FileConfig) debug(a ...interface{}) {
	fmt.Println(a...)
}
