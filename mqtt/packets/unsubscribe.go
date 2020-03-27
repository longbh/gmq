package packets

import (
	"bytes"
	"fmt"
	"io"
	"mqtt/mqtt/utils"
)

//UnsubscribePacket is an internal representation of the fields of the
//Unsubscribe MQTT packet
type UnsubscribePacket struct {
	FixedHeader
	MessageID uint16
	Topics    []string
}

func (u *UnsubscribePacket) String() string {
	str := fmt.Sprintf("%s", u.FixedHeader)
	str += " "
	str += fmt.Sprintf("MessageID: %d", u.MessageID)
	return str
}

func (u *UnsubscribePacket) Write(w io.Writer) error {
	var body bytes.Buffer
	var err error
	body.Write(utils.EncodeUint16(u.MessageID))
	for _, topic := range u.Topics {
		body.Write(utils.EncodeString(topic))
	}
	u.FixedHeader.RemainingLength = body.Len()
	packet := u.FixedHeader.pack()
	packet.Write(body.Bytes())
	_, err = packet.WriteTo(w)

	return err
}

//Unpack decodes the details of a ControlPacket after the fixed
//header has been read
func (u *UnsubscribePacket) Unpack(b io.Reader) error {
	var err error
	u.MessageID, err = utils.DecodeUint16(b)
	if err != nil {
		return err
	}

	for topic, err := utils.DecodeString(b); err == nil && topic != ""; topic, err = utils.DecodeString(b) {
		u.Topics = append(u.Topics, topic)
	}

	return err
}

//Details returns a Details struct containing the Qos and
//MessageID of this ControlPacket
func (u *UnsubscribePacket) Details() Details {
	return Details{Qos: 1, MessageID: u.MessageID}
}

//header
func (pr *UnsubscribePacket) FixHeader() FixedHeader {
	return pr.FixedHeader
}

//response
func (ca *UnsubscribePacket) Process() ControlPacket {
	subBack := NewControlPacket(Unsuback)
	connectPacket := subBack.(*UnsubackPacket)
	connectPacket.MessageID = ca.MessageID
	return connectPacket
}