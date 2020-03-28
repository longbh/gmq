package ext

import (
	"gmq/mqtt/packets"
	"encoding/gob"
	"os"
	"log"
)

//读取文件
func ReadFile(filePath string) map[*packets.PublishPacket]bool {
	var M map[*packets.PublishPacket]bool
	File, _ := os.Open(filePath)
	defer File.Close()
	D := gob.NewDecoder(File)
	D.Decode(&M)
	return M
}

//保存文件
func SaveFile(filePath string,info map[*packets.PublishPacket]bool){
	log.Print("dd",filePath,info)
	File, _ := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0777)
    defer File.Close()
	enc := gob.NewEncoder(File)
	if err := enc.Encode(info); err != nil {
        log.Print("save file error,message:",err)
    }
}