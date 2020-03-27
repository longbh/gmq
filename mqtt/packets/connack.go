package packets

import (
	"bytes"
	"fmt"
	"io"
	"mqtt/mqtt/utils"
)

//ConnackPacket is an internal representation of the fields of the
//Connack MQTT packet
type ConnackPacket struct {
	FixedHeader
	SessionPresent bool
	ReturnCode     byte
}

func (ca *ConnackPacket) String() string {
	str := fmt.Sprintf("%s", ca.FixedHeader)
	str += " "
	str += fmt.Sprintf("sessionpresent: %t returncode: %d", ca.SessionPresent, ca.ReturnCode)
	return str
}

func (ca *ConnackPacket) Write(w io.Writer) error {
	var body bytes.Buffer
	var err error

	body.WriteByte(utils.BoolToByte(ca.SessionPresent))
	body.WriteByte(ca.ReturnCode)
	ca.FixedHeader.RemainingLength = 2
	packet := ca.FixedHeader.pack()
	packet.Write(body.Bytes())
	_, err = packet.WriteTo(w)

	return err
}

//Unpack decodes the details of a ControlPacket after the fixed
//header has been read
func (ca *ConnackPacket) Unpack(b io.Reader) error {
	flags, err := utils.DecodeByte(b)
	if err != nil {
		return err
	}
	ca.SessionPresent = 1&flags > 0
	ca.ReturnCode, err = utils.DecodeByte(b)

	return err
}

//Details returns a Details struct containing the Qos and
//MessageID of this ControlPacket
func (ca *ConnackPacket) Details() Details {
	return Details{Qos: 0, MessageID: 0}
}

//header
func (ca *ConnackPacket) FixHeader() FixedHeader {
	return ca.FixedHeader
}

func (ca *ConnackPacket) Process() ControlPacket {
	control := NewControlPacket(Connack)
	return control
}
