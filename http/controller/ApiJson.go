package controller

type ApiJson struct {
	Status  int         //状态码
	Message string      //详细信息
	Data    interface{} //返回数据
}

const OK int = 200   //成功
const FAIL int = 500 //失败

//返回成功
func Sucess(status int) (apijson *ApiJson) {
	apijson = &ApiJson{Status: status, Data: nil, Message: "OK"}
	return
}

func SucessMessage(status int, message string) (apijson *ApiJson) {
	apijson = &ApiJson{Status: status, Data: nil, Message: message}
	return
}

func SucessData(status int, data interface{}) (apijson *ApiJson) {
	apijson = &ApiJson{Status: status, Data: data, Message: "OK"}
	return
}

func SucessMsgData(status int, message string, data interface{}) (apijson *ApiJson) {
	apijson = &ApiJson{Status: status, Data: data, Message: message}
	return
}

//失败返回
func Fail(status int) (apijson *ApiJson) {
	apijson = &ApiJson{Status: status, Data: nil, Message: "FAIL"}
	return
}

func FailMessage(status int, message string) (apijson *ApiJson) {
	apijson = &ApiJson{Status: status, Data: nil, Message: message}
	return
}
