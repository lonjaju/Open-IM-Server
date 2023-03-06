package register

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/utils"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"sort"
	"strconv"
	"time"
)

var (
	Code2Msg = map[int]string{
		400: "bad request: 请求缺少 secretId 或 businessId",
		401: "forbidden: secretId 或 businessId 错误",
		405: "param error:  请求参数异常",
		410: "signature failure: 签名验证失败，请重新参考demo签名代码",
		420: "request expired: 请求过期",
		421: "contentTypeError:	请求参数错误",
		429: "too many requests: 次数超限",
		430: "replay attack: 重放攻击",
		440: "decode error:	解密错误",
		450: "wrong token:	token错误",
		503: "service unavailable: 接口异常",
		506: "exceed phone send limit: 单手机号发送频率限制",
		507: "balance not enough: 套餐余量不足",
		508: "money is not enough: 试用条数不足",
	}
)

const (
	// 产品密钥ID，产品标识
	secretID = "your_secret_id"
	// 产品私有密钥，服务端生成签名信息使用，请严格保管，避免泄露
	secretKey = "your_secret_key"
	// 业务ID，易盾根据产品业务特点分配
	businessID = "your_business_id"
	// 本机认证服务身份证实人认证在线检测接口地址
	apiURL = "https://sms.dun.163.com/v2/sendsms"
)

type NetEaseSMS struct {
	secretID   string
	secretKey  string
	businessID string
	apiURL     string
}

func (n *NetEaseSMS) SendSms(code int, phoneNumber string) (resp interface{}, err error) {
	// 1.设置公共参数
	params := make(map[string]interface{})
	// 产品秘钥 id ，在网易易盾创建产品时统一分配，产品标识
	params["secretId"] = n.secretID
	// 业务id ，在网易易盾创建产品时统一分配，业务标识
	params["businessId"] = n.businessID
	// 接口版本号，可选值 v2
	params["version"] = "v2"
	// 请求当前 UNIX 时间戳（毫秒单位），请注意服务器时间是否同步
	params["timestamp"] = strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	// 随机字符串，与 timestamp 联合起来，用于防止重放攻击
	params["nonce"] = utils.GetRandomStr(32)

	// 2.设置私有参数
	smsParams := make(map[string]interface{})
	smsParams["code"] = utils.IntToString(code)
	smsParams["time"] = "20180816"

	// params: 该模板ID模板变量中要替换的内容；{"code":"123","time":"20180816"}
	// mobile:
	// 接收短信的手机号
	// 单次调用仅支持一个手机号
	// 发送国际短信，请去掉手机号前的0
	params["mobile"] = phoneNumber
	// 该字段必填：“json”
	params["paramType"] = "json"
	params["params"] = smsParams
	// 模板ID
	params["templateId"] = config.Config.Demo.NetEaseSMS.VerificationCodeTemplateCode
	// 参数templateId传过来的模板ID中的验证码变量名，例：${code}。务必填写正确，否则验证码无法完成替换
	params["codeName"] = "code"
	// 验证码数字个数，支持范围4-10
	params["codeLen"] = 6
	// 验证码有效期，支持范围300-1800秒，单位-秒
	params["codeValidSec"] = 300

	// 3.生成签名信息
	signature := n.genSignature(params)
	params["signature"] = signature

	// 4.发送HTTP请求
	reqBody, err := json.Marshal(params)
	if err != nil {
		return nil, errors.WithMessage(err, "netease.sms.reqBody")
	}
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, errors.WithMessage(err, "netease.sms.req")
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	rawResp, err := client.Do(req)
	if err != nil {
		return nil, errors.WithMessage(err, "netease.sms.do")
	}
	defer rawResp.Body.Close()

	// 5.解析报文返回
	body, err := io.ReadAll(rawResp.Body)
	if err != nil {
		return nil, errors.WithMessage(err, "netease.smd.body")
	}

	var smsResp NetEaseSMSResponse
	err = json.Unmarshal(body, &smsResp)
	if err != nil {
		return nil, errors.WithMessage(err, "netease.sms.resBody")
	}
	b, _ := json.Marshal(smsResp)

	//6.返回结果
	log.Info("netease request id is ", smsResp.Data.RequestID)

	log.Debug("netease send message is ", code, phoneNumber, string(b))

	if smsResp.Code != 200 {
		str := Code2Msg[smsResp.Code]

		return nil, fmt.Errorf("异常返回码: %d %s", smsResp.Code, str)
	}

	return smsResp, nil
}

func NewNetEaseSMS() (*NetEaseSMS, error) {
	cfg := config.Config.Demo.NetEaseSMS

	return &NetEaseSMS{
		apiURL:     "https://sms.dun.163.com/v2/sendsms",
		secretID:   cfg.SecretID,
		secretKey:  cfg.SecretKey,
		businessID: cfg.BusinessID,
	}, nil
}

// NetEaseSMSRequest 短信请求参数
type NetEaseSMSRequest struct {
	SecretID     string                 `json:"secretId"`
	BusinessID   string                 `json:"businessId"`
	Version      string                 `json:"version"`
	Timestamp    string                 `json:"timestamp"`
	Nonce        string                 `json:"nonce"`
	Mobile       string                 `json:"mobile"`
	ParamType    string                 `json:"paramType"`
	Params       map[string]interface{} `json:"params"`
	TemplateID   string                 `json:"templateId"`
	CodeName     string                 `json:"codeName"`
	CodeLen      int                    `json:"codeLen"`
	CodeValidSec int                    `json:"codeValidSec"`
	Signature    string                 `json:"signature"`
}

// NetEaseSMSResponse 短信响应参数
type NetEaseSMSResponse struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
	Data    struct {
		RequestID string `json:"requestId"`
	} `json:"data"`
}

// genSignature 生成签名信息
func (n *NetEaseSMS) genSignature(params map[string]interface{}) string {
	var buf bytes.Buffer
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := fmt.Sprintf("%v", params[k])
		buf.WriteString(k)
		buf.WriteString(v)
	}

	buf.WriteString(secretKey)

	h := md5.New()
	h.Write(buf.Bytes())

	return hex.EncodeToString(h.Sum(nil))
}
