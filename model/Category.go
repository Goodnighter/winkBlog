package model

import (
	"gorm.io/gorm"
	"qiublog/utils/errmsg"
)

// Category 分类表
type Category struct {
	ID       uint      `gorm:"primarykey" json:"id"`
	Name     string    `gorm:"type:varchar(255);not null;unique;comment:分类名" json:"name"`
	Mid      *uint     `gorm:"type:int;comment:菜单子项ID" json:"mid"`
	Homeshow bool      `gorm:"type:bool;comment:主页是否显示;default:true;not null;" json:"homeshow,omitempty"`
	Articles []Article `gorm:"foreignkey:Cid"`
}

// Menuchild 菜单子项表
type Menuchild struct {
	ID       uint       `gorm:"primarykey" json:"id"`
	Sort     uint       `gorm:"comment:排序字段" json:"sort"`
	Name     string     `gorm:"comment:菜单名;not null;unique" json:"name"`
	Ename    string     `gorm:"comment:英文名;not null;unique" json:"ename"`
	Logo     string     `gorm:"type:longtext;comment:图标名;" json:"logo"`
	Link     string     `gorm:"comment:路由名;not null;unique" json:"link"`
	ParentId uint       `gorm:"comment:父级id;not null" json:"parent_id"`
	Cids     []Category `gorm:"foreignkey:Mid"`
}

type SetMenuChild struct {
	Type     string `json:"type"`
	ID       uint   `json:"id,omitempty"`
	Sort     uint   `json:"sort,omitempty"`
	Name     string `json:"name,omitempty"`
	Ename    string `json:"ename,omitempty"`
	Logo     string `json:"logo,omitempty"`
	Link     string `json:"link,omitempty"`
	ParentId uint   `json:"parent_id"` //父级id
}

type SetCategory struct {
	Type     string `json:"type"`
	ID       uint   `json:"id"`
	Name     string `json:"name,omitempty"`
	Mid      uint   `json:"mid,omitempty"`
	Homeshow bool   `json:"homeshow,omitempty"`
}
type GetSingleMenuTy struct {
	ID    uint   `json:"id"`    //菜单id
	Name  string `json:"name"`  //菜单名
	Ename string `json:"ename"` //英文名
	Link  string `json:"link"`  //路由名
	Cids  []struct {
		ID       uint   `json:"id"`       //分类id
		Name     string `json:"name"`     //分类name
		Mid      uint   `json:"mid"`      //分类mid
		Homeshow bool   `json:"homeshow"` //分类首页是否显示
	} `gorm:"foreignkey:Mid" json:"cids"`
}

// SetCate  设置菜单子项
func SetCate(dataL []SetCategory) int {
	tx := Db.Begin()
	for _, data := range dataL {
		switch data.Type {
		case "remove":
			err := tx.Delete(&Category{}, data.ID).Error
			if err != nil {
				tx.Rollback()
				return errmsg.ERROR_SET_RE
			}
		case "new":
			menu := Category{}
			err := tx.Create(&menu).Error
			if err != nil {
				tx.Rollback()
				return errmsg.ERROR_SET_NEW
			}
		case "updata":
		default:
			return errmsg.ERROR_SET_TYPE
		}
	}
	tx.Commit()
	return errmsg.SUCCESS
}

// SetMenu 设置菜单子项
func SetMenu(dataL []SetMenuChild) int {
	tx := Db.Begin()
	for _, data := range dataL {
		switch data.Type {
		case "sort":
			err := tx.Model(&Menuchild{}).Where("id = ?", data.ID).Update("sort", data.Sort).Error
			if err != nil {
				tx.Rollback()
				return errmsg.ERROR_SET_SO
			}
		case "remove":
			err := tx.Delete(&Menuchild{}, data.ID).Error
			if err != nil {
				tx.Rollback()
				return errmsg.ERROR_SET_RE
			}
		case "new":
			menu := Menuchild{Sort: data.Sort, Name: data.Name, Ename: data.Ename, Logo: data.Logo, Link: data.Link}
			err := tx.Create(&menu).Error
			if err != nil {
				tx.Rollback()
				return errmsg.ERROR_SET_NEW
			}
		case "updata":
			err := tx.Model(&Menuchild{}).Where("id = ?", data.ID).Updates(Menuchild{Sort: data.Sort, Name: data.Name, Ename: data.Ename, Logo: data.Logo, Link: data.Link}).Error
			if err != nil {
				tx.Rollback()
				return errmsg.ERROR_SET_UP
			}
		default:
			return errmsg.ERROR_SET_TYPE
		}
	}
	tx.Commit()

	return errmsg.SUCCESS
}

// AddMenu 添加菜单子项
func AddMenu(data *Menuchild) int {
	err := Db.Create(&data).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}

// GetMenu 获取菜单子项
func GetMenu() []Menuchild {
	var data []Menuchild
	err := Db.Order("parent_id asc").Order("sort asc").Find(&data).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil
	}
	return data
}

// GetSingleMenu 获取单菜单项
func GetSingleMenu(link string) (int, *GetSingleMenuTy) {
	var data GetSingleMenuTy
	err := Db.Model(Menuchild{}).Where("link=?", link).Find(&data).Error
	if err != nil {
		return errmsg.ERROR, nil
	}
	err = Db.Model(Category{}).Where("mid=?", data.ID).Find(&data.Cids).Error
	if err != nil {
		return errmsg.ERROR, nil
	}
	return errmsg.SUCCESS, &data
}

// AddCategory 添加分类
func AddCategory(data *Category) (int, *uint) {
	if *data.Mid == 0 {
		data.Mid = nil
	}
	err := Db.Create(&data).Error
	if err != nil {
		return errmsg.ERROR, nil
	}
	return errmsg.SUCCESS, &data.ID
}

type GetCategoryTy struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Mid      uint   `json:"mid"`
	Homeshow bool   `json:"homeshow"`
}

// GetCategory 获取分类
func GetCategory(homeshow bool) []GetCategoryTy {
	var data []GetCategoryTy
	err := Db.Model(&Category{}).
		Where(&Category{Homeshow: homeshow}).
		Find(&data).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil
	}
	return data
}

// GetMidCid 根据mid获取所有cid
func GetMidCid(mid int) []uint {
	var data []Category
	var r []uint
	where := map[string]interface{}{}
	if mid == 0 {
		where["homeshow"] = true
	} else if mid > 0 {
		where["mid"] = mid
	} else if mid == -1 {
		where["homeshow"] = false
	} else {
		return nil
	}
	err := Db.Where(where).Find(&data).Error
	if err != nil {
		return nil
	}
	for _, item := range data {
		r = append(r, item.ID)
	}
	return r
}

// ModifyCategorys 批量修改分类
func ModifyCategorys(data *[]Category) int {
	tx := Db.Begin()
	for _, item := range *data {
		err = tx.Table("category").Select("*").Updates(item).Error
		if err != nil {
			tx.Rollback()
			return errmsg.ERROR
		}
	}
	tx.Commit()
	return errmsg.SUCCESS
}
