package connect

import (
)

//connection interface
type Connect interface {
	Write(data []byte) (int, error)
	Read() ([]byte, error)
	Close() error
}