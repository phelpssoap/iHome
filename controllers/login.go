package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"iHome/models"
	"path"
)

type LoginController struct {
	beego.Controller
}

//得到用户名和密码
type LoginInfo struct {
	Mobile   string `json:"mobile"`
	Password string `json:"password"`
}

type LoginResp struct {
	Errno  string      `json:"errno"`
	Errmsg string      `json:"errmsg"`
	Data   interface{} `json:"data"`
}

type AvatarUrl struct {
	Url string `json:"avatar_url"`
}

type Baseinfo struct {
	User_id   int    `json:"user_id"`
	Name      string `json:"name"`
	Password  string `json:"password"`
	Mobile    string `json:"mobile"`
	Real_name string `json:"real_name"`
	ID_card   string `json:"id_card"`
	Url       string `json:"avatar_url"`
}

//客户端点击修改用户信息时返回给客户端的用户名和密码
type UserinfoResp struct {
	Errno  string   `json:"errno"`
	Errmsg string   `json:"errmsg"`
	Data   Baseinfo `json:"data"`
}

// 上传头像的返回结构
type AvatarResp struct {
	Errno  string    `json:"errno"`
	Errmsg string    `json:"errmsg"`
	Data   AvatarUrl `json:"data"`
}

type AuthcheckResp struct {
	Errno  string   `json:"errno"`
	Errmsg string   `json:"errmsg"`
	Data   Baseinfo `json:"data"`
}

type AuthInfo struct {
	Real_name string `json:"real_name"`
	ID_card   string `json:"id_card"`
}

func (this *LoginController) RetData(resp interface{}) {
	//给客户端返回json数据
	this.Data["json"] = resp
	//将json写回客户端
	this.ServeJSON()
}

func (this *LoginController) Login() {
	//返回客户端
	loginresp := LoginResp{Errno: models.RECODE_OK, Errmsg: models.RecodeText(models.RECODE_OK)}

	defer this.RetData(&loginresp)

	//得到用户信息
	var login_data LoginInfo
	json.Unmarshal(this.Ctx.Input.RequestBody, &login_data)

	fmt.Printf("request data: %+v\n", login_data)

	//校验信息
	if login_data.Mobile == "" || login_data.Password == "" {
		loginresp.Errno = models.RECODE_REQERR
		loginresp.Errmsg = models.RecodeText(loginresp.Errno)
		return
	}

	//查询用户信息保存到userinfo
	var userinfo models.User
	o := orm.NewOrm()
	qs := o.QueryTable("user")
	err := qs.Filter("mobile", login_data.Mobile).Filter("password_hash", login_data.Password).One(&userinfo)

	if err != nil {
		//表示没有任何数据
		loginresp.Errno = models.RECODE_NODATA
		loginresp.Errmsg = models.RecodeText(loginresp.Errno)
		return
	}

	//对比密码
	if userinfo.Password_hash != login_data.Password {
		//密码错误
		loginresp.Errno = models.RECODE_PWDERR
		loginresp.Errmsg = models.RecodeText(loginresp.Errno)
		return
	}

	beego.Info("Login(): ", userinfo)
	loginresp.Data = userinfo

	//将用户名存入session中
	this.SetSession("user_id", userinfo.Id)
	this.SetSession("name", userinfo.Name)
	this.SetSession("mobile", userinfo.Mobile)

	return
}

func (this *LoginController) UserBaseInfo() {

	userinforesp := UserinfoResp{Errno: models.RECODE_OK, Errmsg: models.RecodeText(models.RECODE_OK)}
	defer this.RetData(&userinforesp)

	user_id := this.GetSession("user_id")

	var userinfo models.User
	o := orm.NewOrm()
	qs := o.QueryTable("user")
	err := qs.Filter("id", user_id).One(&userinfo)

	if err != nil {
		//表示没有任何数据
		beego.Debug("UserBaseInfo error: ", err)
		userinforesp.Errno = models.RECODE_NODATA
		userinforesp.Errmsg = models.RecodeText(userinforesp.Errno)
		return
	}

	var baseinfo Baseinfo
	baseinfo.User_id = userinfo.Id
	baseinfo.Name = userinfo.Name
	baseinfo.Password = userinfo.Password_hash
	baseinfo.Mobile = userinfo.Mobile
	baseinfo.Real_name = userinfo.Real_name
	baseinfo.ID_card = userinfo.Id_card
	baseinfo.Url = userinfo.Avatar_url

	userinforesp.Data = baseinfo
	return
}

func (this *LoginController) AuthCheck() {
	authcheckresp := AuthcheckResp{Errno: models.RECODE_OK, Errmsg: models.RecodeText(models.RECODE_OK)}
	defer this.RetData(&authcheckresp)

	user_id := this.GetSession("user_id")

	var userinfo models.User
	o := orm.NewOrm()
	qs := o.QueryTable("user")
	err := qs.Filter("id", user_id).One(&userinfo)

	if err != nil {
		//查询出错
		beego.Debug("AuthCheck error: ", err)
		authcheckresp.Errno = models.RECODE_NODATA
		authcheckresp.Errmsg = models.RecodeText(authcheckresp.Errno)
		return
	}

	if userinfo.Real_name == "" || userinfo.Id_card == "" {
		//表示未认证
		beego.Debug("用户未认证")
		authcheckresp.Errno = models.RECODE_ROLEERR
		authcheckresp.Errmsg = models.RecodeText(authcheckresp.Errno)
		return
	}

	var baseinfo Baseinfo
	baseinfo.User_id = userinfo.Id
	baseinfo.Name = userinfo.Name
	baseinfo.Password = userinfo.Password_hash
	baseinfo.Mobile = userinfo.Mobile
	baseinfo.Real_name = userinfo.Real_name
	baseinfo.ID_card = userinfo.Id_card
	baseinfo.Url = userinfo.Avatar_url

	authcheckresp.Data = baseinfo
	return
}

func (this *LoginController) UpdateAuthinfo() {
	var auth_info AuthInfo
	json.Unmarshal(this.Ctx.Input.RequestBody, &auth_info)

	authinforesp := AuthcheckResp{Errno: models.RECODE_OK, Errmsg: models.RecodeText(models.RECODE_OK)}
	defer this.RetData(&authinforesp)

	user_id := this.GetSession("user_id")
	o := orm.NewOrm()
	_, err := o.QueryTable("user").Filter("id", user_id).Update(orm.Params{"real_name": auth_info.Real_name, "id_card": auth_info.ID_card})

	if err != nil {
		authinforesp.Errno = models.RECODE_DATAERR
		authinforesp.Errmsg = models.RecodeText(authinforesp.Errno)
		return
	}

	return
}

func (this *LoginController) UpdateUsername() {
	var user_name Name
	json.Unmarshal(this.Ctx.Input.RequestBody, &user_name)
	beego.Info("beego info user_name: ", user_name.Name)

	updateresp := LoginResp{Errno: models.RECODE_OK, Errmsg: models.RecodeText(models.RECODE_OK)}
	defer this.RetData(&updateresp)

	oldname := this.GetSession("name")
	beego.Info("beego info oldname: ", oldname.(string))
	o := orm.NewOrm()
	_, err := o.QueryTable("user").Filter("name", oldname.(string)).Update(orm.Params{"name": user_name.Name})

	if err != nil {
		updateresp.Errno = models.RECODE_DBERR
		updateresp.Errmsg = models.RecodeText(updateresp.Errno)
	}

	updateresp.Data = user_name

	this.SetSession("name", user_name.Name)

	return
}

func (this *LoginController) GetAvatar() {
	PicResp := AvatarResp{Errno: models.RECODE_OK, Errmsg: models.RecodeText(models.RECODE_OK)}

	defer this.RetData(&PicResp)

	//获取文件数据
	file, header, err := this.GetFile("avatar")

	if err != nil {
		PicResp.Errno = models.RECODE_SERVERERR
		PicResp.Errmsg = models.RecodeText(PicResp.Errno)
		beego.Info("get file error")
		return
	}
	defer file.Close()

	//创建一个文件的缓冲
	fileBuffer := make([]byte, header.Size)

	if _, err := file.Read(fileBuffer); err != nil {
		PicResp.Errno = models.RECODE_IOERR
		PicResp.Errmsg = models.RecodeText(PicResp.Errno)
		beego.Info("read file error")
		return
	}

	suffix := path.Ext(header.Filename) // suffix = ".jpg"
	groupName, fileId, err1 := models.FDFSUploadByBuffer(fileBuffer, suffix[1:])
	if err1 != nil {
		PicResp.Errno = models.RECODE_IOERR
		PicResp.Errmsg = models.RecodeText(PicResp.Errno)
		beego.Info("fdfs upload file error")
		return
	}

	beego.Info("groupname, ", groupName, " file id: ", fileId)

	//通过session获取当前用户
	user_id := this.GetSession("user_id")

	//添加Avatar_url字段到数据库中
	o := orm.NewOrm()
	user := models.User{Id: user_id.(int), Avatar_url: fileId}

	if _, err := o.Update(&user, "avatar_url"); err != nil {
		PicResp.Errno = models.RECODE_DBERR
		PicResp.Errmsg = models.RecodeText(PicResp.Errno)
		return
	}

	//拼接一个完成的路径
	avatar_url := "http://47.95.219.27:8080/" + fileId
	PicResp.Data.Url = avatar_url
	return
}
