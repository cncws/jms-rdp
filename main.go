package main

import (
	"flag"
	"html/template"
	"log"
	"os"
	"os/exec"
	"runtime"

	"github.com/cncws/jms-rdp/client"
)

const (
	version      = "0.1.2"
	helpTemplate = `
JumpServer RDP Launcher v{{.}}

Support: MacOS(arm64), Windows(amd64)

Usage:
  jms-rdp -url URL -ak AK -sk SK -account ACCOUNT IP
    必需参数
      -url          服务器地址
      -ak           Access Key，在 Web 页面 API Key 列表获取
      -sk           Secret Key，在 Web 页面 API Key 列表获取
      -account      登陆账号
    可选参数
      -fullscreen   是否全屏

`
)

var (
	serverURL            string
	accessKey, secretKey string
	account              string
	fullscreen           string
	cli                  *client.JmsClient
)

func init() {
	log.SetFlags(0)

	if len(os.Args) == 1 {
		tpl, _ := template.New("help").Parse(helpTemplate)
		tpl.Execute(os.Stdout, version)
		os.Exit(0)
	}

	flag.StringVar(&serverURL, "url", "", "server url")
	flag.StringVar(&accessKey, "ak", "", "access key")
	flag.StringVar(&secretKey, "sk", "", "secret key")
	flag.StringVar(&account, "account", "", "登陆账号")
	flag.StringVar(&fullscreen, "fullscreen", "1", "是否全屏")
	flag.Parse()

	if serverURL == "" {
		log.Fatal("需设置 -url 参数")
	}
	if accessKey == "" {
		log.Fatal("需设置 -ak 参数")
	}
	if secretKey == "" {
		log.Fatal("需设置 -sk 参数")
	}
	if account == "" {
		log.Fatal("需设置 -account 参数")
	}
	cli = client.NewJmsClient(serverURL, accessKey, secretKey, 3)
}

func handleRDP(assetID string) error {
	log.Println("下载连接令牌...")
	token, err := cli.GenRDPToken(assetID, account)
	if err != nil {
		log.Println("❌", err)
		return err
	}
	log.Println("✅", token)
	log.Println("下载 RDP 文件...")
	file, err := cli.DownRDP(token, fullscreen)
	if err != nil {
		log.Println("❌", err)
		return err
	}
	log.Println("✅", file)
	log.Println("打开远程桌面...")
	switch runtime.GOOS {
	case "darwin":
		cmd := exec.Command("open", file)
		return cmd.Start()
	case "windows":
		cmd := exec.Command("mstsc.exe", file)
		return cmd.Start()
	default:
		log.Println("❌尚未支持 ", runtime.GOOS)
	}
	return nil
}

func main() {
	ips := flag.Args()
	if len(ips) == 0 {
		log.Fatal("❌未指定连接的资产 IP")
	}

	log.Println("查询资产 ID...")
	id, protocols := cli.QueryIDByIP(ips[0])
	log.Println("✅", id, "supports", protocols)

	for _, p := range protocols {
		switch p {
		case "rdp":
			if handleRDP(id) != nil {
				os.Exit(1)
			}
		default:
			log.Fatal("❌未支持协议", protocols)
		}
	}
}
