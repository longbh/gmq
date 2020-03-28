package master

import (
	"fmt"
	"log"
	"crypto/rand"
    "crypto/tls"
	"gmq/mqtt/config"
	"gmq/mqtt/worker"
	"gmq/mqtt/connect"
	"net/http"
	"net"
	"os"
	"time"
	"github.com/gorilla/websocket"
)

//tcp启动
func NewMaster() {
	udpaddr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf("%s:%d", config.HOST, config.TCP_PORT));
	cashError(err)
	//udp 启动
	tcplisten, _ := net.ListenTCP("tcp", udpaddr)
	//死循环的处理客户端请求
	for {
		//等待客户的连接
		//注意这里是无法并发处理多个请求的
		conn, err3 := tcplisten.Accept();
		//如果有错误直接跳过
		if err3 != nil {
			continue;
		}
		go handlerTcp(conn);
	}
}

//tcp tls启动
func NewTlsMaster() {
	crt, err := tls.LoadX509KeyPair(config.TLSCRTPATH, config.TLSKEYPATH)
	if err != nil {
        log.Fatalln(err.Error())
	}
	tlsConfig := &tls.Config{}
	tlsConfig.Certificates = []tls.Certificate{crt}
	tlsConfig.Time = time.Now
	tlsConfig.Rand = rand.Reader
	//tcp 启动
    tcplisten, err := tls.Listen("tcp", fmt.Sprintf("%s:%d", config.HOST, config.TCP_PORT), tlsConfig)
    cashError(err)
	//死循环的处理客户端请求
	for {
		//等待客户的连接
		//注意这里是无法并发处理多个请求的
		conn, err3 := tcplisten.Accept();
		//如果有错误直接跳过
		if err3 != nil {
			continue;
		}
		go handlerTcp(conn);
	}
}

func NewHttpMaster() {
	http.HandleFunc("/mqtt", handlerWs)
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", config.HOST, config.TCP_PORT),nil)
	//错误信息
	cashError(err)
}

//处理链接
func handlerTcp(conn net.Conn) {
	tcpConnect := connect.TcpConnection{conn}
	worker := &worker.Worker{}
	worker.Init(&tcpConnect)
}

func handlerWs(w http.ResponseWriter, r *http.Request) {
	var (
        wbsCon *websocket.Conn
        err error
	)
	Upgrader := websocket.Upgrader {
        // 读取存储空间大小
        ReadBufferSize:1024,
        // 写入存储空间大小
        WriteBufferSize:1024,
        // 允许跨域
        CheckOrigin: func(r *http.Request) bool {
            return true
		},
		Subprotocols:[]string{r.Header.Get("Sec-WebSocket-Protocol")},  
    }
	// 完成http应答，在httpheader中放下如下参数
    if wbsCon, err = Upgrader.Upgrade(w, r, nil);err != nil {
		cashError(err)
        return // 获取连接失败直接返回
    }
	wsConnect := connect.WsConnection{wbsCon}
	worker := &worker.Worker{}
	worker.Init(&wsConnect)
}

func cashError(err error) {
	if (err != nil) {
		log.Print("crash error",err)
		os.Exit(0);
	}
}
