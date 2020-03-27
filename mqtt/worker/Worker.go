package worker

import (
	"bytes"
	"mqtt/mqtt/packets"
	"mqtt/mqtt/message"
	"mqtt/mqtt/ext"
	"mqtt/mqtt/connect"
	"time"
	"log"
	"mqtt/mqtt/topic"
	"sync/atomic"
)

type Worker struct {
	deltaTime uint32; //心跳检查标记
	clientId  string  //链接标识
 	conn      connect.Connect
	Status    int         //链接状态1=鉴权中 2=链接正常 3=鉴权失败 4=已断开
	timer     *time.Timer
	version	  byte
	versionName string
	idleTime   uint16
	messageIdIncreate  uint16
	messageWindow	*message.Message
	topics	map[string]byte
	idleTimeIncrete uint32	//消息id自增器
}

//初始化
func (worker *Worker) Init(netConnect connect.Connect) {
	worker.conn = netConnect
	worker.Status = 1
	worker.idleTimeIncrete = 0;
	worker.messageWindow = &message.Message{}
	worker.topics = make(map[string]byte)
	defer worker.conn.Close()
	for worker.Status != 4{
		data, err := worker.conn.Read()
		//如果读取数据大小为0或出错则退出
		if err != nil {
			break;
		}
		//去掉两端空白字符==读入读通道
		go worker.Read(data)
	}
}

//读数据
func (worker *Worker) Read(data []byte) {
	control, _ := packets.ReadPacket(bytes.NewReader(data))
	if control == nil{
		worker.Close()
		return;
	}
	log.Print("read data:" + control.String())
	worker.Process(control);
	worker.idleTimeIncrete = 0
}

//写入数据
func (worker *Worker) Write(data []byte) {
	if(worker.conn != nil){
		worker.conn.Write(data)
	}
}

func (worker *Worker) timerLoop() {
	for {
		<-worker.timer.C
		//心跳丢失 检查连接状态，链接丢失则断开连接销毁对象
		atomic.AddUint32(&worker.idleTimeIncrete , 5)
		//检查窗口中是否需要重发或者保存
		data := worker.messageWindow.Resend()
		for k := range data {
			byteBuff := make([] byte, 0,0)
			buffer := bytes.NewBuffer(byteBuff);
			k.Write(buffer)
			worker.Write(buffer.Bytes())
		}
		if worker.idleTimeIncrete > uint32(worker.idleTime + 30){
			worker.Close();
		} else {
			worker.timer.Reset(time.Second * 5)
		}
	}
}

func (worker *Worker) Process(pack packets.ControlPacket)  {
	worker.Status = 2
	ctBack := pack.Process();
	switch pack.FixHeader().MessageType {
	case packets.Connect:
		connectPacket := pack.(*packets.ConnectPacket)
		if(connectPacket.Validate() == packets.Accepted){
			//检查用户密码
			auths := &ext.SdAuths{connectPacket.Username,connectPacket.Password}
			if auths.Login() != packets.Accepted {
				worker.responseData(ctBack)
				worker.Close();
			}
			worker.clientId = connectPacket.GetClientIdentifier()
			worker.version = connectPacket.ProtocolVersion
			worker.versionName = connectPacket.ProtocolName
			worker.idleTime = connectPacket.Keepalive
			NewConnection(connectPacket.GetClientIdentifier(),worker)
			worker.timer = time.NewTimer(time.Second * 5)
			go worker.timerLoop()
			worker.responseData(ctBack)
		} else {
			worker.Close()
		}
		break
	case packets.Disconnect:
		worker.Close();
		break
	case packets.Publish:
		connectPacket := pack.(*packets.PublishPacket)
		topicName := connectPacket.TopicName
		clientIdArr := topic.SearchClientIds(topicName)
		for topic,qos := range clientIdArr{
		// for e := clientIdArr.Front(); e != nil; e = e.Next() {
		//	log.Print("====" + topic) //输出list的值,01234
			connectPacket.Qos = qos
			worker.PushToClients(topic,connectPacket,topicName)
		}
		worker.responseData(ctBack)
		break
	case packets.Puback:
		connectPacket := pack.(*packets.PubackPacket)
		worker.messageWindow.Remove(connectPacket.MessageID)
		break
	case packets.Pubrel:
		worker.responseData(ctBack)
		break
	case packets.Pubrec:
		worker.responseData(ctBack)
		break
	case packets.Pubcomp:
		connectPacket := pack.(*packets.PubcompPacket)
		worker.messageWindow.Remove(connectPacket.MessageID)
		break
	case packets.Subscribe:
		connectPacket := pack.(*packets.SubscribePacket)
		topicName := connectPacket.Topics
		topic.PutTopic(topicName,worker.clientId,connectPacket.Qoss)
		for index,topicItem := range topicName{
			worker.topics[topicItem] = connectPacket.Qoss[index]
		}
		worker.responseData(ctBack)
		break
	case packets.Unsubscribe:
		connectPacket := pack.(*packets.UnsubscribePacket)
		topicName := connectPacket.Topics
		topic.RemoveTopic(topicName,worker.clientId)
		for _,topicItem := range topicName{
			delete(worker.topics,topicItem)
		}
		worker.responseData(ctBack)
		break
	case packets.Pingreq:
		worker.responseData(ctBack)
		break
	}
}

func (worker *Worker) PushToClients(clientIds string,fromPacket *packets.PublishPacket,topic string){
	publishPacket := packets.NewControlPacket(packets.Publish).(*packets.PublishPacket)
	publishPacket.Payload = fromPacket.Payload
	publishPacket.TopicName = topic
	publishPacket.Qos = fromPacket.Qos
	worker.messageWindow.Push(publishPacket)
	byteBuff := make([] byte, 0,0)
	buffer := bytes.NewBuffer(byteBuff);
	publishPacket.Write(buffer)
	go GetConnection(clientIds).Write(buffer.Bytes())
}

func (worker *Worker) responseData(ctBack packets.ControlPacket){
	if(ctBack != nil){
		byteBuff := make([] byte, 0,0)
		buffer := bytes.NewBuffer(byteBuff);
		ctBack.Write(buffer)
		worker.Write(buffer.Bytes())
	}
}

func (worker *Worker) Close() {
	log.Print("connect close:" + worker.clientId)
	worker.Status = 4
	if worker.timer != nil{
		worker.timer.Stop()
	}
	
	//clear all topic suscribe
	for key := range worker.topics{
		topic.RemoveTopicOne(key,worker.clientId)
	}
	//保存contain数据

	worker.conn.Close()
	worker.conn = nil
	CloseConnection(worker.clientId)
}
