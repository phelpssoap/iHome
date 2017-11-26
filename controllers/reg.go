package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"iHome/models"
)

//reg客户端请求的数据
type RegRequest struct {
	Mobile   string `json:"mobile"`
	Password string `json:"password"`
	Sms_code string `json:"sms_code"`
}

//reg业务回复
type RegResp struct {
	Errno  string `json:"errno"`
	Errmsg string `json:"errmsg"`
}

type RegController struct {
	beego.Controller
}

func (this *RegController) RetData(regresp interface{}) {
	//给客户端返回json数据
	this.Data["json"] = regresp
	//将json写回客户端
	this.ServeJSON()
}

//  /api/v1.0/users [post]
func (this *RegController) Reg() {
	regresp := RegResp{Errno: models.RECODE_OK, Errmsg: models.RecodeText(models.RECODE_OK)}

	defer this.RetData(&regresp)

	//得到用户的post请求的数
	//request
	var reg_data RegRequest
	json.Unmarshal(this.Ctx.Input.RequestBody, &reg_data)

	beego.Info("request data: ", reg_data)

	//校验信息
	if reg_data.Mobile == "" || reg_data.Password == "" || reg_data.Sms_code == "" {
		regresp.Errno = models.RECODE_REQERR
		regresp.Errmsg = models.RecodeText(regresp.Errno)
		return
	}

	//对短信进行校验
	//将用户信息入库

	//var user models.User
	user := models.User{}

	user.Mobile = reg_data.Mobile
	user.Password_hash = reg_data.Password
	user.Name = reg_data.Mobile

	o := orm.NewOrm()
	id, err := o.Insert(&user)
	if err != nil {
		regresp.Errno = models.RECODE_DBERR
		regresp.Errmsg = models.RecodeText(regresp.Errno)
		return
	}

	var userinfo models.User
	err = o.QueryTable(user).Filter("mobile", reg_data.Mobile).One(&reg_data)

	if err != nil {
		regresp.Errno = models.RECODE_NODATA
		regresp.Errmsg = models.RecodeText(regresp.Errno)
		return
	}

	user.Id = userinfo.Id

	beego.Info("reg insert succ id = ", id)

	//将用户名存入session中

	this.SetSession("name", user.Name)
	this.SetSession("user_id", user.Id)
	this.SetSession("mobile", user.Mobile)

	return
}
