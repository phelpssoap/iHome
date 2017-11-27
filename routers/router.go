package routers

import (
	"github.com/astaxie/beego"
	"iHome/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/api/v1.0/areas", &controllers.AreaController{}, "get:GetAreas")
	beego.Router("/api/v1.0/session", &controllers.SessionController{}, "get:Getses")
	beego.Router("/api/v1.0/session", &controllers.SessionController{}, "DELETE:Delses")
	beego.Router("/api/v1.0/users", &controllers.RegController{}, "post:Reg")
	beego.Router("/api/v1.0/sessions", &controllers.LoginController{}, "post:Login")
	beego.Router("/api/v1.0/user", &controllers.LoginController{}, "get:UserBaseInfo")
	beego.Router("/api/v1.0/user/name", &controllers.LoginController{}, "put:UpdateUsername")
	beego.Router("/api/v1.0/user/avatar", &controllers.LoginController{}, "post:GetAvatar")
	beego.Router("/api/v1.0/user/auth", &controllers.LoginController{}, "get:AuthCheck")
	beego.Router("/api/v1.0/user/auth", &controllers.LoginController{}, "post:UpdateAuthinfo")
	beego.Router("/api/v1.0/user/houses", &controllers.HouseController{}, "get:GetHouseinfo")
	beego.Router("/api/v1.0/houses", &controllers.HouseController{}, "post:AddHouseinfo")
}
