package message

import(
	"mqtt/mqtt/packets"
	"sync/atomic"
	"sync"
)

type Message struct{
	AutoAddInt uint32	//消息id自增器
	MessageWindow  sync.Map
}

func (m *Message) Push(publishPacket *packets.PublishPacket)  {
	atomic.AddUint32(&m.AutoAddInt, 1)
	publishPacket.MessageID = uint16(m.AutoAddInt)
	if publishPacket.Qos != 0{
		m.MessageWindow.Store(publishPacket.MessageID ,publishPacket)
	}
	if publishPacket.MessageID > 44444{
		m.AutoAddInt = 0
	}
}

func (m *Message) Remove(messageId uint16)  {
	m.MessageWindow.Delete(messageId)
}

func (m *Message) Resend() map[*packets.PublishPacket]bool {
	mapData := make(map[*packets.PublishPacket]bool)
	m.MessageWindow.Range(func(k, v interface{}) bool {
		if v.(*packets.PublishPacket).GetOrUpdateExpire(20){
			v.(*packets.PublishPacket).Dup = true
			mapData[v.(*packets.PublishPacket)] = true;
		}
		return true
	})
	return mapData
}