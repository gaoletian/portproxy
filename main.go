package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"
	"sync"
)

func main() {
	// 通过命令行参数指定端口 --port=1080
	port := flag.Int("p", 1080, "the port to listen on")
	flag.Parse()
	http.HandleFunc("/", handleRequest)

	fmt.Printf("listening on http://127.0.0.1:%d\n", *port)
	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	// 处理 OPTIONS 请求
	// 在处理 OPTIONS 请求时，可以设置 Access-Control-Allow-Headers 头部，
	// 以便客户端在发送带有 Authorization 头部的请求时不会被浏览器拦截。
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// fmt.Printf("%s %s\n", r.Header.Get("Xport"), r.Method)

	var (
		targetPort string
	)

	xport := r.Header.Get("Xport")

	if xport != "" {
		targetPort = xport
	} else {
		// 目标端口  /8080/package.json  ==> 8080
		targetPort = strings.Split(r.URL.Path, "/")[1]
	}

	// 较验端口
	re := regexp.MustCompile(`\d{2,}`)
	if !re.MatchString(targetPort) {
		http.Error(w, "目标端口号到少2位", http.StatusInternalServerError)
		return
	}

	upstream, err := url.Parse(fmt.Sprintf("http://127.0.0.1:%s", targetPort))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 日志输出
	// fmt.Printf("%s [%s] [%s] %s ==> %s\n", r.Header.Get("X-Real-Ip"), r.Header.Get("X-Forwarded-Proto"), r.Method, r.URL, upstream)

	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   upstream.Host,
	})

	// Host重写
	r.Host = upstream.Host
	// path重写
	r.URL.Path = strings.Replace(r.URL.Path, "/"+targetPort, "", 1)

	// 允许跨域
	proxy.ModifyResponse = func(resp *http.Response) error {
		resp.Header.Set("Access-Control-Allow-Origin", "*")
		resp.Header.Set("Access-Control-Allow-Methods", "*")
		resp.Header.Set("Access-Control-Allow-Headers", "*")
		return nil
	}

	// 并发处理
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		proxy.ServeHTTP(w, r)
	}()

	// 等待完成
	wg.Wait()
}
