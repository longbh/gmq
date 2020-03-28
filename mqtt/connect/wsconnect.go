package connect

import (
	"github.com/gorilla/websocket"
	"log"
)

//ws connection implements connection
type WsConnection struct{
	Conn      *websocket.Conn
}

func (ws *WsConnection) Read() ([]byte,error) {
	_, data, err := ws.Conn.ReadMessage()
	if err != nil {
		log.Print(err)
		return nil,err
	}
	//log.Print(utils.ToBase64(data))
	return data, nil
}

func (ws *WsConnection) Write(data []byte) (int,error) {
	err := ws.Conn.WriteMessage(websocket.BinaryMessage, data)
	if err != nil {
		return 0,err
	}
	return len(data),err
}

func (ws *WsConnection) Close() error {
	return ws.Conn.Close()
}