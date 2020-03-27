package config

import (
	"github.com/Unknwon/goconfig"
	"log"
	"os"
	"strconv"
)

var (
	PACKAGE_LENGTH 	int //系统名称
	HOST           	string
	TCP_PORT       	int //tcp端口
	USERNAME	   	string
	PASSWORD		string
	TLSENABLE		bool
	TLSCRTPATH		string
	TLSKEYPATH		string		
)

//配置以及配置文件初始化目录
func InitConf() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	cfg, error := goconfig.LoadConfigFile(dir + "/conf/broker.ini");
	if error != nil {
		log.Println("read broker.ini failed");
		os.Exit(1)
	}
	packageLength, _ := cfg.GetValue("broker", "PACKAGE_LENGTH")
	HOST, _ = cfg.GetValue("broker", "host")
	tcpPort, _ := cfg.GetValue("broker", "tcp_port")
	PACKAGE_LENGTH, _ = strconv.Atoi(packageLength)
	TCP_PORT, _ = strconv.Atoi(tcpPort)
	//tls
	TLSENABLE,_ = cfg.Bool("broker","tls_enalbe")
	TLSCRTPATH,_ = cfg.GetValue("broker","tls_crt_path")
	TLSKEYPATH,_ = cfg.GetValue("broker","tls_key_path")

	//用户信息
	USERNAME,_ = cfg.GetValue("user","username")
	PASSWORD,_ = cfg.GetValue("user","password")

	
}
