package main

import (
	"Open_IM/pkg/utils"
	"flag"
	"fmt"
	"log"
	"strconv"

	"Open_IM/pkg/common/constant"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.DebugMode)

	r := gin.Default()
	r.Use(utils.CorsHandler())

	demoRouterGroup := r.Group("/auth")
	{
		demoRouterGroup.POST("/code", func(c *gin.Context) {

		})
	}

	// 登录页
	{
		// 手机登录注册
		// 重置登录密码
		// 获取验证码
	}

	// tab1 磁场(首页)
	{
		// 附近的人
		// 我的动态
		// 发动态
		// 邀请分享
		// 聊天
		// 快捷用语
		// 语音聊天
		//
	}
	// tab2 圈子
	{
		//
	}
	// tab3 消息
	{
		// 已读 ; 未读
	}
	// tab4 我的
	{
		// 昵称 年龄 距离 会员级别 关注 粉丝 动态
		// 我的动态
		// 我的认证
		// 客服中心
		// 问题与建议
		// 账号中心
		// 我的粉丝
		// 我的关注
		// 用户中心
		{
			// 会员弹窗
			// 普通会员 砖石会员 优惠信息  解锁会员按钮
			// 密码设置
			// 锁屏密码
			// 编辑资料 (上传语音签名 上传头像)
			// 反馈
			// 其他设置 (阅后即焚 隐私模式 邀请码)
			// 我的认证 ( 手机认证 实名认证)
		}
	}

	ginPort := flag.Int("port", 20008, "get ginServerPort from cmd,default 10004 as port")
	flag.Parse()
	address := "0.0.0.0:" + strconv.Itoa(*ginPort)
	fmt.Println("start demo api server address: ", address, ", OpenIM version: ", constant.CurrentVersion, "\n")
	err := r.Run(address)
	if err != nil {
		log.Println("Error", "run failed ", *ginPort, err.Error())
		return
	}
}
