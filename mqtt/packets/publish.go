package packets

import (
	"bytes"
	"fmt"
	"io"
	"time"
	"mqtt/mqtt/utils"
)

//PublishPacket is an internal representation of the fields of the
//Publish MQTT packet
type PublishPacket struct {
	FixedHeader
	TopicName string
	MessageID uint16
	Payload   []byte
	Timestamp int64
}

func (p *PublishPacket) String() string {
	str := fmt.Sprintf("%s", p.FixedHeader)
	str += " "
	str += fmt.Sprintf("topicName: %s MessageID: %d", p.TopicName, p.MessageID)
	str += " "
	str += fmt.Sprintf("payload: %s", string(p.Payload))
	return str
}

func (p *PublishPacket) Write(w io.Writer) error {
	var body bytes.Buffer
	var err error

	body.Write(utils.EncodeString(p.TopicName))
	if p.Qos > 0 {
		body.Write(utils.EncodeUint16(p.MessageID))
	}
	p.FixedHeader.RemainingLength = body.Len() + len(p.Payload)
	packet := p.FixedHeader.pack()
	packet.Write(body.Bytes())
	packet.Write(p.Payload)
	_, err = w.Write(packet.Bytes())

	return err
}

//Unpack decodes the details of a ControlPacket after the fixed
//header has been read
func (p *PublishPacket) Unpack(b io.Reader) error {
	var payloadLength = p.FixedHeader.RemainingLength
	var err error
	p.TopicName, err = utils.DecodeString(b)
	if err != nil {
		return err
	}

	if p.Qos > 0 {
		p.MessageID, err = utils.DecodeUint16(b)
		if err != nil {
			return err
		}
		payloadLength -= len(p.TopicName) + 4
	} else {
		payloadLength -= len(p.TopicName) + 2
	}
	if payloadLength < 0 {
		return fmt.Errorf("Error unpacking publish, payload length < 0")
	}
	p.Payload = make([]byte, payloadLength)
	_, err = b.Read(p.Payload)
	//设置消息时间
	cur := time.Now()
	p.Timestamp = cur.UnixNano() 
	return err
}

//Copy creates a new PublishPacket with the same topic and payload
//but an empty fixed header, useful for when you want to deliver
//a message with different properties such as Qos but the same
//content
func (p *PublishPacket) Copy() *PublishPacket {
	newP := NewControlPacket(Publish).(*PublishPacket)
	newP.TopicName = p.TopicName
	newP.Payload = p.Payload

	return newP
}

//Details returns a Details struct containing the Qos and
//MessageID of this ControlPacket
func (p *PublishPacket) Details() Details {
	return Details{Qos: p.Qos, MessageID: p.MessageID}
}

//header
func (pr *PublishPacket) FixHeader() FixedHeader {
	return pr.FixedHeader
}

func (pr *PublishPacket) GetOrUpdateExpire(expireTime int64) bool {
	cur := time.Now()
	timestamp := cur.UnixNano() 
	value := timestamp - pr.Timestamp > expireTime
	pr.Timestamp = timestamp
	return value
}

//response
func (pr *PublishPacket) Process() ControlPacket {
	qos := pr.FixedHeader.Qos
	if(qos == 0){
			
	} else if(qos == 1){
		control := NewControlPacket(Puback)
		connectPacket := control.(*PubackPacket)
		connectPacket.MessageID = pr.MessageID
		return control
	} else if(qos == 2){
		control := NewControlPacket(Pubrec)
		connectPacket := control.(*PubrecPacket)
		connectPacket.MessageID = pr.MessageID
		return control
	}
	return nil
}