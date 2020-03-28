package controller

import (
	"github.com/kataras/iris"
)

//返回客户端列表
func ClientsList(ctx iris.Context)  {
	ctx.JSON(Sucess(OK))
}

//客户端订阅的topic列表
func ClientsTopicList(ctx iris.Context)  {
	ctx.JSON(Sucess(OK))
}