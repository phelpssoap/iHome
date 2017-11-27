package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"iHome/models"
	"path"
)

type LoginController struct {
	beego.Controller
}

//得到用户名和密码
type LoginData struct {
	Mobile   string `json:"mobile"`
	Password string `json:"password"`
}

type LoginResp struct {
	Errno  string      `json:"errno"`
	Errmsg string      `json:"errmsg"`
	Data   interface{} `json:"data"`
}


func (this *LoginController) RetData(resp interface{}) {
	//给客户端返回json数据
	this.Data["json"] = resp
	//将json写回客户端
	this.ServeJSON()
}

func (this *LoginController) Login() {
	resp := LoginResp{Errno: models.RECODE_OK, Errmsg: models.RecodeText(models.RECODE_OK)}
	defer this.RetData(&resp)

	var data LoginInfo
	json.Unmarshal(this.Ctx.Input.RequestBody, &data)

	beego.Info("Login data:", data)

	//校验信息
	if data.Mobile == "" || data.Password == "" {
		resp.Errno = models.RECODE_REQERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}

	//查询用户信息保存到userinfo
	userinfo := models.User{}
	o := orm.NewOrm()
	qs := o.QueryTable("user")
	err := qs.Filter("mobile", data.Mobile).Filter("password_hash", data.Password).One(&userinfo)

	if err != nil {
		//表示没有任何数据
		resp.Errno = models.RECODE_NODATA
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}

	//对比密码
	if userinfo.Password_hash != data.Password {
		//密码错误
		resp.Errno = models.RECODE_PWDERR
		resp.Errmsg = models.resp(resp.Errno)
		return
	}

	resp.Data = userinfo

	//将用户名存入session中
	this.SetSession("user_id", userinfo.Id)
	this.SetSession("name", userinfo.Name)
	this.SetSession("mobile", userinfo.Mobile)

	return
}
