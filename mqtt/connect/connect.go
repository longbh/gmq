package connect

import (
)

//链接接口
type Connect interface {
	Write(data []byte) (int, error)
	Read() ([]byte, error)
	Close() error
}