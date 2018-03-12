package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/yunnet/gdkxdl/enums"
	"github.com/yunnet/gdkxdl/models"
	"strconv"
	"strings"
	"time"
)

type EquipmentGatewayController struct {
	BaseController
}

func (this *EquipmentGatewayController) Prepare() {
	this.BaseController.Prepare()
	this.checkAuthor("DataGrid", "DataList")
}

func (this *EquipmentGatewayController) Index() {
	this.Data["title"] = "通讯方式"
	this.Data["showMoreQuery"] = true

	this.Data["activeSidebarUrl"] = this.URLFor(this.controllerName + "." + this.actionName)
	this.setTpl()
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["headcssjs"] = "equipmentgateway/index_headcssjs.html"
	this.LayoutSections["footerjs"] = "equipmentgateway/index_footerjs.html"

	//页面里按钮权限控制
	this.Data["canEdit"] = this.checkActionAuthor("EquipmentGatewayController", "Edit")
	this.Data["canDelete"] = this.checkActionAuthor("EquipmentGatewayController", "Delete")
}

func (this *EquipmentGatewayController) DataGrid() {
	var params models.EquipmentGatewayQueryParam
	json.Unmarshal(this.Ctx.Input.RequestBody, &params)
	data, total := models.EquipmentGatewayPageList(&params)

	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data

	this.Data["json"] = result
	this.ServeJSON()
}

func (this *EquipmentGatewayController) DataList() {
	var params = models.EquipmentGatewayQueryParam{}
	data := models.EquipmentGatewayDataList(&params)
	this.jsonResult(enums.JRCodeSucc, "", data)
}

func (this *EquipmentGatewayController) Edit() {
	if this.Ctx.Request.Method == "POST" {
		this.Save()
	}

	Id, _ := this.GetInt(":id", 0)
	m := models.EquipmentGateway{Id: Id}
	if Id > 0 {
		o := orm.NewOrm()
		err := o.Read(&m)
		if err != nil {
			this.pageError("数据无效，请刷新后重试")
		}
	} else {
		m.Used = enums.Enabled
	}

	this.Data["m"] = m
	this.setTpl("equipmentgateway/edit.html", "shared/layout_pullbox.html")
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["footerjs"] = "equipmentgateway/edit_footerjs.html"
}

func (this *EquipmentGatewayController) Save() {
	var err error
	m := models.EquipmentGateway{}

	//获取form里的值
	if err = this.ParseForm(&m); err != nil {
		this.jsonResult(enums.JRCodeFailed, "获取数据失败", m.Id)
	}

	id := this.Input().Get("Id")
	m.Id, _ = strconv.Atoi(id)
	m.GatewayNO = this.GetString("GatewayNO")
	m.GatewayDesc = this.GetString("GatewayDesc")

	o := orm.NewOrm()
	if m.Id == 0 {
		m.CreateUser = this.curUser.RealName
		m.CreateDate = time.Now()

		if _, err = o.Insert(&m); err == nil {
			this.jsonResult(enums.JRCodeSucc, "添加成功", m.Id)
		} else {
			this.jsonResult(enums.JRCodeFailed, "添加失败", m.Id)
		}
	} else {
		m.ChangeUser = this.curUser.RealName
		m.ChangeDate = time.Now()

		if _, err = o.Update(&m, "GatewayNO", "GatewayDesc", "Used", "ChangeUser", "ChangeDate"); err == nil {
			this.jsonResult(enums.JRCodeSucc, "编辑成功", m.Id)
		} else {
			this.jsonResult(enums.JRCodeFailed, "编辑失败", m.Id)
		}
	}
}

func (this *EquipmentGatewayController) Delete() {
	strs := this.GetString("ids")
	ids := make([]int, 0, len(strs))
	for _, str := range strings.Split(strs, ",") {
		if id, err := strconv.Atoi(str); err == nil {
			ids = append(ids, id)
		}
	}

	if num, err := models.EquipmentGatewayBatchDelete(ids); err == nil {
		this.jsonResult(enums.JRCodeSucc, fmt.Sprintf("成功删除 %d 项", num), 0)
	} else {
		this.jsonResult(enums.JRCodeFailed, "删除失败", 0)
	}
}