package routes

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/core/router"
	"mqtt/http/controller"
)

/**
	路由表
 */
func Routes(api *iris.Application) {
	apiParty := api.Party("/api").AllowMethods(iris.MethodOptions)
	//客户端信息
	apiParty.PartyFunc("/clients", func(auths router.Party) {
		auths.Post("/list", controller.ClientsList)
		auths.Post("/topicList", controller.ClientsTopicList)
	})
}
