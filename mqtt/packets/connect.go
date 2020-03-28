package packets

import (
	"bytes"
	"fmt"
	"io"
	"gmq/mqtt/utils"
)

//ConnectPacket is an internal representation of the fields of the
//Connect MQTT packet
type ConnectPacket struct {
	FixedHeader
	ProtocolName    string
	ProtocolVersion byte
	CleanSession    bool
	WillFlag        bool
	WillQos         byte
	WillRetain      bool
	UsernameFlag    bool
	PasswordFlag    bool
	ReservedBit     byte
	Keepalive       uint16

	ClientIdentifier string
	WillTopic        string
	WillMessage      []byte
	Username         string
	Password         []byte
}

func (c *ConnectPacket) String() string {
	str := fmt.Sprintf("%s", c.FixedHeader)
	str += " "
	str += fmt.Sprintf("protocolversion: %d protocolname: %s cleansession: %t willflag: %t WillQos: %d WillRetain: %t Usernameflag: %t Passwordflag: %t keepalive: %d clientId: %s willtopic: %s willmessage: %s Username: %s Password: %s", c.ProtocolVersion, c.ProtocolName, c.CleanSession, c.WillFlag, c.WillQos, c.WillRetain, c.UsernameFlag, c.PasswordFlag, c.Keepalive, c.ClientIdentifier, c.WillTopic, c.WillMessage, c.Username, c.Password)
	return str
}

func (c *ConnectPacket) Write(w io.Writer) error {
	var body bytes.Buffer
	var err error

	body.Write(utils.EncodeString(c.ProtocolName))
	body.WriteByte(c.ProtocolVersion)
	body.WriteByte(utils.BoolToByte(c.CleanSession)<<1 | utils.BoolToByte(c.WillFlag)<<2 | c.WillQos<<3 | utils.BoolToByte(c.WillRetain)<<5 | utils.BoolToByte(c.PasswordFlag)<<6 | utils.BoolToByte(c.UsernameFlag)<<7)
	body.Write(utils.EncodeUint16(c.Keepalive))
	body.Write(utils.EncodeString(c.ClientIdentifier))
	if c.WillFlag {
		body.Write(utils.EncodeString(c.WillTopic))
		body.Write(utils.EncodeBytes(c.WillMessage))
	}
	if c.UsernameFlag {
		body.Write(utils.EncodeString(c.Username))
	}
	if c.PasswordFlag {
		body.Write(utils.EncodeBytes(c.Password))
	}
	c.FixedHeader.RemainingLength = body.Len()
	packet := c.FixedHeader.pack()
	packet.Write(body.Bytes())
	_, err = packet.WriteTo(w)

	return err
}

//Unpack decodes the details of a ControlPacket after the fixed
//header has been read
func (c *ConnectPacket) Unpack(b io.Reader) error {
	var err error
	c.ProtocolName, err = utils.DecodeString(b)
	if err != nil {
		return err
	}
	c.ProtocolVersion, err = utils.DecodeByte(b)
	if err != nil {
		return err
	}
	options, err := utils.DecodeByte(b)
	if err != nil {
		return err
	}
	c.ReservedBit = 1 & options
	c.CleanSession = 1&(options>>1) > 0
	c.WillFlag = 1&(options>>2) > 0
	c.WillQos = 3 & (options >> 3)
	c.WillRetain = 1&(options>>5) > 0
	c.PasswordFlag = 1&(options>>6) > 0
	c.UsernameFlag = 1&(options>>7) > 0
	c.Keepalive, err = utils.DecodeUint16(b)
	if err != nil {
		return err
	}
	c.ClientIdentifier, err = utils.DecodeString(b)
	if err != nil {
		return err
	}
	if c.WillFlag {
		c.WillTopic, err = utils.DecodeString(b)
		if err != nil {
			return err
		}
		c.WillMessage, err = utils.DecodeBytes(b)
		if err != nil {
			return err
		}
	}
	if c.UsernameFlag {
		c.Username, err = utils.DecodeString(b)
		if err != nil {
			return err
		}
	}
	if c.PasswordFlag {
		c.Password, err = utils.DecodeBytes(b)
		if err != nil {
			return err
		}
	}

	return nil
}

//Validate performs validation of the fields of a Connect packet
func (c *ConnectPacket) Validate() byte {
	if c.PasswordFlag && !c.UsernameFlag {
		return ErrRefusedBadUsernameOrPassword
	}
	if c.ReservedBit != 0 {
		//Bad reserved bit
		return ErrProtocolViolation
	}
	if (c.ProtocolName == "MQIsdp" && c.ProtocolVersion != 3) || (c.ProtocolName == "MQTT" && c.ProtocolVersion != 4) {
		//Mismatched or unsupported protocol version
		return ErrRefusedBadProtocolVersion
	}
	if c.ProtocolName != "MQIsdp" && c.ProtocolName != "MQTT" {
		//Bad protocol name
		return ErrProtocolViolation
	}
	if len(c.ClientIdentifier) > 65535 || len(c.Username) > 65535 || len(c.Password) > 65535 {
		//Bad size field
		return ErrProtocolViolation
	}
	
	if len(c.ClientIdentifier) == 0 && !c.CleanSession {
		//Bad client identifier
		return ErrRefusedIDRejected
	}
	return Accepted
}

//Details returns a Details struct containing the Qos and
//MessageID of this ControlPacket
func (c *ConnectPacket) Details() Details {
	return Details{Qos: 0, MessageID: 0}
}

//header
func (ca *ConnectPacket) FixHeader() FixedHeader {
	return ca.FixedHeader
}

//response
func (ca *ConnectPacket) Process() ControlPacket {
	//判断版本是否支持
	resultCode := ca.Validate()
	if(resultCode != Accepted){
		control := NewControlPacket(Connack)
		connectPacket := control.(*ConnackPacket)
		connectPacket.ReturnCode = resultCode
		return control
	}

	control := NewControlPacket(Connack)
	return control
}

func (c *ConnectPacket) GetClientIdentifier() string {
	return c.ClientIdentifier;
}