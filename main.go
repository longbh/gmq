package main

import (
	"mqtt/mqtt/config"
	"mqtt/master"
	"mqtt/http"
	"sync"
)

//启动
func main() {
	config.InitConf();
	//http管理工具接口
	var wg sync.WaitGroup
	wg.Add(1)
	go func(){
		http.InitApiServer()
		defer wg.Done()
	}()

	//连接管理器进程
	wg.Add(1)
	go func(){
		if config.TLSENABLE{
		 	master.NewTlsMaster();
		}else{
		 	master.NewMaster();
		}
		defer wg.Done()
	}()

	wg.Wait()
}
