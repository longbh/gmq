package packets

import (
	"bytes"
	"fmt"
	"io"
	"mqtt/mqtt/utils"
)

//FixedHeader is a struct to hold the decoded information from
//the fixed header of an MQTT ControlPacket
type FixedHeader struct {
	MessageType     byte
	Dup             bool
	Qos             byte
	Retain          bool
	RemainingLength int
}

func (fh FixedHeader) String() string {
	return fmt.Sprintf("%s: dup: %t qos: %d retain: %t rLength: %d", PacketNames[fh.MessageType], fh.Dup, fh.Qos, fh.Retain, fh.RemainingLength)
}

func (fh *FixedHeader) pack() bytes.Buffer {
	var header bytes.Buffer
	header.WriteByte(fh.MessageType<<4 | utils.BoolToByte(fh.Dup)<<3 | fh.Qos<<1 | utils.BoolToByte(fh.Retain))
	header.Write(utils.EncodeLength(fh.RemainingLength))
	return header
}

func (fh *FixedHeader) unpack(typeAndFlags byte, r io.Reader) error {
	fh.MessageType = typeAndFlags >> 4
	fh.Dup = (typeAndFlags>>3)&0x01 > 0
	fh.Qos = (typeAndFlags >> 1) & 0x03
	fh.Retain = typeAndFlags&0x01 > 0

	var err error
	fh.RemainingLength, err = utils.DecodeLength(r)
	return err
}

//Details struct returned by the Details() function called on
//ControlPackets to present details of the Qos and MessageID
//of the ControlPacket
type Details struct {
	Qos       byte
	MessageID uint16
}