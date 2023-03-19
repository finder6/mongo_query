package routers

import (
	"grxx_query/controllers"
	beego "github.com/beego/beego/v2/server/web"
	"grxx_query/dao"
)

func init() {
	dao.Connect()
    beego.Router("/grxxqy", &controllers.MainController{})
	beego.ErrorController(&controllers.ErrorController{})
}
