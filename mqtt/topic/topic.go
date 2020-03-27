package topic

import (
	
)

type Topic struct{
	ppNode 		*Topic
	pNode  		*Topic
	name		string
	Children 	map[string]*Topic
	ClientIds	map[string] byte
}

func (topic *Topic) AddChild(name string)  {
	topic.name = name
}

func (topic *Topic) FindChild(name string) *Topic {
	return topic.Children[name]
}
