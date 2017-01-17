package product

import (
	"encoding/json"
	"goERP/controllers/base"
	md "goERP/models"
	"strconv"
	"strings"
)

type ProductCategoryController struct {
	base.BaseController
}

func (ctl *ProductCategoryController) Post() {
	action := ctl.Input().Get("action")
	switch action {
	case "validator":
		ctl.Validator()
	case "table": //bootstrap table的post请求
		ctl.PostList()
	case "create":
		ctl.PostCreate()
	default:
		ctl.PostList()
	}
}
func (ctl *ProductCategoryController) Put() {
	id := ctl.Ctx.Input.Param(":id")
	ctl.URL = "/product/category/"
	if idInt64, e := strconv.ParseInt(id, 10, 64); e == nil {
		if category, err := md.GetProductCategoryByID(idInt64); err == nil {
			if err := ctl.ParseForm(&category); err == nil {
				if parentID, err := ctl.GetInt64("parent"); err == nil {
					if parent, err := md.GetProductCategoryByID(parentID); err == nil {
						category.Parent = parent
					}
				}
				if err := md.UpdateProductCategoryByID(category); err == nil {
					ctl.Redirect(ctl.URL+id+"?action=detail", 302)
				}
			}
		}
	}
	ctl.Redirect(ctl.URL+id+"?action=edit", 302)

}
func (ctl *ProductCategoryController) Get() {
	ctl.PageName = "产品类别管理"
	action := ctl.Input().Get("action")
	switch action {
	case "create":
		ctl.Create()
	case "edit":
		ctl.Edit()
	case "detail":
		ctl.Detail()
	default:
		ctl.GetList()
	}
	ctl.Data["PageName"] = ctl.PageName + "\\" + ctl.PageAction
	ctl.URL = "/product/category/"
	ctl.Data["URL"] = ctl.URL
	ctl.Data["MenuProductCategoryActive"] = "active"
}
func (ctl *ProductCategoryController) Edit() {
	id := ctl.Ctx.Input.Param(":id")
	categoryInfo := make(map[string]interface{})
	if id != "" {
		if idInt64, e := strconv.ParseInt(id, 10, 64); e == nil {

			if category, err := md.GetProductCategoryByID(idInt64); err == nil {
				ctl.PageAction = category.Name
				categoryInfo["name"] = category.Name
				parent := make(map[string]interface{})
				if category.Parent != nil {
					parent["id"] = category.Parent.ID
					parent["name"] = category.Parent.Name
				}
				categoryInfo["parent"] = parent
			}
		}
	}
	ctl.Data["Action"] = "edit"
	ctl.Data["RecordID"] = id
	ctl.Data["Category"] = categoryInfo
	ctl.Layout = "base/base.html"

	ctl.TplName = "product/product_category_form.html"
}

func (ctl *ProductCategoryController) Detail() {
	//获取信息一样，直接调用Edit
	ctl.Edit()
	ctl.Data["Readonly"] = true
	ctl.Data["Action"] = "detail"
}

//post请求创建产品分类
func (ctl *ProductCategoryController) PostCreate() {
	category := new(md.ProductCategory)
	if err := ctl.ParseForm(category); err == nil {
		if parentID, err := ctl.GetInt64("parent"); err == nil {
			if parent, err := md.GetProductCategoryByID(parentID); err == nil {
				category.Parent = parent
			}
		}
		if id, err := md.AddProductCategory(category); err == nil {
			ctl.Redirect("/product/category/"+strconv.FormatInt(id, 10)+"?action=detail", 302)
		} else {
			ctl.Get()
		}
	} else {
		ctl.Get()
	}
}
func (ctl *ProductCategoryController) Create() {
	ctl.Data["Action"] = "create"
	ctl.Data["Readonly"] = false
	ctl.PageAction = "创建"
	ctl.Layout = "base/base.html"
	ctl.TplName = "product/product_category_form.html"
}
func (ctl *ProductCategoryController) Validator() {
	name := ctl.GetString("name")
	recordID, _ := ctl.GetInt64("recordID")
	name = strings.TrimSpace(name)
	result := make(map[string]bool)
	obj, err := md.GetProductCategoryByName(name)
	if err != nil {
		result["valid"] = true
	} else {
		if obj.Name == name {
			if recordID == obj.ID {
				result["valid"] = true
			} else {
				result["valid"] = false
			}

		} else {
			result["valid"] = true
		}

	}
	ctl.Data["json"] = result
	ctl.ServeJSON()
}

// 获得符合要求的城市数据
func (ctl *ProductCategoryController) productCategoryList(query map[string]string, fields []string, sortby []string, order []string, offset int64, limit int64) (map[string]interface{}, error) {

	var arrs []md.ProductCategory
	paginator, arrs, err := md.GetAllProductCategory(query, fields, sortby, order, offset, limit)
	result := make(map[string]interface{})
	if err == nil {

		// result["recordsFiltered"] = paginator.TotalCount
		tableLines := make([]interface{}, 0, 4)
		for _, line := range arrs {
			oneLine := make(map[string]interface{})
			oneLine["name"] = line.Name
			if line.Parent != nil {
				oneLine["parent"] = line.Parent.Name
			} else {
				oneLine["parent"] = "-"
			}
			oneLine["path"] = line.ParentFullPath
			oneLine["ID"] = line.ID
			oneLine["id"] = line.ID
			tableLines = append(tableLines, oneLine)
		}
		result["data"] = tableLines
		if jsonResult, er := json.Marshal(&paginator); er == nil {
			result["paginator"] = string(jsonResult)
			result["total"] = paginator.TotalCount
		}
	}
	return result, err
}
func (ctl *ProductCategoryController) PostList() {
	query := make(map[string]string)
	fields := make([]string, 0, 0)
	sortby := make([]string, 0, 0)
	order := make([]string, 0, 0)
	offset, _ := ctl.GetInt64("offset")
	limit, _ := ctl.GetInt64("limit")
	if result, err := ctl.productCategoryList(query, fields, sortby, order, offset, limit); err == nil {
		ctl.Data["json"] = result
	}
	ctl.ServeJSON()

}

func (ctl *ProductCategoryController) GetList() {
	viewType := ctl.Input().Get("view")
	if viewType == "" || viewType == "table" {
		ctl.Data["ViewType"] = "table"
	}
	ctl.PageAction = "列表"
	ctl.Data["tableId"] = "table-product-category"
	ctl.Layout = "base/base_list_view.html"
	ctl.TplName = "product/product_category_list_search.html"
}
