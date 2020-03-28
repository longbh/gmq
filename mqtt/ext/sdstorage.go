package ext

import(
	"gmq/mqtt/packets"
)

type SdStorage struct{
}

func (storage *SdStorage) Store(clientIds string,message map[*packets.PublishPacket]bool)  {
	SaveFile("/Users/longbh/Desktop/Finder/project/go/src/gmq/" + clientIds,message)
}

func (storage *SdStorage) Select(clientIds string) map[*packets.PublishPacket]bool  {
	return ReadFile("/Users/longbh/Desktop/Finder/project/go/src/gmq/" + clientIds)
}