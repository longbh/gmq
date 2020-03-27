package ext

type SdStorage struct{
}

func (storage *SdStorage) Store(clientIds string,message []byte) bool  {
	
	return true
}

func (storage *SdStorage) Select(clientIds string) []byte  {
	
	return []byte{}
}