package http

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"mqtt/http/route"
	"mqtt/http/controller"
	"mqtt/mqtt/config"
	"fmt"
)

/**
 * 初始化 iris app
 * @method NewApp
 * @return  {[type]}  api      *iris.Application  [iris app]
 */
func newApp() (api *iris.Application) {
	api = iris.New()
	api.Use(logger.New())
	api.OnErrorCode(iris.StatusNotFound, func(ctx iris.Context) {
		ctx.JSON(controller.SucessMessage(controller.FAIL, "404 Not Found"))
	})
	api.OnErrorCode(iris.StatusInternalServerError, func(ctx iris.Context) {
		ctx.JSON(controller.SucessMessage(controller.FAIL, "error code"))
	})
	return api
}

func InitApiServer() {
	api := newApp();                                                                  //iris初始化
	routes.Routes(api)                                                                //路由初始化
	api.Run(iris.Addr(fmt.Sprintf("0.0.0.0:%d" , config.HTTP_PORT)), iris.WithoutInterruptHandler)
}
