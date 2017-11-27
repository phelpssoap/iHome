package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"iHome/models"
)

type RegController struct {
	beego.Controller
}

//reg客户端请求的数据
type RegData struct {
	Mobile   string `json:"mobile"`
	Password string `json:"password"`
	Sms_code string `json:"sms_code"`
}

//reg业务回复
type RegResp struct {
	Errno  string `json:"errno"`
	Errmsg string `json:"errmsg"`
}


func (this *RegController) RetData(resp interface{}) {
	//给客户端返回json数据
	this.Data["json"] = resp
	//将json写回客户端
	this.ServeJSON()
}

//  /api/v1.0/users [post]
func (this *RegController) Reg() {
	resp := RegResp{Errno: models.RECODE_OK, Errmsg: models.RecodeText(models.RECODE_OK)}
	defer this.RetData(&resp)

	//得到用户的post请求的数
	//request
	var data RegData
	json.Unmarshal(this.Ctx.Input.RequestBody, &data)
	beego.Info("reg data: ", data)

	//校验信息
	if data.Mobile == "" || data.Password == "" || data.Sms_code == "" {
		resp.Errno = models.RECODE_REQERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}

	//对短信进行校验
	//将用户信息入库

	//var user models.User
	userinfo := models.User{}
	userinfo.Name = data.Mobile
	userinfo.Mobile = data.Mobile
	userinfo.Password_hash = data.Password

	o := orm.NewOrm()
	userinfo.Id, err := o.Insert(&userinfo)
	if err != nil {
		resp.Errno = models.RECODE_DBERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}

	beego.Info("reg insert succ id = ", userinfo.Id)

	//将用户名存入session中
	this.SetSession("name", userinfo.Name)
	this.SetSession("user_id", userinfo.Id)
	this.SetSession("mobile", userinfo.Mobile)

	return
}
