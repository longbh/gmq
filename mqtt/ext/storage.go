package ext

import(
	"gmq/mqtt/packets"
)

//消息存储
type StorageData interface {
	//存储
	Store(clientIds string,message map[*packets.PublishPacket]bool)
	//读取历史记录
	Select(clientIds string) map[*packets.PublishPacket]bool
}