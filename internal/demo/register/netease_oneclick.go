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
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

type NetEaseOneClick struct {
	secretID   string
	secretKey  string
	businessID string
	apiURL     string
}

func NewNetEaseOneClick() *NetEaseOneClick {
	cfg := config.Config.Demo.NetEaseOneClick

	return &NetEaseOneClick{
		apiURL:     "https://ye.dun.163.com/v1/oneclick/check",
		secretID:   cfg.SecretID,
		secretKey:  cfg.SecretKey,
		businessID: cfg.BusinessID,
	}
}

// NetEaseOneClickResponse 响应参数
type NetEaseOneClickResponse struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
	Data    struct {
		Phone      string      `json:"phone"`
		ResultCode interface{} `json:"resultCode"`
	} `json:"data"`
}

// genSignature 生成签名信息
func (n *NetEaseOneClick) genSignature(params url.Values) string {
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

func (n *NetEaseOneClick) getPhone(token, accessToken string) (*NetEaseOneClickResponse, error) {
	p := url.Values{}

	p.Set("secretId", n.secretID)
	p.Set("businessId", n.businessID)
	p.Set("version", "v1")
	p.Set("timestamp", strconv.FormatInt(time.Now().UnixNano()/1000000, 10))
	p.Set("nonce", utils.GetRandomStr(32))
	p.Set("signature", n.genSignature(p))

	resp, err := http.Post(n.apiURL, "application/x-www-form-urlencoded", strings.NewReader(p.Encode()))

	if err != nil {
		return nil, fmt.Errorf("调用API接口失败: %s", err)
	}

	defer resp.Body.Close()

	// 5.解析报文返回
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.WithMessage(err, "netease.oneClick.body")
	}

	var clickResp NetEaseOneClickResponse
	err = json.Unmarshal(body, &clickResp)
	if err != nil {
		return nil, errors.WithMessage(err, "netease.oneClick.resBody")
	}
	b, _ := json.Marshal(clickResp)

	//6.返回结果
	log.Info("netease oneclick phone is ", clickResp.Data.Phone)

	log.Debug("netease oneclick response is ", string(b))

	if clickResp.Code != 200 {
		return nil, fmt.Errorf("异常返回码: %d %s", clickResp.Code, "")
	}

	return &clickResp, nil
}
