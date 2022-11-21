package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"qiublog/model"
	"qiublog/utils/ask"
	"qiublog/utils/errmsg"
	"qiublog/utils/tool"
	"strconv"
)

var code int

// GetsArticle 获取文章列表
func GetsArticle(c *gin.Context) {
	pageSize, pageNum := tool.PageTool(c)  //分页最大数,分页偏移量
	cid, _ := strconv.Atoi(c.Query("cid")) //分类ID
	mid, _ := strconv.Atoi(c.Query("mid")) //菜单ID
	cids := model.GetMidCid(mid)
	data, total := model.GetsArticle(pageSize, pageNum, cid, cids)
	code = errmsg.SUCCESS
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"total":   total,
		"message": errmsg.GetErrMsg(code),
	})
}

// GetArticle 获取单文章
func GetArticle(c *gin.Context) {
	aid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ask.ErrParam(c)
		return
	}
	code, data := model.GetArticle(aid)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
		"data":    *data,
	})
}

// ReleaseArticle 发布文章
func ReleaseArticle(c *gin.Context) {
	var data model.Article
	err := c.ShouldBindJSON(&data)
	if err != nil {
		ask.ErrParam(c)
		return
	}
	tx := model.Db.Begin()
	code, Aid := model.CreateArticle(tx, &data)
	if code == errmsg.ERROR {
		tx.Rollback()
	} else if code == errmsg.SUCCESS {
		tx.Commit()
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
		"Aid":     Aid,
	})
}

// ModifyArticle 修改文章
func ModifyArticle(c *gin.Context) {

	var data model.Article
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ask.ErrParam(c)
		return
	}
	err = c.ShouldBindJSON(&data)
	if err != nil {
		ask.ErrParam(c)
		return
	}
	tx := model.Db.Begin()
	code = model.ModifyArticle(tx, id, &data)
	if code == errmsg.ERROR {
		tx.Rollback()
	} else if code == errmsg.SUCCESS {
		tx.Commit()
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}

// DeleteArticle 删除文章
func DeleteArticle(c *gin.Context) {

}

// TagGetArticle  根据标签获取所有文章
func TagGetArticle(c *gin.Context) {
	tagId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ask.ErrParam(c)
		return
	}
	data, total := model.TagGetArticle(tagId)
	code = 200
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"total":   total,
		"message": errmsg.GetErrMsg(code),
	})
}
