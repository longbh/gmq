package ext

import (
	"mqtt/mqtt/config"
	"mqtt/mqtt/packets"
)

type SdAuths struct{
	UserName string
	Password []byte
}

func (auths *SdAuths) Login() byte  {
	if(config.USERNAME == ""){
		return packets.Accepted
	}

	if auths.UserName != config.USERNAME{
		return packets.ErrProtocolViolation
	}
	return packets.Accepted
}

func (auths *SdAuths) SslCheck(path string) bool  {
	
	return true
}