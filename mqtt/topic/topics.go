package topic

import (
	"strings"
)

//链接存储器
var rootTopic = &Topic{nil,nil,"/",make(map[string]*Topic),make(map[string]byte)}

func PutTopic(topicString []string,clientId string,Qos []byte) {
	for index,topicItem := range topicString{
		kov := strings.Split(topicItem, "/")
		var subTopic = rootTopic
		for i:= 0;i < len(kov); i++{
			itemTopic := kov[i]
			if i == len(kov) - 1{
				subTopic.ClientIds[clientId] = Qos[index]
				break
			} else{
				ssubTopic := subTopic.FindChild(itemTopic)
				if ssubTopic != nil{
					continue
				}else{
					ssubTopic = &Topic{nil,nil,"/",make(map[string]*Topic),make(map[string]byte)}
					subTopic.Children[itemTopic] = ssubTopic
				}
				subTopic = ssubTopic
			}
		}
	}
}

func RemoveTopic(topicString []string,clientId string) {
	for _,topicItem := range topicString {
		kov := strings.Split(topicItem, "/")
		var subTopic = rootTopic
		for i:= 0;i < len(kov); i++{
			itemTopic := kov[i]
			if i == len(kov) - 1{
				delete(subTopic.ClientIds, clientId)
				break
			} else{
				ssubTopic := subTopic.FindChild(itemTopic)
				if ssubTopic == nil{
					break
				}
				subTopic = ssubTopic
			}
		}
	}
}

func RemoveTopicOne(topicString string,clientId string) {
	kov := strings.Split(topicString, "/")
	var subTopic = rootTopic
	for i:= 0;i < len(kov); i++{
		itemTopic := kov[i]
		if i == len(kov) - 1{
			delete(subTopic.ClientIds, clientId)
			break
		} else{
			ssubTopic := subTopic.FindChild(itemTopic)
			if ssubTopic == nil{
				break
			}
			subTopic = ssubTopic
		}
	}
}

func SearchClientIds(topicString string) map[string]byte{
	mapData := make(map[string]byte)
	kov := strings.Split(topicString, "/")
	var subTopic = rootTopic
	for i:= 0;i < len(kov); i++{
		if i == len(kov) - 1{
			for k,v := range subTopic.ClientIds { 
				mapData[k] = v
			}
		} else {
			kovItem := kov[i]
			if kovItem == "*"{

			} else if kovItem == "#"{

			}
			ssubTopic := subTopic.FindChild(kovItem)
			if ssubTopic != nil {
				subTopic = ssubTopic
				continue
			} else {
				break
			}
		}
	}
	return mapData
}