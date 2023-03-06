package register

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/utils"

	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
)

var sms SMS

func init() {
	var err error

	cfg := config.Config.Demo
	if cfg.AliSMSVerify.Enable {
		sms, err = NewAliSMS()
		if err != nil {
			panic(err)
		}
	} else if cfg.TencentSMS.Enable {
		sms, err = NewTencentSMS()
		if err != nil {
			panic(err)
		}
	} else {
		sms, err = NewNetEaseSMS()
		if err != nil {
			panic(err)
		}
	}
}

type paramsVerificationCode struct {
	Email          string `json:"email"`
	PhoneNumber    string `json:"phoneNumber"`
	OperationID    string `json:"operationID" binding:"required"`
	UsedFor        int    `json:"usedFor"`
	AreaCode       string `json:"areaCode"`
	InvitationCode string `json:"invitationCode"`
}

// SendVerificationCode 发送验证码
func SendVerificationCode(c *gin.Context) {
	params := paramsVerificationCode{}

	if err := c.BindJSON(&params); err != nil {
		log.NewError("", "BindJSON failed", "err:", err.Error(), "phoneNumber", params.PhoneNumber, "email", params.Email)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": constant.FormattingError, "errMsg": err.Error()})
		return
	}
	operationID := params.OperationID
	if operationID == "" {
		operationID = utils.OperationIDGenerator()
	}
	log.Info(operationID, "SendVerificationCode args: ", "area code: ", params.AreaCode, "Phone Number: ", params.PhoneNumber)
	var account string
	if params.Email != "" {
		account = params.Email
	} else {
		account = params.PhoneNumber
	}
	var accountKey = params.AreaCode + account
	if params.UsedFor == 0 {
		params.UsedFor = constant.VerificationCodeForRegister
	}
	switch params.UsedFor {
	case constant.VerificationCodeForRegister:
		_, err := im_mysql_model.GetRegister(account, params.AreaCode, "")
		if err == nil {
			log.NewError(params.OperationID, "The phone number has been registered", params)
			c.JSON(http.StatusOK, gin.H{"errCode": constant.HasRegistered, "errMsg": "The phone number has been registered"})
			return
		}
		//需要邀请码
		if config.Config.Demo.NeedInvitationCode {
			err = im_mysql_model.CheckInvitationCode(params.InvitationCode)
			if err != nil {
				log.NewError(params.OperationID, "邀请码错误", params)
				c.JSON(http.StatusOK, gin.H{"errCode": constant.InvitationError, "errMsg": "邀请码错误"})
				return
			}
		}
		accountKey = accountKey + "_" + constant.VerificationCodeForRegisterSuffix
		ok, err := db.DB.JudgeAccountEXISTS(accountKey)
		if ok || err != nil {
			log.NewError(params.OperationID, "Repeat send code", params, accountKey)
			c.JSON(http.StatusOK, gin.H{"errCode": constant.RepeatSendCode, "errMsg": "Repeat send code"})
			return
		}

	case constant.VerificationCodeForReset:
		accountKey = accountKey + "_" + constant.VerificationCodeForResetSuffix
		ok, err := db.DB.JudgeAccountEXISTS(accountKey)
		if ok || err != nil {
			log.NewError(params.OperationID, "Repeat send code", params, accountKey)
			c.JSON(http.StatusOK, gin.H{"errCode": constant.RepeatSendCode, "errMsg": "Repeat send code"})
			return
		}
	}
	rand.Seed(time.Now().UnixNano())
	code := 100000 + rand.Intn(900000)
	log.NewInfo(params.OperationID, params.UsedFor, "begin store redis", accountKey, code)
	err := db.DB.SetAccountCode(accountKey, code, config.Config.Demo.CodeTTL)
	if err != nil {
		log.NewError(params.OperationID, "set redis error", accountKey, "err", err.Error())
		c.JSON(http.StatusOK, gin.H{"errCode": constant.SmsSendCodeErr, "errMsg": "Enter the superCode directly in the verification code box, SuperCode can be configured in config.xml"})
		return
	}
	log.NewDebug(params.OperationID, config.Config.Demo)
	if params.Email != "" {
		m := gomail.NewMessage()
		m.SetHeader(`From`, config.Config.Demo.Mail.SenderMail)
		m.SetHeader(`To`, []string{account}...)
		m.SetHeader(`Subject`, config.Config.Demo.Mail.Title)
		m.SetBody(`text/html`, fmt.Sprintf("%d", code))
		if err := gomail.NewDialer(config.Config.Demo.Mail.SmtpAddr, config.Config.Demo.Mail.SmtpPort, config.Config.Demo.Mail.SenderMail, config.Config.Demo.Mail.SenderAuthorizationCode).DialAndSend(m); err != nil {
			log.Error(params.OperationID, "send mail error", account, err.Error())
			c.JSON(http.StatusOK, gin.H{"errCode": constant.MailSendCodeErr, "errMsg": ""})
			return
		}
	} else {
		response, err := sms.SendSms(code, params.AreaCode+params.PhoneNumber)
		if err != nil {
			log.NewError(params.OperationID, "sendSms error", account, "err", err.Error(), response)
			c.JSON(http.StatusOK, gin.H{"errCode": constant.SmsSendCodeErr, "errMsg": "Enter the superCode directly in the verification code box, SuperCode can be configured in config.xml"})
			return
		}
	}
	log.Debug(params.OperationID, "send sms success", code, accountKey)
	data := make(map[string]interface{})
	data["account"] = account
	c.JSON(http.StatusOK, gin.H{"errCode": constant.NoError, "errMsg": "Verification code has been set!", "data": data})
}
