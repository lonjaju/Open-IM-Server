package apiShortPhrases

import (
	"net/http"

	api "Open_IM/pkg/base_info"
	//api "Open_IM/pkg/base_info"
	"github.com/gin-gonic/gin"
)

type AddReq struct {
	Secret      string `json:"secret" binding:"required,max=32"`
	Platform    int32  `json:"platform" binding:"required,min=1,max=12"`
	OperationID string `json:"operationID" binding:"required"`
}

type AddResp struct {
	api.CommResp
}

type ModifyReq struct {
	Secret      string `json:"secret" binding:"required,max=32"`
	Platform    int32  `json:"platform" binding:"required,min=1,max=12"`
	OperationID string `json:"operationID" binding:"required"`
}

type ModifyResp struct {
	api.CommResp
}

type DelfyReq struct {
	Secret      string `json:"secret" binding:"required,max=32"`
	Platform    int32  `json:"platform" binding:"required,min=1,max=12"`
	OperationID string `json:"operationID" binding:"required"`
}

type DelResp struct {
	api.CommResp
}

// @Summary 添加快捷用语
// @Description 添加快捷用语
// @Tags 快捷用语
// @ID ShortPhraseAdd
// @Accept json
func ShortPhraseAdd(c *gin.Context) {
	//params := api.UserRegisterReq{}

	resp := &AddResp{}
	c.JSON(http.StatusOK, resp)
}

func ShortPhraseModify(c *gin.Context) {
	//params := api.UserRegisterReq{}

	resp := &ModifyResp{}
	c.JSON(http.StatusOK, resp)
}

func ShortPhraseDel(c *gin.Context) {
	//params := api.UserRegisterReq{}

	resp := &DelResp{}
	c.JSON(http.StatusOK, resp)
}

func ShortPhraseList(c *gin.Context) {
	//params := api.UserRegisterReq{}

	resp := &AddResp{}
	c.JSON(http.StatusOK, resp)
}
