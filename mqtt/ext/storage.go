package ext

//消息存储
type StorageData interface {
	//存储
	Store(clientIds string,message []byte) bool
	//读取历史记录
	Select(clientIds string) []byte
}