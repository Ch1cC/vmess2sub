package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
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

type user struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var (
	vmessProtocol = "vmess://"
	vmessPath     string
	userPath      string
)

func main() {
	flag.StringVar(&vmessPath, "config", "vmess模板.json", "vmess模板json文件路径")
	flag.StringVar(&userPath, "user", "user模板.json", "user模板json文件路径")
	/*vmessPath := ""
	if len(os.Args) < 2 {
		fmt.Println("请输入文件路径")
		fmt.Scanln(&path)
		vmessPath = path
	} else {
		vmessPath = os.Args[1]
		fmt.Println("读取路径:", vmessPath)
	}*/
	vmesses := formatVmess()
	users := formatUser()
	os.Mkdir("sub", 0644)
	//循环vmess,user对象
	for _, user := range users {
		//字符串拼接
		urlBuilder := strings.Builder{}
		UUID := user.ID
		name := user.Name
		for _, vmess := range vmesses {
			vmess.ID = UUID
			marshal, err := json.Marshal(vmess)
			if err == nil {
				//每个对象都进行base64转换
				vmess := string(marshal)
				base64Url := toVmess(vmess)
				urlBuilder.WriteString(base64Url)
				urlBuilder.WriteString("\n")
			}
		}
		vmessBuilder := urlBuilder.String()
		//最后再base64一次符合小火箭订阅格式
		toString := base64.URLEncoding.EncodeToString([]byte(vmessBuilder))
		os.WriteFile("sub\\"+name, []byte(toString), 0644)
	}
	/*fmt.Println(toString)
	clipboard.WriteAll(toString)
	fmt.Println("已复制到剪切板")*/
	fmt.Printf("文件写入到sub文件夹共%d位用户\n", len(users))
	fmt.Println("回车退出")
	b := make([]byte, 1)
	os.Stdin.Read(b)
}

func toVmess(str string) (base64Url string) {
	return vmessProtocol + base64.StdEncoding.EncodeToString([]byte(str))
}

func formatVmess() []vmess {
	vmessArr := make([]vmess, 0)
	vmess, _ := readJSON(vmessPath)
	//获取json数组
	JSONArr := json.Unmarshal(vmess, &vmessArr)
	if JSONArr != nil {
		fmt.Println("vmess模板.json is error")
		fmt.Println("回车退出")
		b := make([]byte, 1)
		os.Stdin.Read(b)
		panic(JSONArr)
	}
	return vmessArr
}

func formatUser() []user {
	//获取json数组
	userArr := make([]user, 0)
	user, _ := readJSON(userPath)
	JSONArr := json.Unmarshal(user, &userArr)
	if JSONArr != nil {
		fmt.Println("user模板.json is error")
		fmt.Println("回车退出")
		b := make([]byte, 1)
		os.Stdin.Read(b)
		panic(JSONArr)
	}
	return userArr
}
func readJSON(path string) ([]byte, error) {
	open, _ := os.Open(path)
	defer open.Close()
	return io.ReadAll(open)
}
