package worker

import (
	"sync"
)

//链接存储器
var connections sync.Map

//push new connection
func NewConnection(clientId string, worker *Worker) {
	//close old connection
	oldWorlder,_ := connections.Load(clientId)
	if oldWorlder != nil {
		oldWorlder.(*Worker).Close();
		connections.Delete(clientId)
	}
	connections.Store(clientId , worker);
}

func CloseConnection(clientId string) {
	connections.Delete(clientId)
}

func GetConnection(clientId string) *Worker {
	oldWorlder,_ := connections.Load(clientId)
	return oldWorlder.(*Worker);
}