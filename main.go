package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"iHome/models"
	_ "iHome/routers"
	"net/http"
	"strings"
)

func main() {
	ignoreStaticPath()
	beego.Run()
}

func init() {
	orm.RegisterDataBase("default", "mysql", "root:mysql@tcp(127.0.0.1:3306)/iHome?charset=utf8", 30)
	orm.RegisterModel(new(models.User), new(models.Area), new(models.Facility), new(models.House), new(models.HouseImage), new(models.OrderHouse))

	orm.RunSyncdb("default", false, true)

}

//重定向静态路径
func ignoreStaticPath() {
	//透明static
	beego.InsertFilter("/", beego.BeforeRouter, TransparentStatic)
	beego.InsertFilter("/*", beego.BeforeRouter, TransparentStatic)
	//设置一个fastdfs 请求的静态路径
	//http://47.95.219.27:8080/group1/M00/00/00/Zciqq1oaGW-ABnxDAAAHFIcthTk%207176.go
	beego.SetStaticPath("/group1/M00", "fastdfs/storage_data/data")
	beego.SetStaticPath("/down1", "download1")
}

func TransparentStatic(ctx *context.Context) {
	orpath := ctx.Request.URL.Path
	beego.Debug("request url:", orpath)

	//如果请求uri还有API字段，说明是指令应该取消静态资源路径重定向
	if strings.Index(orpath, "api") >= 0 {
		return
	}

	http.ServeFile(ctx.ResponseWriter, ctx.Request, "static/html/"+ctx.Request.URL.Path)
	//将全部的静态资源重定向 加上/static/html路径
	//http://ip:port:8080/index.html----> http://ip:port:8080/static/html/index.html
	//如果restFUL api  那么就取消冲定向
	//http://ip:port:8080/api/v1.0/areas ---> http://ip:port:8080/static/html/api/v1.0/areas

}
