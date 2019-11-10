package xwproxy

import (
	"io"
	"log"
	"net"
	"net/http"
	"strings"
)

type Pxy struct {
	fCfg *FileConfig
	addr string
}

//启动代理服务,可能是多个
func Run(cfg *FileConfig) {
	for addr, _ := range cfg.options{
		go func() {
			pxy := &Pxy{cfg, addr}
			log.Printf("HttpPxoy to runing on %s \n", addr)
			log.Fatalln(http.ListenAndServe(addr, pxy))
		}()
	}

}

// 运行代理服务
func (p *Pxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// http && https
	if req.Method != "CONNECT" {
		// 处理http
		p.HTTP(rw, req)
	} else {
		// 处理https
		// 直通模式不做任何中间处理
		p.HTTPS(rw, req)
	}
}

// http
func (p *Pxy) HTTP(rw http.ResponseWriter, req *http.Request) {

	transport := http.DefaultTransport

	// 新建一个请求outReq
	outReq := new(http.Request)

	// 复制客户端请求到outReq上
	*outReq = *req // 复制请求

	//修改ip
	ip := p.fCfg.GetIP(p.addr, outReq.URL.Hostname())
	if ip != "" {
		port := outReq.URL.Port()
		if port != "" {
			outReq.URL.Host = ip + ":" + outReq.URL.Port()
		} else {
			outReq.URL.Host = ip
		}
		log.Println("http :", outReq.Host, "->", ip)
	}

	//  处理匿名代理
	forwardedIP := p.fCfg.GetOption(p.addr,"forwardedIP")
	if forwardedIP != "0" && forwardedIP != "" { //设置 'X-Forwarded-For'
		ip := ""
		if forwardedIP == "1" { //按真实情况设置
			ip, _, _ = net.SplitHostPort(req.RemoteAddr)
		} else { //设定为指定ip
			ip = forwardedIP
		}
		//执行
		if ip != "" {
			if prior, ok := outReq.Header["X-Forwarded-For"]; ok {
				ip = strings.Join(prior, ", ") + ", " + ip
			}
			outReq.Header.Set("X-Forwarded-For", ip)
		}
	}

	// outReq请求放到传送上
	res, err := transport.RoundTrip(outReq)
	if err != nil {
		rw.WriteHeader(http.StatusBadGateway)
		rw.Write([]byte(err.Error()))
		return
	}

	// 回写http头
	for key, value := range res.Header {
		for _, v := range value {
			rw.Header().Add(key, v)
		}
	}
	// 回写状态码
	rw.WriteHeader(res.StatusCode)
	// 回写body
	io.Copy(rw, res.Body)
	res.Body.Close()
}

// https
func (p *Pxy) HTTPS(rw http.ResponseWriter, req *http.Request) {

	// 拿出host
	host := req.URL.Host
	hij, ok := rw.(http.Hijacker)
	if !ok {
		log.Printf("HTTP Server does not support hijacking")
	}

	client, _, err := hij.Hijack()
	if err != nil {
		return
	}

	//log.Println("https", req.Host)

	//更换ip
	ip := p.fCfg.GetIP(p.addr, req.URL.Hostname())
	if ip != "" {
		port := req.URL.Port()
		if port != "" {
			host = ip + ":" + req.URL.Port()
		} else {
			host = ip
		}
		log.Println("https :", req.Host, "->", ip)
	}

	// 连接远程
	server, err := net.Dial("tcp", host)
	if err != nil {
		return
	}
	client.Write([]byte("HTTP/1.0 200 Connection Established\r\n\r\n"))

	// 直通双向复制
	go io.Copy(server, client)
	go io.Copy(client, server)
}
