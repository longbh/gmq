package controller

import (
	"github.com/kataras/iris"
)

//返回节点状态
func NodeInfo(ctx iris.Context)  {
	ctx.JSON(Sucess(OK))
}