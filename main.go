package main

import (
	"mqtt/mqtt/config"
	"mqtt/master"
)

//启动
func main() {
	config.InitConf();
	//连接管理器进程
	if config.TLSENABLE{
		master.NewTlsMaster();
	}else{
		master.NewMaster();
	}
}
