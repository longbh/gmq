package connect

import (
	"net"
	"mqtt/mqtt/config"
)

type TcpConnection struct{
	Conn      net.Conn
}

func (tcp *TcpConnection) Read() ([]byte,error) {
	data := make([]byte, config.PACKAGE_LENGTH)
	n,err := tcp.Conn.Read(data)
	if n == 0 || err != nil {
		return nil,err
	}
	// log.Printf("data",data)
	return data[0:n],nil
}

func (tcp *TcpConnection) Write(data []byte) (int,error) {
	return tcp.Conn.Write(data)
}

func (tcp *TcpConnection) Close() error {
	err := tcp.Conn.Close()
	if(err != nil){
		tcp.Conn = nil	
	}
	return err
}