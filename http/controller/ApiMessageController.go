package controller

import (
	"github.com/kataras/iris"
)

//返回客户端列表
func MessageList(ctx iris.Context)  {
	ctx.JSON(Sucess(OK))
}

//客户端订阅的topic列表
func SendMessage(ctx iris.Context)  {
	ctx.JSON(Sucess(OK))
}