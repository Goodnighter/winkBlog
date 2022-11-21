package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"qiublog/model"
	"qiublog/utils/ask"
	"qiublog/utils/errmsg"
)

func Register(c *gin.Context) {
	var data model.User
	err := c.ShouldBindJSON(&data)
	if err != nil {
		ask.ErrParam(c)
		return
	}
	code = model.Register(&data)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}
func Login(c *gin.Context) {
	var data model.User
	err := c.ShouldBindJSON(&data)
	if err != nil {
		ask.ErrParam(c)
		return
	}
	code, uid, token := model.CheckLogin(&data)
	c.SetCookie("token", token, 259200, "/", c.Request.Referer(), false, false)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
		"uid":     uid,
	})
}
