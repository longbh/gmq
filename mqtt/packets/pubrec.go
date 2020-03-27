package packets

import (
	"fmt"
	"io"
	"mqtt/mqtt/utils"
)

//PubrecPacket is an internal representation of the fields of the
//Pubrec MQTT packet
type PubrecPacket struct {
	FixedHeader
	MessageID uint16
}

func (pr *PubrecPacket) String() string {
	str := fmt.Sprintf("%s", pr.FixedHeader)
	str += " "
	str += fmt.Sprintf("MessageID: %d", pr.MessageID)
	return str
}

func (pr *PubrecPacket) Write(w io.Writer) error {
	var err error
	pr.FixedHeader.RemainingLength = 2
	packet := pr.FixedHeader.pack()
	packet.Write(utils.EncodeUint16(pr.MessageID))
	_, err = packet.WriteTo(w)

	return err
}

//Unpack decodes the details of a ControlPacket after the fixed
//header has been read
func (pr *PubrecPacket) Unpack(b io.Reader) error {
	var err error
	pr.MessageID, err = utils.DecodeUint16(b)

	return err
}

//Details returns a Details struct containing the Qos and
//MessageID of this ControlPacket
func (pr *PubrecPacket) Details() Details {
	return Details{Qos: pr.Qos, MessageID: pr.MessageID}
}

//header
func (pr *PubrecPacket) FixHeader() FixedHeader {
	return pr.FixedHeader
}

func (ca *PubrecPacket) Process() ControlPacket {
	control := NewControlPacket(Pubrel)
	connectPacket := control.(*PubrelPacket)
	connectPacket.MessageID = ca.MessageID
	return control
}