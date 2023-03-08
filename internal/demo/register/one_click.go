package register

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var oneClick *NetEaseOneClick

func init() {
	oneClick = NewNetEaseOneClick()
}

type paramsOneClick struct {
	Token       string `json:"token"`
	AccessToken string `json:"accessToken"`
	OperationID string `json:"operationID"`
	UsedFor     int    `json:"usedFor"`
	Platform    int    `json:"platform"`
	OS          int    `json:"os"`
}

// OneClickGetPhone 一键登录获取手机号码
func OneClickGetPhone(c *gin.Context) {
	params := paramsOneClick{}

	if err := c.BindJSON(&params); err != nil {
		log.NewError("", "BindJSON failed", "err:", err.Error(), "Token", params.Token, "email", params.AccessToken)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": constant.FormattingError, "errMsg": err.Error()})
		return
	}
	operationID := params.OperationID
	//if operationID == "" {
	//	operationID = utils.OperationIDGenerator()
	//}
	log.Info(operationID, "OneClickGetPhone args: ", "Token", "AccessToken: ", params.AccessToken)

	response, err := oneClick.getPhone(params.Token, params.AccessToken)
	if err != nil {
		log.NewError(params.OperationID, "getPhone error", params.Token, params.AccessToken, "err", err.Error(), response)
		c.JSON(http.StatusOK, gin.H{"errCode": constant.OneClickGetPhoneError, "errMsg": err.Error()})
		return
	}
	log.Debug(params.OperationID, "get phone success", response.Data.Phone, response.Data.ResultCode)
	data := make(map[string]interface{})
	data["phone"] = response.Data.Phone
	c.JSON(http.StatusOK, gin.H{"errCode": constant.NoError, "errMsg": "Success", "data": data})
}
