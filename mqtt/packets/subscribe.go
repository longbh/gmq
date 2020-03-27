package packets

import (
	"bytes"
	"fmt"
	"io"
	"mqtt/mqtt/utils"
)

//SubscribePacket is an internal representation of the fields of the
//Subscribe MQTT packet
type SubscribePacket struct {
	FixedHeader
	MessageID uint16
	Topics    []string
	Qoss      []byte
}

func (s *SubscribePacket) String() string {
	str := fmt.Sprintf("%s", s.FixedHeader)
	str += " "
	str += fmt.Sprintf("MessageID: %d topics: %s", s.MessageID, s.Topics)
	return str
}

func (s *SubscribePacket) Write(w io.Writer) error {
	var body bytes.Buffer
	var err error

	body.Write(utils.EncodeUint16(s.MessageID))
	for i, topic := range s.Topics {
		body.Write(utils.EncodeString(topic))
		body.WriteByte(s.Qoss[i])
	}
	s.FixedHeader.RemainingLength = body.Len()
	packet := s.FixedHeader.pack()
	packet.Write(body.Bytes())
	_, err = packet.WriteTo(w)

	return err
}

//Unpack decodes the details of a ControlPacket after the fixed
//header has been read
func (s *SubscribePacket) Unpack(b io.Reader) error {
	var err error
	s.MessageID, err = utils.DecodeUint16(b)
	if err != nil {
		return err
	}
	payloadLength := s.FixedHeader.RemainingLength - 2
	for payloadLength > 0 {
		topic, err := utils.DecodeString(b)
		if err != nil {
			return err
		}
		s.Topics = append(s.Topics, topic)
		qos, err := utils.DecodeByte(b)
		if err != nil {
			return err
		}
		s.Qoss = append(s.Qoss, qos)
		payloadLength -= 2 + len(topic) + 1 //2 bytes of string length, plus string, plus 1 byte for Qos
	}

	return nil
}

//Details returns a Details struct containing the Qos and
//MessageID of this ControlPacket
func (s *SubscribePacket) Details() Details {
	return Details{Qos: 1, MessageID: s.MessageID}
}

//header
func (pr *SubscribePacket) FixHeader() FixedHeader {
	return pr.FixedHeader
}

//response
func (ca *SubscribePacket) Process() ControlPacket {
	subBack := NewControlPacket(Suback)
	connectPacket := subBack.(*SubackPacket)
	connectPacket.MessageID = ca.MessageID
	return connectPacket
}