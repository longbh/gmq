package packets

import (
	"fmt"
	"io"
	"gmq/mqtt/utils"
)

//PubcompPacket is an internal representation of the fields of the
//Pubcomp MQTT packet
type PubcompPacket struct {
	FixedHeader
	MessageID uint16
}

func (pc *PubcompPacket) String() string {
	str := fmt.Sprintf("%s", pc.FixedHeader)
	str += " "
	str += fmt.Sprintf("MessageID: %d", pc.MessageID)
	return str
}

func (pc *PubcompPacket) Write(w io.Writer) error {
	var err error
	pc.FixedHeader.RemainingLength = 2
	packet := pc.FixedHeader.pack()
	packet.Write(utils.EncodeUint16(pc.MessageID))
	_, err = packet.WriteTo(w)

	return err
}

//Unpack decodes the details of a ControlPacket after the fixed
//header has been read
func (pc *PubcompPacket) Unpack(b io.Reader) error {
	var err error
	pc.MessageID, err = utils.DecodeUint16(b)

	return err
}

//Details returns a Details struct containing the Qos and
//MessageID of this ControlPacket
func (pc *PubcompPacket) Details() Details {
	return Details{Qos: pc.Qos, MessageID: pc.MessageID}
}

//header
func (pr *PubcompPacket) FixHeader() FixedHeader {
	return pr.FixedHeader
}

func (ca *PubcompPacket) Process() ControlPacket {
	return nil
}
