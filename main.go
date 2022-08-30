package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/atotto/clipboard"
	"io"
	"os"
	"strings"
)

type vmess struct {
	V    string `json:"v"`
	Ps   string `json:"ps"`
	Add  string `json:"add"`
	Port string `json:"port"`
	ID   string `json:"id"`
	Aid  string `json:"aid"`
	Scy  string `json:"scy"`
	Net  string `json:"net"`
	Type string `json:"type"`
	Host string `json:"host"`
	Path string `json:"path"`
	TLS  string `json:"tls"`
	Sni  string `json:"sni"`
	Alpn string `json:"alpn"`
}

var (
	vmessProtocol = "vmess://"
	path          string
)

func main() {
	filePath := ""
	if len(os.Args) < 2 {
		fmt.Println("请输入文件路径")
		fmt.Scanln(&path)
		filePath = path
	} else {
		filePath = os.Args[1]
		fmt.Println("读取路径:", filePath)
	}
	//接受json数组
	vmessArr := make([]vmess, 0)
	//字符串拼接
	urlBuilder := strings.Builder{}
	open, _ := os.Open(filePath)
	defer open.Close()
	value, _ := io.ReadAll(open)
	//获取json数组
	jsonErr := json.Unmarshal(value, &vmessArr)
	if jsonErr != nil {
		fmt.Println("jsonErr is error")
		fmt.Println("回车退出")
		b := make([]byte, 1)
		os.Stdin.Read(b)
		panic(jsonErr)
	}
	//循环vmess对象
	for _, v := range vmessArr {
		marshal, err := json.Marshal(v)
		if err == nil {
			//每个对象都进行base64转换
			base64Url := toVmess(string(marshal))
			urlBuilder.WriteString(base64Url)
			urlBuilder.WriteString("\n")
		}
	}
	vmessBuilder := urlBuilder.String()
	//最后再base64一次符合小火箭订阅格式
	toString := base64.URLEncoding.EncodeToString([]byte(vmessBuilder))
	fmt.Println(toString)
	clipboard.WriteAll(toString)
	fmt.Println("已复制到剪切板")
	fmt.Println("回车退出")
	b := make([]byte, 1)
	os.Stdin.Read(b)
}

func toVmess(str string) (base64Url string) {
	return vmessProtocol + base64.StdEncoding.EncodeToString([]byte(str))
}
